package handler

import (
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// LeaseQueryHandler 租约查询处理器
type LeaseQueryHandler struct {
	repo repository.LeaseRepository
}

// NewLeaseQueryHandler 创建租约查询处理器
func NewLeaseQueryHandler(repo repository.LeaseRepository) *LeaseQueryHandler {
	return &LeaseQueryHandler{repo: repo}
}

// HandleGetLease 处理获取租约查询
func (h *LeaseQueryHandler) HandleGetLease(q query.Query) (any, error) {
	getQuery, ok := q.(query.GetLeaseQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	lease, err := h.repo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, model.ErrNotFound
	}

	return &query.LeaseQueryResult{Lease: lease}, nil
}

// HandleListLeases 处理列出租约查询
func (h *LeaseQueryHandler) HandleListLeases(q query.Query) (any, error) {
	listQuery, ok := q.(query.ListLeasesQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	// 构建查询条件
	criteria := repository.LeaseCriteria{
		TenantName:  listQuery.TenantName,
		TenantPhone: listQuery.TenantPhone,
		Status:      listQuery.Status,
		RoomID:      listQuery.RoomID,
		StartDate:   listQuery.StartDate,
		EndDate:     listQuery.EndDate,
	}

	// 设置默认分页大小
	limit := listQuery.Limit
	if limit <= 0 {
		limit = 10 // 默认返回10条
	}

	// 查询数据
	leases, err := h.repo.FindByCriteria(criteria, listQuery.Offset, limit)
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
	if listQuery.Offset > 0 && limit > 0 {
		page = (listQuery.Offset / limit) + 1
	}

	result := &query.LeasesQueryResult{
		Items: leases,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}
