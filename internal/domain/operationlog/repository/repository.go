package repository

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/operationlog/model"
)

// OperationLogRepository 操作日志仓储接口
type OperationLogRepository interface {
	Save(log *model.OperationLog) error
	FindByID(id string) (*model.OperationLog, error)
	FindByAggregateID(aggregateID string) ([]*model.OperationLog, error)
	FindByDomainType(domainType string, offset, limit int) ([]*model.OperationLog, int, error)
	FindByDomainTypeAndAggregateID(domainType, aggregateID string, offset, limit int) ([]*model.OperationLog, int, error)
	FindByTimeRange(start, end time.Time, offset, limit int) ([]*model.OperationLog, int, error)
}
