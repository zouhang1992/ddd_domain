package deposit

import (
	"github.com/zouhang1992/ddd_domain/internal/application/common"
	depositmodel "github.com/zouhang1992/ddd_domain/internal/domain/deposit/model"
	depositrepo "github.com/zouhang1992/ddd_domain/internal/domain/deposit/repository"
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// CommandHandler 押金命令处理器
type CommandHandler struct {
	repo     depositrepo.DepositRepository
	eventBus *event.Bus
}

// NewCommandHandler 创建押金命令处理器
func NewCommandHandler(repo depositrepo.DepositRepository, eventBus *event.Bus) *CommandHandler {
	return &CommandHandler{repo: repo, eventBus: eventBus}
}

// HandleMarkReturning 处理标记押金为待退还命令
func (h *CommandHandler) HandleMarkReturning(cmd common.Command) (any, error) {
	markCmd, ok := cmd.(MarkReturningCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := markCmd.Validate(); err != nil {
		return nil, err
	}

	deposit, err := h.repo.FindByID(markCmd.ID)
	if err != nil {
		return nil, err
	}
	if deposit == nil {
		return nil, domerrors.ErrNotFound
	}

	deposit.MarkReturning()
	if err := h.repo.Save(deposit); err != nil {
		return nil, err
	}

	if h.eventBus != nil {
		for _, evt := range deposit.Events() {
			h.eventBus.PublishAsync(evt)
		}
		deposit.ClearEvents()
	}

	return deposit, nil
}

// HandleMarkReturned 处理标记押金为已退还命令
func (h *CommandHandler) HandleMarkReturned(cmd common.Command) (any, error) {
	markCmd, ok := cmd.(MarkReturnedCommand)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	if err := markCmd.Validate(); err != nil {
		return nil, err
	}

	deposit, err := h.repo.FindByID(markCmd.ID)
	if err != nil {
		return nil, err
	}
	if deposit == nil {
		return nil, domerrors.ErrNotFound
	}

	deposit.MarkReturned()
	if err := h.repo.Save(deposit); err != nil {
		return nil, err
	}

	if h.eventBus != nil {
		for _, evt := range deposit.Events() {
			h.eventBus.PublishAsync(evt)
		}
		deposit.ClearEvents()
	}

	return deposit, nil
}

// QueryHandler 押金查询处理器
type QueryHandler struct {
	repo depositrepo.DepositRepository
}

// NewQueryHandler 创建押金查询处理器
func NewQueryHandler(repo depositrepo.DepositRepository) *QueryHandler {
	return &QueryHandler{repo: repo}
}

// HandleGetDeposit 处理获取押金查询
func (h *QueryHandler) HandleGetDeposit(q common.Query) (any, error) {
	getQuery, ok := q.(GetDepositQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	deposit, err := h.repo.FindByID(getQuery.ID)
	if err != nil {
		return nil, err
	}
	if deposit == nil {
		return nil, domerrors.ErrNotFound
	}

	return &DepositQueryResult{Deposit: deposit}, nil
}

// HandleListDeposits 处理列出押金查询
func (h *QueryHandler) HandleListDeposits(q common.Query) (any, error) {
	listQuery, ok := q.(ListDepositsQuery)
	if !ok {
		return nil, domerrors.ErrInvalidCommand
	}

	var deposits []*depositmodel.Deposit
	var err error

	if listQuery.LeaseID != "" {
		deposit, err := h.repo.FindByLeaseID(listQuery.LeaseID)
		if err != nil {
			return nil, err
		}
		if deposit != nil {
			deposits = []*depositmodel.Deposit{deposit}
		}
	} else {
		deposits, err = h.repo.FindAll()
		if err != nil {
			return nil, err
		}
	}

	// 如果有状态过滤，应用过滤
	if listQuery.Status != "" {
		var filtered []*depositmodel.Deposit
		for _, d := range deposits {
			if string(d.Status) == listQuery.Status {
				filtered = append(filtered, d)
			}
		}
		deposits = filtered
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

	// 简单分页
	var paginated []*depositmodel.Deposit
	total := len(deposits)
	start := listQuery.Offset
	if start < total {
		end := start + limit
		if end > total {
			end = total
		}
		paginated = deposits[start:end]
	}

	result := &DepositsQueryResult{
		Items: paginated,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}
