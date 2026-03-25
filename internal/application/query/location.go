package query

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// GetLocationQuery 获取位置查询
type GetLocationQuery struct {
	ID string
}

// QueryName 实现 Query 接口
func (q GetLocationQuery) QueryName() string {
	return "get_location"
}

// ListLocationsQuery 列出所有位置查询
type ListLocationsQuery struct {
}

// QueryName 实现 Query 接口
func (q ListLocationsQuery) QueryName() string {
	return "list_locations"
}

// LocationQueryResult 位置查询结果
type LocationQueryResult struct {
	*model.Location
}

// LocationsQueryResult 位置列表查询结果
type LocationsQueryResult struct {
	Items []*model.Location
}
