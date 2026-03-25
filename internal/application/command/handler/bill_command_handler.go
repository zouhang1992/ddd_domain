package handler

import (
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// BillCommandHandler 账单命令处理器
type BillCommandHandler struct {
	repo     repository.BillRepository
	eventBus *event.Bus
}

// NewBillCommandHandler 创建账单命令处理器
func NewBillCommandHandler(repo repository.BillRepository, eventBus *event.Bus) *BillCommandHandler {
	return &BillCommandHandler{repo: repo, eventBus: eventBus}
}

// HandleCreateBill 处理创建账单命令
func (h *BillCommandHandler) HandleCreateBill(cmd command.Command) (any, error) {
	createCmd, ok := cmd.(command.CreateBillCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := createCmd.Validate(); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	bill := model.NewBill(id, createCmd.LeaseID, createCmd.Type, createCmd.Amount, createCmd.RentAmount, createCmd.WaterAmount, createCmd.ElectricAmount, createCmd.OtherAmount, createCmd.PaidAt, createCmd.Note)
	if err := h.repo.Save(bill); err != nil {
		return nil, err
	}

	// 发布账单创建事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewBillCreated(bill))
	}

	return bill, nil
}

// HandleUpdateBill 处理更新账单命令
func (h *BillCommandHandler) HandleUpdateBill(cmd command.Command) (any, error) {
	updateCmd, ok := cmd.(command.UpdateBillCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := updateCmd.Validate(); err != nil {
		return nil, err
	}

	bill, err := h.repo.FindByID(updateCmd.ID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, model.ErrNotFound
	}

	bill.Update(updateCmd.Amount, updateCmd.RentAmount, updateCmd.WaterAmount, updateCmd.ElectricAmount, updateCmd.OtherAmount, updateCmd.PaidAt, updateCmd.Note)
	if err := h.repo.Save(bill); err != nil {
		return nil, err
	}

	// 发布账单更新事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewBillUpdated(bill))
	}

	// 如果账单已支付，发布支付事件
	if bill.PaidAt != nil {
		if h.eventBus != nil {
			h.eventBus.PublishAsync(model.NewBillPaid(bill))
		}
	}

	return bill, nil
}

// HandleDeleteBill 处理删除账单命令
func (h *BillCommandHandler) HandleDeleteBill(cmd command.Command) (any, error) {
	deleteCmd, ok := cmd.(command.DeleteBillCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := deleteCmd.Validate(); err != nil {
		return nil, err
	}

	if err := h.repo.Delete(deleteCmd.ID); err != nil {
		return nil, err
	}

	// 发布账单删除事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewBillDeleted(deleteCmd.ID))
	}

	return nil, nil
}

// HandleConfirmBillArrival 处理确认账单到账命令
func (h *BillCommandHandler) HandleConfirmBillArrival(cmd command.Command) (any, error) {
	confirmCmd, ok := cmd.(command.ConfirmBillArrivalCommand)
	if !ok {
		return nil, model.ErrInvalidCommand
	}

	if err := confirmCmd.Validate(); err != nil {
		return nil, err
	}

	bill, err := h.repo.FindByID(confirmCmd.ID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, model.ErrNotFound
	}

	bill.MarkPaid(confirmCmd.PaidAt)
	if err := h.repo.Save(bill); err != nil {
		return nil, err
	}

	// 发布账单更新事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewBillUpdated(bill))
	}

	// 发布账单支付事件
	if h.eventBus != nil {
		h.eventBus.PublishAsync(model.NewBillPaid(bill))
	}

	return bill, nil
}
