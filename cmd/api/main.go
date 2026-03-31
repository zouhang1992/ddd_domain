package main

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/zouhang1992/ddd_domain/internal/application"
	"github.com/zouhang1992/ddd_domain/internal/application/bill"
	"github.com/zouhang1992/ddd_domain/internal/application/common/service"
	"github.com/zouhang1992/ddd_domain/internal/application/deposit"
	"github.com/zouhang1992/ddd_domain/internal/application/landlord"
	"github.com/zouhang1992/ddd_domain/internal/application/lease"
	"github.com/zouhang1992/ddd_domain/internal/application/location"
	"github.com/zouhang1992/ddd_domain/internal/application/operationlog"
	"github.com/zouhang1992/ddd_domain/internal/application/print"
	"github.com/zouhang1992/ddd_domain/internal/application/room"
	"github.com/zouhang1992/ddd_domain/internal/facade"
	busmodule "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busevent "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	logging "github.com/zouhang1992/ddd_domain/internal/infrastructure/logging"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/persistence/sqlite"
)

func main() {
	fx.New(
		// 配置模块
		fx.Options(
			fx.Provide(func() sqlite.Config {
				return sqlite.Config{DSN: "data/ddd.db"}
			}),
		),
		// 各个组件模块
		logging.Module(),
		sqlite.Module,
		busmodule.Module,
		application.Module,
		facade.Module,
		// 注册处理器到总线
		fx.Invoke(registerCommandHandlers),
		fx.Invoke(registerQueryHandlers),
		fx.Invoke(registerEventHandlers),
		// 启动定时任务
		fx.Invoke(startLeaseExpirationScheduler),
		// 设置服务器
		fx.Invoke(startServer),
	).Run()
}

func registerEventHandlers(eventBus *busevent.Bus, leaseRoomHandler *lease.LeaseRoomEventHandler, operationLogHandler *operationlog.OperationLogEventHandler, logger *zap.Logger) {
	// 订阅租约事件以处理房间状态变更
	eventBus.Subscribe("lease.activated", leaseRoomHandler)
	eventBus.Subscribe("lease.checkout", leaseRoomHandler)
	eventBus.Subscribe("lease.expired", leaseRoomHandler)

	logger.Info("Lease room event handler registered",
		zap.String("events", "lease.activated, lease.checkout, lease.expired"))

	// 订阅所有领域事件以记录操作日志
	allEvents := []string{
		// Room events
		"room.created", "room.updated", "room.deleted", "room.rented", "room.available",
		// Lease events
		"lease.created", "lease.activated", "lease.checkout", "lease.expired", "lease.renewed", "lease.deleted",
		// Landlord events
		"landlord.created", "landlord.updated", "landlord.deleted",
		// Bill events
		"bill.created", "bill.updated", "bill.paid", "bill.deleted",
		// Location events
		"location.created", "location.updated", "location.deleted",
		// Deposit events
		"deposit.created", "deposit.returning", "deposit.returned", "deposit.deleted",
	}

	for _, evt := range allEvents {
		eventBus.Subscribe(evt, operationLogHandler)
	}

	logger.Info("Operation log event handler registered",
		zap.Strings("events", allEvents))
}

func registerCommandHandlers(
	bus *buscommand.Bus,
	landlordHandler *landlord.CommandHandler,
	leaseHandler *lease.CommandHandler,
	billHandler *bill.CommandHandler,
	depositHandler *deposit.CommandHandler,
	locationHandler *location.CommandHandler,
	roomHandler *room.CommandHandler,
	printHandler *print.CommandHandler,
) {
	bus.Register("create_landlord", buscommand.HandlerFunc(landlordHandler.HandleCreateLandlord))
	bus.Register("update_landlord", buscommand.HandlerFunc(landlordHandler.HandleUpdateLandlord))
	bus.Register("delete_landlord", buscommand.HandlerFunc(landlordHandler.HandleDeleteLandlord))
	bus.Register("create_lease", buscommand.HandlerFunc(leaseHandler.HandleCreateLease))
	bus.Register("update_lease", buscommand.HandlerFunc(leaseHandler.HandleUpdateLease))
	bus.Register("delete_lease", buscommand.HandlerFunc(leaseHandler.HandleDeleteLease))
	bus.Register("renew_lease", buscommand.HandlerFunc(leaseHandler.HandleRenewLease))
	bus.Register("checkout_lease", buscommand.HandlerFunc(leaseHandler.HandleCheckoutLease))
	bus.Register("checkout_with_bills", buscommand.HandlerFunc(leaseHandler.HandleCheckoutWithBills))
	bus.Register("activate_lease", buscommand.HandlerFunc(leaseHandler.HandleActivateLease))
	bus.Register("create_bill", buscommand.HandlerFunc(billHandler.HandleCreateBill))
	bus.Register("update_bill", buscommand.HandlerFunc(billHandler.HandleUpdateBill))
	bus.Register("delete_bill", buscommand.HandlerFunc(billHandler.HandleDeleteBill))
	bus.Register("confirm_bill_arrival", buscommand.HandlerFunc(billHandler.HandleConfirmBillArrival))
	bus.Register("mark_deposit_returning", buscommand.HandlerFunc(depositHandler.HandleMarkReturning))
	bus.Register("mark_deposit_returned", buscommand.HandlerFunc(depositHandler.HandleMarkReturned))
	bus.Register("create_location", buscommand.HandlerFunc(locationHandler.HandleCreateLocation))
	bus.Register("update_location", buscommand.HandlerFunc(locationHandler.HandleUpdateLocation))
	bus.Register("delete_location", buscommand.HandlerFunc(locationHandler.HandleDeleteLocation))
	bus.Register("create_room", buscommand.HandlerFunc(roomHandler.HandleCreateRoom))
	bus.Register("update_room", buscommand.HandlerFunc(roomHandler.HandleUpdateRoom))
	bus.Register("delete_room", buscommand.HandlerFunc(roomHandler.HandleDeleteRoom))
	bus.Register("print_bill", buscommand.HandlerFunc(printHandler.HandlePrintBill))
	bus.Register("print_lease", buscommand.HandlerFunc(printHandler.HandlePrintLease))
	bus.Register("print_invoice", buscommand.HandlerFunc(printHandler.HandlePrintInvoice))
}

