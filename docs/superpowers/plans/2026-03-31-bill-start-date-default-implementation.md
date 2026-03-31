# 账单开始时间默认值实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现账单开始时间默认值功能，用户选择租约后自动填充开始时间

**Architecture:** 后端添加新的查询API，前端在选择租约后调用该API获取默认开始时间

**Tech Stack:** Go (后端), React + TypeScript + Ant Design (前端), SQLite (数据库)

---

## 任务概述

本计划包含以下任务：
1. 后端 - 添加查询结构
2. 后端 - 更新 QueryHandler 结构和构造函数
3. 后端 - 添加查询处理器
4. 后端 - 在命令总线注册查询处理器
5. 后端 - 添加 HTTP 端点
6. 前端 - 添加 API 函数
7. 前端 - 更新表单逻辑
8. 测试验证

---

### Task 1: 添加查询结构

**Files:**
- Modify: `internal/application/bill/query.go`

- [ ] **Step 1: 读取现有文件内容**

读取 `internal/application/bill/query.go` 的完整内容。

- [ ] **Step 2: 添加 GetNextBillPeriodQuery 结构**

在文件末尾添加：

```go
// GetNextBillPeriodQuery 获取租约下一个账单周期查询
type GetNextBillPeriodQuery struct {
	LeaseID string
}

// QueryName 实现 Query 接口
func (q GetNextBillPeriodQuery) QueryName() string {
	return "get_next_bill_period"
}

// NextBillPeriodQueryResult 下一个账单周期查询结果
type NextBillPeriodQueryResult struct {
	BillStart string `json:"bill_start"`
}
```

- [ ] **Step 3: 验证编译**

Run: `go build ./internal/application/bill`
Expected: 编译成功，无错误

- [ ] **Step 4: Commit**

```bash
git add internal/application/bill/query.go
git commit -m "feat: add GetNextBillPeriodQuery structure"
```

---

### Task 2: 更新 QueryHandler 结构和构造函数

**Files:**
- Modify: `internal/application/bill/handler.go`

- [ ] **Step 1: 读取现有文件内容**

读取 `internal/application/bill/handler.go` 的完整内容。

- [ ] **Step 2: 添加 leaseRepo 导入**

在导入部分添加：
```go
leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
```

- [ ] **Step 3: 更新 QueryHandler 结构**

修改 QueryHandler 结构，添加 leaseRepo 字段：

```go
// QueryHandler 账单查询处理器
type QueryHandler struct {
	billRepo  billrepo.BillRepository
	leaseRepo leaserepo.LeaseRepository
}
```

- [ ] **Step 4: 更新 NewQueryHandler 构造函数**

修改 NewQueryHandler 函数，接受并注入 leaseRepo：

```go
// NewQueryHandler 创建账单查询处理器
func NewQueryHandler(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository) *QueryHandler {
	return &QueryHandler{billRepo: billRepo, leaseRepo: leaseRepo}
}
```

- [ ] **Step 5: 验证编译**

Run: `go build ./internal/application/bill`
Expected: 编译成功，无错误

- [ ] **Step 6: Commit**

```bash
git add internal/application/bill/handler.go
git commit -m "feat: update QueryHandler to inject leaseRepo"
```

---

### Task 3: 添加查询处理器

**Files:**
- Modify: `internal/application/bill/handler.go`

- [ ] **Step 1: 添加时间导入（如需要）**

确保文件顶部有 `time` 导入。

- [ ] **Step 2: 添加 HandleGetNextBillPeriod 方法**

在 QueryHandler 的方法区域末尾添加：

