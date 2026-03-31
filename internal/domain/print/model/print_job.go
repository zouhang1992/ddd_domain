package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
)

// PrintJobStatus 打印作业状态
type PrintJobStatus string

const (
	PrintJobStatusPending    PrintJobStatus = "pending"
	PrintJobStatusProcessing PrintJobStatus = "processing"
	PrintJobStatusCompleted  PrintJobStatus = "completed"
	PrintJobStatusFailed     PrintJobStatus = "failed"
)

// PrintJobType 打印作业类型
type PrintJobType string

const (
	PrintJobTypeBill    PrintJobType = "bill"
	PrintJobTypeLease   PrintJobType = "lease"
	PrintJobTypeInvoice PrintJobType = "invoice"
)

// PrintJob 打印作业领域模型（聚合根）
type PrintJob struct {
	model.BaseAggregateRoot
	Type          PrintJobType  `json:"type"`
	Status        PrintJobStatus `json:"status"`
	ReferenceID   string         `json:"reference_id"`   // 关联的账单ID或租约ID
	TenantName    string         `json:"tenant_name"`    // 租户姓名（冗余字段，方便展示）
	TenantPhone   string         `json:"tenant_phone"`   // 租户电话
	RoomID        string         `json:"room_id"`        // 房间ID
	RoomNumber    string         `json:"room_number"`    // 房间号
	Address       string         `json:"address"`        // 地址
	LandlordName  string         `json:"landlord_name"`  // 房东姓名
	Amount        int64          `json:"amount"`         // 金额（分，仅账单类型有）
	ErrorMsg      string         `json:"error_msg"`      // 错误信息
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
}

// NewPrintJob 创建新打印作业
func NewPrintJob(id string, jobType PrintJobType, referenceID string, tenantName, tenantPhone, roomID, roomNumber, address, landlordName string, amount int64) *PrintJob {
	now := time.Now()
	return &PrintJob{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		Type:              jobType,
		Status:            PrintJobStatusPending,
		ReferenceID:       referenceID,
		TenantName:        tenantName,
		TenantPhone:       tenantPhone,
		RoomID:            roomID,
		RoomNumber:        roomNumber,
		Address:           address,
		LandlordName:      landlordName,
		Amount:            amount,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// MarkProcessing 标记为处理中
func (p *PrintJob) MarkProcessing() {
	p.Status = PrintJobStatusProcessing
	p.UpdatedAt = time.Now()
}

// MarkCompleted 标记为已完成
func (p *PrintJob) MarkCompleted() {
	now := time.Now()
	p.Status = PrintJobStatusCompleted
	p.CompletedAt = &now
	p.UpdatedAt = now
}

// MarkFailed 标记为失败
func (p *PrintJob) MarkFailed(errorMsg string) {
	p.Status = PrintJobStatusFailed
	p.ErrorMsg = errorMsg
	p.UpdatedAt = time.Now()
}

// GetTypeText 获取类型文本
func (p *PrintJob) GetTypeText() string {
	typeMap := map[PrintJobType]string{
		PrintJobTypeBill:    "账单收据",
		PrintJobTypeLease:   "租约合同",
		PrintJobTypeInvoice: "发票",
	}
	if text, ok := typeMap[p.Type]; ok {
		return text
	}
	return string(p.Type)
}

// GetStatusText 获取状态文本
func (p *PrintJob) GetStatusText() string {
	statusMap := map[PrintJobStatus]string{
		PrintJobStatusPending:    "待处理",
		PrintJobStatusProcessing: "处理中",
		PrintJobStatusCompleted:  "已完成",
		PrintJobStatusFailed:     "失败",
	}
	if text, ok := statusMap[p.Status]; ok {
		return text
	}
	return string(p.Status)
}

// Equals 比较打印作业是否相等
func (p *PrintJob) Equals(other *PrintJob) bool {
	return p.ID() == other.ID()
}
