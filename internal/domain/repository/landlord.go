package repository

import "time"

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// LandlordRepository 房东仓储接口
type LandlordRepository interface {
	Save(landlord *model.Landlord) error
	FindByID(id string) (*model.Landlord, error)
	FindAll() ([]*model.Landlord, error)
	FindByCriteria(criteria LandlordCriteria, offset, limit int) ([]*model.Landlord, error)
	CountByCriteria(criteria LandlordCriteria) (int, error)
	Delete(id string) error
	HasLeases(landlordID string) (bool, error)
}

// LandlordCriteria 房东查询条件
type LandlordCriteria struct {
	Name     string
	Phone    string
	StartTime    *time.Time
	EndTime      *time.Time
}
