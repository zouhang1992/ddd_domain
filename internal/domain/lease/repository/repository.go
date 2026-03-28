package repository

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
)

type LeaseRepository interface {
	FindByID(id string) (*model.Lease, error)
	FindAll() ([]*model.Lease, error)
	Save(lease *model.Lease) error
	Delete(id string) error
	FindActiveLeasesExpiringBefore(time time.Time) ([]*model.Lease, error)
	HasBills(leaseID string) (bool, error)
	HasDeposit(leaseID string) (bool, error)
}
