package handler

import (
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// LocationCommandHandler 位置命令处理器
type LocationCommandHandler struct {
	repo     repository.LocationRepository
	eventBus *event.Bus
}

// NewLocationCommandHandler 创建位置命令处理器
func NewLocationCommandHandler(repo repository.LocationRepository, eventBus *event.Bus) *LocationCommandHandler {
	return &LocationCommandHandler{repo: repo, eventBus: eventBus}
}

// HandleCreateLocation 处理创建位置命令
func (h *LocationCommandHandler) HandleCreateLocation(cmd command.Command) (any, error) {
	createCmd, ok := cmd.(command.CreateLocationCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := createCmd.Validate(); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	location := model.NewLocation(id, createCmd.ShortName, createCmd.Detail)
	if err := h.repo.Save(location); err != nil {
		return nil, err
	}

	// 发布位置创建事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLocationCreated(location))
	}

	return location, nil
}

// HandleUpdateLocation 处理更新位置命令
func (h *LocationCommandHandler) HandleUpdateLocation(cmd command.Command) (any, error) {
	updateCmd, ok := cmd.(command.UpdateLocationCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := updateCmd.Validate(); err != nil {
		return nil, err
	}

	location, err := h.repo.FindByID(updateCmd.ID)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, model.ErrNotFound
	}

	location.Update(updateCmd.ShortName, updateCmd.Detail)
	if err := h.repo.Save(location); err != nil {
		return nil, err
	}

	// 发布位置更新事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLocationUpdated(location))
	}

	return location, nil
}

// HandleDeleteLocation 处理删除位置命令
func (h *LocationCommandHandler) HandleDeleteLocation(cmd command.Command) (any, error) {
	deleteCmd, ok := cmd.(command.DeleteLocationCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := deleteCmd.Validate(); err != nil {
		return nil, err
	}

	hasRooms, err := h.repo.HasRooms(deleteCmd.ID)
	if err != nil {
		return nil, err
	}
	if hasRooms {
		return nil, model.ErrCannotDelete
	}

	if err := h.repo.Delete(deleteCmd.ID); err != nil {
		return nil, err
	}

	// 发布位置删除事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLocationDeleted(deleteCmd.ID))
	}

	return nil, nil
}
