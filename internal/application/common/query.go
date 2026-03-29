package common

// Query 定义查询接口
type Query interface {
	// QueryName 返回查询名称
	QueryName() string
}
