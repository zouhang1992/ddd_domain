package facade

import (
	"go.uber.org/fx"
)

// Module provides all HTTP handlers
var Module = fx.Options(
	fx.Provide(NewCQRSLandlordHandler),
	fx.Provide(NewCQRSLeaseHandler),
	fx.Provide(NewCQRSBillHandler),
	fx.Provide(NewCQRSLocationHandler),
	fx.Provide(NewCQRSRoomHandler),
	fx.Provide(NewCQRSPrintHandler),
	fx.Provide(NewAuthHandler),
	fx.Provide(NewIncomeHandler),
	fx.Provide(NewCQRSOperationLogHandler),
)
