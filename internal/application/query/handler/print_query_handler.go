package handler

import (
	"fmt"

	"github.com/zouhang1992/ddd_domain/internal/application/query"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// PrintQueryHandler 打印查询处理器
type PrintQueryHandler struct {
	billRepo  repository.BillRepository
	leaseRepo repository.LeaseRepository
	printRepo repository.PrintJobRepository
}

// NewPrintQueryHandler 创建打印查询处理器
func NewPrintQueryHandler(billRepo repository.BillRepository, leaseRepo repository.LeaseRepository, printRepo repository.PrintJobRepository) *PrintQueryHandler {
	return &PrintQueryHandler{billRepo: billRepo, leaseRepo: leaseRepo, printRepo: printRepo}
}

// HandleGetPrintJob 处理获取打印作业查询
func (h *PrintQueryHandler) HandleGetPrintJob(q query.Query) (any, error) {
	getQuery, ok := q.(query.GetPrintJobQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	job, err := h.printRepo.FindByID(getQuery.JobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, model.ErrNotFound
	}

	return job, nil
}

// HandleListPrintJobs 处理列出打印作业查询
func (h *PrintQueryHandler) HandleListPrintJobs(q query.Query) (any, error) {
	listQuery, ok := q.(query.ListPrintJobsQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	// 构建查询条件
	criteria := repository.PrintJobCriteria{
		Status:     listQuery.Status,
		StartTime:  listQuery.StartDate,
		EndTime:    listQuery.EndDate,
	}

	// 设置默认分页大小
	limit := listQuery.Limit
	if limit <= 0 {
		limit = 10 // 默认返回10条
	}

	// 查询数据
	jobs, err := h.printRepo.FindByCriteria(criteria, listQuery.Offset, limit)
	if err != nil {
		return nil, err
	}

	// 获取总数
	total, err := h.printRepo.CountByCriteria(criteria)
	if err != nil {
		return nil, err
	}

	// 计算页码
	page := 1
	if listQuery.Offset > 0 && limit > 0 {
		page = (listQuery.Offset / limit) + 1
	}

	result := &query.PrintJobsQueryResult{
		Items: jobs,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}

// HandleGetPrintContent 处理获取打印内容查询
func (h *PrintQueryHandler) HandleGetPrintContent(q query.Query) (any, error) {
	getQuery, ok := q.(query.GetPrintContentQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	bill, err := h.billRepo.FindByID(getQuery.BillID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, model.ErrNotFound
	}

	lease, err := h.leaseRepo.FindByID(bill.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, model.ErrNotFound
	}

	// 直接实现打印收据功能（简化版）
	content := []byte(fmt.Sprintf(`
{\rtf1\ansi\deff0{\fonttbl{\f0 Arial;}}
\pard\fs24\b 收据\b0\par
\pard\fs16 账单编号: %s\par
\pard\fs16 租约编号: %s\par
\pard\fs16 租客: %s\par
\pard\fs16 金额: %d 元\par
}`, bill.ID, lease.ID, lease.TenantName, bill.Amount))
	return content, nil
}
