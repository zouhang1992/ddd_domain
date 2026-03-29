package application

import (
	"github.com/zouhang1992/ddd_domain/internal/application/bill"
	"github.com/zouhang1992/ddd_domain/internal/application/common"
	"github.com/zouhang1992/ddd_domain/internal/application/landlord"
	"github.com/zouhang1992/ddd_domain/internal/application/lease"
	"github.com/zouhang1992/ddd_domain/internal/application/location"
	"github.com/zouhang1992/ddd_domain/internal/application/print"
	"github.com/zouhang1992/ddd_domain/internal/application/room"
	"go.uber.org/fx"
)

// Module provides all application components
var Module = fx.Options(
	landlord.Module,
	lease.Module,
	bill.Module,
	room.Module,
	location.Module,
	print.Module,
	common.Module,
)
