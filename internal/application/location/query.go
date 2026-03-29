package location

import locationmodel "github.com/zouhang1992/ddd_domain/internal/domain/location/model"

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
	// 查询条件
	ShortName string // 简称（模糊搜索）
	Detail    string // 详情（模糊搜索）
	// 分页参数
	Offset    int    // 偏移量
	Limit     int    // 每页数量
}

// QueryName 实现 Query 接口
func (q ListLocationsQuery) QueryName() string {
	return "list_locations"
}

// LocationQueryResult 位置查询结果
type LocationQueryResult struct {
	*locationmodel.Location
}

// LocationsQueryResult 位置列表查询结果
type LocationsQueryResult struct {
	Items []*locationmodel.Location `json:"items"`
	Total int                       `json:"total"`
	Page  int                       `json:"page"`
	Limit int                       `json:"limit"`
}
