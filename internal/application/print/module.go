package print

import (
	"go.uber.org/fx"

	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
	locationrepo "github.com/zouhang1992/ddd_domain/internal/domain/location/repository"
	landlordrepo "github.com/zouhang1992/ddd_domain/internal/domain/landlord/repository"
	printservice "github.com/zouhang1992/ddd_domain/internal/domain/print/service"
	printrepo "github.com/zouhang1992/ddd_domain/internal/domain/print/repository"
)

// Module provides print application components
var Module = fx.Options(
	fx.Provide(NewCommandHandler),
	fx.Provide(NewQueryHandler),
	fx.Provide(func(
		billRepo billrepo.BillRepository,
		leaseRepo leaserepo.LeaseRepository,
		roomRepo roomrepo.RoomRepository,
		locationRepo locationrepo.LocationRepository,
		landlordRepo landlordrepo.LandlordRepository,
		printJobRepo printrepo.PrintJobRepository,
	) *printservice.PrintService {
		return printservice.NewPrintService(billRepo, leaseRepo, roomRepo, locationRepo, landlordRepo, printJobRepo)
	}),
)
