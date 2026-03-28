package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

// LeaseStatus 租约状态
type LeaseStatus string

const (
	LeaseStatusPending  LeaseStatus = "pending"
	LeaseStatusActive   LeaseStatus = "active"
	LeaseStatusExpired  LeaseStatus = "expired"
	LeaseStatusCheckout LeaseStatus = "checkout"
)

// Lease 租约领域模型（聚合根）
type Lease struct {
	model.BaseAggregateRoot
	RoomID         string
	LandlordID     string
	TenantName     string
	TenantPhone    string
	StartDate      time.Time
	EndDate        time.Time
	RentAmount     int64
	DepositAmount  int64
	Status         LeaseStatus
	Note           string
	LastChargeAt   *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// 租约事件（本地定义，避免导入循环）
type leaseCreated struct {
	events.BaseEvent
	RoomID     string
	LandlordID string
	TenantName string
}

type leaseActivated struct {
	events.BaseEvent
	RoomID string
}

type leaseCheckout struct {
	events.BaseEvent
	RoomID string
}

type leaseExpired struct {
	events.BaseEvent
	RoomID string
}

type leaseRenewed struct {
	events.BaseEvent
	NewEndDate string
}

// NewLease 创建新租约
func NewLease(id, roomID, landlordID, tenantName, tenantPhone string,
	startDate, endDate time.Time, rentAmount, depositAmount int64, note string) *Lease {
	now := time.Now()
	lease := &Lease{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		RoomID:         roomID,
		LandlordID:     landlordID,
		TenantName:     tenantName,
		TenantPhone:    tenantPhone,
		StartDate:      startDate,
		EndDate:        endDate,
		RentAmount:     rentAmount,
		DepositAmount:  depositAmount,
		Status:         LeaseStatusPending,
		Note:           note,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 创建并记录事件
	evt := leaseCreated{
		BaseEvent:  events.NewBaseEvent("lease.created", lease.ID(), lease.Version()),
		RoomID:     lease.RoomID,
		LandlordID: lease.LandlordID,
		TenantName: lease.TenantName,
	}
	lease.RecordEvent(evt)
	return lease
}

// Activate 激活租约
func (l *Lease) Activate() {
	l.Status = LeaseStatusActive
	l.UpdatedAt = time.Now()
	// 创建并记录事件
	evt := leaseActivated{
		BaseEvent: events.NewBaseEvent("lease.activated", l.ID(), l.Version()),
		RoomID:    l.RoomID,
	}
	l.RecordEvent(evt)
}

// Checkout 退租
func (l *Lease) Checkout() {
	l.Status = LeaseStatusCheckout
	l.UpdatedAt = time.Now()
	// 创建并记录事件
	evt := leaseCheckout{
		BaseEvent: events.NewBaseEvent("lease.checkout", l.ID(), l.Version()),
		RoomID:    l.RoomID,
	}
	l.RecordEvent(evt)
}

// Expire 标记租约为过期状态
func (l *Lease) Expire() {
	l.Status = LeaseStatusExpired
	l.UpdatedAt = time.Now()
	// 创建并记录事件
	evt := leaseExpired{
		BaseEvent: events.NewBaseEvent("lease.expired", l.ID(), l.Version()),
		RoomID:    l.RoomID,
	}
	l.RecordEvent(evt)
}

// Renew 续租
func (l *Lease) Renew(newEndDate time.Time) {
	l.EndDate = newEndDate
	l.UpdatedAt = time.Now()
	// 创建并记录事件
	evt := leaseRenewed{
		BaseEvent:  events.NewBaseEvent("lease.renewed", l.ID(), l.Version()),
		NewEndDate: l.EndDate.Format("2006-01-02"),
	}
	l.RecordEvent(evt)
}
