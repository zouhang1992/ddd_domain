package print

import (
	"fmt"

	"github.com/zouhang1992/ddd_domain/internal/application/common"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
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
}

// NewQueryHandler 创建打印查询处理器
func NewQueryHandler(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository, printService *printservice.PrintService) *QueryHandler {
	return &QueryHandler{billRepo: billRepo, leaseRepo: leaseRepo, printService: printService}
}

// HandleGetPrintJob 处理获取打印作业查询
func (h *QueryHandler) HandleGetPrintJob(q common.Query) (any, error) {
	return nil, nil
}

// HandleListPrintJobs 处理列出打印作业查询
func (h *QueryHandler) HandleListPrintJobs(q common.Query) (any, error) {
	return &PrintJobsQueryResult{
		Items: []interface{}{},
		Total: 0,
		Page:  1,
		Limit: 10,
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
