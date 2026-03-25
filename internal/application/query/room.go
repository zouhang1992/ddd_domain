package query

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// GetRoomQuery 获取房间查询
type GetRoomQuery struct {
	ID string
}

// QueryName 实现 Query 接口
func (q GetRoomQuery) QueryName() string {
	return "get_room"
}

// ListRoomsQuery 列出所有房间查询
type ListRoomsQuery struct {
	LocationID string
	Tags       []string
}

// QueryName 实现 Query 接口
func (q ListRoomsQuery) QueryName() string {
	return "list_rooms"
}

// ListRoomsByLocationQuery 按位置列出房间查询
type ListRoomsByLocationQuery struct {
	LocationID string
}

// QueryName 实现 Query 接口
func (q ListRoomsByLocationQuery) QueryName() string {
	return "list_rooms_by_location"
}

// RoomQueryResult 房间查询结果
type RoomQueryResult struct {
	*model.Room
}

// RoomsQueryResult 房间列表查询结果
type RoomsQueryResult struct {
	Items []*model.Room
}
