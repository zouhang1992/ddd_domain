package events

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

type LocationCreated struct {
	events.BaseEvent
	ShortName string
	Detail    string
}

type LocationUpdated struct {
	events.BaseEvent
	ShortName string
	Detail    string
}

type LocationDeleted struct {
	events.BaseEvent
}

func NewLocationCreated(id string, version int, shortName, detail string) LocationCreated {
	return LocationCreated{
		BaseEvent: events.NewBaseEvent("location.created", id, version),
		ShortName: shortName,
		Detail:    detail,
	}
}

func NewLocationUpdated(id string, version int, shortName, detail string) LocationUpdated {
	return LocationUpdated{
		BaseEvent: events.NewBaseEvent("location.updated", id, version),
		ShortName: shortName,
		Detail:    detail,
	}
}

func NewLocationDeleted(locationID string, version int) LocationDeleted {
	return LocationDeleted{
		BaseEvent: events.NewBaseEvent("location.deleted", locationID, version),
	}
}
