package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
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
	LeaseID    string
	Amount     int64
	Status     DepositStatus
	ReturnedAt *time.Time
	Note       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// 押金事件（本地定义，避免导入循环）
type depositCreated struct {
	events.BaseEvent
	LeaseID string
	Amount  int64
}

type depositReturning struct {
	events.BaseEvent
}

type depositReturned struct {
	events.BaseEvent
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
	// 创建并记录事件
	evt := depositCreated{
		BaseEvent: events.NewBaseEvent("deposit.created", deposit.ID(), deposit.Version()),
		LeaseID:   deposit.LeaseID,
		Amount:    deposit.Amount,
	}
	deposit.RecordEvent(evt)
	return deposit
}

// MarkReturning 标记押金为待退还
func (d *Deposit) MarkReturning() {
	d.Status = DepositStatusReturning
	d.UpdatedAt = time.Now()
	// 创建并记录事件
	evt := depositReturning{
		BaseEvent: events.NewBaseEvent("deposit.returning", d.ID(), d.Version()),
	}
	d.RecordEvent(evt)
}

// MarkReturned 标记押金为已退还
func (d *Deposit) MarkReturned() {
	now := time.Now()
	d.Status = DepositStatusReturned
	d.ReturnedAt = &now
	d.UpdatedAt = now
	// 创建并记录事件
	evt := depositReturned{
		BaseEvent: events.NewBaseEvent("deposit.returned", d.ID(), d.Version()),
	}
	d.RecordEvent(evt)
}
