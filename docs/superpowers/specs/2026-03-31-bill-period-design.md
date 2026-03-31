---
name: Bill Period Design
description: Add billing period (bill_start, bill_end) to bills to track which rental period the bill corresponds to
type: spec
---

# 账单计费周期功能设计

## 概述

为账单添加计费周期信息（bill_start 和 bill_end），用于标识该账单对应的租金时间段。租约可能是一年，但账单可以按月、按季度或任意时间段创建。同时保留 due_date（付款截止日期）用于催收功能。

## 业务规则

### 1. 账单周期自动计算

**账单开始时间的自动计算逻辑：**

1. 查找该租约下的最新账单（按 bill_end 降序）
2. 如果存在已支付的账单：
   - `bill_start` = 最新已支付账单的 `bill_end` + 1天
3. 如果没有已支付账单：
   - `bill_start` = 租约的 `start_date`

**账单结束时间的默认计算：**
- 根据选择的周期类型（月/季）自动计算：
  - 按月：`bill_end` = `bill_start` + 1个月 - 1天
  - 按季：`bill_end` = `bill_start` + 3个月 - 1天

**付款截止日期的默认计算：**
- `due_date` = `bill_end` + 7天（可手动修改）

### 2. 快捷选项

| 选项 | bill_start | bill_end |
|------|------------|----------|
| 本月 | 本月1号 | 本月最后一天 |
| 本季度 | 本季度1号 | 本季度最后一天 |
| 上月 | 上月1号 | 上月最后一天 |
| 上季度 | 上季度1号 | 上季度最后一天 |
| 使用租约周期 | 按业务规则自动计算 | 按选择的月/季自动计算 |

### 3. 数据校验

- `bill_start` 必须 <= `bill_end`
- `bill_start` 和 `bill_end` 不能为空
- 账单周期超出租约范围时：提示警告，但允许创建（特殊情况）

## 数据模型设计

### 领域模型变更

**文件：** `internal/domain/bill/model/bill.go`

```go
type Bill struct {
    model.BaseAggregateRoot
    LeaseID           string     `json:"lease_id"`
    Type              BillType   `json:"type"`
    Status            BillStatus `json:"status"`
    Amount            int64      `json:"amount"`
    RentAmount        int64      `json:"rent_amount"`
    WaterAmount       int64      `json:"water_amount"`
    ElectricAmount    int64      `json:"electric_amount"`
    OtherAmount       int64      `json:"other_amount"`
    RefundDepositAmount int64    `json:"refund_deposit_amount"`
    BillStart         time.Time  `json:"bill_start"`  // 新增：计费周期开始日期
    BillEnd           time.Time  `json:"bill_end"`    // 新增：计费周期结束日期
    DueDate           time.Time  `json:"due_date"`    // 保持不变：付款截止日期
    PaidAt            *time.Time `json:"paid_at"`
    Note              string     `json:"note"`
    CreatedAt         time.Time  `json:"created_at"`
    UpdatedAt         time.Time  `json:"updated_at"`
}
```

### 数据库表变更

**bills 表新增字段：**
- `bill_start DATETIME` - 计费周期开始日期
- `bill_end DATETIME` - 计费周期结束日期

**保留字段：**
- `due_date DATETIME` - 付款截止日期（用于催收功能）

## 后端实现

### 1. Command 层

**文件：** `internal/application/bill/command.go`

```go
type CreateBillCommand struct {
    LeaseID             string
    Type                billmodel.BillType
    Amount              int64
    RentAmount          int64
    WaterAmount         int64
    ElectricAmount      int64
    OtherAmount         int64
    RefundDepositAmount int64
    BillStart           time.Time  // 新增
    BillEnd             time.Time  // 新增
    DueDate             time.Time
    Note                string
}

type UpdateBillCommand struct {
    ID                  string
    Amount              int64
    RentAmount          int64
    WaterAmount         int64
    ElectricAmount      int64
    OtherAmount         int64
    RefundDepositAmount int64
    BillStart           time.Time  // 新增
    BillEnd             time.Time  // 新增
    DueDate             time.Time
    Note                string
}
```

