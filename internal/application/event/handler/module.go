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
	fx.Provide(func(roomRepo repository.RoomRepository, logger *zap.Logger) *LeaseRoomEventHandler {
		return NewLeaseRoomEventHandler(roomRepo, logger)
	}),
)
