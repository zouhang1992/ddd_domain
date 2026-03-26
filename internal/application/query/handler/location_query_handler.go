package handler

import (
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// LocationQueryHandler 位置查询处理器
type LocationQueryHandler struct {
	repo repository.LocationRepository
}

// NewLocationQueryHandler 创建位置查询处理器
func NewLocationQueryHandler(repo repository.LocationRepository) *LocationQueryHandler {
	return &LocationQueryHandler{repo: repo}
}

// HandleGetLocation 处理获取位置查询
func (h *LocationQueryHandler) HandleGetLocation(q query.Query) (any, error) {
	getQuery, ok := q.(query.GetLocationQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	location, err := h.repo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, model.ErrNotFound
	}

	return &query.LocationQueryResult{Location: location}, nil
}

// HandleListLocations 处理列出所有位置查询
func (h *LocationQueryHandler) HandleListLocations(q query.Query) (any, error) {
	listQuery, ok := q.(query.ListLocationsQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	// 构建查询条件
	criteria := repository.LocationCriteria{
		ShortName: listQuery.ShortName,
		Detail:    listQuery.Detail,
	}

	// 设置默认分页大小
	limit := listQuery.Limit
	if limit <= 0 {
		limit = 10 // 默认返回10条
	}

	// 查询数据
	locations, err := h.repo.FindByCriteria(criteria, listQuery.Offset, limit)
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

	result := &query.LocationsQueryResult{
		Items: locations,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}
