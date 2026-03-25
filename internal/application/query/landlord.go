package query

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

// GetLandlordQuery 获取房东查询
type GetLandlordQuery struct {
	ID string
}

// QueryName 实现 Query 接口
func (q GetLandlordQuery) QueryName() string {
	return "get_landlord"
}

// ListLandlordsQuery 列出租东查询
type ListLandlordsQuery struct {
}

// QueryName 实现 Query 接口
func (q ListLandlordsQuery) QueryName() string {
	return "list_landlords"
}

// LandlordQueryResult 房东查询结果
type LandlordQueryResult struct {
	*model.Landlord
}

// LandlordsQueryResult 房东列表查询结果
type LandlordsQueryResult struct {
	Items []*model.Landlord
}
