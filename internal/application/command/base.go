package command

// Command 定义命令接口
type Command interface {
	// CommandName 返回命令名称
	CommandName() string
	// Validate 验证命令
	Validate() error
}

// BaseCommand 基础命令实现
type BaseCommand struct {
	name string
}

// NewBaseCommand 创建基础命令
func NewBaseCommand(name string) BaseCommand {
	return BaseCommand{
		name: name,
	}
}

// CommandName 实现 Command 接口
func (c BaseCommand) CommandName() string {
	return c.name
}
