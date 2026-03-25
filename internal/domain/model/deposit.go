package model

import (
	"time"
)

// DepositStatus 押金状态
type DepositStatus string

const (
	DepositStatusCollected DepositStatus = "collected" // 已收取
	DepositStatusRefunded  DepositStatus = "refunded"  // 已退还
	DepositStatusDeducted  DepositStatus = "deducted"  // 已扣除
)

// Deposit 押金领域模型
type Deposit struct {
	ID         string
	LeaseID    string
	Amount     int64 // 押金金额（分）
	Status     DepositStatus
	RefundedAt *time.Time // 退还时间
	DeductedAt *time.Time // 扣除时间
	Note       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewDeposit 创建新押金记录
func NewDeposit(id, leaseID string, amount int64, note string) *Deposit {
	now := time.Now()
	return &Deposit{
		ID:        id,
		LeaseID:   leaseID,
		Amount:    amount,
		Status:    DepositStatusCollected,
		Note:      note,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Refund 退还押金
func (d *Deposit) Refund(refundAmount int64, refundedAt time.Time, note string) {
	d.Amount = refundAmount
	d.Status = DepositStatusRefunded
	d.RefundedAt = &refundedAt
	d.Note = note
	d.UpdatedAt = time.Now()
}

// Deduct 扣除押金
func (d *Deposit) Deduct(deductAmount int64, deductedAt time.Time, reason, note string) {
	d.Amount = deductAmount
	d.Status = DepositStatusDeducted
	d.DeductedAt = &deductedAt
	if d.Note == "" {
		d.Note = reason
	} else {
		d.Note += "; " + reason
	}
	if note != "" {
		if d.Note == "" {
			d.Note = note
		} else {
			d.Note += "; " + note
		}
	}
	d.UpdatedAt = time.Now()
}

// BindLease 绑定到新租约
func (d *Deposit) BindLease(newLeaseID string) {
	d.LeaseID = newLeaseID
	d.UpdatedAt = time.Now()
}
