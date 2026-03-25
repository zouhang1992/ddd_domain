package repository

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// LeaseRepository 租约仓储接口
type LeaseRepository interface {
	Save(lease *model.Lease) error
	FindByID(id string) (*model.Lease, error)
	FindAll() ([]*model.Lease, error)
	FindByRoomID(roomID string) ([]*model.Lease, error)
	FindByStatus(status model.LeaseStatus) ([]*model.Lease, error)
	FindByRoomIDAndStatus(roomID string, status model.LeaseStatus) ([]*model.Lease, error)
	FindActiveByRoomID(roomID string) (*model.Lease, error)
	Delete(id string) error
	HasBills(leaseID string) (bool, error)
	HasDeposit(leaseID string) (bool, error)
}
