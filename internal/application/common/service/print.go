package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	billmodel "github.com/zouhang1992/ddd_domain/internal/domain/bill/model"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	leasemodel "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
	locationrepo "github.com/zouhang1992/ddd_domain/internal/domain/location/repository"
	landlordrepo "github.com/zouhang1992/ddd_domain/internal/domain/landlord/repository"
)

func formatAmount(amount int64) string {
	return fmt.Sprintf("%.2f", float64(amount)/100)
}

func getBillTypeText(billType billmodel.BillType) string {
	typeMap := map[billmodel.BillType]string{
		billmodel.BillTypeRent:     "租金",
		billmodel.BillTypeWater:    "水费",
		billmodel.BillTypeElectric: "电费",
		billmodel.BillTypeGas:      "燃气费",
		billmodel.BillTypeInternet: "网费",
		billmodel.BillTypeOther:    "其他费用",
		billmodel.BillTypeCharge:   "收账",
		billmodel.BillTypeCheckout: "退租结算",
	}
	if text, ok := typeMap[billType]; ok {
		return text
	}
	return string(billType)
}

func getBillStatusText(status billmodel.BillStatus) string {
	statusMap := map[billmodel.BillStatus]string{
		billmodel.BillStatusPending: "待支付",
		billmodel.BillStatusPaid:    "已支付",
	}
	if text, ok := statusMap[status]; ok {
		return text
	}
	return string(status)
}

// PrintService 打印服务
type PrintService struct {
	billRepo     billrepo.BillRepository
	leaseRepo    leaserepo.LeaseRepository
	roomRepo     roomrepo.RoomRepository
	locationRepo locationrepo.LocationRepository
	landlordRepo landlordrepo.LandlordRepository
}

// NewPrintService 创建打印服务
func NewPrintService(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository, roomRepo roomrepo.RoomRepository, locationRepo locationrepo.LocationRepository, landlordRepo landlordrepo.LandlordRepository) *PrintService {
	return &PrintService{
		billRepo:     billRepo,
		leaseRepo:    leaseRepo,
		roomRepo:     roomRepo,
		locationRepo: locationRepo,
		landlordRepo: landlordRepo,
	}
}

// getRoomAndAddress 获取房间号和地址信息
func (s *PrintService) getRoomAndAddress(roomID string) (roomNumber, address string) {
	if room, err := s.roomRepo.FindByID(roomID); err == nil && room != nil {
		roomNumber = room.RoomNumber
		// 获取位置信息
		if location, err := s.locationRepo.FindByID(room.LocationID); err == nil && location != nil {
			address = location.Detail
			if address == "" {
				address = location.ShortName
			}
		}
	}
	return
}

// getLandlordInfo 获取房东信息
func (s *PrintService) getLandlordInfo(landlordID string) (name, phone string) {
	if landlordID != "" {
		if landlord, err := s.landlordRepo.FindByID(landlordID); err == nil && landlord != nil {
			name = landlord.Name
			phone = landlord.Phone
		}
	}
	return
}

// receiptData 收据模板数据
type receiptData struct {
	Title               string
	PrintDate           string
	BillID              string
	LeaseID             string
	TenantName          string
	TenantPhone         string
	RoomNumber          string
	Address             string
	BillType            string
	IsCheckout          bool
	RentAmount          string
	WaterAmount         string
	ElectricAmount      string
	OtherAmount         string
	RefundDepositAmount string
	TotalAmount         string
	Note                string
	Status              string
	PaidAt              string
	DepositAmount       string
	StartDate           string
	EndDate             string
}

// contractData 合同模板数据
type contractData struct {
	Title         string
	PrintDate     string
	ContractID    string
	LandlordName  string
	LandlordPhone string
	TenantName    string
	TenantPhone   string
	RoomNumber    string
	Address       string
	StartDate     string
	EndDate       string
	RentAmount    string
	DepositAmount string
	Note          string
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

	receiptContent := s.createReceiptHTML(bill, lease)
	return []byte(receiptContent), nil
}

