package facade

import (
	"encoding/json"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"net/http"
)

// CQRSLandlordHandler 基于 CQRS 的房东 HTTP 处理器
type CQRSLandlordHandler struct {
	commandBus *buscommand.Bus
	queryBus   *busquery.Bus
}

// NewCQRSLandlordHandler 创建基于 CQRS 的房东处理器
func NewCQRSLandlordHandler(commandBus *buscommand.Bus, queryBus *busquery.Bus) *CQRSLandlordHandler {
	return &CQRSLandlordHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSLandlordHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /landlords", h.Create)
	mux.HandleFunc("GET /landlords", h.List)
	mux.HandleFunc("GET /landlords/{id}", h.Get)
	mux.HandleFunc("PUT /landlords/{id}", h.Update)
	mux.HandleFunc("DELETE /landlords/{id}", h.Delete)
}

// Create 创建房东
func (h *CQRSLandlordHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Note  string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.CreateLandlordCommand{
		Name:  req.Name,
		Phone: req.Phone,
		Note:  req.Note,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	landlord := result.(any)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(landlord)
}

// List 列出房东
func (h *CQRSLandlordHandler) List(w http.ResponseWriter, r *http.Request) {
	q := query.ListLandlordsQuery{}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult := result.(*query.LandlordsQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Items)
}

// Get 获取房东
func (h *CQRSLandlordHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	q := query.GetLandlordQuery{ID: id}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		if err.Error() == "landlord not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult := result.(*query.LandlordQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Landlord)
}

// Update 更新房东
func (h *CQRSLandlordHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Note  string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.UpdateLandlordCommand{
		ID:    id,
		Name:  req.Name,
		Phone: req.Phone,
		Note:  req.Note,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "landlord not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	landlord := result.(any)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(landlord)
}

// Delete 删除房东
func (h *CQRSLandlordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cmd := command.DeleteLandlordCommand{ID: id}

	if _, err := h.commandBus.Dispatch(cmd); err != nil {
		if err.Error() == "cannot delete landlord with associated leases" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
