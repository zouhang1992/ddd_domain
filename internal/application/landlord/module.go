package landlord

import "go.uber.org/fx"

// Module provides landlord application components
var Module = fx.Options(
	fx.Provide(NewCommandHandler),
	fx.Provide(NewQueryHandler),
)
