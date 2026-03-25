package query

import (
	"errors"
	"fmt"
)

// Bus 查询总线
type Bus struct {
	handlers   map[string]QueryHandler
	middleware []Middleware
}

// NewBus 创建查询总线
func NewBus() *Bus {
	return &Bus{
		handlers:   make(map[string]QueryHandler),
		middleware: make([]Middleware, 0),
	}
}

// Register 注册查询处理器
func (b *Bus) Register(queryName string, handler QueryHandler) {
	b.handlers[queryName] = handler
}

// Use 添加中间件
func (b *Bus) Use(m Middleware) {
	b.middleware = append(b.middleware, m)
}

// Dispatch 分发查询
func (b *Bus) Dispatch(q Query) (any, error) {
	// 验证查询
	if validator, ok := q.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrQueryValidation, err)
		}
	}

	handler, ok := b.handlers[q.QueryName()]
	if !ok {
		return nil, fmt.Errorf("no handler registered for query: %s", q.QueryName())
	}

	// 构建中间件链
	finalHandler := handler
	for i := len(b.middleware) - 1; i >= 0; i-- {
		finalHandler = b.middleware[i](finalHandler)
	}

	return finalHandler.Handle(q)
}

// Middleware 查询中间件类型
type Middleware func(next QueryHandler) QueryHandler

// ErrQueryValidation 查询验证错误
var ErrQueryValidation = errors.New("query validation failed")
