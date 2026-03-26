package repository

import "time"

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// RoomRepository 房间仓储接口
type RoomRepository interface {
	Save(room *model.Room) error
	FindByID(id string) (*model.Room, error)
	FindByLocationIDAndRoomNumber(locationID, roomNumber string) (*model.Room, error)
	FindAll() ([]*model.Room, error)
	FindByCriteria(criteria RoomCriteria, offset, limit int) ([]*model.Room, error)
	CountByCriteria(criteria RoomCriteria) (int, error)
	Delete(id string) error
}

// RoomCriteria 房间查询条件
type RoomCriteria struct {
	LocationID string
	RoomNumber string
	Tags       []string
	StartTime  *time.Time
	EndTime    *time.Time
}
