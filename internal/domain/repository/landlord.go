package repository

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// LandlordRepository 房东仓储接口
type LandlordRepository interface {
	Save(landlord *model.Landlord) error
	FindByID(id string) (*model.Landlord, error)
	FindAll() ([]*model.Landlord, error)
	Delete(id string) error
	HasLeases(landlordID string) (bool, error)
}
