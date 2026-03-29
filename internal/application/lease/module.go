package lease

import "go.uber.org/fx"

// Module provides lease application components
var Module = fx.Options(
	fx.Provide(NewCommandHandler),
	fx.Provide(NewQueryHandler),
	fx.Provide(NewLeaseRoomEventHandler),
)
