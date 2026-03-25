package handler

import (
	"go.uber.org/fx"
)

// Module provides all event handlers
var Module = fx.Options(
	fx.Provide(NewHandler),
)