### 2. Handler 层

**文件：** `internal/application/bill/handler.go`

更新创建和更新账单的处理逻辑，传入新增字段。

### 3. Repository 层

**文件：** `internal/infrastructure/persistence/sqlite/bill_repo.go`

- 更新 Save 方法，保存新字段
- 更新 FindByID、FindAll、FindByLeaseID、FindUnpaidBillsDueBefore 方法，读取新字段

### 4. Facade 层

**文件：** `internal/facade/cqrs_bill_handler.go`

- Create 接口：接收 bill_start、bill_end
- Update 接口：接收 bill_start、bill_end

### 5. 数据库迁移

**文件：** `internal/infrastructure/persistence/sqlite/migration.go`

新增迁移脚本：
```go
type AddBillPeriodToBillsMigration struct{}

func (m *AddBillPeriodToBillsMigration) Version() string {
    return "20260331XXXX" // 格式：YYYYMMDDHHMM
}

func (m *AddBillPeriodToBillsMigration) Up(tx *sql.Tx) error {
    // 检查 bill_start 列是否已存在
    // 检查 bill_end 列是否已存在
    // 添加缺失的列
}
```

将此迁移添加到 `migrations` 切片中。

## 前端实现

### 1. 类型定义

**文件：** `web/src/types/api.ts`

```typescript
export interface Bill {
  id: string;
  leaseId: string;
  type: string;
  status: string;
  amount: number;
  rentAmount: number;
  waterAmount: number;
  electricAmount: number;
  otherAmount: number;
  refundDepositAmount: number;
  billStart: string;  // 新增
  billEnd: string;    // 新增
  dueDate: string;
  paidAt: string | null;
  note: string;
  createdAt: string;
  updatedAt: string;
}
```

### 2. API 层

**文件：** `web/src/api/bill.ts`

更新 create 和 update 接口参数，添加 billStart、billEnd。

### 3. 表单界面

**文件：** `web/src/pages/Bills.tsx`

**新增账单表单变更：**

1. 快捷选项按钮区域：
   - 本月
   - 本季度
   - 上月
   - 上季度
   - 使用租约周期

2. 日期选择器：
   - 开始日期（billStart）
   - 结束日期（billEnd）

3. 选择租约后的自动逻辑：
   - 查询该租约的所有账单
   - 按业务规则计算默认的 billStart 和 billEnd
   - 填充到表单中，用户仍可手动修改

### 4. 列表展示

**文件：** `web/src/pages/Bills.tsx`

- 新增"计费周期"列，显示格式：`2024-01-01 ~ 2024-03-31`
- 保留"租期"列（显示租约的租期）

## 变更文件清单

### 后端文件

1. `internal/domain/bill/model/bill.go` - 领域模型新增字段
2. `internal/application/bill/command.go` - Command 新增字段
3. `internal/application/bill/handler.go` - Handler 更新逻辑
4. `internal/infrastructure/persistence/sqlite/bill_repo.go` - Repository 更新
5. `internal/infrastructure/persistence/sqlite/migration.go` - 新增迁移
6. `internal/facade/cqrs_bill_handler.go` - HTTP Handler 更新

### 前端文件

1. `web/src/types/api.ts` - TypeScript 类型更新
2. `web/src/api/bill.ts` - API 接口更新
3. `web/src/pages/Bills.tsx` - 页面表单和列表更新

## 数据流转

```
用户选择租约
    ↓
查询该租约的所有账单
    ↓
计算默认 bill_start：
  - 有已支付账单？bill_end + 1天
  - 无已支付账单？租约 start_date
    ↓
用户选择快捷选项或手动输入日期
    ↓
计算 bill_end（月/季）
    ↓
due_date = bill_end + 7天（可修改）
    ↓
提交到后端 → 验证 → 保存数据库
    ↓
列表展示：显示计费周期列 + 租期列
```

## 错误处理

1. **账单周期超出租约范围**：
   - 提示警告，但允许创建（因为可能存在特殊情况）

2. **日期校验**：
   - bill_start 必须 <= bill_end
   - bill_start 和 bill_end 不能为空
