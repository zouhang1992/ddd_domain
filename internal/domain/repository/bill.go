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
	Delete(id string) error
}
