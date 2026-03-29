package print

import "time"

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
	StartDate *time.Time // 开始日期范围
	EndDate   *time.Time // 结束日期范围
	// 分页参数
	Offset    int        // 偏移量
	Limit     int        // 每页数量
}

// QueryName 实现 Query 接口
func (q ListPrintJobsQuery) QueryName() string {
	return "list_print_jobs"
}

// PrintJobsQueryResult 打印作业查询结果
type PrintJobsQueryResult struct {
	Items []interface{} `json:"items"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

// GetPrintContentQuery 获取打印内容查询
type GetPrintContentQuery struct {
	BillID string
}

// QueryName 实现 Query 接口
func (q GetPrintContentQuery) QueryName() string {
	return "get_print_content"
}
