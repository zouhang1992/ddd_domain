package handler

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides all event handlers
var Module = fx.Options(
	fx.Provide(func(repo repository.OperationLogRepository, logger *zap.Logger) *Handler {
		return NewHandler(repo, logger)
	}),
)