func (s *PrintService) createReceiptHTML(bill *billmodel.Bill, lease *leasemodel.Lease) string {
	roomNumber, address := s.getRoomAndAddress(lease.RoomID)

	data := receiptData{
		Title:               "收款收据",
		PrintDate:           time.Now().Format("2006年01月02日 15:04:05"),
		BillID:              bill.ID(),
		LeaseID:             lease.ID(),
		TenantName:          lease.TenantName,
		TenantPhone:         lease.TenantPhone,
		RoomNumber:          roomNumber,
		Address:             address,
		BillType:            getBillTypeText(bill.Type),
		IsCheckout:          bill.Type == billmodel.BillTypeCheckout,
		RentAmount:          formatAmount(bill.RentAmount),
		WaterAmount:         formatAmount(bill.WaterAmount),
		ElectricAmount:      formatAmount(bill.ElectricAmount),
		OtherAmount:         formatAmount(bill.OtherAmount),
		RefundDepositAmount: formatAmount(bill.RefundDepositAmount),
		TotalAmount:         formatAmount(bill.Amount),
		Note:                bill.Note,
		Status:              getBillStatusText(bill.Status),
		DepositAmount:       formatAmount(lease.DepositAmount),
		StartDate:           lease.StartDate.Format("2006年01月02日"),
		EndDate:             lease.EndDate.Format("2006年01月02日"),
	}

	if bill.PaidAt != nil {
		data.PaidAt = bill.PaidAt.Format("2006年01月02日 15:04:05")
	}

	const receiptTemplate = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: "Microsoft YaHei", "SimHei", Arial, sans-serif;
            font-size: 14px;
            line-height: 1.6;
            padding: 20px;
            background: #fff;
        }
        .receipt {
            max-width: 800px;
            margin: 0 auto;
            border: 2px solid #333;
            padding: 30px;
            background: #fff;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
            border-bottom: 3px double #333;
            padding-bottom: 20px;
        }
        .title {
            font-size: 28px;
            font-weight: bold;
            color: #333;
        }
        .print-date {
            margin-top: 10px;
            color: #666;
            font-size: 12px;
        }
        .info-section {
            margin-bottom: 25px;
        }
        .info-row {
            display: flex;
            margin-bottom: 10px;
        }
        .info-label {
            width: 100px;
            font-weight: bold;
            color: #555;
        }
        .info-value {
            flex: 1;
            color: #333;
        }
        .details-section {
            margin: 25px 0;
        }
        .section-title {
            font-size: 16px;
            font-weight: bold;
            border-bottom: 2px solid #333;
            padding-bottom: 8px;
            margin-bottom: 15px;
            color: #333;
        }
        .details-table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 15px;
        }
        .details-table th,
        .details-table td {
            border: 1px solid #ccc;
            padding: 12px;
            text-align: left;
        }
        .details-table th {
            background: #f5f5f5;
            font-weight: bold;
            width: 40%;
        }
        .details-table td {
            text-align: right;
        }
        .total-row {
            background: #fff9e6;
            font-weight: bold;
            font-size: 16px;
        }
        .total-row td {
            color: #d32f2f;
        }
        .refund {
            color: #388e3c !important;
        }
        .note-section {
            margin: 20px 0;
            padding: 15px;
            background: #f9f9f9;
            border-left: 4px solid #2196f3;
        }
        .note-title {
            font-weight: bold;
            margin-bottom: 5px;
            color: #333;
        }
        .status-section {
            margin-top: 30px;
            padding: 15px;
            border: 1px solid #ddd;
            background: #fafafa;
        }
        .status-row {
            display: flex;
            justify-content: space-between;
            margin-bottom: 8px;
        }
        .status-label {
            font-weight: bold;
            color: #555;
        }
        .status-value {
            color: #333;
        }
        .status-paid {
            color: #388e3c;
            font-weight: bold;
        }
        .status-unpaid {
            color: #d32f2f;
            font-weight: bold;
        }
        .signature-section {
            margin-top: 50px;
            display: flex;
            justify-content: space-between;
        }
        .signature-box {
            width: 45%;
            text-align: center;
        }
        .signature-line {
            border-bottom: 1px solid #333;
            margin: 40px 0 10px 0;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            color: #999;
            font-size: 12px;
            border-top: 1px solid #eee;
            padding-top: 15px;
        }
        @media print {
            body {
                padding: 0;
            }
            .receipt {
                border: none;
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="receipt">
        <div class="header">
            <div class="title">{{.Title}}</div>
            <div class="print-date">打印时间：{{.PrintDate}}</div>
        </div>

        <div class="info-section">
            <div class="info-row">
                <span class="info-label">账单编号：</span>
                <span class="info-value">{{.BillID}}</span>
            </div>
            <div class="info-row">
                <span class="info-label">租约编号：</span>
                <span class="info-value">{{.LeaseID}}</span>
            </div>
            <div class="info-row">
                <span class="info-label">账单类型：</span>
                <span class="info-value">{{.BillType}}</span>
            </div>
            <div class="info-row">
                <span class="info-label">租客姓名：</span>
                <span class="info-value">{{.TenantName}}</span>
            </div>
            <div class="info-row">
                <span class="info-label">联系电话：</span>
                <span class="info-value">{{.TenantPhone}}</span>
            </div>
            {{if .RoomNumber}}
            <div class="info-row">
                <span class="info-label">房间号：</span>
                <span class="info-value">{{.RoomNumber}}</span>
            </div>
            {{end}}
            {{if .Address}}
            <div class="info-row">
                <span class="info-label">地址：</span>
                <span class="info-value">{{.Address}}</span>
            </div>
            {{end}}
            <div class="info-row">
                <span class="info-label">租赁期限：</span>
                <span class="info-value">{{.StartDate}} 至 {{.EndDate}}</span>
            </div>
        </div>

        <div class="details-section">
            <div class="section-title">费用明细</div>
            <table class="details-table">
                {{if .IsCheckout}}
                <tr>
                    <th>押金金额</th>
                    <td>{{.DepositAmount}}</td>
                </tr>
                <tr>
                    <th>退还押金</th>
                    <td class="refund">- {{.RefundDepositAmount}}</td>
                </tr>
                {{end}}
                {{if ne .RentAmount "0.00"}}
                <tr>
                    <th>{{if .IsCheckout}}退还租金{{else}}租金{{end}}</th>
                    <td>{{if and .IsCheckout (gt .RentAmount 0)}}- {{end}}{{.RentAmount}}</td>
                </tr>
                {{end}}
                {{if ne .WaterAmount "0.00"}}
                <tr>
                    <th>水费</th>
                    <td>{{.WaterAmount}}</td>
                </tr>
                {{end}}
                {{if ne .ElectricAmount "0.00"}}
                <tr>
                    <th>电费</th>
                    <td>{{.ElectricAmount}}</td>
                </tr>
                {{end}}
                {{if ne .OtherAmount "0.00"}}
                <tr>
                    <th>其他费用</th>
                    <td>{{.OtherAmount}}</td>
                </tr>
                {{end}}
                <tr class="total-row">
                    <th>合计金额</th>
                    <td>{{.TotalAmount}}</td>
                </tr>
            </table>
        </div>

        {{if .Note}}
        <div class="note-section">
            <div class="note-title">备注</div>
            <div>{{.Note}}</div>
        </div>
        {{end}}

        <div class="status-section">
            <div class="status-row">
                <span class="status-label">支付状态：</span>
                <span class="status-value {{if eq .Status "已支付"}}status-paid{{else}}status-unpaid{{end}}">{{.Status}}</span>
            </div>
            {{if .PaidAt}}
            <div class="status-row">
                <span class="status-label">支付时间：</span>
                <span class="status-value">{{.PaidAt}}</span>
            </div>
            {{end}}
        </div>

        <div class="signature-section">
            <div class="signature-box">
                <div class="signature-line"></div>
                <div>收款人签字</div>
            </div>
            <div class="signature-box">
                <div class="signature-line"></div>
                <div>付款人签字</div>
            </div>
        </div>

        <div class="footer">
            此收据一式两份，双方各执一份，具有同等法律效力
        </div>
    </div>
</body>
</html>`

	tmpl, err := template.New("receipt").Parse(receiptTemplate)
	if err != nil {
		return fmt.Sprintf("<html><body><h1>模板错误: %v</h1></body></html>", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Sprintf("<html><body><h1>渲染错误: %v</h1></body></html>", err)
	}

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

	contractContent := s.createContractHTML(lease)
	return []byte(contractContent), nil
}

func (s *PrintService) createContractHTML(lease *leasemodel.Lease) string {
	roomNumber, address := s.getRoomAndAddress(lease.RoomID)
	landlordName, landlordPhone := s.getLandlordInfo(lease.LandlordID)

	data := contractData{
		Title:         "房屋租赁合同",
		PrintDate:     time.Now().Format("2006年01月02日 15:04:05"),
		ContractID:    lease.ID(),
		LandlordName:  landlordName,
		LandlordPhone: landlordPhone,
		TenantName:    lease.TenantName,
		TenantPhone:   lease.TenantPhone,
		RoomNumber:    roomNumber,
		Address:       address,
		StartDate:     lease.StartDate.Format("2006年01月02日"),
		EndDate:       lease.EndDate.Format("2006年01月02日"),
		RentAmount:    formatAmount(lease.RentAmount),
		DepositAmount: formatAmount(lease.DepositAmount),
		Note:          lease.Note,
	}

	const contractTemplate = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: "Microsoft YaHei", "SimHei", Arial, sans-serif;
            font-size: 14px;
            line-height: 1.8;
            padding: 30px;
            background: #fff;
        }
        .contract {
            max-width: 800px;
            margin: 0 auto;
            border: 2px solid #333;
            padding: 40px;
            background: #fff;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
            border-bottom: 3px double #333;
            padding-bottom: 20px;
        }
        .title {
            font-size: 32px;
            font-weight: bold;
            color: #333;
        }
        .print-date {
            margin-top: 10px;
            color: #666;
            font-size: 12px;
        }
        .contract-id {
            text-align: right;
            color: #666;
            margin-bottom: 20px;
        }
        .section {
            margin-bottom: 25px;
        }
        .section-title {
            font-size: 18px;
            font-weight: bold;
            border-left: 4px solid #333;
            padding-left: 12px;
            margin-bottom: 15px;
            color: #333;
        }
        .party-info {
            background: #f9f9f9;
            padding: 20px;
            border-radius: 4px;
        }
        .party-row {
            display: flex;
            margin-bottom: 12px;
        }
        .party-label {
            width: 120px;
            font-weight: bold;
            color: #555;
        }
        .party-value {
            flex: 1;
            color: #333;
        }
        .terms-section {
            margin: 30px 0;
        }
        .term-item {
            margin-bottom: 20px;
        }
        .term-title {
            font-weight: bold;
            color: #333;
            margin-bottom: 8px;
        }
        .term-content {
            text-indent: 2em;
            color: #444;
        }
        .info-table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }
        .info-table th,
        .info-table td {
            border: 1px solid #ccc;
            padding: 15px;
            text-align: left;
        }
        .info-table th {
            background: #f5f5f5;
            font-weight: bold;
            width: 35%;
        }
        .note-section {
            margin: 25px 0;
            padding: 15px;
            background: #fff9e6;
            border: 1px solid #ffcc80;
            border-radius: 4px;
        }
        .note-title {
            font-weight: bold;
            margin-bottom: 8px;
            color: #333;
        }
        .signature-section {
            margin-top: 60px;
            display: flex;
            justify-content: space-between;
        }
        .signature-box {
            width: 45%;
            text-align: center;
        }
        .signature-line {
            border-bottom: 1px solid #333;
            margin: 50px 0 10px 0;
        }
        .signature-date {
            margin-top: 10px;
            color: #666;
        }
        .footer {
            margin-top: 40px;
            text-align: center;
            color: #999;
            font-size: 12px;
            border-top: 1px solid #eee;
            padding-top: 20px;
        }
        @media print {
            body {
                padding: 0;
            }
            .contract {
                border: none;
                padding: 30px;
            }
        }
    </style>
</head>
<body>
    <div class="contract">
        <div class="header">
            <div class="title">{{.Title}}</div>
            <div class="print-date">打印时间：{{.PrintDate}}</div>
        </div>

        <div class="contract-id">
            合同编号：{{.ContractID}}
        </div>

        <div class="section">
            <div class="section-title">双方当事人信息</div>
            <div class="party-info">
                <div class="party-row">
                    <span class="party-label">甲方（出租方）：</span>
                    <span class="party-value">{{if .LandlordName}}{{.LandlordName}}{{else}}房东{{end}}</span>
                </div>
                {{if .LandlordPhone}}
                <div class="party-row">
                    <span class="party-label">联系电话：</span>
                    <span class="party-value">{{.LandlordPhone}}</span>
                </div>
                {{end}}
            </div>
            <div style="height: 15px;"></div>
            <div class="party-info">
                <div class="party-row">
                    <span class="party-label">乙方（承租方）：</span>
                    <span class="party-value">{{.TenantName}}</span>
                </div>
                <div class="party-row">
                    <span class="party-label">联系电话：</span>
                    <span class="party-value">{{.TenantPhone}}</span>
                </div>
            </div>
        </div>

        <div class="section">
            <div class="section-title">租赁房屋信息</div>
            <table class="info-table">
                {{if .RoomNumber}}
                <tr>
                    <th>房间号</th>
                    <td>{{.RoomNumber}}</td>
                </tr>
                {{end}}
                {{if .Address}}
                <tr>
                    <th>房屋地址</th>
                    <td>{{.Address}}</td>
                </tr>
                {{end}}
                <tr>
                    <th>租赁期限</th>
                    <td>{{.StartDate}} 至 {{.EndDate}}</td>
                </tr>
                <tr>
                    <th>月租金</th>
                    <td>{{.RentAmount}} 元/月</td>
                </tr>
                <tr>
                    <th>押金金额</th>
                    <td>{{.DepositAmount}} 元</td>
                </tr>
            </table>
        </div>

        <div class="terms-section">
            <div class="section-title">合同条款</div>

            <div class="term-item">
                <div class="term-title">第一条 租赁用途</div>
                <div class="term-content">
                    乙方租赁该房屋仅作为居住使用，不得擅自改变房屋用途，不得利用该房屋从事违法活动。
                </div>
            </div>

            <div class="term-item">
                <div class="term-title">第二条 租金支付方式</div>
                <div class="term-content">
                    乙方应按约定时间支付租金。租金支付方式为：[请在此处填写具体支付方式，如：月付、季付等]。
                    乙方应在每期租金到期前[请在此处填写天数]天内支付下期租金。
                </div>
            </div>

            <div class="term-item">
                <div class="term-title">第三条 押金</div>
                <div class="term-content">
                    乙方应向甲方支付押金 {{.DepositAmount}} 元。租赁期满或合同解除后，扣除应由乙方承担的费用后，
                    甲方应将押金余额无息退还给乙方。
                </div>
            </div>

            <div class="term-item">
                <div class="term-title">第四条 房屋使用与维护</div>
                <div class="term-content">
                    乙方应合理使用并爱护该房屋及其附属设施。因乙方保管不当或不合理使用，致使该房屋及其附属设施
                    发生损坏或故障的，乙方应负责维修或承担赔偿责任。
                </div>
            </div>

            <div class="term-item">
                <div class="term-title">第五条 水费、电费等费用</div>
                <div class="term-content">
                    租赁期间，乙方应承担水费、电费、燃气费、物业管理费等因使用该房屋所产生的费用。
                    乙方应按时缴纳上述费用。
                </div>
            </div>

            <div class="term-item">
                <div class="term-title">第六条 合同解除</div>
                <div class="term-content">
                    经甲乙双方协商一致，可以解除本合同。因不可抗力导致本合同无法继续履行的，本合同自行解除。
                </div>
            </div>

            <div class="term-item">
                <div class="term-title">第七条 违约责任</div>
                <div class="term-content">
                    双方应严格履行本合同约定的义务。任何一方违约，应承担相应的违约责任。
                    具体违约责任由双方另行约定或按照相关法律法规执行。
                </div>
            </div>

            <div class="term-item">
                <div class="term-title">第八条 争议解决</div>
                <div class="term-content">
                    本合同履行过程中发生的争议，由双方协商解决；协商不成的，任何一方均有权向房屋所在地人民法院提起诉讼。
                </div>
            </div>

            <div class="term-item">
                <div class="term-title">第九条 其他约定</div>
                <div class="term-content">
                    本合同未尽事宜，可由双方协商补充。补充协议与本合同具有同等法律效力。
                    本合同一式两份，甲乙双方各执一份，自双方签字或盖章之日起生效。
                </div>
            </div>
        </div>

        {{if .Note}}
        <div class="note-section">
            <div class="note-title">补充说明</div>
            <div>{{.Note}}</div>
        </div>
        {{end}}

        <div class="signature-section">
            <div class="signature-box">
                <div class="signature-line"></div>
                <div>甲方（签字或盖章）</div>
                <div class="signature-date">日期：________________</div>
            </div>
            <div class="signature-box">
                <div class="signature-line"></div>
                <div>乙方（签字或盖章）</div>
                <div class="signature-date">日期：________________</div>
            </div>
        </div>

        <div class="footer">
            本合同一式两份，甲乙双方各执一份，具有同等法律效力
        </div>
    </div>
</body>
</html>`

	tmpl, err := template.New("contract").Parse(contractTemplate)
	if err != nil {
		return fmt.Sprintf("<html><body><h1>模板错误: %v</h1></body></html>", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Sprintf("<html><body><h1>渲染错误: %v</h1></body></html>", err)
	}

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
