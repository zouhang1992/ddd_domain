package model

import (
	"time"
)

// LeaseStatus 租约状态
type LeaseStatus string

const (
	LeaseStatusPending  LeaseStatus = "pending"  // 待生效
	LeaseStatusActive   LeaseStatus = "active"   // 生效中
	LeaseStatusExpired  LeaseStatus = "expired"  // 已过期
	LeaseStatusCheckout LeaseStatus = "checkout" // 已退租
)

// Lease 租约领域模型
type Lease struct {
	ID             string       `json:"id"`
	RoomID         string       `json:"roomId"`
	LandlordID     string       `json:"landlordId"`
	TenantName     string       `json:"tenantName"`
	TenantPhone    string       `json:"tenantPhone"`
	StartDate      time.Time    `json:"startDate"`
	EndDate        time.Time    `json:"endDate"`
	RentAmount     int64        `json:"rentAmount"` // 租金（分）
	DepositAmount  int64        `json:"depositAmount"` // 押金金额（分）
	Status         LeaseStatus  `json:"status"`
	Note           string       `json:"note"`
	LastChargeAt   *time.Time   `json:"lastChargeAt"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}

// NewLease 创建新租约
func NewLease(id, roomID, landlordID, tenantName, tenantPhone string,
	startDate, endDate time.Time, rentAmount, depositAmount int64, note string) *Lease {
	now := time.Now()
	return &Lease{
		ID:            id,
		RoomID:        roomID,
		LandlordID:    landlordID,
		TenantName:    tenantName,
		TenantPhone:   tenantPhone,
		StartDate:     startDate,
		EndDate:       endDate,
		RentAmount:    rentAmount,
		DepositAmount: depositAmount,
		Status:        LeaseStatusPending,
		Note:          note,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// Update 更新租约信息
func (l *Lease) Update(tenantName, tenantPhone string, startDate, endDate time.Time,
	rentAmount int64, note string) {
	l.TenantName = tenantName
	l.TenantPhone = tenantPhone
	l.StartDate = startDate
	l.EndDate = endDate
	l.RentAmount = rentAmount
	l.Note = note
	l.UpdatedAt = time.Now()
}

// SetDepositAmount 设置押金金额
func (l *Lease) SetDepositAmount(amount int64) {
	l.DepositAmount = amount
	l.UpdatedAt = time.Now()
}

// Activate 激活租约
func (l *Lease) Activate() {
	l.Status = LeaseStatusActive
	l.UpdatedAt = time.Now()
}

// Checkout 退租
func (l *Lease) Checkout() {
	l.Status = LeaseStatusCheckout
	l.UpdatedAt = time.Now()
}

// UpdateLastChargeAt 更新最后收账时间
func (l *Lease) UpdateLastChargeAt(chargeAt time.Time) {
	l.LastChargeAt = &chargeAt
	l.UpdatedAt = time.Now()
}

// IsActive 检查租约是否生效中
func (l *Lease) IsActive() bool {
	return l.Status == LeaseStatusActive
}

// Renew 续租
func (l *Lease) Renew(newStartDate, newEndDate time.Time, newRentAmount int64, note string) error {
	if l.Status != LeaseStatusActive && l.Status != LeaseStatusExpired {
		return ErrInvalidState
	}
	l.StartDate = newStartDate
	l.EndDate = newEndDate
	l.RentAmount = newRentAmount
	l.Note = note
	l.Status = LeaseStatusActive
	l.UpdatedAt = time.Now()
	return nil
}

// CanDelete 检查租约是否可以删除
func (l *Lease) CanDelete(hasBills, hasDeposit bool) bool {
	if hasBills || hasDeposit {
		return false
	}
	if l.IsActive() {
		return false
	}
	return true
}
