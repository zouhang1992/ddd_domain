package lease

import (
	leasemodel "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
	"time"
)

// GetLeaseQuery 获取租约查询
type GetLeaseQuery struct {
	ID string
}

// QueryName 实现 Query 接口
func (q GetLeaseQuery) QueryName() string {
	return "get_lease"
}

// ListLeasesQuery 列出租约查询
type ListLeasesQuery struct {
	// 查询条件
	TenantName string     // 租户姓名（模糊搜索）
	TenantPhone string    // 租户电话（模糊搜索）
	Status      string     // 状态
	LocationID  string     // 位置ID
	RoomID      string     // 房间ID
	StartDate   *time.Time // 开始日期范围
	EndDate     *time.Time // 结束日期范围
	// 分页参数
	Offset      int        // 偏移量
	Limit       int        // 每页数量
}

// QueryName 实现 Query 接口
func (q ListLeasesQuery) QueryName() string {
	return "list_leases"
}

// LeaseQueryResult 租约查询结果
type LeaseQueryResult struct {
	*leasemodel.Lease
}

// LeasesQueryResult 租约列表查询结果
type LeasesQueryResult struct {
	Items []*leasemodel.Lease `json:"items"`
	Total int                  `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}
