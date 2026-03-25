package handler

import (
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// RoomCommandHandler 房间命令处理器
type RoomCommandHandler struct {
	repo     repository.RoomRepository
	eventBus *event.Bus
}

// NewRoomCommandHandler 创建房间命令处理器
func NewRoomCommandHandler(repo repository.RoomRepository, eventBus *event.Bus) *RoomCommandHandler {
	return &RoomCommandHandler{repo: repo, eventBus: eventBus}
}

// HandleCreateRoom 处理创建房间命令
func (h *RoomCommandHandler) HandleCreateRoom(cmd command.Command) (any, error) {
	createCmd, ok := cmd.(command.CreateRoomCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := createCmd.Validate(); err != nil {
		return nil, err
	}

	// 检查房间号在同一位置下是否已存在
	existingRoom, err := h.repo.FindByLocationIDAndRoomNumber(createCmd.LocationID, createCmd.RoomNumber)
	if err != nil {
		return nil, err
	}
	if existingRoom != nil {
		return nil, model.ErrRoomNumberExists
	}

	id := uuid.NewString()
	room := model.NewRoom(id, createCmd.LocationID, createCmd.RoomNumber, createCmd.Tags)
	if err := h.repo.Save(room); err != nil {
		return nil, err
	}

	// 发布房间创建事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewRoomCreated(room))
	}

	return room, nil
}

// HandleUpdateRoom 处理更新房间命令
func (h *RoomCommandHandler) HandleUpdateRoom(cmd command.Command) (any, error) {
	updateCmd, ok := cmd.(command.UpdateRoomCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := updateCmd.Validate(); err != nil {
		return nil, err
	}

	room, err := h.repo.FindByID(updateCmd.ID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, model.ErrNotFound
	}

	// 如果位置或房间号发生了变化，检查新值是否已存在
	if updateCmd.LocationID != "" && updateCmd.RoomNumber != "" && (updateCmd.LocationID != room.LocationID || updateCmd.RoomNumber != room.RoomNumber) {
		existingRoom, err := h.repo.FindByLocationIDAndRoomNumber(updateCmd.LocationID, updateCmd.RoomNumber)
		if err != nil {
			return nil, err
		}
		if existingRoom != nil && existingRoom.ID != updateCmd.ID {
			return nil, model.ErrRoomNumberExists
		}
	}

	room.Update(updateCmd.LocationID, updateCmd.RoomNumber, updateCmd.Tags)
	if err := h.repo.Save(room); err != nil {
		return nil, err
	}

	// 发布房间更新事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewRoomUpdated(room))
	}

	return room, nil
}

// HandleDeleteRoom 处理删除房间命令
func (h *RoomCommandHandler) HandleDeleteRoom(cmd command.Command) (any, error) {
	deleteCmd, ok := cmd.(command.DeleteRoomCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := deleteCmd.Validate(); err != nil {
		return nil, err
	}

	if err := h.repo.Delete(deleteCmd.ID); err != nil {
		return nil, err
	}

	// 发布房间删除事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewRoomDeleted(deleteCmd.ID))
	}

	return nil, nil
}
