package facade

import (
	"encoding/json"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	"github.com/zouhang1992/ddd_domain/internal/application/service"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"net/http"
	"time"
)

// CQRSBillHandler 基于 CQRS 的账单 HTTP 处理器
type CQRSBillHandler struct {
	commandBus   *buscommand.Bus
	queryBus     *busquery.Bus
	printService *service.PrintService
}

// NewCQRSBillHandler 创建基于 CQRS 的账单处理器
func NewCQRSBillHandler(commandBus *buscommand.Bus, queryBus *busquery.Bus, printService *service.PrintService) *CQRSBillHandler {
	return &CQRSBillHandler{
		commandBus:   commandBus,
		queryBus:     queryBus,
		printService: printService,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSBillHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /bills", h.Create)
	mux.HandleFunc("GET /bills", h.List)
	mux.HandleFunc("GET /bills/{id}", h.Get)
	mux.HandleFunc("PUT /bills/{id}", h.Update)
	mux.HandleFunc("DELETE /bills/{id}", h.Delete)
	mux.HandleFunc("GET /bills/{id}/receipt", h.PrintReceipt)
	mux.HandleFunc("POST /bills/{id}/confirm-arrival", h.ConfirmArrival)
}

// Create 创建账单
func (h *CQRSBillHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LeaseID        string     `json:"lease_id"`
		Type           string     `json:"type"`
		Amount         int64      `json:"amount"`
		RentAmount     int64      `json:"rent_amount"`
		WaterAmount    int64      `json:"water_amount"`
		ElectricAmount int64      `json:"electric_amount"`
		OtherAmount    int64      `json:"other_amount"`
		PaidAt         *time.Time `json:"paid_at"`
		Note           string     `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.CreateBillCommand{
		LeaseID:        req.LeaseID,
		Type:           model.BillType(req.Type),
		Amount:         req.Amount,
		RentAmount:     req.RentAmount,
		WaterAmount:    req.WaterAmount,
		ElectricAmount: req.ElectricAmount,
		OtherAmount:    req.OtherAmount,
		PaidAt:         req.PaidAt,
		Note:           req.Note,
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

	bill := result.(any)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(bill)
}

// List 列出账单
func (h *CQRSBillHandler) List(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	q := query.ListBillsQuery{
		LeaseID: r.URL.Query().Get("lease_id"),
		RoomID:  r.URL.Query().Get("room_id"),
		Month:   r.URL.Query().Get("month"),
	}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// Get 获取账单
func (h *CQRSBillHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	q := query.GetBillQuery{ID: id}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		if err.Error() == "bill not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult := result.(*query.BillQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Bill)
}

// Update 更新账单
func (h *CQRSBillHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Amount         int64      `json:"amount"`
		RentAmount     int64      `json:"rent_amount"`
		WaterAmount    int64      `json:"water_amount"`
		ElectricAmount int64      `json:"electric_amount"`
		OtherAmount    int64      `json:"other_amount"`
		PaidAt         *time.Time `json:"paid_at"`
		Note           string     `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.UpdateBillCommand{
		ID:             id,
		Amount:         req.Amount,
		RentAmount:     req.RentAmount,
		WaterAmount:    req.WaterAmount,
		ElectricAmount: req.ElectricAmount,
		OtherAmount:    req.OtherAmount,
		PaidAt:         req.PaidAt,
		Note:           req.Note,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "bill not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bill := result.(any)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(bill)
}

// Delete 删除账单
func (h *CQRSBillHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cmd := command.DeleteBillCommand{ID: id}

	if _, err := h.commandBus.Dispatch(cmd); err != nil {
		if err.Error() == "cannot delete this bill" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PrintReceipt 打印收据
func (h *CQRSBillHandler) PrintReceipt(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	content, err := h.printService.PrintReceipt(id)
	if err != nil {
		if err.Error() == "bill not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rtf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"receipt.rtf\"")
	w.Write(content)
}

// ConfirmArrival 确认账单到账
func (h *CQRSBillHandler) ConfirmArrival(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		PaidAt *time.Time `json:"paid_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paidAt := req.PaidAt
	if paidAt == nil {
		now := time.Now()
		paidAt = &now
	}

	cmd := command.ConfirmBillArrivalCommand{
		ID:     id,
		PaidAt: *paidAt,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		if err.Error() == "bill not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bill := result.(any)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(bill)
}
