package command

import (
	"errors"
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/model"
)

// CreateBillCommand 创建账单命令
type CreateBillCommand struct {
	LeaseID        string
	Type           model.BillType
	Amount         int64
	RentAmount     int64
	WaterAmount    int64
	ElectricAmount int64
	OtherAmount    int64
	PaidAt         *time.Time
	Note           string
}

// CommandName 实现 Command 接口
func (c CreateBillCommand) CommandName() string {
	return "create_bill"
}

// Validate 验证命令
func (c CreateBillCommand) Validate() error {
	if c.LeaseID == "" {
		return errors.New("lease_id is required")
	}
	if c.Type == "" {
		return errors.New("type is required")
	}
	if c.Type != model.BillTypeCharge && c.Type != model.BillTypeCheckout {
		return errors.New("invalid bill type")
	}
	if c.Amount < 0 {
		return errors.New("amount cannot be negative")
	}
	if c.RentAmount < 0 {
		return errors.New("rent_amount cannot be negative")
	}
	if c.WaterAmount < 0 {
		return errors.New("water_amount cannot be negative")
	}
	if c.ElectricAmount < 0 {
		return errors.New("electric_amount cannot be negative")
	}
	if c.OtherAmount < 0 {
		return errors.New("other_amount cannot be negative")
	}
	return nil
}

// UpdateBillCommand 更新账单命令
type UpdateBillCommand struct {
	ID             string
	Amount         int64
	RentAmount     int64
	WaterAmount    int64
	ElectricAmount int64
	OtherAmount    int64
	PaidAt         *time.Time
	Note           string
}

// CommandName 实现 Command 接口
func (c UpdateBillCommand) CommandName() string {
	return "update_bill"
}

// Validate 验证命令
func (c UpdateBillCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	if c.Amount < 0 {
		return errors.New("amount cannot be negative")
	}
	if c.RentAmount < 0 {
		return errors.New("rent_amount cannot be negative")
	}
	if c.WaterAmount < 0 {
		return errors.New("water_amount cannot be negative")
	}
	if c.ElectricAmount < 0 {
		return errors.New("electric_amount cannot be negative")
	}
	if c.OtherAmount < 0 {
		return errors.New("other_amount cannot be negative")
	}
	return nil
}

// DeleteBillCommand 删除账单命令
type DeleteBillCommand struct {
	ID string
}

// CommandName 实现 Command 接口
func (c DeleteBillCommand) CommandName() string {
	return "delete_bill"
}

// Validate 验证命令
func (c DeleteBillCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	return nil
}

// ConfirmBillArrivalCommand 确认账单到账命令
type ConfirmBillArrivalCommand struct {
	ID string
	PaidAt time.Time
}

// CommandName 实现 Command 接口
func (c ConfirmBillArrivalCommand) CommandName() string {
	return "confirm_bill_arrival"
}

// Validate 验证命令
func (c ConfirmBillArrivalCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	return nil
}
