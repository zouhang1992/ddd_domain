# 账单开始时间默认值设计

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为新增账单表单添加开始时间默认值逻辑，提升用户体验

**Architecture:** 后端添加新的查询API，前端在选择租约后调用该API获取默认开始时间

**Tech Stack:** Go (后端), React + TypeScript + Ant Design (前端), SQLite (数据库)

---

## 1. 需求概述

### 功能需求
- 用户在新增账单时选择租约后，自动填充计费开始日期的默认值
- 如果该租约已有账单，开始时间 = 该租约所有账单中最大的 billEnd + 1天
- 如果该租约没有账单，开始时间 = 租约的 startDate
- 计费结束日期不设置默认值，由用户手动选择

### 非功能需求
- 响应速度快，用户选择租约后应立即显示默认值
- 数据准确，基于数据库中的真实账单数据计算
- 代码可维护，业务逻辑集中在后端

---

## 2. 后端设计

### 2.1 查询结构

**文件:** `internal/application/bill/query.go`

新增查询结构：

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

### 2.2 查询处理器

**文件:** `internal/application/bill/handler.go`

新增查询处理器：

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

注意：QueryHandler 需要注入 leaseRepo，需要更新 NewQueryHandler 构造函数。

### 2.3 HTTP 端点

**文件:** `internal/facade/cqrs_bill_handler.go`

新增 HTTP 路由和处理器：

```go
// RegisterRoutes 注册路由
func (h *CQRSBillHandler) RegisterRoutes(mux *http.ServeMux) {
	// ... 现有路由 ...
	mux.HandleFunc("GET /leases/{leaseId}/next-bill-period", h.GetNextBillPeriod)
}

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

### 2.4 依赖注入更新

需要更新 QueryHandler 的构造函数，注入 leaseRepo：

```go
// QueryHandler 账单查询处理器
type QueryHandler struct {
	billRepo  billrepo.BillRepository
	leaseRepo leaserepo.LeaseRepository
}

// NewQueryHandler 创建账单查询处理器
func NewQueryHandler(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository) *QueryHandler {
	return &QueryHandler{billRepo: billRepo, leaseRepo: leaseRepo}
}
```

同时需要更新依赖注入配置。

---

## 3. 前端设计

### 3.1 API 函数

**文件:** `web/src/api/bill.ts`

新增 API 函数：

```typescript
export const billApi = {
  // ... 现有 API ...

  getNextBillPeriod: async (leaseId: string) => {
    const response = await apiClient.get<{ billStart: string }>(
      `/leases/${leaseId}/next-bill-period`
    );
    return response.data;
  },
};
```

### 3.2 表单逻辑更新

**文件:** `web/src/pages/Bills.tsx`

修改点：

1. 在租约选择器添加 `onChange` 回调
2. 当用户选择租约时，调用 API 获取 billStart
3. 将返回的 billStart 设置到表单中
4. billEnd 保持为空，不设置默认值

代码示例：

```typescript
// 在租约选择器的 Form.Item 中
<Form.Item
  name="leaseId"
  label="租约"
  rules={[{ required: true, message: '请选择租约' }]}
>
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
    {/* ... 选项 ... */}
  </Select>
</Form.Item>
```

---

## 4. 数据流程

```
1. 用户打开新增账单表单
2. 用户选择租约
   ↓
3. 前端触发 Select onChange 事件
   ↓
4. 前端调用 GET /leases/{leaseId}/next-bill-period
   ↓
5. 后端接收请求
   ↓
6. 后端查询该租约的所有账单 (billRepo.FindByLeaseID)
   ↓
7. 如果有账单:
   - 找到最大的 billEnd
   - 计算 billStart = maxBillEnd + 1天

   如果没有账单:
   - 查询租约信息 (leaseRepo.FindByID)
   - billStart = 租约.startDate
   ↓
8. 后端返回 { "billStart": "2024-01-01" }
   ↓
9. 前端接收响应
   ↓
10. 前端将 billStart 设置到表单的 DatePicker 组件
    ↓
11. billEnd 保持为空，用户手动选择
```

---

## 5. 错误处理

### 后端错误
- 租约不存在 → 返回 404 "lease not found"
- 数据库查询失败 → 返回 500 内部错误

### 前端错误
- API 调用失败 → 显示错误提示 "获取默认计费开始时间失败"
- 用户仍可以手动选择开始时间

---

## 6. 文件清单

### 需要创建/修改的文件

**后端:**
- `internal/application/bill/query.go` - 添加查询结构
- `internal/application/bill/handler.go` - 添加查询处理器，更新构造函数
- `internal/facade/cqrs_bill_handler.go` - 添加HTTP端点和路由
- 依赖注入配置 - 更新 QueryHandler 的依赖注入

**前端:**
- `web/src/api/bill.ts` - 添加新的API函数
- `web/src/pages/Bills.tsx` - 添加租约选择器的onChange逻辑

---

## 7. 验收标准

- [ ] 用户选择租约后，计费开始日期自动填充
- [ ] 租约有账单时，开始时间 = 最大billEnd + 1天
- [ ] 租约无账单时，开始时间 = 租约startDate
- [ ] 计费结束日期不设置默认值，保持为空
- [ ] API调用失败时显示错误提示，但不阻塞用户操作
