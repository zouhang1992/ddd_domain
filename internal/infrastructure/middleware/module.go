package middleware

import "go.uber.org/fx"

// Module Fx 模块
var Module = fx.Options(
	fx.Provide(
		NewAuthMiddleware,
		NewRBACMiddleware,
	),
)
