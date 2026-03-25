package model

import (
	"time"
)

// Landlord 房东领域模型
type Landlord struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewLandlord 创建新房东
func NewLandlord(id, name, phone, note string) *Landlord {
	now := time.Now()
	return &Landlord{
		ID:        id,
		Name:      name,
		Phone:     phone,
		Note:      note,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update 更新房东信息
func (l *Landlord) Update(name, phone, note string) {
	l.Name = name
	l.Phone = phone
	l.Note = note
	l.UpdatedAt = time.Now()
}
