package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// ==================== 房东相关事件 ====================

// LandlordCreated 房东创建事件
type LandlordCreated struct {
	event.BaseEvent
	LandlordID string
	Name       string
	Phone      string
}

// NewLandlordCreated 创建房东创建事件
func NewLandlordCreated(landlord *Landlord) *LandlordCreated {
	return &LandlordCreated{
		BaseEvent:  event.NewBaseEvent("landlord.created"),
		LandlordID: landlord.ID,
		Name:       landlord.Name,
		Phone:      landlord.Phone,
	}
}

// LandlordUpdated 房东更新事件
type LandlordUpdated struct {
	event.BaseEvent
	LandlordID string
	Name       string
	Phone      string
}

// NewLandlordUpdated 创建房东更新事件
func NewLandlordUpdated(landlord *Landlord) *LandlordUpdated {
	return &LandlordUpdated{
		BaseEvent:  event.NewBaseEvent("landlord.updated"),
		LandlordID: landlord.ID,
		Name:       landlord.Name,
		Phone:      landlord.Phone,
	}
}

// LandlordDeleted 房东删除事件
type LandlordDeleted struct {
	event.BaseEvent
	LandlordID string
}

// NewLandlordDeleted 创建房东删除事件
func NewLandlordDeleted(landlordID string) *LandlordDeleted {
	return &LandlordDeleted{
		BaseEvent:  event.NewBaseEvent("landlord.deleted"),
		LandlordID: landlordID,
	}
}

// ==================== 租约相关事件 ====================

// LeaseCreated 租约创建事件
type LeaseCreated struct {
	event.BaseEvent
	LeaseID    string
	RoomID     string
	LandlordID string
	TenantName string
}

// NewLeaseCreated 创建租约创建事件
func NewLeaseCreated(lease *Lease) *LeaseCreated {
	return &LeaseCreated{
		BaseEvent:  event.NewBaseEvent("lease.created"),
		LeaseID:    lease.ID,
		RoomID:     lease.RoomID,
		LandlordID: lease.LandlordID,
		TenantName: lease.TenantName,
	}
}

// LeaseUpdated 租约更新事件
type LeaseUpdated struct {
	event.BaseEvent
	LeaseID    string
	TenantName string
}

// NewLeaseUpdated 创建租约更新事件
func NewLeaseUpdated(lease *Lease) *LeaseUpdated {
	return &LeaseUpdated{
		BaseEvent:  event.NewBaseEvent("lease.updated"),
		LeaseID:    lease.ID,
		TenantName: lease.TenantName,
	}
}

// LeaseRenewed 租约续租事件
type LeaseRenewed struct {
	event.BaseEvent
	LeaseID    string
	NewEndDate string
}

// NewLeaseRenewed 创建租约续租事件
func NewLeaseRenewed(lease *Lease) *LeaseRenewed {
	return &LeaseRenewed{
		BaseEvent:  event.NewBaseEvent("lease.renewed"),
		LeaseID:    lease.ID,
		NewEndDate: lease.EndDate.Format("2006-01-02"),
	}
}

// LeaseCheckout 租约退租事件
type LeaseCheckout struct {
	event.BaseEvent
	LeaseID string
}

// NewLeaseCheckout 创建租约退租事件
func NewLeaseCheckout(lease *Lease) *LeaseCheckout {
	return &LeaseCheckout{
		BaseEvent: event.NewBaseEvent("lease.checkout"),
		LeaseID:   lease.ID,
	}
}

// LeaseActivated 租约生效事件
type LeaseActivated struct {
	event.BaseEvent
	LeaseID string
	RoomID  string
}

// NewLeaseActivated 创建租约生效事件
func NewLeaseActivated(lease *Lease) *LeaseActivated {
	return &LeaseActivated{
		BaseEvent: event.NewBaseEvent("lease.activated"),
		LeaseID:   lease.ID,
		RoomID:    lease.RoomID,
	}
}

// LeaseDeleted 租约删除事件
type LeaseDeleted struct {
	event.BaseEvent
	LeaseID string
}

// NewLeaseDeleted 创建租约删除事件
func NewLeaseDeleted(leaseID string) *LeaseDeleted {
	return &LeaseDeleted{
		BaseEvent: event.NewBaseEvent("lease.deleted"),
		LeaseID:   leaseID,
	}
}

// ==================== 账单相关事件 ====================

// BillCreated 账单创建事件
type BillCreated struct {
	event.BaseEvent
	BillID  string
	LeaseID string
	Type    string
	Amount  int64
}

// NewBillCreated 创建账单创建事件
func NewBillCreated(bill *Bill) *BillCreated {
	return &BillCreated{
		BaseEvent: event.NewBaseEvent("bill.created"),
		BillID:    bill.ID,
		LeaseID:   bill.LeaseID,
		Type:      string(bill.Type),
		Amount:    bill.Amount,
	}
}

// BillUpdated 账单更新事件
type BillUpdated struct {
	event.BaseEvent
	BillID string
	Amount int64
}

// NewBillUpdated 创建账单更新事件
func NewBillUpdated(bill *Bill) *BillUpdated {
	return &BillUpdated{
		BaseEvent: event.NewBaseEvent("bill.updated"),
		BillID:    bill.ID,
		Amount:    bill.Amount,
	}
}

// BillPaid 账单支付事件
type BillPaid struct {
	event.BaseEvent
	BillID string
	PaidAt string
}

// NewBillPaid 创建账单支付事件
func NewBillPaid(bill *Bill) *BillPaid {
	paidAt := ""
	if bill.PaidAt != nil {
		paidAt = bill.PaidAt.Format("2006-01-02")
	}
	return &BillPaid{
		BaseEvent: event.NewBaseEvent("bill.paid"),
		BillID:    bill.ID,
		PaidAt:    paidAt,
	}
}

