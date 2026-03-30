package bill

import (
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/common"
	billmodel "github.com/zouhang1992/ddd_domain/internal/domain/bill/model"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// CommandHandler 账单命令处理器
type CommandHandler struct {
	repo     billrepo.BillRepository
	eventBus *event.Bus
}

// NewCommandHandler 创建账单命令处理器
func NewCommandHandler(repo billrepo.BillRepository, eventBus *event.Bus) *CommandHandler {
	return &CommandHandler{repo: repo, eventBus: eventBus}
}

// HandleCreateBill 处理创建账单命令
func (h *CommandHandler) HandleCreateBill(cmd common.Command) (any, error) {
	createCmd, ok := cmd.(CreateBillCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := createCmd.Validate(); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	var bill *billmodel.Bill

	// If detailed amounts are provided, use the new constructor
	if createCmd.RentAmount > 0 || createCmd.WaterAmount > 0 || createCmd.ElectricAmount > 0 || createCmd.OtherAmount > 0 {
		bill = billmodel.NewBillWithDetails(
			id,
			createCmd.LeaseID,
			createCmd.Type,
			createCmd.RentAmount,
			createCmd.WaterAmount,
			createCmd.ElectricAmount,
			createCmd.OtherAmount,
			createCmd.DueDate,
			createCmd.Note,
		)
	} else {
		// Fall back to traditional single amount
		bill = billmodel.NewBill(id, createCmd.LeaseID, createCmd.Type, createCmd.Amount, createCmd.DueDate, createCmd.Note)
	}

	if err := h.repo.Save(bill); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range bill.Events() {
			h.eventBus.PublishAsync(evt)
		}
		bill.ClearEvents()
	}

	return bill, nil
}

// HandleUpdateBill 处理更新账单命令
func (h *CommandHandler) HandleUpdateBill(cmd common.Command) (any, error) {
	updateCmd, ok := cmd.(UpdateBillCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := updateCmd.Validate(); err != nil {
		return nil, err
	}

	bill, err := h.repo.FindByID(updateCmd.ID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, domerrors.ErrNotFound
	}

	// If detailed amounts are provided, use the new update method
	if updateCmd.RentAmount > 0 || updateCmd.WaterAmount > 0 || updateCmd.ElectricAmount > 0 || updateCmd.OtherAmount > 0 {
		bill.UpdateWithDetails(
			updateCmd.RentAmount,
			updateCmd.WaterAmount,
			updateCmd.ElectricAmount,
			updateCmd.OtherAmount,
			updateCmd.DueDate,
			updateCmd.Note,
		)
	} else {
		// Fall back to traditional single amount update
		bill.Update(updateCmd.Amount, updateCmd.DueDate, updateCmd.Note)
	}

	if err := h.repo.Save(bill); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range bill.Events() {
			h.eventBus.PublishAsync(evt)
		}
		bill.ClearEvents()
	}

	return bill, nil
}

// HandleDeleteBill 处理删除账单命令
func (h *CommandHandler) HandleDeleteBill(cmd common.Command) (any, error) {
	deleteCmd, ok := cmd.(DeleteBillCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := deleteCmd.Validate(); err != nil {
		return nil, err
	}

	if err := h.repo.Delete(deleteCmd.ID); err != nil {
		return nil, err
	}

	return nil, nil
}

// HandleConfirmBillArrival 处理确认账单到账命令
func (h *CommandHandler) HandleConfirmBillArrival(cmd common.Command) (any, error) {
	confirmCmd, ok := cmd.(ConfirmBillArrivalCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := confirmCmd.Validate(); err != nil {
		return nil, err
	}

	bill, err := h.repo.FindByID(confirmCmd.ID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, domerrors.ErrNotFound
	}

	bill.MarkPaid()
	if err := h.repo.Save(bill); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range bill.Events() {
			h.eventBus.PublishAsync(evt)
		}
		bill.ClearEvents()
	}

	return bill, nil
}

// QueryHandler 账单查询处理器
type QueryHandler struct {
	billRepo  billrepo.BillRepository
	leaseRepo leaserepo.LeaseRepository
}

// NewQueryHandler 创建账单查询处理器
func NewQueryHandler(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository) *QueryHandler {
	return &QueryHandler{billRepo: billRepo, leaseRepo: leaseRepo}
}

// HandleGetBill 处理获取账单查询
func (h *QueryHandler) HandleGetBill(q common.Query) (any, error) {
	getQuery, ok := q.(GetBillQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	bill, err := h.billRepo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, domerrors.ErrNotFound
	}

	return &BillQueryResult{Bill: bill}, nil
}

// HandleListBills 处理列出账单查询
func (h *QueryHandler) HandleListBills(q common.Query) (any, error) {
	listQuery, ok := q.(ListBillsQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	bills, err := h.billRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 设置默认分页大小
	limit := listQuery.Limit
	if limit <= 0 {
		limit = 10
	}

	// 计算页码
	page := 1
	if listQuery.Offset > 0 && limit > 0 {
		page = (listQuery.Offset / limit) + 1
	}

	// Simple pagination
	var paginated []*billmodel.Bill
	total := len(bills)
	start := listQuery.Offset
	if start < total {
		end := start + limit
		if end > total {
			end = total
		}
		paginated = bills[start:end]
	}

	result := &BillsQueryResult{
		Items: paginated,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}

// HandleIncomeReport 处理收入报表查询
func (h *QueryHandler) HandleIncomeReport(q common.Query) (any, error) {
	reportQuery, ok := q.(IncomeReportQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	bills, err := h.billRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var result IncomeReportQueryResult
	result.Year = reportQuery.Year
	result.Month = int(reportQuery.Month)

	for _, bill := range bills {
		if bill.PaidAt != nil {
			result.Total += bill.Amount
		}
	}

	return &result, nil
}
