package lease

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/zouhang1992/ddd_domain/internal/application/common"
	leasemodel "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	leaseservice "github.com/zouhang1992/ddd_domain/internal/domain/lease/service"
	depositrepo "github.com/zouhang1992/ddd_domain/internal/domain/deposit/repository"
	depositmodel "github.com/zouhang1992/ddd_domain/internal/domain/deposit/model"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	billmodel "github.com/zouhang1992/ddd_domain/internal/domain/bill/model"
	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// CommandHandler 租约命令处理器
type CommandHandler struct {
	repo         leaserepo.LeaseRepository
	depositRepo  depositrepo.DepositRepository
	billRepo     billrepo.BillRepository
	roomRepo     roomrepo.RoomRepository
	eventBus     *event.Bus
	leaseService *leaseservice.LeaseService
	log          *zap.Logger
}

// NewCommandHandler 创建租约命令处理器
func NewCommandHandler(repo leaserepo.LeaseRepository, depositRepo depositrepo.DepositRepository, billRepo billrepo.BillRepository, roomRepo roomrepo.RoomRepository, eventBus *event.Bus, leaseService *leaseservice.LeaseService, log *zap.Logger) *CommandHandler {
	return &CommandHandler{repo: repo, depositRepo: depositRepo, billRepo: billRepo, roomRepo: roomRepo, eventBus: eventBus, leaseService: leaseService, log: log}
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

	result, err := h.leaseService.CreateLease(
		createCmd.RoomID,
		createCmd.LandlordID,
		createCmd.TenantName,
		createCmd.TenantPhone,
		createCmd.StartDate,
		createCmd.EndDate,
		createCmd.RentAmount,
		createCmd.DepositAmount,
		createCmd.Note,
		createCmd.DepositNote,
	)
	if err != nil {
		return nil, err
	}

	if err := h.repo.Save(result.Lease); err != nil {
		return nil, err
	}

	if result.Deposit != nil {
		if err := h.depositRepo.Save(result.Deposit); err != nil {
			return nil, err
		}
	}

	if h.eventBus != nil {
		for _, evt := range result.Lease.Events() {
			h.eventBus.PublishAsync(evt)
		}
		result.Lease.ClearEvents()
	}

	return result.Lease, nil
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

	if err := h.leaseService.ValidateDelete(deleteCmd.ID); err != nil {
		return nil, err
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

	// 直接更新房间状态为可出租
	room, err := h.roomRepo.FindByID(lease.RoomID)
	if err == nil && room != nil {
		room.MarkAvailable()
		if err := h.roomRepo.Save(room); err != nil {
			h.log.Error("Failed to update room status to available",
				zap.String("room_id", room.ID()),
				zap.Error(err))
		}

		if h.eventBus != nil {
			for _, evt := range room.Events() {
				if err := h.eventBus.Publish(evt); err != nil {
					h.log.Error("Failed to publish event",
						zap.String("event", evt.EventName()),
						zap.Error(err))
				}
			}
			room.ClearEvents()
		}
	}

	if h.eventBus != nil {
		for _, evt := range lease.Events() {
			if err := h.eventBus.Publish(evt); err != nil {
				h.log.Error("Failed to publish event",
					zap.String("event", evt.EventName()),
					zap.Error(err))
			}
		}
		lease.ClearEvents()
	}

	return lease, nil
}

// CheckoutWithBillsResult 退租并创建结算账单结果
type CheckoutWithBillsResult struct {
	Lease        *leasemodel.Lease   `json:"lease"`
	CheckoutBill *billmodel.Bill      `json:"checkout_bill"`
	Deposit      *depositmodel.Deposit `json:"deposit,omitempty"`
}

// HandleCheckoutWithBills 处理退租并创建结算账单命令
func (h *CommandHandler) HandleCheckoutWithBills(cmd common.Command) (any, error) {
	checkoutCmd, ok := cmd.(CheckoutWithBillsCommand)
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

	// 1. Checkout the lease
	lease.Checkout()
	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// 2. Create checkout bill (refund amounts are negative, charges are positive)
	// Refunds are negative because they're money going out to the tenant
	// Charges are positive because they're money coming in from the tenant
	totalAmount := -checkoutCmd.RefundRentAmount - checkoutCmd.RefundDepositAmount + checkoutCmd.WaterAmount + checkoutCmd.ElectricAmount + checkoutCmd.OtherAmount

	checkoutBillID := uuid.NewString()
	checkoutBill := billmodel.NewBillWithDetails(
		checkoutBillID,
		lease.ID(),
		billmodel.BillTypeCheckout,
		checkoutCmd.WaterAmount,    // Positive: tenant owes water
		checkoutCmd.ElectricAmount, // Positive: tenant owes electric
		0,                           // Other utility charges go here
		checkoutCmd.OtherAmount,    // Other charges
		time.Now(),
		checkoutCmd.Note,
	)
	// Set refund amounts as negative rent refund
	checkoutBill.RentAmount = -checkoutCmd.RefundRentAmount // Negative: refund to tenant
	checkoutBill.Amount = totalAmount

	if err := h.billRepo.Save(checkoutBill); err != nil {
		return nil, err
	}

	// 3. Find and update deposit status - always mark as returned directly on checkout
	var deposit *depositmodel.Deposit
	deposit, err = h.depositRepo.FindByLeaseID(lease.ID())
	if err == nil && deposit != nil {
		// 退租时直接将押金状态改为退还
		deposit.MarkReturned()
		if err := h.depositRepo.Save(deposit); err != nil {
			return nil, err
		}
	}

	// 直接更新房间状态为可出租
	room, err := h.roomRepo.FindByID(lease.RoomID)
	if err == nil && room != nil {
		room.MarkAvailable()
		if err := h.roomRepo.Save(room); err != nil {
			h.log.Error("Failed to update room status to available",
				zap.String("room_id", room.ID()),
				zap.Error(err))
		}

		if h.eventBus != nil {
			for _, evt := range room.Events() {
				if err := h.eventBus.Publish(evt); err != nil {
					h.log.Error("Failed to publish event",
						zap.String("event", evt.EventName()),
						zap.Error(err))
				}
			}
			room.ClearEvents()
		}
	}

	// Publish events from aggregates
	if h.eventBus != nil {
		for _, evt := range lease.Events() {
			if err := h.eventBus.Publish(evt); err != nil {
				h.log.Error("Failed to publish event",
					zap.String("event", evt.EventName()),
					zap.Error(err))
			}
		}
		lease.ClearEvents()
		for _, evt := range checkoutBill.Events() {
			if err := h.eventBus.Publish(evt); err != nil {
				h.log.Error("Failed to publish event",
					zap.String("event", evt.EventName()),
					zap.Error(err))
			}
		}
		checkoutBill.ClearEvents()
		if deposit != nil {
			for _, evt := range deposit.Events() {
				if err := h.eventBus.Publish(evt); err != nil {
					h.log.Error("Failed to publish event",
						zap.String("event", evt.EventName()),
						zap.Error(err))
				}
			}
			deposit.ClearEvents()
		}
	}

	return &CheckoutWithBillsResult{
		Lease:        lease,
		CheckoutBill: checkoutBill,
		Deposit:      deposit,
	}, nil
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

	room, err := h.roomRepo.FindByID(lease.RoomID)
	if err != nil {
		return nil, err
	}

	if err := h.leaseService.ValidateActivate(lease, room); err != nil {
		return nil, err
	}

	lease.Activate()
	if err := h.repo.Save(lease); err != nil {
		return nil, err
	}

	// 直接更新房间状态为已出租
	room.MarkRented()
	if err := h.roomRepo.Save(room); err != nil {
		h.log.Error("Failed to update room status to rented",
			zap.String("room_id", room.ID()),
			zap.Error(err))
		// 不返回错误，因为租约已经激活成功了
	}

	if h.eventBus != nil {
		for _, evt := range lease.Events() {
			if err := h.eventBus.Publish(evt); err != nil {
				h.log.Error("Failed to publish event",
					zap.String("event", evt.EventName()),
					zap.Error(err))
			}
		}
		lease.ClearEvents()
		for _, evt := range room.Events() {
			if err := h.eventBus.Publish(evt); err != nil {
				h.log.Error("Failed to publish event",
					zap.String("event", evt.EventName()),
					zap.Error(err))
			}
		}
		room.ClearEvents()
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

	limit := listQuery.Limit
	if limit <= 0 {
		limit = 10
	}

	page := 1
	if listQuery.Offset > 0 && limit > 0 {
		page = (listQuery.Offset / limit) + 1
	}

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
