package handler

import (
	"go.uber.org/fx"
)

// Module provides all query handlers
var Module = fx.Options(
	fx.Provide(NewLandlordQueryHandler),
	fx.Provide(NewLeaseQueryHandler),
	fx.Provide(NewBillQueryHandler),
	fx.Provide(NewLocationQueryHandler),
	fx.Provide(NewRoomQueryHandler),
	fx.Provide(NewPrintQueryHandler),
	fx.Provide(NewOperationLogQueryHandler),
)
