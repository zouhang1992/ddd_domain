package query

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"time"
)

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
	// 查询条件
	LocationID  string     // 位置ID
	RoomNumber  string     // 房间号（模糊搜索）
	Tags        []string   // 标签
	StartDate   *time.Time // 创建开始日期
	EndDate     *time.Time // 创建结束日期
	// 分页参数
	Offset      int        // 偏移量
	Limit       int        // 每页数量
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
	Items []*model.Room `json:"items"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}