// BillDeleted 账单删除事件
type BillDeleted struct {
	event.BaseEvent
	BillID string
}

// NewBillDeleted 创建账单删除事件
func NewBillDeleted(billID string) *BillDeleted {
	return &BillDeleted{
		BaseEvent: event.NewBaseEvent("bill.deleted"),
		BillID:    billID,
	}
}

// ==================== 位置相关事件 ====================

// LocationCreated 位置创建事件
type LocationCreated struct {
	event.BaseEvent
	LocationID string
	ShortName  string
	Detail     string
}

// NewLocationCreated 创建位置创建事件
func NewLocationCreated(location *Location) *LocationCreated {
	return &LocationCreated{
		BaseEvent:  event.NewBaseEvent("location.created"),
		LocationID: location.ID,
		ShortName:  location.ShortName,
		Detail:     location.Detail,
	}
}

// LocationUpdated 位置更新事件
type LocationUpdated struct {
	event.BaseEvent
	LocationID string
	ShortName  string
	Detail     string
}

// NewLocationUpdated 创建位置更新事件
func NewLocationUpdated(location *Location) *LocationUpdated {
	return &LocationUpdated{
		BaseEvent:  event.NewBaseEvent("location.updated"),
		LocationID: location.ID,
		ShortName:  location.ShortName,
		Detail:     location.Detail,
	}
}

// LocationDeleted 位置删除事件
type LocationDeleted struct {
	event.BaseEvent
	LocationID string
}

// NewLocationDeleted 创建位置删除事件
func NewLocationDeleted(locationID string) *LocationDeleted {
	return &LocationDeleted{
		BaseEvent:  event.NewBaseEvent("location.deleted"),
		LocationID: locationID,
	}
}

// ==================== 房间相关事件 ====================

// RoomCreated 房间创建事件
type RoomCreated struct {
	event.BaseEvent
	RoomID     string
	LocationID string
	RoomNumber string
	Tags       []string
}

// NewRoomCreated 创建房间创建事件
func NewRoomCreated(room *Room) *RoomCreated {
	return &RoomCreated{
		BaseEvent:  event.NewBaseEvent("room.created"),
		RoomID:     room.ID,
		LocationID: room.LocationID,
		RoomNumber: room.RoomNumber,
		Tags:       room.Tags,
	}
}

// RoomUpdated 房间更新事件
type RoomUpdated struct {
	event.BaseEvent
	RoomID     string
	LocationID string
	RoomNumber string
	Tags       []string
}

// NewRoomUpdated 创建房间更新事件
func NewRoomUpdated(room *Room) *RoomUpdated {
	return &RoomUpdated{
		BaseEvent:  event.NewBaseEvent("room.updated"),
		RoomID:     room.ID,
		LocationID: room.LocationID,
		RoomNumber: room.RoomNumber,
		Tags:       room.Tags,
	}
}

// RoomDeleted 房间删除事件
type RoomDeleted struct {
	event.BaseEvent
	RoomID string
}

// NewRoomDeleted 创建房间删除事件
func NewRoomDeleted(roomID string) *RoomDeleted {
	return &RoomDeleted{
		BaseEvent: event.NewBaseEvent("room.deleted"),
		RoomID:    roomID,
	}
}

// ==================== 打印相关事件 ====================

// BillPrinted 账单打印事件
type BillPrinted struct {
	event.BaseEvent
	JobID     string
	BillID    string
	PrintedAt string
	Content   []byte
}

// NewBillPrinted 创建账单打印事件
func NewBillPrinted(jobID, billID string, content []byte) *BillPrinted {
	return &BillPrinted{
		BaseEvent:  event.NewBaseEvent("bill.printed"),
		JobID:      jobID,
		BillID:     billID,
		PrintedAt:  time.Now().Format("2006-01-02 15:04:05"),
		Content:    content,
	}
}

// LeasePrinted 租约打印事件
type LeasePrinted struct {
	event.BaseEvent
	JobID     string
	LeaseID   string
	PrintedAt string
	Content   []byte
}

// NewLeasePrinted 创建租约打印事件
func NewLeasePrinted(jobID, leaseID string, content []byte) *LeasePrinted {
	return &LeasePrinted{
		BaseEvent:  event.NewBaseEvent("lease.printed"),
		JobID:      jobID,
		LeaseID:    leaseID,
		PrintedAt:  time.Now().Format("2006-01-02 15:04:05"),
		Content:    content,
	}
}

// InvoicePrinted 发票打印事件
type InvoicePrinted struct {
	event.BaseEvent
	JobID     string
	BillID    string
	PrintedAt string
	Content   []byte
}

// NewInvoicePrinted 创建发票打印事件
func NewInvoicePrinted(jobID, billID string, content []byte) *InvoicePrinted {
	return &InvoicePrinted{
		BaseEvent:  event.NewBaseEvent("invoice.printed"),
		JobID:      jobID,
		BillID:     billID,
		PrintedAt:  time.Now().Format("2006-01-02 15:04:05"),
		Content:    content,
	}
}

// PrintJobFailed 打印作业失败事件
type PrintJobFailed struct {
	event.BaseEvent
	JobID    string
	BillID   string
	LeaseID  string
	FailedAt string
	Error    string
}

// NewPrintJobFailed 创建打印作业失败事件
func NewPrintJobFailed(jobID, billID, leaseID, err string) *PrintJobFailed {
	return &PrintJobFailed{
		BaseEvent:  event.NewBaseEvent("print.failed"),
		JobID:      jobID,
		BillID:     billID,
		LeaseID:    leaseID,
		FailedAt:   time.Now().Format("2006-01-02 15:04:05"),
		Error:      err,
	}
}
