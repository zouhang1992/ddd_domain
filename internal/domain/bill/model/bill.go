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
	LeaseID   string     `json:"lease_id"`
	Type      BillType   `json:"type"`
	Status    BillStatus `json:"status"`
	Amount    int64      `json:"amount"`
	RentAmount    int64  `json:"rent_amount"`
	WaterAmount   int64  `json:"water_amount"`
	ElectricAmount int64 `json:"electric_amount"`
	OtherAmount   int64  `json:"other_amount"`
	PaidAt    *time.Time `json:"paid_at"`
	DueDate   time.Time  `json:"due_date"`
	Note      string     `json:"note"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
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
	dueDate time.Time, note string) *Bill {
	now := time.Now()
	bill := &Bill{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		LeaseID:           leaseID,
		Type:              billType,
		Status:            BillStatusPending,
		Amount:            amount,
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
	rentAmount, waterAmount, electricAmount, otherAmount int64,
	dueDate time.Time, note string) *Bill {
	now := time.Now()
	// Calculate total amount
	totalAmount := rentAmount + waterAmount + electricAmount + otherAmount

	bill := &Bill{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		LeaseID:           leaseID,
		Type:              billType,
		Status:            BillStatusPending,
		Amount:            totalAmount,
		RentAmount:        rentAmount,
		WaterAmount:       waterAmount,
		ElectricAmount:    electricAmount,
		OtherAmount:       otherAmount,
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
func (b *Bill) Update(amount int64, dueDate time.Time, note string) {
	b.Amount = amount
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
func (b *Bill) UpdateWithDetails(rentAmount, waterAmount, electricAmount, otherAmount int64, dueDate time.Time, note string) {
	b.RentAmount = rentAmount
	b.WaterAmount = waterAmount
	b.ElectricAmount = electricAmount
	b.OtherAmount = otherAmount
	b.Amount = rentAmount + waterAmount + electricAmount + otherAmount
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