func registerQueryHandlers(
	queryBus *busquery.Bus,
	landlordQueryHandler *landlord.QueryHandler,
	leaseQueryHandler *lease.QueryHandler,
	billQueryHandler *bill.QueryHandler,
	depositQueryHandler *deposit.QueryHandler,
	locationQueryHandler *location.QueryHandler,
	roomQueryHandler *room.QueryHandler,
	printQueryHandler *print.QueryHandler,
) {
	queryBus.Register("get_landlord", busquery.HandlerFunc(landlordQueryHandler.HandleGetLandlord))
	queryBus.Register("list_landlords", busquery.HandlerFunc(landlordQueryHandler.HandleListLandlords))
	queryBus.Register("get_lease", busquery.HandlerFunc(leaseQueryHandler.HandleGetLease))
	queryBus.Register("list_leases", busquery.HandlerFunc(leaseQueryHandler.HandleListLeases))
	queryBus.Register("get_bill", busquery.HandlerFunc(billQueryHandler.HandleGetBill))
	queryBus.Register("list_bills", busquery.HandlerFunc(billQueryHandler.HandleListBills))
	queryBus.Register("income_report", busquery.HandlerFunc(billQueryHandler.HandleIncomeReport))
	queryBus.Register("get_next_bill_period", busquery.HandlerFunc(billQueryHandler.HandleGetNextBillPeriod))
	queryBus.Register("get_deposit", busquery.HandlerFunc(depositQueryHandler.HandleGetDeposit))
	queryBus.Register("list_deposits", busquery.HandlerFunc(depositQueryHandler.HandleListDeposits))
	queryBus.Register("get_location", busquery.HandlerFunc(locationQueryHandler.HandleGetLocation))
	queryBus.Register("list_locations", busquery.HandlerFunc(locationQueryHandler.HandleListLocations))
	queryBus.Register("get_room", busquery.HandlerFunc(roomQueryHandler.HandleGetRoom))
	queryBus.Register("list_rooms", busquery.HandlerFunc(roomQueryHandler.HandleListRooms))
	queryBus.Register("list_rooms_by_location", busquery.HandlerFunc(roomQueryHandler.HandleListRoomsByLocation))
	queryBus.Register("get_print_job", busquery.HandlerFunc(printQueryHandler.HandleGetPrintJob))
	queryBus.Register("list_print_jobs", busquery.HandlerFunc(printQueryHandler.HandleListPrintJobs))
	queryBus.Register("get_print_content", busquery.HandlerFunc(printQueryHandler.HandleGetPrintContent))
}

func startServer(
	logger *zap.Logger,
	locationHandler *facade.CQRSLocationHandler,
	roomHandler *facade.CQRSRoomHandler,
	landlordHandler *facade.CQRSLandlordHandler,
	leaseHandler *facade.CQRSLeaseHandler,
	billHandler *facade.CQRSBillHandler,
	depositHandler *facade.CQRSDepositHandler,
	printHandler *facade.CQRSPrintHandler,
	authHandler *facade.AuthHandler,
	incomeHandler *facade.IncomeHandler,
	operationLogHandler *facade.OperationLogHandler,
	authService *service.AuthService,
	printService *service.PrintService,
) {
	mux := http.NewServeMux()

	// CORS 中间件
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 允许所有源（生产环境应该限制特定域）
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// 处理预检请求
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// 健康检查路由
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// 注册业务路由
	locationHandler.RegisterRoutes(mux)
	roomHandler.RegisterRoutes(mux)
	landlordHandler.RegisterRoutes(mux)
	leaseHandler.RegisterRoutes(mux)
	billHandler.RegisterRoutes(mux)
	depositHandler.RegisterRoutes(mux)
	printHandler.RegisterRoutes(mux)
	authHandler.RegisterRoutes(mux)
	incomeHandler.RegisterRoutes(mux)
	operationLogHandler.RegisterRoutes(mux)

	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", corsHandler(mux)); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}

func startLeaseExpirationScheduler(scheduler *lease.LeaseExpirationScheduler, logger *zap.Logger) {
	if err := scheduler.Start(); err != nil {
		logger.Error("Failed to start lease expiration scheduler", zap.Error(err))
	}
}
