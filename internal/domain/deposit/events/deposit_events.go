package events

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

type DepositCreated struct {
	events.BaseEvent
	LeaseID string
	Amount  int64
}

type DepositReturning struct {
	events.BaseEvent
}

type DepositReturned struct {
	events.BaseEvent
}

type DepositDeleted struct {
	events.BaseEvent
}

func NewDepositCreated(id string, version int, leaseID string, amount int64) DepositCreated {
	return DepositCreated{
		BaseEvent: events.NewBaseEvent("deposit.created", id, version),
		LeaseID:   leaseID,
		Amount:    amount,
	}
}

func NewDepositReturning(id string, version int) DepositReturning {
	return DepositReturning{
		BaseEvent: events.NewBaseEvent("deposit.returning", id, version),
	}
}

func NewDepositReturned(id string, version int) DepositReturned {
	return DepositReturned{
		BaseEvent: events.NewBaseEvent("deposit.returned", id, version),
	}
}

func NewDepositDeleted(depositID string, version int) DepositDeleted {
	return DepositDeleted{
		BaseEvent: events.NewBaseEvent("deposit.deleted", depositID, version),
	}
}
