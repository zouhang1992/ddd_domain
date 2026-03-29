package room

import "go.uber.org/fx"

// Module provides room application components
var Module = fx.Options(
	fx.Provide(NewCommandHandler),
	fx.Provide(NewQueryHandler),
)
