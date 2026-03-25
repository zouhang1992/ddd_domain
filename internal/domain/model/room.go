package model

import (
	"strings"
	"time"
)

// Room 房间领域模型
type Room struct {
	ID         string    `json:"id"`
	LocationID string    `json:"locationId"`
	RoomNumber string    `json:"roomNumber"`
	Tags       []string  `json:"tags"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// NewRoom 创建新房间
func NewRoom(id, locationID, roomNumber string, tags []string) *Room {
	now := time.Now()
	return &Room{
		ID:         id,
		LocationID: locationID,
		RoomNumber: roomNumber,
		Tags:       tags,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Update 更新房间信息
func (r *Room) Update(locationID, roomNumber string, tags []string) {
	r.LocationID = locationID
	r.RoomNumber = roomNumber
	r.Tags = tags
	r.UpdatedAt = time.Now()
}

// TagsString 获取标签字符串（逗号分隔）
func (r *Room) TagsString() string {
	if len(r.Tags) == 0 {
		return ""
	}
	return strings.Join(r.Tags, ",")
}

// ParseTags 解析标签字符串
func ParseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	parts := strings.Split(tagsStr, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
