package bill

import "go.uber.org/fx"

// Module provides bill application components
var Module = fx.Options(
	fx.Provide(NewCommandHandler),
	fx.Provide(NewQueryHandler),
)
