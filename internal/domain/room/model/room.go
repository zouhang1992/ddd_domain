package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

// RoomStatus 房间状态
type RoomStatus string

const (
	RoomStatusAvailable RoomStatus = "available"
	RoomStatusRented    RoomStatus = "rented"
	RoomStatusMaintain  RoomStatus = "maintain"
)

// Room 房间领域模型（聚合根）
type Room struct {
	model.BaseAggregateRoot
	LocationID string
	RoomNumber string
	Status     RoomStatus
	Tags       []string
	Note       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// 房间事件（本地定义，避免导入循环）
type roomCreated struct {
	events.BaseEvent
	LocationID string
	RoomNumber string
	Tags       []string
}

type roomUpdated struct {
	events.BaseEvent
	LocationID string
	RoomNumber string
	Tags       []string
}

// NewRoom 创建新房间
func NewRoom(id, locationID, roomNumber string, tags []string, note string) *Room {
	now := time.Now()
	room := &Room{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		LocationID:        locationID,
		RoomNumber:        roomNumber,
		Status:            RoomStatusAvailable,
		Tags:              tags,
		Note:              note,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	// 创建并记录事件
	evt := roomCreated{
		BaseEvent:  events.NewBaseEvent("room.created", room.ID(), room.Version()),
		LocationID: room.LocationID,
		RoomNumber: room.RoomNumber,
		Tags:       room.Tags,
	}
	room.RecordEvent(evt)
	return room
}

// Update 更新房间信息
func (r *Room) Update(locationID, roomNumber string, tags []string, note string) {
	r.LocationID = locationID
	r.RoomNumber = roomNumber
	r.Tags = tags
	r.Note = note
	r.UpdatedAt = time.Now()
	// 创建并记录事件
	evt := roomUpdated{
		BaseEvent:  events.NewBaseEvent("room.updated", r.ID(), r.Version()),
		LocationID: r.LocationID,
		RoomNumber: r.RoomNumber,
		Tags:       r.Tags,
	}
	r.RecordEvent(evt)
}

// MarkRented 标记房间为已出租
func (r *Room) MarkRented() {
	r.Status = RoomStatusRented
	r.UpdatedAt = time.Now()
}

// MarkAvailable 标记房间为可出租
func (r *Room) MarkAvailable() {
	r.Status = RoomStatusAvailable
	r.UpdatedAt = time.Now()
}
