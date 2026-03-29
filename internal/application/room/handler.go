package room

import (
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/common"
	roommodel "github.com/zouhang1992/ddd_domain/internal/domain/room/model"
	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// CommandHandler 房间命令处理器
type CommandHandler struct {
	repo     roomrepo.RoomRepository
	eventBus *event.Bus
}

// NewCommandHandler 创建房间命令处理器
func NewCommandHandler(repo roomrepo.RoomRepository, eventBus *event.Bus) *CommandHandler {
	return &CommandHandler{repo: repo, eventBus: eventBus}
}

// HandleCreateRoom 处理创建房间命令
func (h *CommandHandler) HandleCreateRoom(cmd common.Command) (any, error) {
	createCmd, ok := cmd.(CreateRoomCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := createCmd.Validate(); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	room := roommodel.NewRoom(id, createCmd.LocationID, createCmd.RoomNumber, createCmd.Tags, "")
	if err := h.repo.Save(room); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range room.Events() {
			h.eventBus.PublishAsync(evt)
		}
		room.ClearEvents()
	}

	return room, nil
}

// HandleUpdateRoom 处理更新房间命令
func (h *CommandHandler) HandleUpdateRoom(cmd common.Command) (any, error) {
	updateCmd, ok := cmd.(UpdateRoomCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := updateCmd.Validate(); err != nil {
		return nil, err
	}

	room, err := h.repo.FindByID(updateCmd.ID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, domerrors.ErrNotFound
	}

	room.Update(updateCmd.LocationID, updateCmd.RoomNumber, updateCmd.Tags, "")
	if err := h.repo.Save(room); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range room.Events() {
			h.eventBus.PublishAsync(evt)
		}
		room.ClearEvents()
	}

	return room, nil
}

// HandleDeleteRoom 处理删除房间命令
func (h *CommandHandler) HandleDeleteRoom(cmd common.Command) (any, error) {
	deleteCmd, ok := cmd.(DeleteRoomCommand)
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

// QueryHandler 房间查询处理器
type QueryHandler struct {
	repo roomrepo.RoomRepository
}

// NewQueryHandler 创建房间查询处理器
func NewQueryHandler(repo roomrepo.RoomRepository) *QueryHandler {
	return &QueryHandler{repo: repo}
}

// HandleGetRoom 处理获取房间查询
func (h *QueryHandler) HandleGetRoom(q common.Query) (any, error) {
	getQuery, ok := q.(GetRoomQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	room, err := h.repo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, domerrors.ErrNotFound
	}

	return &RoomQueryResult{Room: room}, nil
}

// HandleListRooms 处理列出所有房间查询
func (h *QueryHandler) HandleListRooms(q common.Query) (any, error) {
	listQuery, ok := q.(ListRoomsQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	rooms, err := h.repo.FindAll()
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
	var paginated []*roommodel.Room
	total := len(rooms)
	start := listQuery.Offset
	if start < total {
		end := start + limit
		if end > total {
			end = total
		}
		paginated = rooms[start:end]
	}

	result := &RoomsQueryResult{
		Items: paginated,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}

// HandleListRoomsByLocation 处理按位置列出房间查询
func (h *QueryHandler) HandleListRoomsByLocation(q common.Query) (any, error) {
	listQuery, ok := q.(ListRoomsByLocationQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	allRooms, err := h.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var filtered []*roommodel.Room
	for _, room := range allRooms {
		if room.LocationID == listQuery.LocationID {
			filtered = append(filtered, room)
		}
	}

	return &RoomsQueryResult{Items: filtered}, nil
}
