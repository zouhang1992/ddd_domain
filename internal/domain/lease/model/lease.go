package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
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

	// 暂时注释掉，先解决导入循环问题
	// 记录创建事件
	// lease.RecordEvent(events.NewLeaseCreated(lease.ID(), lease.Version(), lease.RoomID, lease.LandlordID, lease.TenantName))
	return lease
}

// Activate 激活租约
func (l *Lease) Activate() {
	l.Status = LeaseStatusActive
	l.UpdatedAt = time.Now()
	// 暂时注释掉，先解决导入循环问题
	// l.RecordEvent(events.NewLeaseActivated(l.ID(), l.Version(), l.RoomID))
}

// Checkout 退租
func (l *Lease) Checkout() {
	l.Status = LeaseStatusCheckout
	l.UpdatedAt = time.Now()
	// 暂时注释掉，先解决导入循环问题
	// l.RecordEvent(events.NewLeaseCheckout(l.ID(), l.Version(), l.RoomID))
}

// Expire 标记租约为过期状态
func (l *Lease) Expire() {
	l.Status = LeaseStatusExpired
	l.UpdatedAt = time.Now()
	// 暂时注释掉，先解决导入循环问题
	// l.RecordEvent(events.NewLeaseExpired(l.ID(), l.Version(), l.RoomID))
}

// Renew 续租
func (l *Lease) Renew(newEndDate time.Time) {
	l.EndDate = newEndDate
	l.UpdatedAt = time.Now()
	// 暂时注释掉，先解决导入循环问题
	// l.RecordEvent(events.NewLeaseRenewed(l.ID(), l.Version(), l.EndDate.Format("2006-01-02")))
}
