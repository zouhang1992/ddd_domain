package repository

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/room/model"
)

type RoomRepository interface {
	FindByID(id string) (*model.Room, error)
	FindAll() ([]*model.Room, error)
	Save(room *model.Room) error
	Delete(id string) error
}
