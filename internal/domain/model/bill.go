package model

import (
	"time"
)

// BillType 账单类型
type BillType string

const (
	BillTypeCharge   BillType = "charge"   // 收账
	BillTypeCheckout BillType = "checkout" // 退租结算
)

// BillStatus 账单状态
type BillStatus string

const (
	BillStatusPaid    BillStatus = "paid"    // 已到账
	BillStatusPending BillStatus = "pending" // 待到账
)

// Bill 账单领域模型
type Bill struct {
	ID             string       `json:"id"`
	LeaseID        string       `json:"leaseId"`
	Type           BillType     `json:"type"`
	Status         BillStatus   `json:"status"`
	Amount         int64        `json:"amount"`      // 总金额（分）
	RentAmount     int64        `json:"rentAmount"`  // 租金金额（分）
	WaterAmount    int64        `json:"waterAmount"` // 水费金额（分）
	ElectricAmount int64        `json:"electricAmount"` // 电费金额（分）
	OtherAmount    int64        `json:"otherAmount"` // 其他金额（分）
	PaidAt         *time.Time   `json:"paidAt"`      // 到账时间
	Note           string       `json:"note"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}

// NewBill 创建新账单
func NewBill(id, leaseID string, billType BillType, amount, rentAmount, waterAmount,
	electricAmount, otherAmount int64, paidAt *time.Time, note string) *Bill {
	now := time.Now()
	status := BillStatusPending
	if paidAt != nil {
		status = BillStatusPaid
	}
	return &Bill{
		ID:             id,
		LeaseID:        leaseID,
		Type:           billType,
		Status:         status,
		Amount:         amount,
		RentAmount:     rentAmount,
		WaterAmount:    waterAmount,
		ElectricAmount: electricAmount,
		OtherAmount:    otherAmount,
		PaidAt:         paidAt,
		Note:           note,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Update 更新账单信息
func (b *Bill) Update(amount, rentAmount, waterAmount, electricAmount, otherAmount int64,
	paidAt *time.Time, note string) {
	b.Amount = amount
	b.RentAmount = rentAmount
	b.WaterAmount = waterAmount
	b.ElectricAmount = electricAmount
	b.OtherAmount = otherAmount
	b.PaidAt = paidAt
	if paidAt != nil {
		b.Status = BillStatusPaid
	} else {
		b.Status = BillStatusPending
	}
	b.Note = note
	b.UpdatedAt = time.Now()
}

// MarkPaid 标记账单为已到账
func (b *Bill) MarkPaid(paidAt time.Time) {
	b.PaidAt = &paidAt
	b.Status = BillStatusPaid
	b.UpdatedAt = time.Now()
}

// CanDelete 检查账单是否可以删除
func (b *Bill) CanDelete() bool {
	return true // 目前所有类型的账单都可以删除，但 checkout 类型删除会回滚状态
}
