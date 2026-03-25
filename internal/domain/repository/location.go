package repository

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// LocationRepository 位置仓储接口
type LocationRepository interface {
	Save(location *model.Location) error
	FindByID(id string) (*model.Location, error)
	FindAll() ([]*model.Location, error)
	Delete(id string) error
	HasRooms(locationID string) (bool, error)
}
