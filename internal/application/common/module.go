package common

import (
	"github.com/zouhang1992/ddd_domain/internal/application/common/service"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
	locationrepo "github.com/zouhang1992/ddd_domain/internal/domain/location/repository"
	landlordrepo "github.com/zouhang1992/ddd_domain/internal/domain/landlord/repository"
	"go.uber.org/fx"
	"time"
)

// Module provides common application components
var Module = fx.Options(
	fx.Provide(func() *service.AuthService {
		return service.NewAuthService("", 7*24*time.Hour)
	}),
	fx.Provide(func(billRepo billrepo.BillRepository, leaseRepo leaserepo.LeaseRepository, roomRepo roomrepo.RoomRepository, locationRepo locationrepo.LocationRepository, landlordRepo landlordrepo.LandlordRepository) *service.PrintService {
		return service.NewPrintService(billRepo, leaseRepo, roomRepo, locationRepo, landlordRepo)
	}),
)
