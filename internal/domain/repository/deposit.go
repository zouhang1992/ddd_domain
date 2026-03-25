package repository

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// DepositRepository 押金仓储接口
type DepositRepository interface {
	Save(deposit *model.Deposit) error
	FindByID(id string) (*model.Deposit, error)
	FindByLeaseID(leaseID string) (*model.Deposit, error)
	FindAll() ([]*model.Deposit, error)
	Delete(id string) error
}
