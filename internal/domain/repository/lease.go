package repository

import "time"

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
	FindByCriteria(criteria LeaseCriteria, offset, limit int) ([]*model.Lease, error)
	CountByCriteria(criteria LeaseCriteria) (int, error)
	Delete(id string) error
	HasBills(leaseID string) (bool, error)
	HasDeposit(leaseID string) (bool, error)
}

// LeaseCriteria 租约查询条件
type LeaseCriteria struct {
	TenantName  string
	TenantPhone string
	Status      string
	RoomID      string
	StartDate   *time.Time
	EndDate     *time.Time
}
