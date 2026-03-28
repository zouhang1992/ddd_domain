package events

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

type RoomCreated struct {
	events.BaseEvent
	LocationID string
	RoomNumber string
	Tags       []string
}

type RoomUpdated struct {
	events.BaseEvent
	LocationID string
	RoomNumber string
	Tags       []string
}

type RoomDeleted struct {
	events.BaseEvent
}

func NewRoomCreated(id string, version int, locationID, roomNumber string, tags []string) RoomCreated {
	return RoomCreated{
		BaseEvent:  events.NewBaseEvent("room.created", id, version),
		LocationID: locationID,
		RoomNumber: roomNumber,
		Tags:       tags,
	}
}

func NewRoomUpdated(id string, version int, locationID, roomNumber string, tags []string) RoomUpdated {
	return RoomUpdated{
		BaseEvent:  events.NewBaseEvent("room.updated", id, version),
		LocationID: locationID,
		RoomNumber: roomNumber,
		Tags:       tags,
	}
}

func NewRoomDeleted(roomID string, version int) RoomDeleted {
	return RoomDeleted{
		BaseEvent: events.NewBaseEvent("room.deleted", roomID, version),
	}
}
