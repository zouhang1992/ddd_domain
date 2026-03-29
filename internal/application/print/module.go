package print

import "go.uber.org/fx"

// Module provides print application components
var Module = fx.Options(
	fx.Provide(NewCommandHandler),
	fx.Provide(NewQueryHandler),
)
