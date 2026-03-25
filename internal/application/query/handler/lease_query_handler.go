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

	var leases []*model.Lease
	var err error

	if listQuery.Status != "" && listQuery.RoomID != "" {
		leases, err = h.repo.FindByRoomIDAndStatus(listQuery.RoomID, model.LeaseStatus(listQuery.Status))
	} else if listQuery.Status != "" {
		leases, err = h.repo.FindByStatus(model.LeaseStatus(listQuery.Status))
	} else if listQuery.RoomID != "" {
		leases, err = h.repo.FindByRoomID(listQuery.RoomID)
	} else {
		leases, err = h.repo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	return &query.LeasesQueryResult{Items: leases}, nil
}
