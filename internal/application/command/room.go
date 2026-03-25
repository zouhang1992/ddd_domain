package command

import "errors"

// CreateRoomCommand 创建房间命令
type CreateRoomCommand struct {
	LocationID string
	RoomNumber string
	Tags       []string
}

// CommandName 实现 Command 接口
func (c CreateRoomCommand) CommandName() string {
	return "create_room"
}

// Validate 验证命令
func (c CreateRoomCommand) Validate() error {
	if c.LocationID == "" {
		return errors.New("location id is required")
	}
	if c.RoomNumber == "" {
		return errors.New("room number is required")
	}
	return nil
}

// UpdateRoomCommand 更新房间命令
type UpdateRoomCommand struct {
	ID         string
	LocationID string
	RoomNumber string
	Tags       []string
}

// CommandName 实现 Command 接口
func (c UpdateRoomCommand) CommandName() string {
	return "update_room"
}

// Validate 验证命令
func (c UpdateRoomCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	if c.LocationID == "" && c.RoomNumber == "" && len(c.Tags) == 0 {
		return errors.New("at least one field is required")
	}
	return nil
}

// DeleteRoomCommand 删除房间命令
type DeleteRoomCommand struct {
	ID string
}

// CommandName 实现 Command 接口
func (c DeleteRoomCommand) CommandName() string {
	return "delete_room"
}

// Validate 验证命令
func (c DeleteRoomCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	return nil
}
