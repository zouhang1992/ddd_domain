package mysql

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/zouhang1992/ddd_domain/internal/application/config"
)

// Module provides all persistence components
var Module = fx.Options(
	fx.Provide(func(cfg config.Config, logger *zap.Logger) (*Connection, error) {
		return NewConnection(cfg.Database, logger)
	}),
	fx.Provide(NewLandlordRepository),
	fx.Provide(NewLeaseRepository),
	fx.Provide(NewBillRepository),
	fx.Provide(NewDepositRepository),
	fx.Provide(NewLocationRepository),
	fx.Provide(NewRoomRepository),
	fx.Provide(NewOperationLogRepository),
	fx.Provide(NewPrintJobRepository),
	fx.Provide(NewSessionRepository),
)
