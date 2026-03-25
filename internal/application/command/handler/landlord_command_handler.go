package handler

import (
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// LandlordCommandHandler 房东命令处理器
type LandlordCommandHandler struct {
	repo     repository.LandlordRepository
	eventBus *event.Bus
}

// NewLandlordCommandHandler 创建房东命令处理器
func NewLandlordCommandHandler(repo repository.LandlordRepository, eventBus *event.Bus) *LandlordCommandHandler {
	return &LandlordCommandHandler{repo: repo, eventBus: eventBus}
}

// HandleCreateLandlord 处理创建房东命令
func (h *LandlordCommandHandler) HandleCreateLandlord(cmd command.Command) (any, error) {
	createCmd, ok := cmd.(command.CreateLandlordCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := createCmd.Validate(); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	landlord := model.NewLandlord(id, createCmd.Name, createCmd.Phone, createCmd.Note)
	if err := h.repo.Save(landlord); err != nil {
		return nil, err
	}

	// 发布房东创建事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLandlordCreated(landlord))
	}

	return landlord, nil
}

// HandleUpdateLandlord 处理更新房东命令
func (h *LandlordCommandHandler) HandleUpdateLandlord(cmd command.Command) (any, error) {
	updateCmd, ok := cmd.(command.UpdateLandlordCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := updateCmd.Validate(); err != nil {
		return nil, err
	}

	landlord, err := h.repo.FindByID(updateCmd.ID)
	if err != nil {
		return nil, err
	}
	if landlord == nil {
		return nil, model.ErrNotFound
	}

	landlord.Update(updateCmd.Name, updateCmd.Phone, updateCmd.Note)
	if err := h.repo.Save(landlord); err != nil {
		return nil, err
	}

	// 发布房东更新事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLandlordUpdated(landlord))
	}

	return landlord, nil
}

// HandleDeleteLandlord 处理删除房东命令
func (h *LandlordCommandHandler) HandleDeleteLandlord(cmd command.Command) (any, error) {
	deleteCmd, ok := cmd.(command.DeleteLandlordCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := deleteCmd.Validate(); err != nil {
		return nil, err
	}

	hasLeases, err := h.repo.HasLeases(deleteCmd.ID)
	if err != nil {
		return nil, err
	}
	if hasLeases {
		return nil, model.ErrCannotDelete
	}

	if err := h.repo.Delete(deleteCmd.ID); err != nil {
		return nil, err
	}

	// 发布房东删除事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLandlordDeleted(deleteCmd.ID))
	}

	return nil, nil
}
