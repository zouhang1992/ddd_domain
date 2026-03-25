package handler

import (
	"time"

	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// LeaseCommandHandler 租约命令处理器
type LeaseCommandHandler struct {
	repo        repository.LeaseRepository
	depositRepo repository.DepositRepository
	billRepo    repository.BillRepository
	eventBus    *event.Bus
}

// NewLeaseCommandHandler 创建租约命令处理器
func NewLeaseCommandHandler(repo repository.LeaseRepository, depositRepo repository.DepositRepository, billRepo repository.BillRepository, eventBus *event.Bus) *LeaseCommandHandler {
	return &LeaseCommandHandler{repo: repo, depositRepo: depositRepo, billRepo: billRepo, eventBus: eventBus}
}

// HandleCreateLease 处理创建租约命令
func (h *LeaseCommandHandler) HandleCreateLease(cmd command.Command) (any, error) {
	createCmd, ok := cmd.(command.CreateLeaseCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := createCmd.Validate(); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	lease := model.NewLease(id, createCmd.RoomID, createCmd.LandlordID, createCmd.TenantName, createCmd.TenantPhone, createCmd.StartDate, createCmd.EndDate, createCmd.RentAmount, createCmd.DepositAmount, createCmd.Note)
	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// 创建押金
	if createCmd.DepositAmount > 0 {
		depositID := uuid.NewString()
		deposit := model.NewDeposit(depositID, id, createCmd.DepositAmount, createCmd.DepositNote)
		if err := h.depositRepo.Save(deposit); err != nil {
			return nil, err
		}
	}

	// 发布租约创建事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLeaseCreated(lease))
	}

	return lease, nil
}

// HandleUpdateLease 处理更新租约命令
func (h *LeaseCommandHandler) HandleUpdateLease(cmd command.Command) (any, error) {
	updateCmd, ok := cmd.(command.UpdateLeaseCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := updateCmd.Validate(); err != nil {
		return nil, err
	}

	lease, err := h.repo.FindByID(updateCmd.ID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, model.ErrNotFound
	}

	lease.Update(updateCmd.TenantName, updateCmd.TenantPhone, updateCmd.StartDate, updateCmd.EndDate, updateCmd.RentAmount, updateCmd.Note)
	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// 发布租约更新事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLeaseUpdated(lease))
	}

	return lease, nil
}

// HandleDeleteLease 处理删除租约命令
func (h *LeaseCommandHandler) HandleDeleteLease(cmd command.Command) (any, error) {
	deleteCmd, ok := cmd.(command.DeleteLeaseCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := deleteCmd.Validate(); err != nil {
		return nil, err
	}

	hasBills, err := h.repo.HasBills(deleteCmd.ID)
	if err != nil {
		return nil, err
	}
	if hasBills {
		return nil, model.ErrCannotDelete
	}

	hasDeposit, err := h.repo.HasDeposit(deleteCmd.ID)
	if err != nil {
		return nil, err
	}
	if hasDeposit {
		return nil, model.ErrCannotDelete
	}

	if err := h.repo.Delete(deleteCmd.ID); err != nil {
		return nil, err
	}

	// 发布租约删除事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLeaseDeleted(deleteCmd.ID))
	}

	return nil, nil
}

// HandleRenewLease 处理续租命令
func (h *LeaseCommandHandler) HandleRenewLease(cmd command.Command) (any, error) {
	renewCmd, ok := cmd.(command.RenewLeaseCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := renewCmd.Validate(); err != nil {
		return nil, err
	}

	lease, err := h.repo.FindByID(renewCmd.ID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, model.ErrNotFound
	}

	if err := lease.Renew(renewCmd.NewStartDate, renewCmd.NewEndDate, renewCmd.NewRentAmount, renewCmd.Note); err != nil {
		return nil, err
	}

	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// 发布租约续租事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLeaseRenewed(lease))
	}

	return lease, nil
}

// HandleCheckoutLease 处理退租命令
func (h *LeaseCommandHandler) HandleCheckoutLease(cmd command.Command) (any, error) {
	checkoutCmd, ok := cmd.(command.CheckoutLeaseCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := checkoutCmd.Validate(); err != nil {
		return nil, err
	}

	lease, err := h.repo.FindByID(checkoutCmd.ID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, model.ErrNotFound
	}

	lease.Checkout()
	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// 发布租约退租事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLeaseCheckout(lease))
	}

	return lease, nil
}

// HandleActivateLease 处理租约生效命令
func (h *LeaseCommandHandler) HandleActivateLease(cmd command.Command) (any, error) {
	activateCmd, ok := cmd.(command.ActivateLeaseCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := activateCmd.Validate(); err != nil {
		return nil, err
	}

	lease, err := h.repo.FindByID(activateCmd.ID)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, model.ErrNotFound
	}

	if lease.Status != model.LeaseStatusPending {
		return nil, model.ErrInvalidState
	}

	// 检查开始日期是否已到
	if lease.StartDate.After(time.Now()) {
		return nil, model.ErrInvalidState
	}

	lease.Activate()
	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// 发布租约生效事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewLeaseActivated(lease))
	}

	return lease, nil
}
