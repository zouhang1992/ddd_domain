package query

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/model"
)

// ==================== 查询操作日志列表 ====================

// ListOperationLogsQuery 查询操作日志列表
type ListOperationLogsQuery struct {
	BaseQuery
	DomainType  string     // 领域类型
	EventName   string     // 事件名称（模糊匹配）
	AggregateID string     // 聚合根ID
	OperatorID  string     // 操作人ID
	StartTime   *time.Time // 开始时间
	EndTime     *time.Time // 结束时间
	Offset      int        // 偏移量
	Limit       int        // 每页数量
}

// QueryName 实现 Query 接口
func (q ListOperationLogsQuery) QueryName() string {
	return "list_operation_logs"
}

// OperationLogItem 操作日志项（用于查询结果）
type OperationLogItem struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	EventName   string                 `json:"eventName"`
	DomainType  string                 `json:"domainType"`
	AggregateID string                 `json:"aggregateId"`
	OperatorID  string                 `json:"operatorId"`
	Action      string                 `json:"action"`
	Details     map[string]interface{} `json:"details,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
}

// OperationLogsQueryResult 操作日志查询结果
type OperationLogsQueryResult struct {
	Items []*OperationLogItem `json:"items"`
	Total int                 `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}

// ==================== 获取单条操作日志 ====================

// GetOperationLogQuery 获取单条操作日志
type GetOperationLogQuery struct {
	BaseQuery
	ID string
}

// QueryName 实现 Query 接口
func (q GetOperationLogQuery) QueryName() string {
	return "get_operation_log"
}

// OperationLogQueryResult 单条操作日志查询结果
type OperationLogQueryResult struct {
	OperationLog *model.OperationLog `json:"operationLog"`
}

// ToOperationLogItem 将 model.OperationLog 转换为 OperationLogItem
func ToOperationLogItem(log *model.OperationLog) *OperationLogItem {
	if log == nil {
		return nil
	}
	return &OperationLogItem{
		ID:          log.ID(),
		Timestamp:   log.Timestamp(),
		EventName:   log.EventName(),
		DomainType:  log.DomainType(),
		AggregateID: log.AggregateID(),
		OperatorID:  log.OperatorID(),
		Action:      log.Action(),
		Details:     log.Details(),
		CreatedAt:   log.CreatedAt(),
	}
}

// ToOperationLogItems 将 model.OperationLog 切片转换为 OperationLogItem 切片
func ToOperationLogItems(logs []*model.OperationLog) []*OperationLogItem {
	items := make([]*OperationLogItem, len(logs))
	for i, log := range logs {
		items[i] = ToOperationLogItem(log)
	}
	return items
}
