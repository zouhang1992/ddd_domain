package handler

import (
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
	"time"
)

// BillQueryHandler 账单查询处理器
type BillQueryHandler struct {
	billRepo  repository.BillRepository
	leaseRepo repository.LeaseRepository
}

// NewBillQueryHandler 创建账单查询处理器
func NewBillQueryHandler(billRepo repository.BillRepository, leaseRepo repository.LeaseRepository) *BillQueryHandler {
	return &BillQueryHandler{billRepo: billRepo, leaseRepo: leaseRepo}
}

// HandleGetBill 处理获取账单查询
func (h *BillQueryHandler) HandleGetBill(q query.Query) (any, error) {
	getQuery, ok := q.(query.GetBillQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	bill, err := h.billRepo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, model.ErrNotFound
	}

	return &query.BillQueryResult{Bill: bill}, nil
}

// HandleListBills 处理列出账单查询
func (h *BillQueryHandler) HandleListBills(q query.Query) (any, error) {
	listQuery, ok := q.(query.ListBillsQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	// 构建查询条件
	criteria := repository.BillCriteria{
		Type:        listQuery.Type,
		Status:      listQuery.Status,
		LeaseID:     listQuery.LeaseID,
		RoomID:      listQuery.RoomID,
		Month:       listQuery.Month,
		MinAmount:   listQuery.MinAmount,
		MaxAmount:   listQuery.MaxAmount,
		StartDate:   listQuery.StartDate,
		EndDate:     listQuery.EndDate,
	}

	// 设置默认分页大小
	limit := listQuery.Limit
	if limit <= 0 {
		limit = 10 // 默认返回10条
	}

	// 查询数据
	bills, err := h.billRepo.FindByCriteria(criteria, listQuery.Offset, limit)
	if err != nil {
		return nil, err
	}

	// 获取总数
	total, err := h.billRepo.CountByCriteria(criteria)
	if err != nil {
		return nil, err
	}

	// 计算页码
	page := 1
	if listQuery.Offset > 0 && limit > 0 {
		page = (listQuery.Offset / limit) + 1
	}

	result := &query.BillsQueryResult{
		Items: bills,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}

// HandleIncomeReport 处理收入报表查询
func (h *BillQueryHandler) HandleIncomeReport(q query.Query) (any, error) {
	reportQuery, ok := q.(query.IncomeReportQuery)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	bills, err := h.billRepo.FindByMonth(reportQuery.Year, reportQuery.Month)
	if err != nil {
		return nil, err
	}

	var result query.IncomeReportQueryResult
	result.Year = reportQuery.Year
	result.Month = int(reportQuery.Month)

	for _, bill := range bills {
		if bill.Status == model.BillStatusPaid && bill.PaidAt != nil {
			result.RentTotal += bill.RentAmount
			result.WaterTotal += bill.WaterAmount
			result.ElectricTotal += bill.ElectricAmount
			result.OtherTotal += bill.OtherAmount
			result.Total += bill.Amount
		}
	}

	return &result, nil
}

// parseMonth 解析月份（简单实现）
func parseMonth(monthStr string, year *int, mon *time.Month) (bool, error) {
	if len(monthStr) != 7 {
		return false, nil
	}

	t, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return false, err
	}

	*year = t.Year()
	*mon = t.Month()
	return true, nil
}
