package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
)

// Location 位置领域模型（聚合根）
type Location struct {
	model.BaseAggregateRoot
	ShortName string
	Detail    string
	CreatedAt time.Time
	UpdatedAt time.Time
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
	// 暂时注释掉，先解决导入循环问题
	// location.RecordEvent(events.NewLocationCreated(location.ID(), location.Version(), location.ShortName, location.Detail))
	return location
}

// Update 更新位置信息
func (l *Location) Update(shortName, detail string) {
	l.ShortName = shortName
	l.Detail = detail
	l.UpdatedAt = time.Now()
	// 暂时注释掉，先解决导入循环问题
	// l.RecordEvent(events.NewLocationUpdated(l.ID(), l.Version(), l.ShortName, l.Detail))
}
