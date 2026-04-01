package auth

import (
	"go.uber.org/fx"

	"github.com/zouhang1992/ddd_domain/internal/application/config"
)

// Module Fx 模块
var Module = fx.Options(
	fx.Provide(func(cfg config.Config) Config {
		return cfg.OIDC
	}),
	fx.Provide(
		NewOIDCService,
		NewRBACService,
	),
)
