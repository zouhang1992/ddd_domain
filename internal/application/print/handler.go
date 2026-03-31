package print

import (
	"fmt"

	"github.com/zouhang1992/ddd_domain/internal/application/common"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	printmodel "github.com/zouhang1992/ddd_domain/internal/domain/print/model"
	printrepo "github.com/zouhang1992/ddd_domain/internal/domain/print/repository"
	printservice "github.com/zouhang1992/ddd_domain/internal/domain/print/service"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// CommandHandler 打印命令处理器
type CommandHandler struct {
	billRepo     billrepo.BillRepository
	leaseRepo    leaserepo.LeaseRepository
	eventBus     *event.Bus
	printService *printservice.PrintService
}

// NewCommandHandler 创建打印命令处理器
func NewCommandHandler(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository, eventBus *event.Bus, printService *printservice.PrintService) *CommandHandler {
	return &CommandHandler{billRepo: billRepo, leaseRepo: leaseRepo, eventBus: eventBus, printService: printService}
}

// HandlePrintBill 处理打印账单命令
func (h *CommandHandler) HandlePrintBill(cmd common.Command) (any, error) {
	printCmd, ok := cmd.(PrintBillCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := printCmd.Validate(); err != nil {
		return nil, err
	}

	return h.printService.CreateBillPrintJob(printCmd.BillID)
}

// HandlePrintLease 处理打印租约命令
func (h *CommandHandler) HandlePrintLease(cmd common.Command) (any, error) {
	printCmd, ok := cmd.(PrintLeaseCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := printCmd.Validate(); err != nil {
		return nil, err
	}

	return h.printService.CreateLeasePrintJob(printCmd.LeaseID)
}

// HandlePrintInvoice 处理打印发票命令
func (h *CommandHandler) HandlePrintInvoice(cmd common.Command) (any, error) {
	printCmd, ok := cmd.(PrintInvoiceCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := printCmd.Validate(); err != nil {
		return nil, err
	}

	return h.printService.CreateInvoicePrintJob(printCmd.BillID)
}

// QueryHandler 打印查询处理器
type QueryHandler struct {
	billRepo     billrepo.BillRepository
	leaseRepo    leaserepo.LeaseRepository
	printService *printservice.PrintService
	printJobRepo printrepo.PrintJobRepository
}

// NewQueryHandler 创建打印查询处理器
func NewQueryHandler(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository, printService *printservice.PrintService, printJobRepo printrepo.PrintJobRepository) *QueryHandler {
	return &QueryHandler{billRepo: billRepo, leaseRepo: leaseRepo, printService: printService, printJobRepo: printJobRepo}
}

// HandleGetPrintJob 处理获取打印作业查询
func (h *QueryHandler) HandleGetPrintJob(q common.Query) (any, error) {
	getQuery, ok := q.(GetPrintJobQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	job, err := h.printJobRepo.FindByID(getQuery.JobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, fmt.Errorf("print job not found")
	}

	return toPrintJobResult(job), nil
}

// HandleListPrintJobs 处理列出打印作业查询
func (h *QueryHandler) HandleListPrintJobs(q common.Query) (any, error) {
	listQuery, ok := q.(ListPrintJobsQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	var status printmodel.PrintJobStatus
	if listQuery.Status != "" {
		status = printmodel.PrintJobStatus(listQuery.Status)
	}

	var jobType printmodel.PrintJobType
	if listQuery.Type != "" {
		jobType = printmodel.PrintJobType(listQuery.Type)
	}

	offset := listQuery.Offset
	limit := listQuery.Limit
	if limit <= 0 {
		limit = 10
	}

	jobs, total, err := h.printJobRepo.FindByFilters(status, jobType, listQuery.StartDate, listQuery.EndDate, offset, limit)
	if err != nil {
		return nil, err
	}

	results := make([]*PrintJobResult, len(jobs))
	for i, job := range jobs {
		results[i] = toPrintJobResult(job)
	}

	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}

	return &PrintJobsQueryResult{
		Items: results,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// HandleGetPrintContent 处理获取打印内容查询
func (h *QueryHandler) HandleGetPrintContent(q common.Query) (any, error) {
	getQuery, ok := q.(GetPrintContentQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	bill, err := h.billRepo.FindByID(getQuery.BillID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, fmt.Errorf("bill not found")
	}

	lease, err := h.leaseRepo.FindByID(bill.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, fmt.Errorf("lease not found")
	}

	return h.printService.GenerateInvoiceContent(bill, lease), nil
}

// toPrintJobResult 转换领域模型为查询结果
func toPrintJobResult(job *printmodel.PrintJob) *PrintJobResult {
	amountYuan := ""
	if job.Amount > 0 {
		amountYuan = fmt.Sprintf("%.2f", float64(job.Amount)/100)
	}

	return &PrintJobResult{
		ID:           job.ID(),
		Type:         job.Type,
		Status:       job.Status,
		ReferenceID:  job.ReferenceID,
		TenantName:   job.TenantName,
		TenantPhone:  job.TenantPhone,
		RoomID:       job.RoomID,
		RoomNumber:   job.RoomNumber,
		Address:      job.Address,
		LandlordName: job.LandlordName,
		Amount:       job.Amount,
		ErrorMsg:     job.ErrorMsg,
		CreatedAt:    job.CreatedAt,
		UpdatedAt:    job.UpdatedAt,
		CompletedAt:  job.CompletedAt,
		TypeText:     job.GetTypeText(),
		StatusText:   job.GetStatusText(),
		AmountYuan:   amountYuan,
	}
}
