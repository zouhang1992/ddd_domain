package command

import (
	appcommand "github.com/zouhang1992/ddd_domain/internal/application/command"
)

// Command 定义命令接口（使用 application 包的 Command 接口）
type Command = appcommand.Command

// CommandHandler 定义命令处理器接口
type CommandHandler interface {
	// Handle 处理命令
	Handle(cmd Command) (any, error)
}

// HandlerFunc 函数类型适配器
type HandlerFunc func(cmd Command) (any, error)

// Handle 实现 CommandHandler 接口
func (f HandlerFunc) Handle(cmd Command) (any, error) {
	return f(cmd)
}
