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
	// 查询条件
	Type        string     // 账单类型
	Status      string     // 账单状态
	LeaseID     string     // 租约ID
	RoomID      string     // 房间ID
	Month       string     // 月份 (格式: "2006-01")
	MinAmount   int64      // 最小金额（分）
	MaxAmount   int64      // 最大金额（分）
	StartDate   *time.Time // 开始日期范围
	EndDate     *time.Time // 结束日期范围
	// 分页参数
	Offset      int        // 偏移量
	Limit       int        // 每页数量
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
	Items []*model.Bill `json:"items"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
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
