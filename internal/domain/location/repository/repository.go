package repository

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/location/model"
)

type LocationRepository interface {
	FindByID(id string) (*model.Location, error)
	FindAll() ([]*model.Location, error)
	Save(location *model.Location) error
	Delete(id string) error
}
