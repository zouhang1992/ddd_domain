package facade

import (
	"encoding/json"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"net/http"
	"strconv"
	"time"
)

// CQRSLeaseHandler 基于 CQRS 的租约 HTTP 处理器
type CQRSLeaseHandler struct {
	commandBus *buscommand.Bus
	queryBus   *busquery.Bus
}

// NewCQRSLeaseHandler 创建基于 CQRS 的租约处理器
func NewCQRSLeaseHandler(commandBus *buscommand.Bus, queryBus *busquery.Bus) *CQRSLeaseHandler {
	return &CQRSLeaseHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSLeaseHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /leases", h.Create)
	mux.HandleFunc("GET /leases", h.List)
	mux.HandleFunc("GET /leases/{id}", h.Get)
	mux.HandleFunc("PUT /leases/{id}", h.Update)
	mux.HandleFunc("DELETE /leases/{id}", h.Delete)
	mux.HandleFunc("POST /leases/{id}/renew", h.Renew)
	mux.HandleFunc("POST /leases/{id}/checkout", h.Checkout)
	mux.HandleFunc("PUT /leases/{id}/activate", h.Activate)
}

// Create 创建租约
func (h *CQRSLeaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID        string `json:"room_id"`
		LandlordID    string `json:"landlord_id"`
		TenantName    string `json:"tenant_name"`
		TenantPhone   string `json:"tenant_phone"`
		StartDate     string `json:"start_date"`
		EndDate       string `json:"end_date"`
		RentAmount    int64  `json:"rent_amount"`
		Note          string `json:"note"`
		DepositAmount int64  `json:"deposit_amount"`
		DepositNote   string `json:"deposit_note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		http.Error(w, "invalid end_date format", http.StatusBadRequest)
		return
	}

	cmd := command.CreateLeaseCommand{
		RoomID:        req.RoomID,
		LandlordID:    req.LandlordID,
		TenantName:    req.TenantName,
		TenantPhone:   req.TenantPhone,
		StartDate:     startDate,
		EndDate:       endDate,
		RentAmount:    req.RentAmount,
		Note:          req.Note,
		DepositAmount: req.DepositAmount,
		DepositNote:   req.DepositNote,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "room already has an active lease" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lease := result.(any)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(lease)
}

// List 列出租约
func (h *CQRSLeaseHandler) List(w http.ResponseWriter, r *http.Request) {
	q := query.ListLeasesQuery{
		TenantName:  r.URL.Query().Get("tenant_name"),
		TenantPhone: r.URL.Query().Get("tenant_phone"),
		Status:      r.URL.Query().Get("status"),
		RoomID:      r.URL.Query().Get("room_id"),
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// Get 获取租约
func (h *CQRSLeaseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	q := query.GetLeaseQuery{ID: id}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		if err.Error() == "lease not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult := result.(*query.LeaseQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Lease)
}

// Update 更新租约
func (h *CQRSLeaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		TenantName  string `json:"tenant_name"`
		TenantPhone string `json:"tenant_phone"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		RentAmount  int64  `json:"rent_amount"`
		Note        string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		http.Error(w, "invalid end_date format", http.StatusBadRequest)
		return
	}

	cmd := command.UpdateLeaseCommand{
		ID:          id,
		TenantName:  req.TenantName,
		TenantPhone: req.TenantPhone,
		StartDate:   startDate,
		EndDate:     endDate,
		RentAmount:  req.RentAmount,
		Note:        req.Note,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "lease not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lease := result.(any)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(lease)
}

// Delete 删除租约
func (h *CQRSLeaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cmd := command.DeleteLeaseCommand{ID: id}

	if _, err := h.commandBus.Dispatch(cmd); err != nil {
		if err.Error() == "cannot delete lease with bills or deposit, or active lease" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Renew 续租
func (h *CQRSLeaseHandler) Renew(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		NewStartDate  string `json:"new_start_date"`
		NewEndDate    string `json:"new_end_date"`
		NewRentAmount int64  `json:"new_rent_amount"`
		Note          string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newStartDate, err := time.Parse("2006-01-02", req.NewStartDate)
	if err != nil {
		http.Error(w, "invalid new_start_date format", http.StatusBadRequest)
		return
	}

	newEndDate, err := time.Parse("2006-01-02", req.NewEndDate)
	if err != nil {
		http.Error(w, "invalid new_end_date format", http.StatusBadRequest)
		return
	}

	cmd := command.RenewLeaseCommand{
		ID:            id,
		NewStartDate:  newStartDate,
		NewEndDate:    newEndDate,
		NewRentAmount: req.NewRentAmount,
		Note:          req.Note,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "lease not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newLease := result.(any)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(newLease)
}

// Checkout 退租
func (h *CQRSLeaseHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cmd := command.CheckoutLeaseCommand{ID: id}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "lease not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lease := result.(any)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(lease)
}

// Activate 租约生效
func (h *CQRSLeaseHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cmd := command.ActivateLeaseCommand{ID: id}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "lease not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "invalid lease state" || err.Error() == "lease start date has not arrived" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lease := result.(any)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(lease)
}
