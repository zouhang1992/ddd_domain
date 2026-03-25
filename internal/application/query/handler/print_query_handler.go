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
}

// NewPrintQueryHandler 创建打印查询处理器
func NewPrintQueryHandler(billRepo repository.BillRepository, leaseRepo repository.LeaseRepository) *PrintQueryHandler {
	return &PrintQueryHandler{billRepo: billRepo, leaseRepo: leaseRepo}
}

// HandleGetPrintJob 处理获取打印作业查询
func (h *PrintQueryHandler) HandleGetPrintJob(q query.Query) (any, error) {
	// 暂时返回空实现，实际项目中需要有打印作业仓储
	return nil, nil
}

// HandleListPrintJobs 处理列出打印作业查询
func (h *PrintQueryHandler) HandleListPrintJobs(q query.Query) (any, error) {
	// 暂时返回空实现
	return nil, nil
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
