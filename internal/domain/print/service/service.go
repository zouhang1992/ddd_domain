package service

import (
	"fmt"

	"github.com/google/uuid"

	billmodel "github.com/zouhang1992/ddd_domain/internal/domain/bill/model"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	leasemodel "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
)

// PrintService 打印领域服务
type PrintService struct {
	billRepo  billrepo.BillRepository
	leaseRepo leaserepo.LeaseRepository
}

// NewPrintService 创建打印领域服务
func NewPrintService(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository) *PrintService {
	return &PrintService{
		billRepo:  billRepo,
		leaseRepo: leaseRepo,
	}
}

// CreateBillPrintJob 创建账单打印作业
func (s *PrintService) CreateBillPrintJob(billID string) (string, error) {
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return "", err
	}
	if bill == nil {
		return "", domerrors.ErrNotFound
	}

	lease, err := s.leaseRepo.FindByID(bill.LeaseID)
	if err != nil {
		return "", err
	}
	if lease == nil {
		return "", domerrors.ErrNotFound
	}

	jobID := uuid.NewString()
	return jobID, nil
}

// CreateLeasePrintJob 创建租约打印作业
func (s *PrintService) CreateLeasePrintJob(leaseID string) (string, error) {
	lease, err := s.leaseRepo.FindByID(leaseID)
	if err != nil {
		return "", err
	}
	if lease == nil {
		return "", domerrors.ErrNotFound
	}

	jobID := uuid.NewString()
	return jobID, nil
}

// CreateInvoicePrintJob 创建发票打印作业
func (s *PrintService) CreateInvoicePrintJob(billID string) (string, error) {
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return "", err
	}
	if bill == nil {
		return "", domerrors.ErrNotFound
	}

	lease, err := s.leaseRepo.FindByID(bill.LeaseID)
	if err != nil {
		return "", err
	}
	if lease == nil {
		return "", domerrors.ErrNotFound
	}

	jobID := uuid.NewString()
	return jobID, nil
}

// GenerateBillContent 生成账单打印内容
func (s *PrintService) GenerateBillContent(bill *billmodel.Bill, lease *leasemodel.Lease) []byte {
	content := []byte(fmt.Sprintf(`
{\rtf1\ansi\deff0{\fonttbl{\f0 Arial;}}
\pard\fs24\b 账单\b0\par
\pard\fs16 账单编号: %s\par
\pard\fs16 租约编号: %s\par
\pard\fs16 租客: %s\par
\pard\fs16 金额: %d 元\par
}`, bill.ID(), lease.ID(), lease.TenantName, bill.Amount))
	return content
}

// GenerateLeaseContent 生成租约打印内容
func (s *PrintService) GenerateLeaseContent(lease *leasemodel.Lease) []byte {
	content := []byte(fmt.Sprintf(`
{\rtf1\ansi\deff0{\fonttbl{\f0 Arial;}}
\pard\fs24\b 租约\b0\par
\pard\fs16 租约编号: %s\par
\pard\fs16 租客: %s\par
\pard\fs16 租期: %s 至 %s\par
\pard\fs16 租金: %d 元\par
}`, lease.ID(), lease.TenantName,
		lease.StartDate.Format("2006-01-02"),
		lease.EndDate.Format("2006-01-02"),
		lease.RentAmount))
	return content
}

// GenerateInvoiceContent 生成发票打印内容
func (s *PrintService) GenerateInvoiceContent(bill *billmodel.Bill, lease *leasemodel.Lease) []byte {
	content := []byte(fmt.Sprintf(`
{\rtf1\ansi\deff0{\fonttbl{\f0 Arial;}}
\pard\fs24\b 收据\b0\par
\pard\fs16 账单编号: %s\par
\pard\fs16 租约编号: %s\par
\pard\fs16 租客: %s\par
\pard\fs16 金额: %d 元\par
}`, bill.ID(), lease.ID(), lease.TenantName, bill.Amount))
	return content
}
