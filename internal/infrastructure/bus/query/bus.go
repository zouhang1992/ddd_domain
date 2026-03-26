package query

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

// Bus 查询总线
type Bus struct {
	handlers   map[string]QueryHandler
	middleware []Middleware
	log        *zap.Logger
}

// NewBus 创建查询总线
func NewBus(logger *zap.Logger) *Bus {
	return &Bus{
		handlers:   make(map[string]QueryHandler),
		middleware: make([]Middleware, 0),
		log:        logger,
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
	queryName := q.QueryName()
	b.log.Info("Dispatching query", zap.String("query", queryName))

	// 验证查询
	if validator, ok := q.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			b.log.Warn("Query validation failed",
				zap.String("query", queryName),
				zap.Error(err))
			return nil, fmt.Errorf("%w: %v", ErrQueryValidation, err)
		}
	}

	handler, ok := b.handlers[queryName]
	if !ok {
		b.log.Error("No handler registered for query",
			zap.String("query", queryName))
		return nil, fmt.Errorf("no handler registered for query: %s", queryName)
	}

	// 构建中间件链
	finalHandler := handler
	for i := len(b.middleware) - 1; i >= 0; i-- {
		finalHandler = b.middleware[i](finalHandler)
	}

	result, err := finalHandler.Handle(q)
	if err != nil {
		b.log.Error("Query execution failed",
			zap.String("query", queryName),
			zap.Error(err))
	} else {
		b.log.Info("Query executed successfully",
			zap.String("query", queryName))
	}

	return result, err
}

// Middleware 查询中间件类型
type Middleware func(next QueryHandler) QueryHandler

// ErrQueryValidation 查询验证错误
var ErrQueryValidation = errors.New("query validation failed")
