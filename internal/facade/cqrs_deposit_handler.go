package facade

import (
	"encoding/json"
	"github.com/zouhang1992/ddd_domain/internal/application/deposit"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/logging"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/middleware"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// CQRSDepositHandler 基于 CQRS 的押金 HTTP 处理器
type CQRSDepositHandler struct {
	commandBus      *buscommand.Bus
	queryBus        *busquery.Bus
	authMiddleware  *middleware.AuthMiddleware
}

// NewCQRSDepositHandler 创建基于 CQRS 的押金处理器
func NewCQRSDepositHandler(
	commandBus *buscommand.Bus,
	queryBus *busquery.Bus,
	authMiddleware *middleware.AuthMiddleware,
) *CQRSDepositHandler {
	return &CQRSDepositHandler{
		commandBus:     commandBus,
		queryBus:       queryBus,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSDepositHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /deposits", h.authMiddleware.RequireAuth(h.List))
	mux.HandleFunc("GET /deposits/{id}", h.authMiddleware.RequireAuth(h.Get))
	mux.HandleFunc("POST /deposits/{id}/mark-returning", h.authMiddleware.RequireAuth(h.MarkReturning))
	mux.HandleFunc("POST /deposits/{id}/mark-returned", h.authMiddleware.RequireAuth(h.MarkReturned))
}

// List 列出押金
func (h *CQRSDepositHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logging.Ctx(ctx).Info("Listing deposits")

	query := r.URL.Query()

	offsetStr := query.Get("offset")
	limitStr := query.Get("limit")
	pageStr := query.Get("page")

	offset := 0
	limit := 10

	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
		if page < 1 {
			page = 1
		}
		offset = (page - 1) * limit
	} else if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err == nil && o >= 0 {
			offset = o
		}
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	q := deposit.ListDepositsQuery{
		LeaseID: query.Get("lease_id"),
		Status:  query.Get("status"),
		Offset:  offset,
		Limit:   limit,
	}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		logging.Ctx(ctx).Error("Failed to list deposits", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Ctx(ctx).Info("Deposits listed successfully")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// Get 获取押金
func (h *CQRSDepositHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	logging.Ctx(ctx).Info("Getting deposit", zap.String("id", id))

	q := deposit.GetDepositQuery{ID: id}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		if err.Error() == "not found" {
			logging.Ctx(ctx).Warn("Deposit not found", zap.String("id", id))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logging.Ctx(ctx).Error("Failed to get deposit", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Ctx(ctx).Info("Deposit retrieved successfully", zap.String("id", id))
	queryResult := result.(*deposit.DepositQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Deposit)
}

// MarkReturning 标记押金为待退还
func (h *CQRSDepositHandler) MarkReturning(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	logging.Ctx(ctx).Info("Marking deposit as returning", zap.String("id", id))

	cmd := deposit.MarkReturningCommand{ID: id}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "not found" {
			logging.Ctx(ctx).Warn("Deposit not found for marking returning", zap.String("id", id))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logging.Ctx(ctx).Error("Failed to mark deposit as returning", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Ctx(ctx).Info("Deposit marked as returning successfully", zap.String("id", id))
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// MarkReturned 标记押金为已退还
func (h *CQRSDepositHandler) MarkReturned(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	logging.Ctx(ctx).Info("Marking deposit as returned", zap.String("id", id))

	cmd := deposit.MarkReturnedCommand{ID: id}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "not found" {
			logging.Ctx(ctx).Warn("Deposit not found for marking returned", zap.String("id", id))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logging.Ctx(ctx).Error("Failed to mark deposit as returned", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Ctx(ctx).Info("Deposit marked as returned successfully", zap.String("id", id))
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
