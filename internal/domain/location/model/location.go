package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

// Location 位置领域模型（聚合根）
type Location struct {
	model.BaseAggregateRoot
	ShortName string    `json:"short_name"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 位置事件（本地定义，避免导入循环）
type locationCreated struct {
	events.BaseEvent
	ShortName string
	Detail    string
}

type locationUpdated struct {
	events.BaseEvent
	ShortName string
	Detail    string
}

type locationDeleted struct {
	events.BaseEvent
}

// NewLocation 创建新位置
func NewLocation(id, shortName, detail string) *Location {
	now := time.Now()
	location := &Location{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		ShortName:         shortName,
		Detail:            detail,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	// 创建并记录事件
	evt := locationCreated{
		BaseEvent: events.NewBaseEvent("location.created", location.ID(), location.Version()),
		ShortName: location.ShortName,
		Detail:    location.Detail,
	}
	location.RecordEvent(evt)
	return location
}

// Update 更新位置信息
func (l *Location) Update(shortName, detail string) {
	l.ShortName = shortName
	l.Detail = detail
	l.UpdatedAt = time.Now()
	// 创建并记录事件
	evt := locationUpdated{
		BaseEvent: events.NewBaseEvent("location.updated", l.ID(), l.Version()),
		ShortName: l.ShortName,
		Detail:    l.Detail,
	}
	l.RecordEvent(evt)
}

// Equals 比较位置是否相等
func (l *Location) Equals(other *Location) bool {
	return l.ID() == other.ID()
}
