package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
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
	// 暂时注释掉，先解决导入循环问题
	// room.RecordEvent(events.NewRoomCreated(room.ID(), room.Version(), room.LocationID, room.RoomNumber, room.Tags))
	return room
}

// Update 更新房间信息
func (r *Room) Update(locationID, roomNumber string, tags []string, note string) {
	r.LocationID = locationID
	r.RoomNumber = roomNumber
	r.Tags = tags
	r.Note = note
	r.UpdatedAt = time.Now()
	// 暂时注释掉，先解决导入循环问题
	// r.RecordEvent(events.NewRoomUpdated(r.ID(), r.Version(), r.LocationID, r.RoomNumber, r.Tags))
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
