package common

// Command 定义命令接口
type Command interface {
	// CommandName 返回命令名称
	CommandName() string
	// Validate 验证命令
	Validate() error
}
