package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
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
	// 暂时注释掉，先解决导入循环问题
	// bill.RecordEvent(events.NewBillCreated(bill.ID(), bill.LeaseID, bill.Version(), string(bill.Type), bill.Amount))
	return bill
}

// MarkPaid 标记账单为已支付
func (b *Bill) MarkPaid() {
	now := time.Now()
	b.PaidAt = &now
	b.UpdatedAt = now
	// 暂时注释掉，先解决导入循环问题
	// paidAt := ""
	// if b.PaidAt != nil {
	// 	paidAt = b.PaidAt.Format("2006-01-02")
	// }
	// b.RecordEvent(events.NewBillPaid(b.ID(), b.Version(), paidAt))
}

// Update 更新账单信息
func (b *Bill) Update(amount int64, dueDate time.Time, note string) {
	b.Amount = amount
	b.DueDate = dueDate
	b.Note = note
	b.UpdatedAt = time.Now()
	// 暂时注释掉，先解决导入循环问题
	// b.RecordEvent(events.NewBillUpdated(b.ID(), b.Version(), b.Amount))
}
