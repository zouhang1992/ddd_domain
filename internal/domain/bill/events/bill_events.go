package events

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

type BillCreated struct {
	events.BaseEvent
	LeaseID string
	Type    string
	Amount  int64
}

type BillUpdated struct {
	events.BaseEvent
	Amount int64
}

type BillPaid struct {
	events.BaseEvent
	PaidAt string
}

type BillDeleted struct {
	events.BaseEvent
}

func NewBillCreated(id, leaseID string, version int, billType string, amount int64) BillCreated {
	return BillCreated{
		BaseEvent: events.NewBaseEvent("bill.created", id, version),
		LeaseID:   leaseID,
		Type:      billType,
		Amount:    amount,
	}
}

func NewBillUpdated(id string, version int, amount int64) BillUpdated {
	return BillUpdated{
		BaseEvent: events.NewBaseEvent("bill.updated", id, version),
		Amount:    amount,
	}
}

func NewBillPaid(id string, version int, paidAt string) BillPaid {
	return BillPaid{
		BaseEvent: events.NewBaseEvent("bill.paid", id, version),
		PaidAt:    paidAt,
	}
}

func NewBillDeleted(billID string, version int) BillDeleted {
	return BillDeleted{
		BaseEvent: events.NewBaseEvent("bill.deleted", billID, version),
	}
}
