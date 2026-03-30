package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

// Landlord 房东领域模型（聚合根）
type Landlord struct {
	model.BaseAggregateRoot
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 房东事件（本地定义，避免导入循环）
type landlordCreated struct {
	events.BaseEvent
	Name  string
	Phone string
	Note  string
}

type landlordUpdated struct {
	events.BaseEvent
	Name  string
	Phone string
	Note  string
}

type landlordDeleted struct {
	events.BaseEvent
}

// NewLandlord 创建新房东
func NewLandlord(id, name, phone, note string) *Landlord {
	now := time.Now()
	landlord := &Landlord{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		Name:              name,
		Phone:             phone,
		Note:              note,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	// 创建并记录事件
	evt := landlordCreated{
		BaseEvent: events.NewBaseEvent("landlord.created", landlord.ID(), landlord.Version()),
		Name:      landlord.Name,
		Phone:     landlord.Phone,
		Note:      landlord.Note,
	}
	landlord.RecordEvent(evt)
	return landlord
}

// Update 更新房东信息
func (l *Landlord) Update(name, phone, note string) {
	l.Name = name
	l.Phone = phone
	l.Note = note
	l.UpdatedAt = time.Now()
	// 创建并记录事件
	evt := landlordUpdated{
		BaseEvent: events.NewBaseEvent("landlord.updated", l.ID(), l.Version()),
		Name:      l.Name,
		Phone:     l.Phone,
		Note:      l.Note,
	}
	l.RecordEvent(evt)
}

// Equals 比较房东是否相等
func (l *Landlord) Equals(other *Landlord) bool {
	return l.ID() == other.ID()
}
