package repository

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/landlord/model"
)

type LandlordRepository interface {
	FindByID(id string) (*model.Landlord, error)
	FindAll() ([]*model.Landlord, error)
	Save(landlord *model.Landlord) error
	Delete(id string) error
}
