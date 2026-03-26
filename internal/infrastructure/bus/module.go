package bus

import (
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busevent "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides all bus components
var Module = fx.Options(
	fx.Provide(func(logger *zap.Logger) *buscommand.Bus {
		return buscommand.NewBus(logger)
	}),
	fx.Provide(func(logger *zap.Logger) *busevent.Bus {
		return busevent.NewBus(logger)
	}),
	fx.Provide(func(logger *zap.Logger) *busquery.Bus {
		return busquery.NewBus(logger)
	}),
)
