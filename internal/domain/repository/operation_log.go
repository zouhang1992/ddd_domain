package repository

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/model"
)

// OperationLogRepository 操作日志仓储接口
type OperationLogRepository interface {
	Save(log *model.OperationLog) error
	FindByID(id string) (*model.OperationLog, error)
	FindAll(offset, limit int) ([]*model.OperationLog, error)
	FindByDomainType(domainType string, offset, limit int) ([]*model.OperationLog, error)
	FindByTimeRange(startTime, endTime time.Time, offset, limit int) ([]*model.OperationLog, error)
	FindByAggregateID(aggregateID string, offset, limit int) ([]*model.OperationLog, error)
	FindByOperatorID(operatorID string, offset, limit int) ([]*model.OperationLog, error)
	FindByCriteria(criteria OperationLogCriteria, offset, limit int) ([]*model.OperationLog, error)
	CountByCriteria(criteria OperationLogCriteria) (int, error)
}

// OperationLogCriteria 操作日志查询条件
type OperationLogCriteria struct {
	DomainType   string
	EventName    string
	AggregateID  string
	OperatorID   string
	StartTime    *time.Time
	EndTime      *time.Time
}
