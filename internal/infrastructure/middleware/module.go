package middleware

import (
	"github.com/zouhang1992/ddd_domain/internal/application/auth"
	"github.com/zouhang1992/ddd_domain/internal/application/config"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/persistence/sqlite"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module Fx 模块
var Module = fx.Options(
	fx.Provide(func(cfg config.Config, repo *sqlite.SessionRepository, svc *auth.OIDCService, log *zap.Logger) *AuthMiddleware {
		return NewAuthMiddleware(repo, svc, cfg.OIDC.DevMode, log)
	}),
	fx.Provide(
		NewRBACMiddleware,
	),
)
