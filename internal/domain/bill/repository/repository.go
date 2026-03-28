package repository

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/bill/model"
)

type BillRepository interface {
	FindByID(id string) (*model.Bill, error)
	FindByLeaseID(leaseID string) ([]*model.Bill, error)
	FindAll() ([]*model.Bill, error)
	Save(bill *model.Bill) error
	Delete(id string) error
	FindUnpaidBillsDueBefore(dueDate time.Time) ([]*model.Bill, error)
}
