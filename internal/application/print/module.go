package print

import (
	"go.uber.org/fx"

	printservice "github.com/zouhang1992/ddd_domain/internal/domain/print/service"
)

// Module provides print application components
var Module = fx.Options(
	fx.Provide(NewCommandHandler),
	fx.Provide(NewQueryHandler),
	fx.Provide(printservice.NewPrintService),
)
