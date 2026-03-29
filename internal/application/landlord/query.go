package landlord

import landlordmodel "github.com/zouhang1992/ddd_domain/internal/domain/landlord/model"

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
	// 查询条件
	Name   string // 姓名（模糊搜索）
	Phone  string // 电话（模糊搜索）
	// 分页参数
	Offset int    // 偏移量
	Limit  int    // 每页数量
}

// QueryName 实现 Query 接口
func (q ListLandlordsQuery) QueryName() string {
	return "list_landlords"
}

// LandlordQueryResult 房东查询结果
type LandlordQueryResult struct {
	*landlordmodel.Landlord
}

// LandlordsQueryResult 房东列表查询结果
type LandlordsQueryResult struct {
	Items []*landlordmodel.Landlord `json:"items"`
	Total int                        `json:"total"`
	Page  int                        `json:"page"`
	Limit int                        `json:"limit"`
}
