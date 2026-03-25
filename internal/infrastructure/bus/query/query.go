package query

import (
	appquery "github.com/zouhang1992/ddd_domain/internal/application/query"
)

// Query 定义查询接口（使用 application 包的 Query 接口）
type Query = appquery.Query

// QueryHandler 定义查询处理器接口
type QueryHandler interface {
	// Handle 处理查询
	Handle(query Query) (any, error)
}

// HandlerFunc 函数类型适配器
type HandlerFunc func(query Query) (any, error)

// Handle 实现 QueryHandler 接口
func (f HandlerFunc) Handle(query Query) (any, error) {
	return f(query)
}