```go
// HandleGetNextBillPeriod 处理获取下一个账单周期查询
func (h *QueryHandler) HandleGetNextBillPeriod(q common.Query) (any, error) {
	getQuery, ok := q.(GetNextBillPeriodQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	// 获取租约的所有账单
	bills, err := h.billRepo.FindByLeaseID(getQuery.LeaseID)
	if err != nil {
		return nil, err
	}

	// 如果有账单，找到最大的 billEnd + 1天
	if len(bills) > 0 {
		var maxBillEnd time.Time
		for _, bill := range bills {
			if bill.BillEnd.After(maxBillEnd) {
				maxBillEnd = bill.BillEnd
			}
		}
		nextBillStart := maxBillEnd.AddDate(0, 0, 1)
		return &NextBillPeriodQueryResult{
			BillStart: nextBillStart.Format("2006-01-02"),
		}, nil
	}

	// 如果没有账单，获取租约信息，返回租约的 startDate
	lease, err := h.leaseRepo.FindByID(getQuery.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, domerrors.ErrNotFound
	}

	return &NextBillPeriodQueryResult{
		BillStart: lease.StartDate.Format("2006-01-02"),
	}, nil
}
```

- [ ] **Step 3: 验证编译**

Run: `go build ./internal/application/bill`
Expected: 编译成功，无错误

- [ ] **Step 4: Commit**

```bash
git add internal/application/bill/handler.go
git commit -m "feat: add HandleGetNextBillPeriod query handler"
```

---

### Task 4: 在查询总线注册查询处理器

**Files:**
- Modify: `internal/application/module.go` (或依赖注入配置文件)

- [ ] **Step 1: 查找依赖注入配置文件**

先查找依赖注入配置的位置：
```bash
ls internal/application/*.go
```

找到设置查询处理器的文件，通常是 `module.go` 或类似名称。

- [ ] **Step 2: 更新 QueryHandler 的构造函数调用**

在依赖注入配置中，更新 `NewQueryHandler` 的调用，添加 leaseRepo 参数。

参考模式：
```go
fx.Provide(func(
	billRepo billrepo.BillRepository,
	leaseRepo leaserepo.LeaseRepository,
) *bill.QueryHandler {
	return bill.NewQueryHandler(billRepo, leaseRepo)
}),
```

- [ ] **Step 3: 注册查询处理器到查询总线**

在查询总线注册部分添加：
```go
bus.Register("get_next_bill_period", busquery.HandlerFunc(billQueryHandler.HandleGetNextBillPeriod))
```

- [ ] **Step 4: 验证编译**

Run: `go build ./internal/application`
Expected: 编译成功，无错误

- [ ] **Step 5: Commit**

```bash
git add internal/application/module.go
git commit -m "feat: register GetNextBillPeriod query handler"
```

---

### Task 5: 添加 HTTP 端点

**Files:**
- Modify: `internal/facade/cqrs_bill_handler.go`

- [ ] **Step 1: 读取现有文件内容**

读取 `internal/facade/cqrs_bill_handler.go` 的完整内容。

- [ ] **Step 2: 在 RegisterRoutes 中添加新路由**

在 `RegisterRoutes` 方法中添加：
```go
mux.HandleFunc("GET /leases/{leaseId}/next-bill-period", h.GetNextBillPeriod)
```

- [ ] **Step 3: 添加 GetNextBillPeriod HTTP 处理器方法**

在文件末尾添加：

```go
// GetNextBillPeriod 获取租约的下一个账单周期
func (h *CQRSBillHandler) GetNextBillPeriod(w http.ResponseWriter, r *http.Request) {
	leaseId := r.PathValue("leaseId")
	q := bill.GetNextBillPeriodQuery{LeaseID: leaseId}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		if err.Error() == "lease not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
```

- [ ] **Step 4: 验证编译**

Run: `go build ./internal/facade`
Expected: 编译成功，无错误

- [ ] **Step 5: Commit**

```bash
git add internal/facade/cqrs_bill_handler.go
git commit -m "feat: add HTTP endpoint for next bill period"
```

---

### Task 6: 前端添加 API 函数

**Files:**
- Modify: `web/src/api/bill.ts`

- [ ] **Step 1: 读取现有文件内容**

