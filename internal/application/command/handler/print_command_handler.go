package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// PrintCommandHandler 打印命令处理器
type PrintCommandHandler struct {
	billRepo  repository.BillRepository
	leaseRepo repository.LeaseRepository
	eventBus  *event.Bus
}

// NewPrintCommandHandler 创建打印命令处理器
func NewPrintCommandHandler(billRepo repository.BillRepository, leaseRepo repository.LeaseRepository, eventBus *event.Bus) *PrintCommandHandler {
	return &PrintCommandHandler{billRepo: billRepo, leaseRepo: leaseRepo, eventBus: eventBus}
}

// HandlePrintBill 处理打印账单命令
func (h *PrintCommandHandler) HandlePrintBill(cmd command.Command) (any, error) {
	printCmd, ok := cmd.(command.PrintBillCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := printCmd.Validate(); err != nil {
		return nil, err
	}

	bill, err := h.billRepo.FindByID(printCmd.BillID)
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

	// 这里需要创建一个打印内容的方法，复用 PrintService 的逻辑
	printService := &PrintService{billRepo: h.billRepo, leaseRepo: h.leaseRepo}
	content, err := printService.PrintReceipt(bill.ID)
	if err != nil {
		return nil, err
	}

	// 生成打印作业ID
	jobID := uuid.NewString()

	// 发布账单打印事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewBillPrinted(jobID, bill.ID, content))
	}

	return jobID, nil
}

// HandlePrintLease 处理打印租约命令
func (h *PrintCommandHandler) HandlePrintLease(cmd command.Command) (any, error) {
	printCmd, ok := cmd.(command.PrintLeaseCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := printCmd.Validate(); err != nil {
		return nil, err
	}

	lease, err := h.leaseRepo.FindByID(printCmd.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, model.ErrNotFound
	}

	printService := &PrintService{billRepo: h.billRepo, leaseRepo: h.leaseRepo}
	content, err := printService.PrintContract(lease.ID)
	if err != nil {
		return nil, err
	}

	jobID := uuid.NewString()

	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLeasePrinted(jobID, lease.ID, content))
	}

	return jobID, nil
}

// HandlePrintInvoice 处理打印发票命令
func (h *PrintCommandHandler) HandlePrintInvoice(cmd command.Command) (any, error) {
	printCmd, ok := cmd.(command.PrintInvoiceCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := printCmd.Validate(); err != nil {
		return nil, err
	}

	bill, err := h.billRepo.FindByID(printCmd.BillID)
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

	printService := &PrintService{billRepo: h.billRepo, leaseRepo: h.leaseRepo}
	content, err := printService.PrintInvoice(bill.ID)
	if err != nil {
		return nil, err
	}

	jobID := uuid.NewString()

	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewInvoicePrinted(jobID, bill.ID, content))
	}

	return jobID, nil
}

// PrintService 内部使用的打印服务，复用原来的逻辑
type PrintService struct {
	billRepo  repository.BillRepository
	leaseRepo repository.LeaseRepository
}

// PrintReceipt 打印收据（RTF格式）
func (s *PrintService) PrintReceipt(billID string) ([]byte, error) {
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, model.ErrNotFound
	}

	lease, err := s.leaseRepo.FindByID(bill.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, model.ErrNotFound
	}

	return s.createReceiptRTF(bill, lease), nil
}

// PrintContract 打印合同（RTF格式）
func (s *PrintService) PrintContract(leaseID string) ([]byte, error) {
	lease, err := s.leaseRepo.FindByID(leaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, model.ErrNotFound
	}

	return s.createContractRTF(lease), nil
}

// PrintInvoice 打印发票
func (s *PrintService) PrintInvoice(billID string) ([]byte, error) {
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, model.ErrNotFound
	}

	lease, err := s.leaseRepo.FindByID(bill.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, model.ErrNotFound
	}

	// 简单的JSON格式发票
	return []byte(`{"BillID":"` + bill.ID + `","Amount":` + fmt.Sprintf("%d", bill.Amount) + `}`), nil
}

// createReceiptRTF 创建收据RTF内容（简化版，实际项目中可能需要更复杂的格式）
func (s *PrintService) createReceiptRTF(bill *model.Bill, lease *model.Lease) []byte {
	content := []byte(`
{\rtf1\ansi\deff0{\fonttbl{\f0 Arial;}}
\pard\fs24\b 收据\b0\par
\pard\fs16 账单编号: ` + bill.ID + `\par
\pard\fs16 租约编号: ` + lease.ID + `\par
\pard\fs16 租客: ` + lease.TenantName + `\par
\pard\fs16 金额: ` + fmt.Sprintf("%d", bill.Amount) + ` 元\par
}`)
	return content
}

// createContractRTF 创建合同RTF内容（简化版）
func (s *PrintService) createContractRTF(lease *model.Lease) []byte {
	content := []byte(`
{\rtf1\ansi\deff0{\fonttbl{\f0 Arial;}}
\pard\fs24\b 房屋租赁合同\b0\par
\pard\fs16 合同编号: ` + lease.ID + `\par
\pard\fs16 租客: ` + lease.TenantName + `\par
\pard\fs16 房间编号: ` + lease.RoomID + `\par
\pard\fs16 租赁期限: ` + lease.StartDate.Format("2006-01-02") + ` 至 ` + lease.EndDate.Format("2006-01-02") + `\par
}`)
	return content
}
