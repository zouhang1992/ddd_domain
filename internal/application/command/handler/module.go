package handler

import (
	"go.uber.org/fx"
)

// Module provides all command handlers
var Module = fx.Options(
	fx.Provide(NewLandlordCommandHandler),
	fx.Provide(NewLeaseCommandHandler),
	fx.Provide(NewBillCommandHandler),
	fx.Provide(NewLocationCommandHandler),
	fx.Provide(NewRoomCommandHandler),
	fx.Provide(NewPrintCommandHandler),
)
