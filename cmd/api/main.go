package main

import (
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/zouhang1992/ddd_domain/internal/application/command/handler"
	eventhandler "github.com/zouhang1992/ddd_domain/internal/application/event/handler"
	queryhandler "github.com/zouhang1992/ddd_domain/internal/application/query/handler"
	"github.com/zouhang1992/ddd_domain/internal/application/service"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
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
		handler.Module,
		eventhandler.Module,
		queryhandler.Module,
		facade.Module,
		// 应用服务
		fx.Options(
			fx.Provide(func() *service.AuthService {
				return service.NewAuthService("", 7*24*time.Hour)
			}),
			fx.Provide(func(billRepo repository.BillRepository, leaseRepo repository.LeaseRepository) *service.PrintService {
				return service.NewPrintService(billRepo, leaseRepo)
			}),
		),
		// 注册处理器到总线
		fx.Invoke(registerCommandHandlers),
		fx.Invoke(registerQueryHandlers),
		fx.Invoke(registerEventHandlers),
		// 设置服务器
		fx.Invoke(startServer),
	).Run()
}

func registerEventHandlers(eventBus *busevent.Bus, logHandler *eventhandler.Handler) {
	logHandler.SubscribeToAllEvents(eventBus)
}

func registerCommandHandlers(
	bus *buscommand.Bus,
	landlordHandler *handler.LandlordCommandHandler,
	leaseHandler *handler.LeaseCommandHandler,
	billHandler *handler.BillCommandHandler,
	locationHandler *handler.LocationCommandHandler,
	roomHandler *handler.RoomCommandHandler,
	printHandler *handler.PrintCommandHandler,
) {
	bus.Register("create_landlord", buscommand.HandlerFunc(landlordHandler.HandleCreateLandlord))
	bus.Register("update_landlord", buscommand.HandlerFunc(landlordHandler.HandleUpdateLandlord))
	bus.Register("delete_landlord", buscommand.HandlerFunc(landlordHandler.HandleDeleteLandlord))
	bus.Register("create_lease", buscommand.HandlerFunc(leaseHandler.HandleCreateLease))
	bus.Register("update_lease", buscommand.HandlerFunc(leaseHandler.HandleUpdateLease))
	bus.Register("delete_lease", buscommand.HandlerFunc(leaseHandler.HandleDeleteLease))
	bus.Register("renew_lease", buscommand.HandlerFunc(leaseHandler.HandleRenewLease))
	bus.Register("checkout_lease", buscommand.HandlerFunc(leaseHandler.HandleCheckoutLease))
	bus.Register("activate_lease", buscommand.HandlerFunc(leaseHandler.HandleActivateLease))
	bus.Register("create_bill", buscommand.HandlerFunc(billHandler.HandleCreateBill))
	bus.Register("update_bill", buscommand.HandlerFunc(billHandler.HandleUpdateBill))
	bus.Register("delete_bill", buscommand.HandlerFunc(billHandler.HandleDeleteBill))
	bus.Register("confirm_bill_arrival", buscommand.HandlerFunc(billHandler.HandleConfirmBillArrival))
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
	landlordQueryHandler *queryhandler.LandlordQueryHandler,
	leaseQueryHandler *queryhandler.LeaseQueryHandler,
	billQueryHandler *queryhandler.BillQueryHandler,
	locationQueryHandler *queryhandler.LocationQueryHandler,
	roomQueryHandler *queryhandler.RoomQueryHandler,
	printQueryHandler *queryhandler.PrintQueryHandler,
	operationLogQueryHandler *queryhandler.OperationLogQueryHandler,
) {
	queryBus.Register("get_landlord", busquery.HandlerFunc(landlordQueryHandler.HandleGetLandlord))
	queryBus.Register("list_landlords", busquery.HandlerFunc(landlordQueryHandler.HandleListLandlords))
	queryBus.Register("get_lease", busquery.HandlerFunc(leaseQueryHandler.HandleGetLease))
	queryBus.Register("list_leases", busquery.HandlerFunc(leaseQueryHandler.HandleListLeases))
	queryBus.Register("get_bill", busquery.HandlerFunc(billQueryHandler.HandleGetBill))
	queryBus.Register("list_bills", busquery.HandlerFunc(billQueryHandler.HandleListBills))
	queryBus.Register("income_report", busquery.HandlerFunc(billQueryHandler.HandleIncomeReport))
	queryBus.Register("get_location", busquery.HandlerFunc(locationQueryHandler.HandleGetLocation))
	queryBus.Register("list_locations", busquery.HandlerFunc(locationQueryHandler.HandleListLocations))
	queryBus.Register("get_room", busquery.HandlerFunc(roomQueryHandler.HandleGetRoom))
	queryBus.Register("list_rooms", busquery.HandlerFunc(roomQueryHandler.HandleListRooms))
	queryBus.Register("list_rooms_by_location", busquery.HandlerFunc(roomQueryHandler.HandleListRoomsByLocation))
	queryBus.Register("get_print_job", busquery.HandlerFunc(printQueryHandler.HandleGetPrintJob))
	queryBus.Register("list_print_jobs", busquery.HandlerFunc(printQueryHandler.HandleListPrintJobs))
	queryBus.Register("get_print_content", busquery.HandlerFunc(printQueryHandler.HandleGetPrintContent))
	// 操作日志查询
	queryBus.Register("list_operation_logs", busquery.HandlerFunc(operationLogQueryHandler.HandleListOperationLogs))
	queryBus.Register("get_operation_log", busquery.HandlerFunc(operationLogQueryHandler.HandleGetOperationLog))
}

func startServer(
	logger *zap.Logger,
	locationHandler *facade.CQRSLocationHandler,
	roomHandler *facade.CQRSRoomHandler,
	landlordHandler *facade.CQRSLandlordHandler,
	leaseHandler *facade.CQRSLeaseHandler,
	billHandler *facade.CQRSBillHandler,
	printHandler *facade.CQRSPrintHandler,
	authHandler *facade.AuthHandler,
	incomeHandler *facade.IncomeHandler,
	operationLogHandler *facade.CQRSOperationLogHandler,
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
	printHandler.RegisterRoutes(mux)
	authHandler.RegisterRoutes(mux)
	incomeHandler.RegisterRoutes(mux)
	operationLogHandler.RegisterRoutes(mux)

	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", corsHandler(mux)); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
