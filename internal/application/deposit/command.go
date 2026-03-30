package deposit

import (
	domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
)

// DepositCommand 押金命令接口
type DepositCommand interface {
	CommandName() string
	Validate() error
}

// MarkReturningCommand 标记押金为待退还命令
type MarkReturningCommand struct {
	ID string
}

func (c MarkReturningCommand) CommandName() string {
	return "mark_deposit_returning"
}

func (c MarkReturningCommand) Validate() error {
	if c.ID == "" {
		return domerrors.ErrInvalidCommand
	}
	return nil
}

// MarkReturnedCommand 标记押金为已退还命令
type MarkReturnedCommand struct {
	ID string
}

func (c MarkReturnedCommand) CommandName() string {
	return "mark_deposit_returned"
}

func (c MarkReturnedCommand) Validate() error {
	if c.ID == "" {
		return domerrors.ErrInvalidCommand
	}
	return nil
}
