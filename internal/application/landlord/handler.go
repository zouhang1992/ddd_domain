package landlord

import (
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/common"
	landlordmodel "github.com/zouhang1992/ddd_domain/internal/domain/landlord/model"
	landlordrepo "github.com/zouhang1992/ddd_domain/internal/domain/landlord/repository"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// CommandHandler 房东命令处理器
type CommandHandler struct {
	repo     landlordrepo.LandlordRepository
	eventBus *event.Bus
}

// NewCommandHandler 创建房东命令处理器
func NewCommandHandler(repo landlordrepo.LandlordRepository, eventBus *event.Bus) *CommandHandler {
	return &CommandHandler{repo: repo, eventBus: eventBus}
}

// HandleCreateLandlord 处理创建房东命令
func (h *CommandHandler) HandleCreateLandlord(cmd common.Command) (any, error) {
	createCmd, ok := cmd.(CreateLandlordCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := createCmd.Validate(); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	landlord := landlordmodel.NewLandlord(id, createCmd.Name, createCmd.Phone, createCmd.Note)
	if err := h.repo.Save(landlord); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range landlord.Events() {
			h.eventBus.PublishAsync(evt)
		}
		landlord.ClearEvents()
	}

	return landlord, nil
}

// HandleUpdateLandlord 处理更新房东命令
func (h *CommandHandler) HandleUpdateLandlord(cmd common.Command) (any, error) {
	updateCmd, ok := cmd.(UpdateLandlordCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := updateCmd.Validate(); err != nil {
		return nil, err
	}

	landlord, err := h.repo.FindByID(updateCmd.ID)
	if err != nil {
		return nil, err
	}
	if landlord == nil {
		return nil, domerrors.ErrNotFound
	}

	landlord.Update(updateCmd.Name, updateCmd.Phone, updateCmd.Note)
	if err := h.repo.Save(landlord); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range landlord.Events() {
			h.eventBus.PublishAsync(evt)
		}
		landlord.ClearEvents()
	}

	return landlord, nil
}

// HandleDeleteLandlord 处理删除房东命令
func (h *CommandHandler) HandleDeleteLandlord(cmd common.Command) (any, error) {
	deleteCmd, ok := cmd.(DeleteLandlordCommand)
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

// QueryHandler 房东查询处理器
type QueryHandler struct {
	repo landlordrepo.LandlordRepository
}

// NewQueryHandler 创建房东查询处理器
func NewQueryHandler(repo landlordrepo.LandlordRepository) *QueryHandler {
	return &QueryHandler{repo: repo}
}

// HandleGetLandlord 处理获取房东查询
func (h *QueryHandler) HandleGetLandlord(q common.Query) (any, error) {
	getQuery, ok := q.(GetLandlordQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	landlord, err := h.repo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if landlord == nil {
		return nil, domerrors.ErrNotFound
	}

	return &LandlordQueryResult{Landlord: landlord}, nil
}

// HandleListLandlords 处理列出租东查询
func (h *QueryHandler) HandleListLandlords(q common.Query) (any, error) {
	listQuery, ok := q.(ListLandlordsQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	landlords, err := h.repo.FindAll()
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
	var paginated []*landlordmodel.Landlord
	total := len(landlords)
	start := listQuery.Offset
	if start < total {
		end := start + limit
		if end > total {
			end = total
		}
		paginated = landlords[start:end]
	}

	result := &LandlordsQueryResult{
		Items: paginated,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}
