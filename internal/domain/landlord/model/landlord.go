package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
)

// Landlord 房东领域模型（聚合根）
type Landlord struct {
	model.BaseAggregateRoot
	Name      string
	Phone     string
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
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
	// 暂时注释掉，先解决导入循环问题
	// landlord.RecordEvent(events.NewLandlordCreated(landlord.ID(), landlord.Version(), landlord.Name, landlord.Phone, landlord.Note))
	return landlord
}

// Update 更新房东信息
func (l *Landlord) Update(name, phone, note string) {
	l.Name = name
	l.Phone = phone
	l.Note = note
	l.UpdatedAt = time.Now()
	// 暂时注释掉，先解决导入循环问题
	// l.RecordEvent(events.NewLandlordUpdated(l.ID(), l.Version(), l.Name, l.Phone, l.Note))
}

// Equals 比较房东是否相等
func (l *Landlord) Equals(other *Landlord) bool {
	return l.ID() == other.ID()
}
