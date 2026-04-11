package persistence

import (
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/zouhang1992/ddd_domain/internal/application/config"
	billrepo "github.com/zouhang1992/ddd_domain/internal/domain/bill/repository"
	depositrepo "github.com/zouhang1992/ddd_domain/internal/domain/deposit/repository"
	landlordrepo "github.com/zouhang1992/ddd_domain/internal/domain/landlord/repository"
	leaserepo "github.com/zouhang1992/ddd_domain/internal/domain/lease/repository"
	locationrepo "github.com/zouhang1992/ddd_domain/internal/domain/location/repository"
	operationlogrepo "github.com/zouhang1992/ddd_domain/internal/domain/operationlog/repository"
	printrepo "github.com/zouhang1992/ddd_domain/internal/domain/print/repository"
	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/persistence/mysql"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/persistence/sqlite"
)

// Module provides persistence components based on configuration
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideConnection,
				fx.As(new(any)),
			),
		),
		fx.Provide(
			fx.Annotate(
				provideLandlordRepository,
				fx.As(new(landlordrepo.LandlordRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				provideLeaseRepository,
				fx.As(new(leaserepo.LeaseRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				provideBillRepository,
				fx.As(new(billrepo.BillRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				provideDepositRepository,
				fx.As(new(depositrepo.DepositRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				provideLocationRepository,
				fx.As(new(locationrepo.LocationRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				provideRoomRepository,
				fx.As(new(roomrepo.RoomRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				provideOperationLogRepository,
				fx.As(new(operationlogrepo.OperationLogRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				providePrintJobRepository,
				fx.As(new(printrepo.PrintJobRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				provideSessionRepository,
				fx.As(new(any)),
			),
		),
		fx.Provide(
			fx.Annotate(
				provideToClaimsFunc,
				fx.As(new(ToClaimsFunc)),
			),
		),
		fx.Provide(
			fx.Annotate(
				provideFromClaimsFunc,
				fx.As(new(FromClaimsFunc)),
			),
		),
	)
}

// Connection interface that both SQLite and MySQL connections implement
type Connection interface {
	DB() any
	Close() error
}

// Repository provider functions
func provideConnection(cfg config.Config, logger *zap.Logger) (any, error) {
	switch cfg.Database.Type {
	case "mysql":
		logger.Info("Using MySQL database")
		return mysql.NewConnection(cfg.Database, logger)
	case "sqlite", "":
		logger.Info("Using SQLite database")
		return sqlite.NewConnection(cfg.Database, logger)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}
}

func provideLandlordRepository(conn any) (landlordrepo.LandlordRepository, error) {
	switch c := conn.(type) {
	case *sqlite.Connection:
		return sqlite.NewLandlordRepository(c), nil
	case *mysql.Connection:
		return mysql.NewLandlordRepository(c), nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}

func provideLeaseRepository(conn any) (leaserepo.LeaseRepository, error) {
	switch c := conn.(type) {
	case *sqlite.Connection:
		return sqlite.NewLeaseRepository(c), nil
	case *mysql.Connection:
		return mysql.NewLeaseRepository(c), nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}

func provideBillRepository(conn any) (billrepo.BillRepository, error) {
	switch c := conn.(type) {
	case *sqlite.Connection:
		return sqlite.NewBillRepository(c), nil
	case *mysql.Connection:
		return mysql.NewBillRepository(c), nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}

func provideDepositRepository(conn any) (depositrepo.DepositRepository, error) {
	switch c := conn.(type) {
	case *sqlite.Connection:
		return sqlite.NewDepositRepository(c), nil
	case *mysql.Connection:
		return mysql.NewDepositRepository(c), nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}

func provideLocationRepository(conn any) (locationrepo.LocationRepository, error) {
	switch c := conn.(type) {
	case *sqlite.Connection:
		return sqlite.NewLocationRepository(c), nil
	case *mysql.Connection:
		return mysql.NewLocationRepository(c), nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}

func provideRoomRepository(conn any) (roomrepo.RoomRepository, error) {
	switch c := conn.(type) {
	case *sqlite.Connection:
		return sqlite.NewRoomRepository(c), nil
	case *mysql.Connection:
		return mysql.NewRoomRepository(c), nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}

func provideOperationLogRepository(conn any) (operationlogrepo.OperationLogRepository, error) {
	switch c := conn.(type) {
	case *sqlite.Connection:
		return sqlite.NewOperationLogRepository(c), nil
	case *mysql.Connection:
		return mysql.NewOperationLogRepository(c), nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}

func providePrintJobRepository(conn any) (printrepo.PrintJobRepository, error) {
	switch c := conn.(type) {
	case *sqlite.Connection:
		return sqlite.NewPrintJobRepository(c), nil
	case *mysql.Connection:
		return mysql.NewPrintJobRepository(c), nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}

func provideSessionRepository(conn any) (any, error) {
	switch c := conn.(type) {
	case *sqlite.Connection:
		return sqlite.NewSessionRepository(c), nil
	case *mysql.Connection:
		return mysql.NewSessionRepository(c), nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}

func provideToClaimsFunc(conn any) (ToClaimsFunc, error) {
	switch conn.(type) {
	case *sqlite.Connection:
		return sqlite.ToClaims, nil
	case *mysql.Connection:
		return mysql.ToClaims, nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}

func provideFromClaimsFunc(conn any) (FromClaimsFunc, error) {
	switch conn.(type) {
	case *sqlite.Connection:
		return sqlite.FromClaims, nil
	case *mysql.Connection:
		return mysql.FromClaims, nil
	default:
		return nil, fmt.Errorf("unsupported connection type")
	}
}
