package query

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
	Status     string
	StartDate  string
	EndDate    string
}

// QueryName 实现 Query 接口
func (q ListPrintJobsQuery) QueryName() string {
	return "list_print_jobs"
}

// GetPrintContentQuery 获取打印内容查询
type GetPrintContentQuery struct {
	BillID string
}

// QueryName 实现 Query 接口
func (q GetPrintContentQuery) QueryName() string {
	return "get_print_content"
}
