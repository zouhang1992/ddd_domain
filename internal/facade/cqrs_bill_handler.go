package facade

import (
	"encoding/json"
	"github.com/zouhang1992/ddd_domain/internal/application/bill"
	"github.com/zouhang1992/ddd_domain/internal/application/common/service"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/middleware"
	billmodel "github.com/zouhang1992/ddd_domain/internal/domain/bill/model"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"net/http"
	"time"
)

// CQRSBillHandler 基于 CQRS 的账单 HTTP 处理器
type CQRSBillHandler struct {
	commandBus      *buscommand.Bus
	queryBus        *busquery.Bus
	printService    *service.PrintService
	authMiddleware  *middleware.AuthMiddleware
}

// NewCQRSBillHandler 创建基于 CQRS 的账单处理器
func NewCQRSBillHandler(
	commandBus *buscommand.Bus, 
	queryBus *busquery.Bus, 
	printService *service.PrintService,
	authMiddleware *middleware.AuthMiddleware,
) *CQRSBillHandler {
	return &CQRSBillHandler{
		commandBus:     commandBus,
		queryBus:       queryBus,
		printService:   printService,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSBillHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /bills", h.authMiddleware.RequireAuth(h.Create))
	mux.HandleFunc("GET /bills", h.authMiddleware.RequireAuth(h.List))
	mux.HandleFunc("GET /bills/{id}", h.authMiddleware.RequireAuth(h.Get))
	mux.HandleFunc("PUT /bills/{id}", h.authMiddleware.RequireAuth(h.Update))
	mux.HandleFunc("DELETE /bills/{id}", h.authMiddleware.RequireAuth(h.Delete))
	mux.HandleFunc("GET /bills/{id}/receipt", h.authMiddleware.RequireAuth(h.PrintReceipt))
	mux.HandleFunc("POST /bills/{id}/confirm-arrival", h.authMiddleware.RequireAuth(h.ConfirmArrival))
	mux.HandleFunc("GET /leases/{leaseId}/next-bill-period", h.authMiddleware.RequireAuth(h.GetNextBillPeriod))
}

// Create 创建账单
func (h *CQRSBillHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LeaseID             string     `json:"lease_id"`
		Type                string     `json:"type"`
		Amount              int64      `json:"amount"`
		RentAmount          int64      `json:"rent_amount"`
		WaterAmount         int64      `json:"water_amount"`
		ElectricAmount      int64      `json:"electric_amount"`
		OtherAmount         int64      `json:"other_amount"`
		RefundDepositAmount int64      `json:"refund_deposit_amount"`
		BillStart           string     `json:"bill_start"`
		BillEnd             string     `json:"bill_end"`
		DueDate             string     `json:"due_date"`
		Note                string     `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse bill start date
	billStart, err := time.Parse("2006-01-02", req.BillStart)
	if err != nil {
		billStart = time.Now()
	}

	// Parse bill end date
	billEnd, err := time.Parse("2006-01-02", req.BillEnd)
	if err != nil {
		billEnd = time.Now().AddDate(0, 1, 0)
	}

	// Parse due date
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		dueDate = time.Now().AddDate(0, 1, 0) // Default to 1 month from now
	}

	cmd := bill.CreateBillCommand{
		LeaseID:             req.LeaseID,
		Type:                billmodel.BillType(req.Type),
		Amount:              req.Amount,
		RentAmount:          req.RentAmount,
		WaterAmount:         req.WaterAmount,
		ElectricAmount:      req.ElectricAmount,
		OtherAmount:         req.OtherAmount,
		RefundDepositAmount: req.RefundDepositAmount,
		BillStart:           billStart,
		BillEnd:             billEnd,
		DueDate:             dueDate,
		Note:                req.Note,
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
	q := bill.ListBillsQuery{
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
	q := bill.GetBillQuery{ID: id}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		if err.Error() == "bill not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult := result.(*bill.BillQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Bill)
}

// Update 更新账单
func (h *CQRSBillHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Amount              int64      `json:"amount"`
		RentAmount          int64      `json:"rent_amount"`
		WaterAmount         int64      `json:"water_amount"`
		ElectricAmount      int64      `json:"electric_amount"`
		OtherAmount         int64      `json:"other_amount"`
		RefundDepositAmount int64      `json:"refund_deposit_amount"`
		BillStart           string     `json:"bill_start"`
		BillEnd             string     `json:"bill_end"`
		DueDate             string     `json:"due_date"`
		Note                string     `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse bill start date
	billStart, err := time.Parse("2006-01-02", req.BillStart)
	if err != nil {
		billStart = time.Now()
	}

	// Parse bill end date
	billEnd, err := time.Parse("2006-01-02", req.BillEnd)
	if err != nil {
		billEnd = time.Now().AddDate(0, 1, 0)
	}

	// Parse due date
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		dueDate = time.Now().AddDate(0, 1, 0) // Default to 1 month from now
	}

	cmd := bill.UpdateBillCommand{
		ID:                  id,
		Amount:              req.Amount,
		RentAmount:          req.RentAmount,
		WaterAmount:         req.WaterAmount,
		ElectricAmount:      req.ElectricAmount,
		OtherAmount:         req.OtherAmount,
		RefundDepositAmount: req.RefundDepositAmount,
		BillStart:           billStart,
		BillEnd:             billEnd,
		DueDate:             dueDate,
		Note:                req.Note,
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
	cmd := bill.DeleteBillCommand{ID: id}

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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"receipt.html\"")
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

	cmd := bill.ConfirmBillArrivalCommand{
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

// GetNextBillPeriod 获取租约的下一个账单周期
func (h *CQRSBillHandler) GetNextBillPeriod(w http.ResponseWriter, r *http.Request) {
	leaseId := r.PathValue("leaseId")
	q := bill.GetNextBillPeriodQuery{LeaseID: leaseId}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		if err.Error() == "lease not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
