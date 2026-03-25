package facade

import (
	"encoding/json"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"net/http"
)

// CQRSLocationHandler 基于 CQRS 的位置 HTTP 处理器
type CQRSLocationHandler struct {
	commandBus *buscommand.Bus
	queryBus   *busquery.Bus
}

// NewCQRSLocationHandler 创建基于 CQRS 的位置处理器
func NewCQRSLocationHandler(commandBus *buscommand.Bus, queryBus *busquery.Bus) *CQRSLocationHandler {
	return &CQRSLocationHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSLocationHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /locations", h.Create)
	mux.HandleFunc("GET /locations", h.List)
	mux.HandleFunc("GET /locations/{id}", h.Get)
	mux.HandleFunc("PUT /locations/{id}", h.Update)
	mux.HandleFunc("DELETE /locations/{id}", h.Delete)
}

// Create 创建位置
func (h *CQRSLocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ShortName string `json:"shortName"`
		Detail    string `json:"detail"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.CreateLocationCommand{
		ShortName: req.ShortName,
		Detail:    req.Detail,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	location := result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(location)
}

// List 列出位置
func (h *CQRSLocationHandler) List(w http.ResponseWriter, r *http.Request) {
	q := query.ListLocationsQuery{}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult := result.(*query.LocationsQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Items)
}

// Get 获取位置
func (h *CQRSLocationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	q := query.GetLocationQuery{ID: id}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult := result.(*query.LocationQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Location)
}

// Update 更新位置
func (h *CQRSLocationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		ShortName string `json:"shortName"`
		Detail    string `json:"detail"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.UpdateLocationCommand{
		ID:        id,
		ShortName: req.ShortName,
		Detail:    req.Detail,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	location := result
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(location)
}

// Delete 删除位置
func (h *CQRSLocationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cmd := command.DeleteLocationCommand{ID: id}

	if _, err := h.commandBus.Dispatch(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
