# 账单管理、退租结算和收入汇总调整设计文档

**日期**: 2026-03-30
**项目**: ddd_domain

## 概述

本文档描述了三个主要功能调整：
1. 账单页面调整 - 根据账单类型动态显示表单字段
2. 退租结算表单优化 - 调整分组布局和默认值
3. 收入汇总逻辑调整 - 正确处理退租结算中的租金退还

## 1. 账单页面调整

### 1.1 功能描述

在账单管理页面，当用户点击"新增账单"时，根据选择的账单类型动态显示不同的表单字段：

- **收租类类型**（rent, water, electric, gas, internet, other, charge）：显示普通账单表单
- **退租结算类型**（checkout）：可以选择隐藏或显示退租结算专用表单

### 1.2 设计细节

#### 1.2.1 类型选择

保留现有的账单类型下拉框：
- rent: 租金
- water: 水费
- electric: 电费
- gas: 燃气费
- internet: 网费
- other: 其他
- charge: 收账
- checkout: 退租结算

#### 1.2.2 动态字段显示

**收租类表单字段：**
- 租约选择（必填）
- 类型选择（已选择）
- 金额区域：
  - 租金金额（分）
  - 水费金额（分）
  - 电费金额（分）
  - 其他金额（分）
- 到账时间（可选）
- 备注（可选）

**退租结算表单字段（可选）：**
考虑到退租结算是通过租约列表触发的，可以选择：
- 选项1：在账单页面隐藏"退租结算"类型选项
- 选项2：显示但要求选择租约，并且只读模式显示退租信息

### 1.3 涉及文件修改

- `web/src/pages/Bills.tsx`

---

## 2. 退租结算表单优化

### 2.1 功能描述

在租约列表页面点击"退租"按钮时，弹出优化后的退租结算表单，包括：
- 分组显示字段（退还金额、收取费用）
- 合理的默认值设置
- 完善的验证规则

### 2.2 设计细节

#### 2.2.1 表单布局

```
┌─────────────────────────────────────────┐
│  租约信息（灰色背景区域）              │
│  - 租户：{tenantName}                    │
│  - 押金金额：¥{depositAmount}          │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│  退还金额（标题 + 边框区域）           │
│  [ ] 退还租金（分）：[0        ]       │
│  [ ] 退还押金（分）：[全额押金  ]       │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│  收取费用（标题 + 边框区域）           │
│  [ ] 水费（分）：[0        ]           │
│  [ ] 电费（分）：[0        ]           │
│  [ ] 其他费用（分）：[0        ]       │
└─────────────────────────────────────────┘

[备注输入框]
```

#### 2.2.2 默认值设置

| 字段 | 默认值 | 说明 |
|------|--------|------|
| 退还租金 | 0 | 默认不退还租金 |
| 退还押金 | 租约押金金额 | 默认全部退还押金 |
| 水费 | 0 | 默认不收水费 |
| 电费 | 0 | 默认不收电费 |
| 其他费用 | 0 | 默认不收其他费用 |
| 备注 | 空 | 可选填写 |

#### 2.2.3 验证规则

- 退还租金 ≥ 0
- 0 ≤ 退还押金 ≤ 租约押金金额
- 水费 ≥ 0
- 电费 ≥ 0
- 其他费用 ≥ 0

### 2.3 涉及文件修改

- `web/src/pages/Leases.tsx`

---

## 3. 收入汇总逻辑调整

### 3.1 功能描述

调整收入汇总的计算逻辑，确保：
- 正确处理退租结算账单中的租金退还（作为支出）
- 区分收入和支出
- 显示总收入、总支出、净收入

### 3.2 设计细节

#### 3.2.1 后端数据结构调整

```go
// IncomeReport 收入报告
type IncomeReport struct {
    Year                  int    `json:"year"`
    Month                 int    `json:"month"`

    // 收入部分
    RentIncome            int64  `json:"rent_income"`
    WaterIncome           int64  `json:"water_income"`
    ElectricIncome        int64  `json:"electric_income"`
    OtherIncome           int64  `json:"other_income"`
    DepositIncome         int64  `json:"deposit_income"`

    // 支出部分
    RentExpense           int64  `json:"rent_expense"`
    DepositExpense        int64  `json:"deposit_expense"`

    // 计算结果
    TotalIncome           int64  `json:"total_income"`
    TotalExpense          int64  `json:"total_expense"`
    NetIncome             int64  `json:"net_income"`

    // 格式化字符串
    RentIncomeFormatted     string `json:"rent_income_formatted"`
    WaterIncomeFormatted    string `json:"water_income_formatted"`
    ElectricIncomeFormatted string `json:"electric_income_formatted"`
    OtherIncomeFormatted    string `json:"other_income_formatted"`
    DepositIncomeFormatted  string `json:"deposit_income_formatted"`
    RentExpenseFormatted    string `json:"rent_expense_formatted"`
    DepositExpenseFormatted string `json:"deposit_expense_formatted"`
    TotalIncomeFormatted    string `json:"total_income_formatted"`
    TotalExpenseFormatted   string `json:"total_expense_formatted"`
    NetIncomeFormatted      string `json:"net_income_formatted"`
}
```

