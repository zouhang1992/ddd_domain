package handler

import (
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// LandlordQueryHandler 房东查询处理器
type LandlordQueryHandler struct {
	repo repository.LandlordRepository
}

// NewLandlordQueryHandler 创建房东查询处理器
func NewLandlordQueryHandler(repo repository.LandlordRepository) *LandlordQueryHandler {
	return &LandlordQueryHandler{repo: repo}
}

// HandleGetLandlord 处理获取房东查询
func (h *LandlordQueryHandler) HandleGetLandlord(q query.Query) (any, error) {
	getQuery, ok := q.(query.GetLandlordQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	landlord, err := h.repo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if landlord == nil {
		return nil, model.ErrNotFound
	}

	return &query.LandlordQueryResult{Landlord: landlord}, nil
}

// HandleListLandlords 处理列出租东查询
func (h *LandlordQueryHandler) HandleListLandlords(q query.Query) (any, error) {
	listQuery, ok := q.(query.ListLandlordsQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	// 构建查询条件
	criteria := repository.LandlordCriteria{
		Name:  listQuery.Name,
		Phone: listQuery.Phone,
	}

	// 设置默认分页大小
	limit := listQuery.Limit
	if limit <= 0 {
		limit = 10 // 默认返回10条
	}

	// 查询数据
	landlords, err := h.repo.FindByCriteria(criteria, listQuery.Offset, limit)
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

	result := &query.LandlordsQueryResult{
		Items: landlords,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}
