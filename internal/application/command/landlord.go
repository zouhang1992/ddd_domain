package command

import (
	"errors"
)

// CreateLandlordCommand 创建房东命令
type CreateLandlordCommand struct {
	Name  string
	Phone string
	Note  string
}

// CommandName 实现 Command 接口
func (c CreateLandlordCommand) CommandName() string {
	return "create_landlord"
}

// Validate 验证命令
func (c CreateLandlordCommand) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// UpdateLandlordCommand 更新房东命令
type UpdateLandlordCommand struct {
	ID    string
	Name  string
	Phone string
	Note  string
}

// CommandName 实现 Command 接口
func (c UpdateLandlordCommand) CommandName() string {
	return "update_landlord"
}

// Validate 验证命令
func (c UpdateLandlordCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	if c.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// DeleteLandlordCommand 删除房东命令
type DeleteLandlordCommand struct {
	ID string
}

// CommandName 实现 Command 接口
func (c DeleteLandlordCommand) CommandName() string {
	return "delete_landlord"
}

// Validate 验证命令
func (c DeleteLandlordCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	return nil
}
