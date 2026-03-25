package handler

import (
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// RoomQueryHandler 房间查询处理器
type RoomQueryHandler struct {
	repo repository.RoomRepository
}

// NewRoomQueryHandler 创建房间查询处理器
func NewRoomQueryHandler(repo repository.RoomRepository) *RoomQueryHandler {
	return &RoomQueryHandler{repo: repo}
}

// HandleGetRoom 处理获取房间查询
func (h *RoomQueryHandler) HandleGetRoom(q query.Query) (any, error) {
	getQuery, ok := q.(query.GetRoomQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	room, err := h.repo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, model.ErrNotFound
	}

	return &query.RoomQueryResult{Room: room}, nil
}

// HandleListRooms 处理列出所有房间查询
func (h *RoomQueryHandler) HandleListRooms(q query.Query) (any, error) {
	listQuery, ok := q.(query.ListRoomsQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	// 先获取所有房间
	allRooms, err := h.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// 应用过滤
	var filtered []*model.Room
	for _, room := range allRooms {
		// 按位置过滤
		if listQuery.LocationID != "" && room.LocationID != listQuery.LocationID {
			continue
		}
		// 按标签过滤
		if len(listQuery.Tags) > 0 && !hasAnyTag(room.Tags, listQuery.Tags) {
			continue
		}
		filtered = append(filtered, room)
	}

	return &query.RoomsQueryResult{Items: filtered}, nil
}

// HandleListRoomsByLocation 处理按位置列出房间查询
func (h *RoomQueryHandler) HandleListRoomsByLocation(q query.Query) (any, error) {
	listQuery, ok := q.(query.ListRoomsByLocationQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	// 先获取所有房间，再过滤
	allRooms, err := h.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var filtered []*model.Room
	for _, room := range allRooms {
		if room.LocationID == listQuery.LocationID {
			filtered = append(filtered, room)
		}
	}

	return &query.RoomsQueryResult{Items: filtered}, nil
}

// hasAnyTag 检查房间是否有任意一个指定标签
func hasAnyTag(roomTags, queryTags []string) bool {
	tagSet := make(map[string]bool)
	for _, t := range roomTags {
		tagSet[t] = true
	}
	for _, t := range queryTags {
		if tagSet[t] {
			return true
		}
	}
	return false
}
