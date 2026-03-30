package lease

import (
	"go.uber.org/fx"

	leaseservice "github.com/zouhang1992/ddd_domain/internal/domain/lease/service"
)

// Module provides lease application components
var Module = fx.Options(
	fx.Provide(NewCommandHandler),
	fx.Provide(NewQueryHandler),
	fx.Provide(NewLeaseRoomEventHandler),
	fx.Provide(NewLeaseExpirationScheduler),
	fx.Provide(leaseservice.NewLeaseService),
)
