package operationlog

import (
	"go.uber.org/fx"
)

// Module provides operationlog application components
var Module = fx.Options(
	fx.Provide(NewOperationLogEventHandler),
)
