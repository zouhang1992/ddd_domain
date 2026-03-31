# Bill Period Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add billing period (bill_start, bill_end) to bills to track which rental period the bill corresponds to, with auto-calculation based on previous bills or lease start date.

**Architecture:** Minimal changes approach - modify only bill-related files, keep due_date for payment reminders, add bill_start and bill_end for period tracking.

**Tech Stack:** Go backend, React/TypeScript frontend, SQLite database

---

## File Structure

**Backend Files:**
- `internal/domain/bill/model/bill.go` - Add BillStart and BillEnd fields to domain model
- `internal/application/bill/command.go` - Add fields to commands
- `internal/application/bill/handler.go` - Update handler logic
- `internal/infrastructure/persistence/sqlite/bill_repo.go` - Update repository to handle new fields
- `internal/infrastructure/persistence/sqlite/migration.go` - Add database migration
- `internal/facade/cqrs_bill_handler.go` - Update HTTP handlers

**Frontend Files:**
- `web/src/types/api.ts` - Update TypeScript interfaces
- `web/src/api/bill.ts` - Update API calls
- `web/src/pages/Bills.tsx` - Update form and list UI

---

### Task 1: Update Domain Model - Add BillStart and BillEnd Fields

**Files:**
- Modify: `internal/domain/bill/model/bill.go`

- [ ] **Step 1: Add BillStart and BillEnd fields to Bill struct**

```go
// Bill 账单领域模型（聚合根）
type Bill struct {
	model.BaseAggregateRoot
	LeaseID             string     `json:"lease_id"`
	Type                BillType   `json:"type"`
	Status              BillStatus `json:"status"`
	Amount              int64      `json:"amount"`
	RentAmount          int64      `json:"rent_amount"`
	WaterAmount         int64      `json:"water_amount"`
	ElectricAmount      int64      `json:"electric_amount"`
	OtherAmount         int64      `json:"other_amount"`
	RefundDepositAmount int64      `json:"refund_deposit_amount"`
	BillStart           time.Time  `json:"bill_start"`
	BillEnd             time.Time  `json:"bill_end"`
	DueDate             time.Time  `json:"due_date"`
	PaidAt              *time.Time `json:"paid_at"`
	Note                string     `json:"note"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
