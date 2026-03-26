package repository

import "time"

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// LocationRepository 位置仓储接口
type LocationRepository interface {
	Save(location *model.Location) error
	FindByID(id string) (*model.Location, error)
	FindAll() ([]*model.Location, error)
	FindByCriteria(criteria LocationCriteria, offset, limit int) ([]*model.Location, error)
	CountByCriteria(criteria LocationCriteria) (int, error)
	Delete(id string) error
	HasRooms(locationID string) (bool, error)
}

// LocationCriteria 位置查询条件
type LocationCriteria struct {
	ShortName string
	Detail    string
	StartTime *time.Time
	EndTime   *time.Time
}