#### 3.2.2 计算逻辑

**步骤1：遍历所有账单**

对于每个已支付（PaidAt != nil）且在指定月份的账单：

- **如果是 checkout 类型（退租结算）**：
  - `rent_amount < 0`：`RentExpense += |rent_amount|`（租金支出）
  - `rent_amount > 0`：`RentIncome += rent_amount`（租金收入）
  - `water_amount`：`WaterIncome += water_amount`
  - `electric_amount`：`ElectricIncome += electric_amount`
  - `other_amount`：`OtherIncome += other_amount`

- **如果是其他类型**：
  - `rent_amount`：`RentIncome += rent_amount`
  - `water_amount`：`WaterIncome += water_amount`
  - `electric_amount`：`ElectricIncome += electric_amount`
  - `other_amount`：`OtherIncome += other_amount`

**步骤2：遍历所有押金**

对于每个押金：

- `created_at` 在指定月份：`DepositIncome += amount`
- `refunded_at` 在指定月份：`DepositExpense += amount`

**步骤3：计算总计**

```
TotalIncome = RentIncome + WaterIncome + ElectricIncome + OtherIncome + DepositIncome
TotalExpense = RentExpense + DepositExpense
NetIncome = TotalIncome - TotalExpense
```

#### 3.2.3 前端显示调整

**顶部统计卡片：**

1. **总收入**（绿色）
   - 显示：TotalIncome / 100 元

2. **总支出**（红色）
   - 显示：TotalExpense / 100 元

3. **净收入**（蓝色/红色）
   - 显示：NetIncome / 100 元
   - 正数用蓝色，负数用红色

**明细表调整：**

- 收入项：正常颜色显示
- 支出项：红色显示

### 3.3 涉及文件修改

- `internal/facade/income_handler.go`
- `web/src/pages/Income.tsx`
- `web/src/api/income.ts`（如需更新类型定义）

---

## 4. 数据流程

### 4.1 退租结算流程

```
用户点击"退租"按钮
    ↓
弹出退租结算表单（已优化布局）
    ↓
用户填写退还金额和收取费用
    ↓
提交表单
    ↓
后端处理：
  - 租约状态改为 checkout
  - 创建 checkout 类型账单
  - 押金状态改为 returned
    ↓
前端刷新列表
```

### 4.2 收入汇总流程

```
用户选择月份（或默认当前月）
    ↓
前端请求 /income?month=YYYY-MM
    ↓
后端查询：
  - 所有已支付账单
  - 所有押金记录
    ↓
后端计算：
  - 按规则分类收入和支出
  - 计算总收入、总支出、净收入
    ↓
返回 IncomeReport
    ↓
前端显示统计卡片和明细表
```

---

## 5. 注意事项

### 5.1 账单金额符号规则

- **正数**：表示向租户收取（收入）
- **负数**：表示退还给租户（支出）

### 5.2 退租结算账单

退租结算账单（type = "checkout"）的金额计算公式：
```
amount = -refundRentAmount - refundDepositAmount + waterAmount + electricAmount + otherAmount
```

### 5.3 时间范围判断

收入汇总的时间判断：
- 账单：使用 PaidAt 字段
- 押金收入：使用 CreatedAt 字段
- 押金支出：使用 RefundedAt 字段

---

## 6. 实现优先级

1. **高优先级**：
   - 退租结算表单优化（布局、默认值、验证）
   - 收入汇总逻辑调整（正确处理租金支出）

2. **中优先级**：
   - 账单页面动态表单字段

---

## 7. 验收标准

- [ ] 退租结算表单分组显示退还金额和收取费用
- [ ] 退租结算表单默认押金为全额退还
- [ ] 收入汇总正确显示总收入、总支出、净收入
- [ ] 退租结算中的租金退还正确计入支出
- [ ] 所有验证规则正常工作
