package bus

import (
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busevent "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"go.uber.org/fx"
)

// Module provides all bus components
var Module = fx.Options(
	fx.Provide(buscommand.NewBus),
	fx.Provide(busevent.NewBus),
	fx.Provide(busquery.NewBus),
)
