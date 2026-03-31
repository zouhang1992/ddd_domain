package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

// BillType 账单类型
type BillType string

const (
	BillTypeRent     BillType = "rent"
	BillTypeWater    BillType = "water"
	BillTypeElectric BillType = "electric"
	BillTypeGas      BillType = "gas"
	BillTypeInternet BillType = "internet"
	BillTypeOther    BillType = "other"
	BillTypeCharge   BillType = "charge"
	BillTypeCheckout BillType = "checkout"
)

// BillStatus 账单状态
type BillStatus string

const (
	BillStatusPending BillStatus = "pending"
	BillStatusPaid    BillStatus = "paid"
)

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
	RefundDepositAmount int64      `json:"refund_deposit_amount"` // 退还押金金额（分，正数表示退还）
	BillStart           time.Time  `json:"bill_start"`
	BillEnd             time.Time  `json:"bill_end"`
	DueDate             time.Time  `json:"due_date"`
	PaidAt              *time.Time `json:"paid_at"`
	Note                string     `json:"note"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// 账单事件（本地定义，避免导入循环）
type billCreated struct {
	events.BaseEvent
	LeaseID string
	Type    string
	Amount  int64
}

type billUpdated struct {
	events.BaseEvent
	Amount int64
}

type billPaid struct {
	events.BaseEvent
	PaidAt string
}

type billDeleted struct {
	events.BaseEvent
}

// NewBill 创建新账单
func NewBill(id, leaseID string, billType BillType, amount int64,
	billStart, billEnd time.Time, dueDate time.Time, note string) *Bill {
	now := time.Now()
	bill := &Bill{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		LeaseID:           leaseID,
		Type:              billType,
		Status:            BillStatusPending,
		Amount:            amount,
		BillStart:         billStart,
		BillEnd:           billEnd,
		DueDate:           dueDate,
		Note:              note,
		CreatedAt:         now,
		UpdatedAt:         now,
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

// MarkPaid 标记账单为已支付
func (b *Bill) MarkPaid() {
	now := time.Now()
	b.Status = BillStatusPaid
	b.PaidAt = &now
	b.UpdatedAt = now
	// 创建并记录事件
	paidAt := ""
	if b.PaidAt != nil {
		paidAt = b.PaidAt.Format("2006-01-02")
	}
	evt := billPaid{
		BaseEvent: events.NewBaseEvent("bill.paid", b.ID(), b.Version()),
		PaidAt:    paidAt,
	}
	b.RecordEvent(evt)
}

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

// Equals 比较账单是否相等
func (b *Bill) Equals(other *Bill) bool {
	return b.ID() == other.ID()
}
