package middleware

import (
	"github.com/zouhang1992/ddd_domain/internal/application/auth"
	"github.com/zouhang1992/ddd_domain/internal/application/config"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/persistence"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module Fx 模块
var Module = fx.Options(
	fx.Provide(func(cfg config.Config, repo any, toClaims persistence.ToClaimsFunc, svc *auth.OIDCService, log *zap.Logger) *AuthMiddleware {
		return NewAuthMiddleware(repo, toClaims, svc, cfg.OIDC.DevMode, log)
	}),
	fx.Provide(
		NewRBACMiddleware,
	),
)
