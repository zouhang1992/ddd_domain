package handler

import (
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// OperationLogQueryHandler 操作日志查询处理器
type OperationLogQueryHandler struct {
	repo repository.OperationLogRepository
}

// NewOperationLogQueryHandler 创建操作日志查询处理器
func NewOperationLogQueryHandler(repo repository.OperationLogRepository) *OperationLogQueryHandler {
	return &OperationLogQueryHandler{repo: repo}
}

// HandleListOperationLogs 处理列表查询
func (h *OperationLogQueryHandler) HandleListOperationLogs(q query.Query) (any, error) {
	listQuery, ok := q.(query.ListOperationLogsQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	// 构建查询条件
	criteria := repository.OperationLogCriteria{
		DomainType:   listQuery.DomainType,
		EventName:    listQuery.EventName,
		AggregateID:  listQuery.AggregateID,
		OperatorID:   listQuery.OperatorID,
		StartTime:    listQuery.StartTime,
		EndTime:      listQuery.EndTime,
	}

	// 查询数据
	logs, err := h.repo.FindByCriteria(criteria, listQuery.Offset, listQuery.Limit)
	if err != nil {
		return nil, err
	}

	// 获取总数
	total, err := h.repo.CountByCriteria(criteria)
	if err != nil {
		return nil, err
	}

	// 计算页码
	page := 1
	if listQuery.Offset > 0 && listQuery.Limit > 0 {
		page = (listQuery.Offset / listQuery.Limit) + 1
	}

	// 转换结果格式
	items := query.ToOperationLogItems(logs)

	result := &query.OperationLogsQueryResult{
		Items: items,
		Total: total,
		Page:  page,
		Limit: listQuery.Limit,
	}

	return result, nil
}

// HandleGetOperationLog 处理获取单条操作日志查询
func (h *OperationLogQueryHandler) HandleGetOperationLog(q query.Query) (any, error) {
	getQuery, ok := q.(query.GetOperationLogQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	log, err := h.repo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}

	if log == nil {
		return nil, model.ErrNotFound
	}

	return &query.OperationLogQueryResult{OperationLog: log}, nil
}
