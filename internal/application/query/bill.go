package query

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"time"
)

// GetBillQuery 获取账单查询
type GetBillQuery struct {
	ID string
}

// QueryName 实现 Query 接口
func (q GetBillQuery) QueryName() string {
	return "get_bill"
}

// ListBillsQuery 列出账单查询
type ListBillsQuery struct {
	LeaseID string
	RoomID  string
	Month   string
}

// QueryName 实现 Query 接口
func (q ListBillsQuery) QueryName() string {
	return "list_bills"
}

// BillQueryResult 账单查询结果
type BillQueryResult struct {
	*model.Bill
}

// BillsQueryResult 账单列表查询结果
type BillsQueryResult struct {
	Items []*model.Bill
}

// IncomeReportQuery 收入报表查询
type IncomeReportQuery struct {
	Year  int
	Month time.Month
}

// QueryName 实现 Query 接口
func (q IncomeReportQuery) QueryName() string {
	return "income_report"
}

// IncomeReportQueryResult 收入报表查询结果
type IncomeReportQueryResult struct {
	Year          int
	Month         int
	RentTotal     int64
	WaterTotal    int64
	ElectricTotal int64
	OtherTotal    int64
	Total         int64
}
