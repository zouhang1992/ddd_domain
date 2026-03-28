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
)

// Bill 账单领域模型（聚合根）
type Bill struct {
	model.BaseAggregateRoot
	LeaseID   string
	Type      BillType
	Amount    int64
	PaidAt    *time.Time
	DueDate   time.Time
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
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

// NewBill 创建新账单
func NewBill(id, leaseID string, billType BillType, amount int64,
	dueDate time.Time, note string) *Bill {
	now := time.Now()
	bill := &Bill{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		LeaseID:         leaseID,
		Type:            billType,
		Amount:          amount,
		DueDate:         dueDate,
		Note:            note,
		CreatedAt:       now,
		UpdatedAt:       now,
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
