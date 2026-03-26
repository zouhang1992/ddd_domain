package sqlite

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides all persistence components
var Module = fx.Options(
	fx.Provide(func(cfg Config, logger *zap.Logger) (*Connection, error) {
		return NewConnection(cfg, logger)
	}),
	fx.Provide(NewLandlordRepository),
	fx.Provide(NewLeaseRepository),
	fx.Provide(NewBillRepository),
	fx.Provide(NewDepositRepository),
	fx.Provide(NewLocationRepository),
	fx.Provide(NewRoomRepository),
	fx.Provide(NewOperationLogRepository),
	fx.Provide(NewPrintJobRepository),
)
