package auth

import (
	"go.uber.org/fx"
)

// Module Fx 模块
var Module = fx.Options(
	fx.Provide(
		NewOIDCService,
		NewRBACService,
	),
)
