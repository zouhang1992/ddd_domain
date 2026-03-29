package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	billmodel "github.com/zouhang1992/ddd_domain/internal/domain/bill/model"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	leasemodel "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
)

// PrintService 打印服务
type PrintService struct {
	billRepo  billrepo.BillRepository
	leaseRepo leaserepo.LeaseRepository
}

// NewPrintService 创建打印服务
func NewPrintService(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository) *PrintService {
	return &PrintService{
		billRepo:  billRepo,
		leaseRepo: leaseRepo,
	}
}

// PrintReceipt 打印收据
func (s *PrintService) PrintReceipt(billID string) ([]byte, error) {
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, fmt.Errorf("bill not found")
	}

	lease, err := s.leaseRepo.FindByID(bill.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, fmt.Errorf("lease not found")
	}

	receiptContent := s.createReceiptRTF(bill, lease)
	return []byte(receiptContent), nil
}

func (s *PrintService) createReceiptRTF(bill *billmodel.Bill, lease *leasemodel.Lease) string {
	var buf bytes.Buffer

	buf.WriteString(`{\rtf1\ansi\deff0{\fonttbl{\f0 Arial;}}`)
	buf.WriteString(`\pard\fs24\b 收据\b0\par`)
	buf.WriteString(`\pard\fs16 日期：` + time.Now().Format("2006年01月02日") + `\par`)
	buf.WriteString(`\pard\fs16 账单编号：` + bill.ID() + `\par`)
	buf.WriteString(`\pard\fs16 租约编号：` + lease.ID() + `\par`)
	buf.WriteString(`\pard\fs16 租客：` + lease.TenantName + `\par`)
	buf.WriteString(`\pard\fs16 电话：` + lease.TenantPhone + `\par`)
	buf.WriteString(`\pard\fs16 \par`)

	buf.WriteString(`{\pard\fs18\b 费用详情：\b0\par`)
	buf.WriteString(fmt.Sprintf(`\pard\fs16 合计：%.2f 元\par`, float64(bill.Amount)/100))
	buf.WriteString(`\pard\fs16 \par`)

	if bill.Note != "" {
		buf.WriteString(`\pard\fs16 备注：` + bill.Note + `\par`)
		buf.WriteString(`\pard\fs16 \par`)
	}

	buf.WriteString(`\pard\fs16 \line \line \line`)
	buf.WriteString(`\pard\fs16 ------------------------------\par`)
	buf.WriteString(`\pard\fs16 \line \line \line`)

	buf.WriteString(`\pard\fs18\b 支付信息：\b0\par`)
	if bill.PaidAt != nil {
		buf.WriteString(`\pard\fs16 支付状态：已支付\par`)
		buf.WriteString(`\pard\fs16 支付时间：` + bill.PaidAt.Format("2006年01月02日 15:04:05") + `\par`)
	} else {
		buf.WriteString(`\pard\fs16 支付状态：待支付\par`)
	}

	buf.WriteString(`}`)
	return buf.String()
}

// PrintContract 打印合同
func (s *PrintService) PrintContract(leaseID string) ([]byte, error) {
	lease, err := s.leaseRepo.FindByID(leaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, fmt.Errorf("lease not found")
	}

	contractContent := s.createContractRTF(lease)
	return []byte(contractContent), nil
}

func (s *PrintService) createContractRTF(lease *leasemodel.Lease) string {
	var buf bytes.Buffer

	buf.WriteString(`{\rtf1\ansi\deff0{\fonttbl{\f0 Arial;}}`)
	buf.WriteString(`\pard\fs24\b 房屋租赁合同\b0\par`)
	buf.WriteString(`\pard\fs16 合同编号：` + lease.ID() + `\par`)
	buf.WriteString(`\pard\fs16 签订日期：` + time.Now().Format("2006年01月02日") + `\par`)
	buf.WriteString(`\pard\fs16 \par`)

	buf.WriteString(`\pard\fs18\b 出租方（甲方）：\b0\par`)
	buf.WriteString(`\pard\fs16 姓名：房东\par`)
	buf.WriteString(`\pard\fs16 联系方式：\par`)
	buf.WriteString(`\pard\fs16 \par`)

	buf.WriteString(`\pard\fs18\b 承租方（乙方）：\b0\par`)
	buf.WriteString(`\pard\fs16 姓名：` + lease.TenantName + `\par`)
	buf.WriteString(`\pard\fs16 联系方式：` + lease.TenantPhone + `\par`)
	buf.WriteString(`\pard\fs16 \par`)

	buf.WriteString(`\pard\fs18\b 租赁房屋：\b0\par`)
	buf.WriteString(`\pard\fs16 房间编号：` + lease.RoomID + `\par`)
	buf.WriteString(`\pard\fs16 租赁期限：` + lease.StartDate.Format("2006年01月02日") +
		` 至 ` + lease.EndDate.Format("2006年01月02日") + `\par`)
	buf.WriteString(`\pard\fs16 租金：` + fmt.Sprintf(`%.2f 元/月\par`, float64(lease.RentAmount)/100))
	buf.WriteString(`\pard\fs16 \par`)

	buf.WriteString(`\pard\fs16 双方签字：\par`)
	buf.WriteString(`\pard\fs16 \line \line \line`)
	buf.WriteString(`\pard\fs16 甲方签字：__________________________\par`)
	buf.WriteString(`\pard\fs16 乙方签字：__________________________\par`)
	buf.WriteString(`}`)

	return buf.String()
}

// PrintInvoice 打印发票
func (s *PrintService) PrintInvoice(billID string) ([]byte, error) {
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, fmt.Errorf("bill not found")
	}

	lease, err := s.leaseRepo.FindByID(bill.LeaseID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, fmt.Errorf("lease not found")
	}

	invoice := map[string]interface{}{
		"BillID":      billID,
		"LeaseID":     lease.ID(),
		"Amount":      bill.Amount,
		"Type":        bill.Type,
		"TenantName":  lease.TenantName,
		"TenantPhone": lease.TenantPhone,
		"PaidAt":      bill.PaidAt,
		"Note":        bill.Note,
	}

	data, err := json.MarshalIndent(invoice, "", "  ")
	if err != nil {
		return nil, err
	}
	return data, nil
}
