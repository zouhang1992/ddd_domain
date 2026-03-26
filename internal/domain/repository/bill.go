package repository

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"time"
)

// BillRepository 账单仓储接口
type BillRepository interface {
	Save(bill *model.Bill) error
	FindByID(id string) (*model.Bill, error)
	FindAll() ([]*model.Bill, error)
	FindByLeaseID(leaseID string) ([]*model.Bill, error)
	FindByRoomID(roomID string) ([]*model.Bill, error)
	FindByMonth(year int, month time.Month) ([]*model.Bill, error)
	FindByCriteria(criteria BillCriteria, offset, limit int) ([]*model.Bill, error)
	CountByCriteria(criteria BillCriteria) (int, error)
	Delete(id string) error
}

// BillCriteria 账单查询条件
type BillCriteria struct {
	Type        string
	Status      string
	LeaseID     string
	RoomID      string
	Month       string
	MinAmount   int64
	MaxAmount   int64
	StartDate   *time.Time
	EndDate     *time.Time
}
