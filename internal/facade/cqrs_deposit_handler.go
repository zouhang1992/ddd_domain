package facade

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zouhang1992/ddd_domain/internal/application/deposit"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
)

// CQRSDepositHandler 基于 CQRS 的押金 HTTP 处理器
type CQRSDepositHandler struct {
	commandBus *buscommand.Bus
	queryBus   *busquery.Bus
}

// NewCQRSDepositHandler 创建基于 CQRS 的押金处理器
func NewCQRSDepositHandler(commandBus *buscommand.Bus, queryBus *busquery.Bus) *CQRSDepositHandler {
	return &CQRSDepositHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSDepositHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /deposits", h.List)
	mux.HandleFunc("GET /deposits/{id}", h.Get)
	mux.HandleFunc("POST /deposits/{id}/mark-returning", h.MarkReturning)
	mux.HandleFunc("POST /deposits/{id}/mark-returned", h.MarkReturned)
}

// List 列出押金
func (h *CQRSDepositHandler) List(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// Get 获取押金
func (h *CQRSDepositHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	q := deposit.GetDepositQuery{ID: id}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult := result.(*deposit.DepositQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Deposit)
}

// MarkReturning 标记押金为待退还
func (h *CQRSDepositHandler) MarkReturning(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cmd := deposit.MarkReturningCommand{ID: id}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// MarkReturned 标记押金为已退还
func (h *CQRSDepositHandler) MarkReturned(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cmd := deposit.MarkReturnedCommand{ID: id}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
