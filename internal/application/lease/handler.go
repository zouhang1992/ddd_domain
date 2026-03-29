package lease

import (
	"time"

	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/common"
	leasemodel "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	depositmodel "github.com/zouhang1992/ddd_domain/internal/domain/deposit/model"
	depositrepo "github.com/zouhang1992/ddd_domain/internal/domain/deposit/repository"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// CommandHandler 租约命令处理器
type CommandHandler struct {
	repo        leaserepo.LeaseRepository
	depositRepo depositrepo.DepositRepository
	billRepo    billrepo.BillRepository
	eventBus    *event.Bus
}

// NewCommandHandler 创建租约命令处理器
func NewCommandHandler(repo leaserepo.LeaseRepository, depositRepo depositrepo.DepositRepository, billRepo billrepo.BillRepository, eventBus *event.Bus) *CommandHandler {
	return &CommandHandler{repo: repo, depositRepo: depositRepo, billRepo: billRepo, eventBus: eventBus}
}

// HandleCreateLease 处理创建租约命令
func (h *CommandHandler) HandleCreateLease(cmd common.Command) (any, error) {
	createCmd, ok := cmd.(CreateLeaseCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := createCmd.Validate(); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	lease := leasemodel.NewLease(id, createCmd.RoomID, createCmd.LandlordID, createCmd.TenantName, createCmd.TenantPhone, createCmd.StartDate, createCmd.EndDate, createCmd.RentAmount, createCmd.DepositAmount, createCmd.Note)
	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// 创建押金
	if createCmd.DepositAmount > 0 {
		depositID := uuid.NewString()
		deposit := depositmodel.NewDeposit(depositID, id, createCmd.DepositAmount, createCmd.DepositNote)
		if err := h.depositRepo.Save(deposit); err != nil {
			return nil, err
		}
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range lease.Events() {
			h.eventBus.PublishAsync(evt)
		}
		lease.ClearEvents()
	}

	return lease, nil
}

// HandleUpdateLease 处理更新租约命令
func (h *CommandHandler) HandleUpdateLease(cmd common.Command) (any, error) {
	updateCmd, ok := cmd.(UpdateLeaseCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := updateCmd.Validate(); err != nil {
		return nil, err
	}

	lease, err := h.repo.FindByID(updateCmd.ID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, domerrors.ErrNotFound
	}

	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range lease.Events() {
			h.eventBus.PublishAsync(evt)
		}
		lease.ClearEvents()
	}

	return lease, nil
}

// HandleDeleteLease 处理删除租约命令
func (h *CommandHandler) HandleDeleteLease(cmd common.Command) (any, error) {
	deleteCmd, ok := cmd.(DeleteLeaseCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := deleteCmd.Validate(); err != nil {
		return nil, err
	}

	hasBills, err := h.repo.HasBills(deleteCmd.ID)
	if err != nil {
		return nil, err
	}
	if hasBills {
		return nil, domerrors.ErrCannotDelete
	}

	hasDeposit, err := h.repo.HasDeposit(deleteCmd.ID)
	if err != nil {
		return nil, err
	}
	if hasDeposit {
		return nil, domerrors.ErrCannotDelete
	}

	if err := h.repo.Delete(deleteCmd.ID); err != nil {
		return nil, err
	}

	return nil, nil
}

// HandleRenewLease 处理续租命令
func (h *CommandHandler) HandleRenewLease(cmd common.Command) (any, error) {
	renewCmd, ok := cmd.(RenewLeaseCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := renewCmd.Validate(); err != nil {
		return nil, err
	}

	lease, err := h.repo.FindByID(renewCmd.ID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, domerrors.ErrNotFound
	}

	lease.Renew(renewCmd.NewEndDate)
	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range lease.Events() {
			h.eventBus.PublishAsync(evt)
		}
		lease.ClearEvents()
	}

	return lease, nil
}

// HandleCheckoutLease 处理退租命令
func (h *CommandHandler) HandleCheckoutLease(cmd common.Command) (any, error) {
	checkoutCmd, ok := cmd.(CheckoutLeaseCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := checkoutCmd.Validate(); err != nil {
		return nil, err
	}

	lease, err := h.repo.FindByID(checkoutCmd.ID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, domerrors.ErrNotFound
	}

	lease.Checkout()
	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range lease.Events() {
			h.eventBus.PublishAsync(evt)
		}
		lease.ClearEvents()
	}

	return lease, nil
}

// HandleActivateLease 处理租约生效命令
func (h *CommandHandler) HandleActivateLease(cmd common.Command) (any, error) {
	activateCmd, ok := cmd.(ActivateLeaseCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := activateCmd.Validate(); err != nil {
		return nil, err
	}

	lease, err := h.repo.FindByID(activateCmd.ID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, domerrors.ErrNotFound
	}

	if lease.Status != leasemodel.LeaseStatusPending {
		return nil, domerrors.ErrInvalidState
	}

	// 检查开始日期是否已到
	if lease.StartDate.After(time.Now()) {
		return nil, domerrors.ErrInvalidState
	}

	lease.Activate()
	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// Publish events from aggregate
	if h.eventBus != nil {
		for _, evt := range lease.Events() {
			h.eventBus.PublishAsync(evt)
		}
		lease.ClearEvents()
	}

	return lease, nil
}

// QueryHandler 租约查询处理器
type QueryHandler struct {
	repo leaserepo.LeaseRepository
}

// NewQueryHandler 创建租约查询处理器
func NewQueryHandler(repo leaserepo.LeaseRepository) *QueryHandler {
	return &QueryHandler{repo: repo}
}

// HandleGetLease 处理获取租约查询
func (h *QueryHandler) HandleGetLease(q common.Query) (any, error) {
	getQuery, ok := q.(GetLeaseQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	lease, err := h.repo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, domerrors.ErrNotFound
	}

	return &LeaseQueryResult{Lease: lease}, nil
}

// HandleListLeases 处理列出租约查询
func (h *QueryHandler) HandleListLeases(q common.Query) (any, error) {
	listQuery, ok := q.(ListLeasesQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	leases, err := h.repo.FindAll()
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
	var paginated []*leasemodel.Lease
	total := len(leases)
	start := listQuery.Offset
	if start < total {
		end := start + limit
		if end > total {
			end = total
		}
		paginated = leases[start:end]
	}

	result := &LeasesQueryResult{
		Items: paginated,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}
