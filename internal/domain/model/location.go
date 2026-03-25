package model

import "time"

// Location 位置领域模型
type Location struct {
	ID        string    `json:"id"`
	ShortName string    `json:"shortName"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewLocation 创建新位置
func NewLocation(id, shortName, detail string) *Location {
	now := time.Now()
	return &Location{
		ID:        id,
		ShortName: shortName,
		Detail:    detail,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update 更新位置信息
func (l *Location) Update(shortName, detail string) {
	l.ShortName = shortName
	l.Detail = detail
	l.UpdatedAt = time.Now()
}
