package events

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

type LandlordCreated struct {
	events.BaseEvent
	Name  string
	Phone string
	Note  string
}

type LandlordUpdated struct {
	events.BaseEvent
	Name  string
	Phone string
	Note  string
}

type LandlordDeleted struct {
	events.BaseEvent
}

func NewLandlordCreated(id string, version int, name, phone, note string) LandlordCreated {
	return LandlordCreated{
		BaseEvent: events.NewBaseEvent("landlord.created", id, version),
		Name:      name,
		Phone:     phone,
		Note:      note,
	}
}

func NewLandlordUpdated(id string, version int, name, phone, note string) LandlordUpdated {
	return LandlordUpdated{
		BaseEvent: events.NewBaseEvent("landlord.updated", id, version),
		Name:      name,
		Phone:     phone,
		Note:      note,
	}
}

func NewLandlordDeleted(landlordID string, version int) LandlordDeleted {
	return LandlordDeleted{
		BaseEvent: events.NewBaseEvent("landlord.deleted", landlordID, version),
	}
}
