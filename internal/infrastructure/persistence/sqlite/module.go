package sqlite

import (
	"go.uber.org/fx"
)

// Module provides all persistence components
var Module = fx.Options(
	fx.Provide(NewConnection),
	fx.Provide(NewLandlordRepository),
	fx.Provide(NewLeaseRepository),
	fx.Provide(NewBillRepository),
	fx.Provide(NewDepositRepository),
	fx.Provide(NewLocationRepository),
	fx.Provide(NewRoomRepository),
	fx.Provide(NewOperationLogRepository),
)
