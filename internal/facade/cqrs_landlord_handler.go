package facade

import (
	"encoding/json"
	"github.com/zouhang1992/ddd_domain/internal/application/landlord"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/logging"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/middleware"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// CQRSLandlordHandler 基于 CQRS 的房东 HTTP 处理器
type CQRSLandlordHandler struct {
	commandBus      *buscommand.Bus
	queryBus        *busquery.Bus
	authMiddleware  *middleware.AuthMiddleware
}

// NewCQRSLandlordHandler 创建基于 CQRS 的房东处理器
func NewCQRSLandlordHandler(
	commandBus *buscommand.Bus,
	queryBus *busquery.Bus,
	authMiddleware *middleware.AuthMiddleware,
) *CQRSLandlordHandler {
	return &CQRSLandlordHandler{
		commandBus:     commandBus,
		queryBus:       queryBus,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSLandlordHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /landlords", h.authMiddleware.RequireAuth(h.Create))
	mux.HandleFunc("GET /landlords", h.authMiddleware.RequireAuth(h.List))
	mux.HandleFunc("GET /landlords/{id}", h.authMiddleware.RequireAuth(h.Get))
	mux.HandleFunc("PUT /landlords/{id}", h.authMiddleware.RequireAuth(h.Update))
	mux.HandleFunc("DELETE /landlords/{id}", h.authMiddleware.RequireAuth(h.Delete))
}

// Create 创建房东
func (h *CQRSLandlordHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logging.Ctx(ctx).Info("Creating landlord")

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Note  string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Ctx(ctx).Error("Failed to decode request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := landlord.CreateLandlordCommand{
		Name:  req.Name,
		Phone: req.Phone,
		Note:  req.Note,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		logging.Ctx(ctx).Error("Failed to create landlord", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Ctx(ctx).Info("Landlord created successfully")
	landlord := result.(any)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(landlord)
}

// List 列出房东
func (h *CQRSLandlordHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logging.Ctx(ctx).Info("Listing landlords")

	// 解析查询参数
	q := landlord.ListLandlordsQuery{
		Name:  r.URL.Query().Get("name"),
		Phone: r.URL.Query().Get("phone"),
	}

	// 解析分页参数
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			q.Offset = offset
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			q.Limit = limit
		}
	}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		logging.Ctx(ctx).Error("Failed to list landlords", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Ctx(ctx).Info("Landlords listed successfully")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// Get 获取房东
func (h *CQRSLandlordHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	logging.Ctx(ctx).Info("Getting landlord", zap.String("id", id))

	q := landlord.GetLandlordQuery{ID: id}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		if err.Error() == "landlord not found" {
			logging.Ctx(ctx).Warn("Landlord not found", zap.String("id", id))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logging.Ctx(ctx).Error("Failed to get landlord", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Ctx(ctx).Info("Landlord retrieved successfully", zap.String("id", id))
	queryResult := result.(*landlord.LandlordQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Landlord)
}

// Update 更新房东
func (h *CQRSLandlordHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	logging.Ctx(ctx).Info("Updating landlord", zap.String("id", id))

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Note  string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Ctx(ctx).Error("Failed to decode request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := landlord.UpdateLandlordCommand{
		ID:    id,
		Name:  req.Name,
		Phone: req.Phone,
		Note:  req.Note,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "landlord not found" {
			logging.Ctx(ctx).Warn("Landlord not found for update", zap.String("id", id))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logging.Ctx(ctx).Error("Failed to update landlord", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Ctx(ctx).Info("Landlord updated successfully", zap.String("id", id))
	landlord := result.(any)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(landlord)
}

// Delete 删除房东
func (h *CQRSLandlordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	logging.Ctx(ctx).Info("Deleting landlord", zap.String("id", id))

	cmd := landlord.DeleteLandlordCommand{ID: id}

	if _, err := h.commandBus.Dispatch(cmd); err != nil {
		if err.Error() == "cannot delete landlord with associated leases" {
			logging.Ctx(ctx).Warn("Cannot delete landlord with associated leases", zap.String("id", id))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logging.Ctx(ctx).Error("Failed to delete landlord", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Ctx(ctx).Info("Landlord deleted successfully", zap.String("id", id))
	w.WriteHeader(http.StatusNoContent)
}
