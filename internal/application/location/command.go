package location

import "errors"

// CreateLocationCommand 创建位置命令
type CreateLocationCommand struct {
	ShortName string
	Detail    string
}

// CommandName 实现 Command 接口
func (c CreateLocationCommand) CommandName() string {
	return "create_location"
}

// Validate 验证命令
func (c CreateLocationCommand) Validate() error {
	if c.ShortName == "" {
		return errors.New("short name is required")
	}
	return nil
}

// UpdateLocationCommand 更新位置命令
type UpdateLocationCommand struct {
	ID        string
	ShortName string
	Detail    string
}

// CommandName 实现 Command 接口
func (c UpdateLocationCommand) CommandName() string {
	return "update_location"
}

// Validate 验证命令
func (c UpdateLocationCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	if c.ShortName == "" && c.Detail == "" {
		return errors.New("at least one field is required")
	}
	return nil
}

// DeleteLocationCommand 删除位置命令
type DeleteLocationCommand struct {
	ID string
}

// CommandName 实现 Command 接口
func (c DeleteLocationCommand) CommandName() string {
	return "delete_location"
}

// Validate 验证命令
func (c DeleteLocationCommand) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}
	return nil
}