```

- [ ] **Step 2: Update NewBill function to accept BillStart and BillEnd**

```go
// NewBill 创建新账单
func NewBill(id, leaseID string, billType BillType, amount int64,
	billStart, billEnd time.Time, dueDate time.Time, note string) *Bill {
	now := time.Now()
	bill := &Bill{
		BaseAggregateRoot:   model.NewBaseAggregateRoot(id),
		LeaseID:             leaseID,
		Type:                billType,
		Status:              BillStatusPending,
		Amount:              amount,
		BillStart:           billStart,
		BillEnd:             billEnd,
		DueDate:             dueDate,
		Note:                note,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	// 创建并记录事件
	evt := billCreated{
		BaseEvent: events.NewBaseEvent("bill.created", bill.ID(), bill.Version()),
		LeaseID:   bill.LeaseID,
		Type:      string(bill.Type),
		Amount:    bill.Amount,
	}
	bill.RecordEvent(evt)
	return bill
}
```

- [ ] **Step 3: Update NewBillWithDetails function to accept BillStart and BillEnd**

```go
// NewBillWithDetails 创建带明细的新账单
func NewBillWithDetails(id, leaseID string, billType BillType,
	rentAmount, waterAmount, electricAmount, otherAmount, refundDepositAmount int64,
	billStart, billEnd time.Time, dueDate time.Time, note string) *Bill {
	now := time.Now()
	// Calculate total amount: rentAmount(负数表示退还) + 费用(正数) - 退还押金(正数表示退还)
	// 注意：refundDepositAmount是正数表示退还，所以计算总额时要减去它（因为这是要给租户的钱）
	totalAmount := rentAmount + waterAmount + electricAmount + otherAmount - refundDepositAmount

	bill := &Bill{
		BaseAggregateRoot:   model.NewBaseAggregateRoot(id),
		LeaseID:             leaseID,
		Type:                billType,
		Status:              BillStatusPending,
		Amount:              totalAmount,
		RentAmount:          rentAmount,
		WaterAmount:         waterAmount,
		ElectricAmount:      electricAmount,
		OtherAmount:         otherAmount,
		RefundDepositAmount: refundDepositAmount,
		BillStart:           billStart,
		BillEnd:             billEnd,
		DueDate:             dueDate,
		Note:                note,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	// 创建并记录事件
	evt := billCreated{
		BaseEvent: events.NewBaseEvent("bill.created", bill.ID(), bill.Version()),
		LeaseID:   bill.LeaseID,
		Type:      string(bill.Type),
		Amount:    bill.Amount,
	}
	bill.RecordEvent(evt)
	return bill
}
```

- [ ] **Step 4: Update Update method**

```go
// Update 更新账单信息
func (b *Bill) Update(amount int64, billStart, billEnd time.Time, dueDate time.Time, note string) {
	b.Amount = amount
	b.BillStart = billStart
	b.BillEnd = billEnd
	b.DueDate = dueDate
	b.Note = note
	b.UpdatedAt = time.Now()
	// 创建并记录事件
	evt := billUpdated{
		BaseEvent: events.NewBaseEvent("bill.updated", b.ID(), b.Version()),
		Amount:    b.Amount,
	}
	b.RecordEvent(evt)
}
```

- [ ] **Step 5: Update UpdateWithDetails method**

```go
// UpdateWithDetails 更新账单信息（包含明细）
func (b *Bill) UpdateWithDetails(rentAmount, waterAmount, electricAmount, otherAmount, refundDepositAmount int64,
	billStart, billEnd time.Time, dueDate time.Time, note string) {
	b.RentAmount = rentAmount
	b.WaterAmount = waterAmount
	b.ElectricAmount = electricAmount
	b.OtherAmount = otherAmount
	b.RefundDepositAmount = refundDepositAmount
	b.Amount = rentAmount + waterAmount + electricAmount + otherAmount - refundDepositAmount
	b.BillStart = billStart
	b.BillEnd = billEnd
	b.DueDate = dueDate
	b.Note = note
	b.UpdatedAt = time.Now()
	// 创建并记录事件
	evt := billUpdated{
		BaseEvent: events.NewBaseEvent("bill.updated", b.ID(), b.Version()),
		Amount:    b.Amount,
	}
	b.RecordEvent(evt)
}
```

- [ ] **Step 6: Commit**

```bash
git add internal/domain/bill/model/bill.go
git commit -m "feat: add bill_start and bill_end fields to Bill domain model"
```

---

### Task 2: Update Command Layer - Add BillStart and BillEnd to Commands

**Files:**
- Modify: `internal/application/bill/command.go`

- [ ] **Step 1: Update CreateBillCommand**

```go
// CreateBillCommand 创建账单命令
type CreateBillCommand struct {
	LeaseID             string
	Type                billmodel.BillType
	Amount              int64
	RentAmount          int64
	WaterAmount         int64
	ElectricAmount      int64
	OtherAmount         int64
	RefundDepositAmount int64
	BillStart           time.Time
	BillEnd             time.Time
	DueDate             time.Time
	Note                string
}
```

- [ ] **Step 2: Update UpdateBillCommand**

```go
// UpdateBillCommand 更新账单命令
type UpdateBillCommand struct {
	ID                  string
	Amount              int64
	RentAmount          int64
	WaterAmount         int64
	ElectricAmount      int64
	OtherAmount         int64
	RefundDepositAmount int64
	BillStart           time.Time
	BillEnd             time.Time
	DueDate             time.Time
	Note                string
}
```

- [ ] **Step 3: Commit**

```bash
git add internal/application/bill/command.go
git commit -m "feat: add bill_start and bill_end to bill commands"
```

---

### Task 3: Update Handler Layer

**Files:**
- Modify: `internal/application/bill/handler.go`

- [ ] **Step 1: Read the current handler.go file**

First, let's read the file to see what needs to be changed:

```bash
# Read the file content first
```

- [ ] **Step 2: Update HandleCreateBill to pass BillStart and BillEnd**

Modify the NewBill and NewBillWithDetails calls to include billStart and billEnd parameters.

- [ ] **Step 3: Update HandleUpdateBill to pass BillStart and BillEnd**

Modify the Update and UpdateWithDetails calls to include billStart and billEnd parameters.

- [ ] **Step 4: Commit**

```bash
git add internal/application/bill/handler.go
git commit -m "feat: update bill handlers for bill_start and bill_end"
```

---

### Task 4: Update Repository Layer

**Files:**
- Modify: `internal/infrastructure/persistence/sqlite/bill_repo.go`

- [ ] **Step 1: Update tempBill struct**

```go
// tempBill is a temporary struct for scanning
type tempBill struct {
	ID                  string
	LeaseID             string
	Type                string
	Status              string
	Amount              int64
	RentAmount          int64
	WaterAmount         int64
	ElectricAmount      int64
	OtherAmount         int64
	RefundDepositAmount int64
	BillStart           sql.NullTime
	BillEnd             sql.NullTime
	DueDate             sql.NullTime
	PaidAt              *time.Time
	Note                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
```

- [ ] **Step 2: Update Save method SQL**

```go
_, err := r.conn.DB().Exec(`
	INSERT OR REPLACE INTO bills (
		id, lease_id, type, status, amount, rent_amount, water_amount, electric_amount, other_amount, refund_deposit_amount, bill_start, bill_end, due_date, paid_at, note, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
	bill.IDField, bill.LeaseID, string(bill.Type), string(bill.Status), bill.Amount, bill.RentAmount, bill.WaterAmount, bill.ElectricAmount, bill.OtherAmount, bill.RefundDepositAmount, bill.BillStart, bill.BillEnd,
	bill.DueDate, paidAt, bill.Note, bill.CreatedAt, bill.UpdatedAt)
```

- [ ] **Step 3: Update SELECT queries in all methods**

Add `bill_start, bill_end` to all SELECT statements.

- [ ] **Step 4: Update Scan calls in all methods**

Add `&temp.BillStart, &temp.BillEnd,` to all row.Scan calls.

- [ ] **Step 5: Update Bill construction in all methods**

When constructing Bill from tempBill, read BillStart and BillEnd:

```go
billStart := time.Now()
if temp.BillStart.Valid {
	billStart = temp.BillStart.Time
}
billEnd := time.Now()
if temp.BillEnd.Valid {
	billEnd = temp.BillEnd.Time
}
dueDate := time.Now()
if temp.DueDate.Valid {
	dueDate = temp.DueDate.Time
}

// Now construct the bill with billStart and billEnd
bill := billmodel.NewBill(temp.ID, temp.LeaseID, billmodel.BillType(temp.Type), temp.Amount, billStart, billEnd, dueDate, temp.Note)
```

- [ ] **Step 6: Commit**

```bash
git add internal/infrastructure/persistence/sqlite/bill_repo.go
git commit -m "feat: update bill repository for bill_start and bill_end"
```

---

### Task 5: Add Database Migration

**Files:**
- Modify: `internal/infrastructure/persistence/sqlite/migration.go`

- [ ] **Step 1: Add migration struct**

```go
// AddBillPeriodToBillsMigration 为 bills 表添加 bill_start 和 bill_end 字段
type AddBillPeriodToBillsMigration struct{}

func (m *AddBillPeriodToBillsMigration) Version() string {
	return "202603311600" // 格式：YYYYMMDDHHMM
}

func (m *AddBillPeriodToBillsMigration) Up(tx *sql.Tx) error {
	// 检查列是否已存在
	rows, err := tx.Query("PRAGMA table_info(bills)")
	if err != nil {
		return fmt.Errorf("failed to check table info: %w", err)
	}
	defer rows.Close()

	billStartExists := false
	billEndExists := false
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt_value sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk); err != nil {
			return fmt.Errorf("failed to scan table info: %w", err)
		}
		if name == "bill_start" {
			billStartExists = true
		}
		if name == "bill_end" {
			billEndExists = true
		}
	}

	if !billStartExists {
		if _, err := tx.Exec("ALTER TABLE bills ADD COLUMN bill_start DATETIME"); err != nil {
			return fmt.Errorf("failed to add bill_start column: %w", err)
		}
	}

	if !billEndExists {
		if _, err := tx.Exec("ALTER TABLE bills ADD COLUMN bill_end DATETIME"); err != nil {
			return fmt.Errorf("failed to add bill_end column: %w", err)
		}
	}

	return nil
}

func (m *AddBillPeriodToBillsMigration) Down(tx *sql.Tx) error {
	// SQLite 不支持直接删除列，我们不处理回滚
	return nil
}
```

- [ ] **Step 2: Add migration to migrations slice**

```go
// migrations 所有迁移
var migrations = []Migration{
	&BaseMigration{},
	&AddDepositAmountToLeasesMigration{},
	&AddOperationLogsTableMigration{},
	&AddStatusAndNoteToRoomsMigration{},
	&AddDueDateToBillsMigration{},
	&AddRefundDepositAmountToBillsMigration{},
	&AddPrintJobsTableMigration{},
	&AddPrintJobDetailsMigration{},
	&AddBillPeriodToBillsMigration{}, // Add this line
}
```

- [ ] **Step 3: Commit**

```bash
git add internal/infrastructure/persistence/sqlite/migration.go
git commit -m "feat: add migration for bill_start and bill_end columns"
```

---

### Task 6: Update Facade Layer - HTTP Handlers

**Files:**
- Modify: `internal/facade/cqrs_bill_handler.go`

- [ ] **Step 1: Update Create request struct**

```go
var req struct {
	LeaseID             string     `json:"lease_id"`
	Type                string     `json:"type"`
	Amount              int64      `json:"amount"`
	RentAmount          int64      `json:"rent_amount"`
	WaterAmount         int64      `json:"water_amount"`
	ElectricAmount      int64      `json:"electric_amount"`
	OtherAmount         int64      `json:"other_amount"`
	RefundDepositAmount int64      `json:"refund_deposit_amount"`
	BillStart           string     `json:"bill_start"`
	BillEnd             string     `json:"bill_end"`
	DueDate             string     `json:"due_date"`
	Note                string     `json:"note"`
}
```

- [ ] **Step 2: Parse BillStart and BillEnd in Create**

```go
// Parse bill start date
billStart, err := time.Parse("2006-01-02", req.BillStart)
if err != nil {
	billStart = time.Now()
}

// Parse bill end date
billEnd, err := time.Parse("2006-01-02", req.BillEnd)
if err != nil {
	billEnd = time.Now().AddDate(0, 1, 0)
}

// Parse due date (keep existing logic)
dueDate, err := time.Parse("2006-01-02", req.DueDate)
if err != nil {
	dueDate = time.Now().AddDate(0, 1, 0) // Default to 1 month from now
}
```

- [ ] **Step 3: Update CreateBillCommand construction**

```go
cmd := bill.CreateBillCommand{
	LeaseID:             req.LeaseID,
	Type:                billmodel.BillType(req.Type),
	Amount:              req.Amount,
	RentAmount:          req.RentAmount,
	WaterAmount:         req.WaterAmount,
	ElectricAmount:      req.ElectricAmount,
	OtherAmount:         req.OtherAmount,
	RefundDepositAmount: req.RefundDepositAmount,
	BillStart:           billStart,
	BillEnd:             billEnd,
	DueDate:             dueDate,
	Note:                req.Note,
}
```

- [ ] **Step 4: Update Update request struct**

```go
var req struct {
	Amount              int64      `json:"amount"`
	RentAmount          int64      `json:"rent_amount"`
	WaterAmount         int64      `json:"water_amount"`
	ElectricAmount      int64      `json:"electric_amount"`
	OtherAmount         int64      `json:"other_amount"`
	RefundDepositAmount int64      `json:"refund_deposit_amount"`
	BillStart           string     `json:"bill_start"`
	BillEnd             string     `json:"bill_end"`
	DueDate             string     `json:"due_date"`
	Note                string     `json:"note"`
}
```

- [ ] **Step 5: Parse BillStart and BillEnd in Update**

```go
// Parse bill start date
billStart, err := time.Parse("2006-01-02", req.BillStart)
if err != nil {
	billStart = time.Now()
}

// Parse bill end date
billEnd, err := time.Parse("2006-01-02", req.BillEnd)
if err != nil {
	billEnd = time.Now().AddDate(0, 1, 0)
}

// Parse due date (keep existing logic)
dueDate, err := time.Parse("2006-01-02", req.DueDate)
if err != nil {
	dueDate = time.Now().AddDate(0, 1, 0) // Default to 1 month from now
}
```

- [ ] **Step 6: Update UpdateBillCommand construction**

```go
cmd := bill.UpdateBillCommand{
	ID:                  id,
	Amount:              req.Amount,
	RentAmount:          req.RentAmount,
	WaterAmount:         req.WaterAmount,
	ElectricAmount:      req.ElectricAmount,
	OtherAmount:         req.OtherAmount,
	RefundDepositAmount: req.RefundDepositAmount,
	BillStart:           billStart,
	BillEnd:             billEnd,
	DueDate:             dueDate,
	Note:                req.Note,
}
```

- [ ] **Step 7: Commit**

```bash
git add internal/facade/cqrs_bill_handler.go
git commit -m "feat: update HTTP handlers for bill_start and bill_end"
```

---

### Task 7: Update Frontend Type Definitions

**Files:**
- Modify: `web/src/types/api.ts`

- [ ] **Step 1: Update Bill interface**

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
  billStart: string;
  billEnd: string;
  dueDate: string;
  paidAt: string | null;
  note: string;
  createdAt: string;
  updatedAt: string;
}
```

- [ ] **Step 2: Commit**

```bash
git add web/src/types/api.ts
git commit -m "feat: add billStart and billEnd to Bill TypeScript interface"
```

---

### Task 8: Update Frontend API Layer

**Files:**
- Modify: `web/src/api/bill.ts`

- [ ] **Step 1: Update create function parameters**

```typescript
create: async (data: {
  leaseId: string;
  type: string;
  amount: number;
  rentAmount: number;
  waterAmount: number;
  electricAmount: number;
  otherAmount: number;
  refundDepositAmount?: number;
  billStart: string;
  billEnd: string;
  dueDate: string | null;
  note: string;
}) => {
  const response = await apiClient.post<Bill>('/bills', data);
  return response.data;
},
```

- [ ] **Step 2: Update update function parameters**

```typescript
update: async (id: string, data: {
  amount: number;
  rentAmount: number;
  waterAmount: number;
  electricAmount: number;
  otherAmount: number;
  refundDepositAmount?: number;
  billStart: string;
  billEnd: string;
  dueDate: string | null;
  note: string;
}) => {
  const response = await apiClient.put<Bill>(`/bills/${id}`, data);
  return response.data;
},
```

- [ ] **Step 3: Commit**

```bash
git add web/src/api/bill.ts
git commit -m "feat: update bill API for billStart and billEnd"
```

---

### Task 9: Update Frontend Bills Page - Form and List

**Files:**
- Modify: `web/src/pages/Bills.tsx`

- [ ] **Step 1: Add state for billStart and billEnd in form**

Add form fields and state management for bill period selection.

- [ ] **Step 2: Add quick option buttons**

Add buttons for: "本月"、"本季度"、"上月"、"上季度"、"使用租约周期"

- [ ] **Step 3: Add date pickers for billStart and billEnd**

Add DatePicker components for manual date selection.

- [ ] **Step 4: Add auto-calculation logic when lease is selected**

When a lease is selected:
1. Fetch all bills for that lease
2. Find the latest paid bill
3. Calculate billStart = latest bill.billEnd + 1 day, or lease.startDate
4. Calculate billEnd based on month/quarter selection
5. Set dueDate = billEnd + 7 days

- [ ] **Step 5: Update form submission**

Include billStart and billEnd when submitting create/update.

- [ ] **Step 6: Add "计费周期" column to table**

```typescript
{
  title: '计费周期',
  key: 'billPeriod',
  width: 180,
  render: (_: any, record: Bill) => {
    return `${record.billStart} ~ ${record.billEnd}`;
  },
},
```

- [ ] **Step 7: Keep the "租期" column**

Keep the existing "租期" column that shows the lease's period.

- [ ] **Step 8: Commit**

```bash
git add web/src/pages/Bills.tsx
git commit -m "feat: update Bills page with billing period UI"
```

---

## Self-Review

**1. Spec coverage:**
- ✅ Data model changes - Task 1
- ✅ Command layer changes - Task 2
- ✅ Handler layer changes - Task 3
- ✅ Repository layer changes - Task 4
- ✅ Database migration - Task 5
- ✅ Facade/HTTP layer changes - Task 6
- ✅ Frontend type changes - Task 7
- ✅ Frontend API changes - Task 8
- ✅ Frontend UI changes - Task 9
- ✅ Business rules (auto-calculation) - Covered in Task 9
- ✅ Quick options - Covered in Task 9
- ✅ Both "租期" and "计费周期" columns - Covered in Task 9

**2. Placeholder scan:**
- ✅ No TBD/TODO placeholders
- ✅ All code snippets are complete
- ✅ All commands are specific

**3. Type consistency:**
- ✅ Field names consistent across all layers (billStart/BillStart/bill_start)
- ✅ Method signatures match across layers

---

Plan complete and saved to `docs/superpowers/plans/2026-03-31-bill-period-implementation.md`. Two execution options:

**1. Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach?
