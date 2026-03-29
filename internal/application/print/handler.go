package print

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/common"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// CommandHandler 打印命令处理器
type CommandHandler struct {
	billRepo  billrepo.BillRepository
	leaseRepo leaserepo.LeaseRepository
	eventBus  *event.Bus
}

// NewCommandHandler 创建打印命令处理器
func NewCommandHandler(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository, eventBus *event.Bus) *CommandHandler {
	return &CommandHandler{billRepo: billRepo, leaseRepo: leaseRepo, eventBus: eventBus}
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

	bill, err := h.billRepo.FindByID(printCmd.BillID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, domerrors.ErrNotFound
	}

	lease, err := h.leaseRepo.FindByID(bill.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, domerrors.ErrNotFound
	}

	jobID := uuid.NewString()
	return jobID, nil
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

	lease, err := h.leaseRepo.FindByID(printCmd.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, domerrors.ErrNotFound
	}

	jobID := uuid.NewString()
	return jobID, nil
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

	bill, err := h.billRepo.FindByID(printCmd.BillID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, domerrors.ErrNotFound
	}

	lease, err := h.leaseRepo.FindByID(bill.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, domerrors.ErrNotFound
	}

	jobID := uuid.NewString()
	return jobID, nil
}

// QueryHandler 打印查询处理器
type QueryHandler struct {
	billRepo  billrepo.BillRepository
	leaseRepo leaserepo.LeaseRepository
}

// NewQueryHandler 创建打印查询处理器
func NewQueryHandler(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository) *QueryHandler {
	return &QueryHandler{billRepo: billRepo, leaseRepo: leaseRepo}
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

	content := []byte(fmt.Sprintf(`
{\rtf1\ansi\deff0{\fonttbl{\f0 Arial;}}
\pard\fs24\b 收据\b0\par
\pard\fs16 账单编号: %s\par
\pard\fs16 租约编号: %s\par
\pard\fs16 租客: %s\par
\pard\fs16 金额: %d 元\par
}`, bill.ID(), lease.ID(), lease.TenantName, bill.Amount))
	return content, nil
}
