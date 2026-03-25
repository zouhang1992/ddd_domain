package command

import (
	"errors"
	"time"
)

// CreateLeaseCommand 创建租约命令
type CreateLeaseCommand struct {
	RoomID        string
	LandlordID    string
	TenantName    string
	TenantPhone   string
	StartDate     time.Time
	EndDate       time.Time
	RentAmount    int64
	Note          string
	DepositAmount int64
	DepositNote   string
}

// CommandName 实现 Command 接口
func (c CreateLeaseCommand) CommandName() string {
	return "create_lease"
}

// Validate 验证命令
func (c CreateLeaseCommand) Validate() error {
	if c.RoomID == "" {
		return errors.New("room_id is required")
	}
	if c.TenantName == "" {
		return errors.New("tenant_name is required")
	}
	if c.StartDate.IsZero() {
		return errors.New("start_date is required")
	}
	if c.EndDate.IsZero() {
		return errors.New("end_date is required")
	}
	if c.StartDate.After(c.EndDate) {
		return errors.New("start_date must be before end_date")
	}
	if c.RentAmount < 0 {
		return errors.New("rent_amount cannot be negative")
	}
	return nil
}

// UpdateLeaseCommand 更新租约命令
type UpdateLeaseCommand struct {
	ID          string
	TenantName  string
	TenantPhone string
	StartDate   time.Time
	EndDate     time.Time
	RentAmount  int64
	Note        string
}

// CommandName 实现 Command 接口
func (c UpdateLeaseCommand) CommandName() string {
	return "update_lease"
}

// Validate 验证命令
func (c UpdateLeaseCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	if c.TenantName == "" {
		return errors.New("tenant_name is required")
	}
	if c.StartDate.IsZero() {
		return errors.New("start_date is required")
	}
	if c.EndDate.IsZero() {
		return errors.New("end_date is required")
	}
	if c.StartDate.After(c.EndDate) {
		return errors.New("start_date must be before end_date")
	}
	if c.RentAmount < 0 {
		return errors.New("rent_amount cannot be negative")
	}
	return nil
}

// DeleteLeaseCommand 删除租约命令
type DeleteLeaseCommand struct {
	ID string
}

// CommandName 实现 Command 接口
func (c DeleteLeaseCommand) CommandName() string {
	return "delete_lease"
}

// Validate 验证命令
func (c DeleteLeaseCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	return nil
}

// RenewLeaseCommand 续租命令
type RenewLeaseCommand struct {
	ID            string
	NewStartDate  time.Time
	NewEndDate    time.Time
	NewRentAmount int64
	Note          string
}

// CommandName 实现 Command 接口
func (c RenewLeaseCommand) CommandName() string {
	return "renew_lease"
}

// Validate 验证命令
func (c RenewLeaseCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	if c.NewStartDate.IsZero() {
		return errors.New("new_start_date is required")
	}
	if c.NewEndDate.IsZero() {
		return errors.New("new_end_date is required")
	}
	if c.NewStartDate.After(c.NewEndDate) {
		return errors.New("new_start_date must be before new_end_date")
	}
	if c.NewRentAmount < 0 {
		return errors.New("new_rent_amount cannot be negative")
	}
	return nil
}

// CheckoutLeaseCommand 退租命令
type CheckoutLeaseCommand struct {
	ID string
}

// CommandName 实现 Command 接口
func (c CheckoutLeaseCommand) CommandName() string {
	return "checkout_lease"
}

// Validate 验证命令
func (c CheckoutLeaseCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	return nil
}

// ActivateLeaseCommand 租约生效命令
type ActivateLeaseCommand struct {
	ID string
}

// CommandName 实现 Command 接口
func (c ActivateLeaseCommand) CommandName() string {
	return "activate_lease"
}

// Validate 验证命令
func (c ActivateLeaseCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	return nil
}