读取 `web/src/api/bill.ts` 的完整内容。

- [ ] **Step 2: 添加 getNextBillPeriod 函数**

在 billApi 对象中添加：

```typescript
getNextBillPeriod: async (leaseId: string) => {
  const response = await apiClient.get<{ billStart: string }>(
    `/leases/${leaseId}/next-bill-period`
  );
  return response.data;
},
```

- [ ] **Step 3: 验证类型检查**

Run: `cd web && npm run build` (或类似的类型检查命令)
Expected: 无类型错误

- [ ] **Step 4: Commit**

```bash
git add web/src/api/bill.ts
git commit -m "feat: add getNextBillPeriod API function"
```

---

### Task 7: 前端更新表单逻辑

**Files:**
- Modify: `web/src/pages/Bills.tsx`

- [ ] **Step 1: 读取现有文件内容**

读取 `web/src/pages/Bills.tsx`，重点关注租约选择器部分（约第 595-622 行）。

- [ ] **Step 2: 添加 dayjs 导入（如需要）**

确认文件顶部已有 `dayjs` 导入。

- [ ] **Step 3: 在租约选择器添加 onChange 回调**

修改租约选择器的 Select 组件，添加 onChange 回调：

```typescript
<Select
  placeholder="请选择租约"
  showSearch
  optionFilterProp="children"
  onChange={async (leaseId) => {
    if (leaseId) {
      try {
        const data = await billApi.getNextBillPeriod(leaseId);
        form.setFieldValue('billStart', dayjs(data.billStart));
        // 清空 billEnd，让用户自己选择
        form.setFieldValue('billEnd', null);
      } catch {
        message.error('获取默认计费开始时间失败');
      }
    }
  }}
>
  {/* ... 现有选项 ... */}
</Select>
```

- [ ] **Step 4: 验证类型检查**

Run: `cd web && npm run build`
Expected: 无类型错误

- [ ] **Step 5: Commit**

```bash
git add web/src/pages/Bills.tsx
git commit -m "feat: auto-fill billStart when selecting lease"
```

---

### Task 8: 测试验证

**Files:** 手动测试

- [ ] **Step 1: 启动后端服务**

Run: `go run cmd/api/main.go`
Expected: 服务成功启动在 8080 端口

- [ ] **Step 2: 启动前端服务**

Run: `cd web && npm start`
Expected: 前端成功启动

- [ ] **Step 3: 测试无账单的租约**

1. 选择一个没有账单的租约
2. 验证 billStart 自动填充为租约的 startDate
3. 验证 billEnd 保持为空

- [ ] **Step 4: 测试有账单的租约**

1. 选择一个有账单的租约
2. 验证 billStart 自动填充为最大 billEnd + 1天
3. 验证 billEnd 保持为空

- [ ] **Step 5: 测试错误处理**

1. 模拟 API 失败（如断开后端连接）
2. 选择租约
3. 验证显示错误提示 "获取默认计费开始时间失败"
4. 验证用户仍可手动选择开始时间

---

## 验收标准核对清单

- [ ] 用户选择租约后，计费开始日期自动填充
- [ ] 租约有账单时，开始时间 = 最大billEnd + 1天
- [ ] 租约无账单时，开始时间 = 租约startDate
- [ ] 计费结束日期不设置默认值，保持为空
- [ ] API调用失败时显示错误提示，但不阻塞用户操作

---

## 文件汇总

**修改的文件：**
1. `internal/application/bill/query.go` - 添加查询结构
2. `internal/application/bill/handler.go` - 添加查询处理器，更新构造函数
3. `internal/application/module.go` - 更新依赖注入和查询总线注册
4. `internal/facade/cqrs_bill_handler.go` - 添加HTTP端点
5. `web/src/api/bill.ts` - 添加API函数
6. `web/src/pages/Bills.tsx` - 添加租约选择器的onChange逻辑
