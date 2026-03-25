package query

import "github.com/zouhang1992/ddd_domain/internal/domain/model"

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
	Status string
	RoomID string
}

// QueryName 实现 Query 接口
func (q ListLeasesQuery) QueryName() string {
	return "list_leases"
}

// LeaseQueryResult 租约查询结果
type LeaseQueryResult struct {
	*model.Lease
}

// LeasesQueryResult 租约列表查询结果
type LeasesQueryResult struct {
	Items []*model.Lease
}
