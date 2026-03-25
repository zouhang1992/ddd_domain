package command

import (
	"errors"
	"fmt"
)

// Bus 命令总线
type Bus struct {
	handlers   map[string]CommandHandler
	middleware []Middleware
}

// NewBus 创建命令总线
func NewBus() *Bus {
	return &Bus{
		handlers:   make(map[string]CommandHandler),
		middleware: make([]Middleware, 0),
	}
}

// Register 注册命令处理器
func (b *Bus) Register(cmdName string, handler CommandHandler) {
	b.handlers[cmdName] = handler
}

// Use 添加中间件
func (b *Bus) Use(m Middleware) {
	b.middleware = append(b.middleware, m)
}

// Dispatch 分发命令
func (b *Bus) Dispatch(cmd Command) (any, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCommandValidation, err)
	}

	handler, ok := b.handlers[cmd.CommandName()]
	if !ok {
		return nil, fmt.Errorf("no handler registered for command: %s", cmd.CommandName())
	}

	// 构建中间件链
	finalHandler := handler
	for i := len(b.middleware) - 1; i >= 0; i-- {
		finalHandler = b.middleware[i](finalHandler)
	}

	return finalHandler.Handle(cmd)
}

// Middleware 命令中间件类型
type Middleware func(next CommandHandler) CommandHandler

// ErrCommandValidation 命令验证错误
var ErrCommandValidation = errors.New("command validation failed")
