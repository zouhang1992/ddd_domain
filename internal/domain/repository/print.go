package repository

import "time"

// PrintJobRepository 打印作业仓储接口
type PrintJobRepository interface {
	FindByID(id string) (interface{}, error)
	FindAll() ([]interface{}, error)
	FindByCriteria(criteria PrintJobCriteria, offset, limit int) ([]interface{}, error)
	CountByCriteria(criteria PrintJobCriteria) (int, error)
	Save(job interface{}) error
	Delete(id string) error
}

// PrintJobCriteria 打印作业查询条件
type PrintJobCriteria struct {
	Status     string
	StartTime  *time.Time
	EndTime    *time.Time
}
