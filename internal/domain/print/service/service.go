package service

import (
	"fmt"

	"github.com/google/uuid"

	billmodel "github.com/zouhang1992/ddd_domain/internal/domain/bill/model"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	leasemodel "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	printmodel "github.com/zouhang1992/ddd_domain/internal/domain/print/model"
	printrepo "github.com/zouhang1992/ddd_domain/internal/domain/print/repository"
	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
	locationrepo "github.com/zouhang1992/ddd_domain/internal/domain/location/repository"
	landlordrepo "github.com/zouhang1992/ddd_domain/internal/domain/landlord/repository"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
)

// PrintService 打印领域服务
type PrintService struct {
	billRepo     billrepo.BillRepository
	leaseRepo    leaserepo.LeaseRepository
	roomRepo     roomrepo.RoomRepository
	locationRepo locationrepo.LocationRepository
	landlordRepo landlordrepo.LandlordRepository
	printJobRepo printrepo.PrintJobRepository
}

// NewPrintService 创建打印领域服务
func NewPrintService(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository, roomRepo roomrepo.RoomRepository, locationRepo locationrepo.LocationRepository, landlordRepo landlordrepo.LandlordRepository, printJobRepo printrepo.PrintJobRepository) *PrintService {
	return &PrintService{
		billRepo:     billRepo,
		leaseRepo:    leaseRepo,
		roomRepo:     roomRepo,
		locationRepo: locationRepo,
		landlordRepo: landlordRepo,
		printJobRepo: printJobRepo,
	}
}

// getPrintJobDetails 获取打印作业详细信息
func (s *PrintService) getPrintJobDetails(lease *leasemodel.Lease) (tenantPhone, roomID, roomNumber, address, landlordName string) {
	tenantPhone = lease.TenantPhone
	roomID = lease.RoomID

	// 获取房间信息
	if room, err := s.roomRepo.FindByID(lease.RoomID); err == nil && room != nil {
		roomNumber = room.RoomNumber
		// 获取位置信息
		if location, err := s.locationRepo.FindByID(room.LocationID); err == nil && location != nil {
			address = location.Detail
			if address == "" {
				address = location.ShortName
			}
		}
	}

	// 获取房东信息
	if lease.LandlordID != "" {
		if landlord, err := s.landlordRepo.FindByID(lease.LandlordID); err == nil && landlord != nil {
			landlordName = landlord.Name
		}
	}

	return
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

	tenantPhone, roomID, roomNumber, address, landlordName := s.getPrintJobDetails(lease)

	jobID := uuid.NewString()
	job := printmodel.NewPrintJob(
		jobID,
		printmodel.PrintJobTypeBill,
		billID,
		lease.TenantName,
		tenantPhone,
		roomID,
		roomNumber,
		address,
		landlordName,
		bill.Amount,
	)
	job.MarkProcessing()

	if err := s.printJobRepo.Save(job); err != nil {
		return "", err
	}

	// 模拟打印完成
	job.MarkCompleted()
	if err := s.printJobRepo.Save(job); err != nil {
		return "", err
	}

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

	tenantPhone, roomID, roomNumber, address, landlordName := s.getPrintJobDetails(lease)

	jobID := uuid.NewString()
	job := printmodel.NewPrintJob(
		jobID,
		printmodel.PrintJobTypeLease,
		leaseID,
		lease.TenantName,
		tenantPhone,
		roomID,
		roomNumber,
		address,
		landlordName,
		0,
	)
	job.MarkProcessing()

	if err := s.printJobRepo.Save(job); err != nil {
		return "", err
	}

	// 模拟打印完成
	job.MarkCompleted()
	if err := s.printJobRepo.Save(job); err != nil {
		return "", err
	}

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

	tenantPhone, roomID, roomNumber, address, landlordName := s.getPrintJobDetails(lease)

	jobID := uuid.NewString()
	job := printmodel.NewPrintJob(
		jobID,
		printmodel.PrintJobTypeInvoice,
		billID,
		lease.TenantName,
		tenantPhone,
		roomID,
		roomNumber,
		address,
		landlordName,
		bill.Amount,
	)
	job.MarkProcessing()

	if err := s.printJobRepo.Save(job); err != nil {
		return "", err
	}

	// 模拟打印完成
	job.MarkCompleted()
	if err := s.printJobRepo.Save(job); err != nil {
		return "", err
	}

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
