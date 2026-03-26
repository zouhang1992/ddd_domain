package command

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

// Bus 命令总线
type Bus struct {
	handlers   map[string]CommandHandler
	middleware []Middleware
	log        *zap.Logger
}

// NewBus 创建命令总线
func NewBus(logger *zap.Logger) *Bus {
	return &Bus{
		handlers:   make(map[string]CommandHandler),
		middleware: make([]Middleware, 0),
		log:        logger,
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
	cmdName := cmd.CommandName()
	b.log.Info("Dispatching command", zap.String("command", cmdName))

	// 验证命令
	if err := cmd.Validate(); err != nil {
		b.log.Warn("Command validation failed",
			zap.String("command", cmdName),
			zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrCommandValidation, err)
	}

	handler, ok := b.handlers[cmdName]
	if !ok {
		b.log.Error("No handler registered for command",
			zap.String("command", cmdName))
		return nil, fmt.Errorf("no handler registered for command: %s", cmdName)
	}

	// 构建中间件链
	finalHandler := handler
	for i := len(b.middleware) - 1; i >= 0; i-- {
		finalHandler = b.middleware[i](finalHandler)
	}

	result, err := finalHandler.Handle(cmd)
	if err != nil {
		b.log.Error("Command execution failed",
			zap.String("command", cmdName),
			zap.Error(err))
	} else {
		b.log.Info("Command executed successfully",
			zap.String("command", cmdName))
	}

	return result, err
}

// Middleware 命令中间件类型
type Middleware func(next CommandHandler) CommandHandler

// ErrCommandValidation 命令验证错误
var ErrCommandValidation = errors.New("command validation failed")
