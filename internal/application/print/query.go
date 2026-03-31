package print

import (
	"time"

	printmodel "github.com/zouhang1992/ddd_domain/internal/domain/print/model"
)

// GetPrintJobQuery 获取打印作业查询
type GetPrintJobQuery struct {
	JobID string
}

// QueryName 实现 Query 接口
func (q GetPrintJobQuery) QueryName() string {
	return "get_print_job"
}

// ListPrintJobsQuery 列出打印作业查询
type ListPrintJobsQuery struct {
	// 查询条件
	Status    string     // 状态
	Type      string     // 类型
	StartDate *time.Time // 开始日期范围
	EndDate   *time.Time // 结束日期范围
	// 分页参数
	Offset int // 偏移量
	Limit  int // 每页数量
}

// QueryName 实现 Query 接口
func (q ListPrintJobsQuery) QueryName() string {
	return "list_print_jobs"
}

// PrintJobResult 打印作业查询结果
type PrintJobResult struct {
	ID           string                     `json:"id"`
	Type         printmodel.PrintJobType   `json:"type"`
	Status       printmodel.PrintJobStatus `json:"status"`
	ReferenceID  string                     `json:"reference_id"`
	TenantName   string                     `json:"tenant_name"`
	TenantPhone  string                     `json:"tenant_phone"`
	RoomID       string                     `json:"room_id"`
	RoomNumber   string                     `json:"room_number"`
	Address      string                     `json:"address"`
	LandlordName string                     `json:"landlord_name"`
	Amount       int64                      `json:"amount"`
	ErrorMsg     string                     `json:"error_msg,omitempty"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at"`
	CompletedAt  *time.Time                 `json:"completed_at,omitempty"`
	TypeText     string                     `json:"type_text"`
	StatusText   string                     `json:"status_text"`
	AmountYuan   string                     `json:"amount_yuan"`
}

// PrintJobsQueryResult 打印作业列表查询结果
type PrintJobsQueryResult struct {
	Items []*PrintJobResult `json:"items"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

// GetPrintContentQuery 获取打印内容查询
type GetPrintContentQuery struct {
	BillID string
}

// QueryName 实现 Query 接口
func (q GetPrintContentQuery) QueryName() string {
	return "get_print_content"
}
