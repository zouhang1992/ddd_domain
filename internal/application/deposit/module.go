package deposit

import (
	"go.uber.org/fx"
)

// Module provides deposit components
var Module = fx.Options(
	fx.Provide(NewCommandHandler),
	fx.Provide(NewQueryHandler),
)
