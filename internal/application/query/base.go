package query

// Query 定义查询接口
type Query interface {
	// QueryName 返回查询名称
	QueryName() string
}

// BaseQuery 基础查询结构
type BaseQuery struct{}
