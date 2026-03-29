package location

import "go.uber.org/fx"

// Module provides location application components
var Module = fx.Options(
	fx.Provide(NewCommandHandler),
	fx.Provide(NewQueryHandler),
)
