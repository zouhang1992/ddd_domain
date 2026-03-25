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
	_, ok := q.(query.ListLandlordsQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	landlords, err := h.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return &query.LandlordsQueryResult{Items: landlords}, nil
}
