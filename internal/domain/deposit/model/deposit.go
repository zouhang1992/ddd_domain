package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
)

// DepositStatus 押金状态
type DepositStatus string

const (
	DepositStatusPaid      DepositStatus = "paid"
	DepositStatusReturning DepositStatus = "returning"
	DepositStatusReturned  DepositStatus = "returned"
)

// Deposit 押金领域模型（聚合根）
type Deposit struct {
	model.BaseAggregateRoot
	LeaseID  string
	Amount   int64
	Status   DepositStatus
	ReturnedAt *time.Time
	Note     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewDeposit 创建新押金
func NewDeposit(id, leaseID string, amount int64, note string) *Deposit {
	now := time.Now()
	deposit := &Deposit{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		LeaseID:         leaseID,
		Amount:          amount,
		Status:          DepositStatusPaid,
		Note:            note,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	// 暂时注释掉，先解决导入循环问题
	// deposit.RecordEvent(events.NewDepositCreated(deposit.ID(), deposit.Version(), deposit.LeaseID, deposit.Amount))
	return deposit
}

// MarkReturning 标记押金为待退还
func (d *Deposit) MarkReturning() {
	d.Status = DepositStatusReturning
	d.UpdatedAt = time.Now()
	// 暂时注释掉，先解决导入循环问题
	// d.RecordEvent(events.NewDepositReturning(d.ID(), d.Version()))
}

// MarkReturned 标记押金为已退还
func (d *Deposit) MarkReturned() {
	now := time.Now()
	d.Status = DepositStatusReturned
	d.ReturnedAt = &now
	d.UpdatedAt = now
	// 暂时注释掉，先解决导入循环问题
	// d.RecordEvent(events.NewDepositReturned(d.ID(), d.Version()))
}
