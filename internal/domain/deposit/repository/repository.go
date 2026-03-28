package repository

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/deposit/model"
)

type DepositRepository interface {
	FindByID(id string) (*model.Deposit, error)
	FindByLeaseID(leaseID string) (*model.Deposit, error)
	FindAll() ([]*model.Deposit, error)
	Save(deposit *model.Deposit) error
	Delete(id string) error
}
