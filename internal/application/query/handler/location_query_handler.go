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
	_, ok := q.(query.ListLocationsQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	locations, err := h.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return &query.LocationsQueryResult{Items: locations}, nil
}
