package repository

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// RoomRepository 房间仓储接口
type RoomRepository interface {
	Save(room *model.Room) error
	FindByID(id string) (*model.Room, error)
	FindByLocationIDAndRoomNumber(locationID, roomNumber string) (*model.Room, error)
	FindAll() ([]*model.Room, error)
	Delete(id string) error
}
