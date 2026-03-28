package events

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

type LeaseCreated struct {
	events.BaseEvent
	RoomID     string
	LandlordID string
	TenantName string
}

type LeaseActivated struct {
	events.BaseEvent
	RoomID string
}

type LeaseCheckout struct {
	events.BaseEvent
	RoomID string
}

type LeaseExpired struct {
	events.BaseEvent
	RoomID string
}

type LeaseRenewed struct {
	events.BaseEvent
	NewEndDate string
}

type LeaseDeleted struct {
	events.BaseEvent
}

func NewLeaseCreated(id string, version int, roomID, landlordID, tenantName string) LeaseCreated {
	return LeaseCreated{
		BaseEvent:  events.NewBaseEvent("lease.created", id, version),
		RoomID:     roomID,
		LandlordID: landlordID,
		TenantName: tenantName,
	}
}

func NewLeaseActivated(id string, version int, roomID string) LeaseActivated {
	return LeaseActivated{
		BaseEvent: events.NewBaseEvent("lease.activated", id, version),
		RoomID:    roomID,
	}
}

func NewLeaseCheckout(id string, version int, roomID string) LeaseCheckout {
	return LeaseCheckout{
		BaseEvent: events.NewBaseEvent("lease.checkout", id, version),
		RoomID:    roomID,
	}
}

func NewLeaseExpired(id string, version int, roomID string) LeaseExpired {
	return LeaseExpired{
		BaseEvent: events.NewBaseEvent("lease.expired", id, version),
		RoomID:    roomID,
	}
}

func NewLeaseRenewed(id string, version int, newEndDate string) LeaseRenewed {
	return LeaseRenewed{
		BaseEvent:  events.NewBaseEvent("lease.renewed", id, version),
		NewEndDate: newEndDate,
	}
}

func NewLeaseDeleted(leaseID string, version int) LeaseDeleted {
	return LeaseDeleted{
		BaseEvent: events.NewBaseEvent("lease.deleted", leaseID, version),
	}
}
