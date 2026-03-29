package location

import (
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/common"
	locationmodel "github.com/zouhang1992/ddd_domain/internal/domain/location/model"
	locationrepo "github.com/zouhang1992/ddd_domain/internal/domain/location/repository"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// CommandHandler 位置命令处理器
type CommandHandler struct {
	repo     locationrepo.LocationRepository
	eventBus *event.Bus
}

// NewCommandHandler 创建位置命令处理器
func NewCommandHandler(repo locationrepo.LocationRepository, eventBus *event.Bus) *CommandHandler {
	return &CommandHandler{repo: repo, eventBus: eventBus}
}

// HandleCreateLocation 处理创建位置命令
func (h *CommandHandler) HandleCreateLocation(cmd common.Command) (any, error) {
	createCmd, ok := cmd.(CreateLocationCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := createCmd.Validate(); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	location := locationmodel.NewLocation(id, createCmd.ShortName, createCmd.Detail)
	if err := h.repo.Save(location); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range location.Events() {
			h.eventBus.PublishAsync(evt)
		}
		location.ClearEvents()
	}

	return location, nil
}

// HandleUpdateLocation 处理更新位置命令
func (h *CommandHandler) HandleUpdateLocation(cmd common.Command) (any, error) {
	updateCmd, ok := cmd.(UpdateLocationCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := updateCmd.Validate(); err != nil {
		return nil, err
	}

	location, err := h.repo.FindByID(updateCmd.ID)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, domerrors.ErrNotFound
	}

	location.Update(updateCmd.ShortName, updateCmd.Detail)
	if err := h.repo.Save(location); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range location.Events() {
			h.eventBus.PublishAsync(evt)
		}
		location.ClearEvents()
	}

	return location, nil
}

// HandleDeleteLocation 处理删除位置命令
func (h *CommandHandler) HandleDeleteLocation(cmd common.Command) (any, error) {
	deleteCmd, ok := cmd.(DeleteLocationCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := deleteCmd.Validate(); err != nil {
		return nil, err
	}

	if err := h.repo.Delete(deleteCmd.ID); err != nil {
		return nil, err
	}

	return nil, nil
}

// QueryHandler 位置查询处理器
type QueryHandler struct {
	repo locationrepo.LocationRepository
}

// NewQueryHandler 创建位置查询处理器
func NewQueryHandler(repo locationrepo.LocationRepository) *QueryHandler {
	return &QueryHandler{repo: repo}
}

// HandleGetLocation 处理获取位置查询
func (h *QueryHandler) HandleGetLocation(q common.Query) (any, error) {
	getQuery, ok := q.(GetLocationQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	location, err := h.repo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, domerrors.ErrNotFound
	}

	return &LocationQueryResult{Location: location}, nil
}

// HandleListLocations 处理列出所有位置查询
func (h *QueryHandler) HandleListLocations(q common.Query) (any, error) {
	listQuery, ok := q.(ListLocationsQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	locations, err := h.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// 设置默认分页大小
	limit := listQuery.Limit
	if limit <= 0 {
		limit = 10
	}

	// 计算页码
	page := 1
	if listQuery.Offset > 0 && limit > 0 {
		page = (listQuery.Offset / limit) + 1
	}

	// Simple pagination
	var paginated []*locationmodel.Location
	total := len(locations)
	start := listQuery.Offset
	if start < total {
		end := start + limit
		if end > total {
			end = total
		}
		paginated = locations[start:end]
	}

	result := &LocationsQueryResult{
		Items: paginated,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}
