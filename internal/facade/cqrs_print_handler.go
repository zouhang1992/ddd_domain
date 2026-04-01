package facade

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/zouhang1992/ddd_domain/internal/application/print"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/middleware"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"net/http"
)

// CQRSPrintHandler 基于 CQRS 的打印 HTTP 处理器
type CQRSPrintHandler struct {
	commandBus      *buscommand.Bus
	queryBus        *busquery.Bus
	authMiddleware  *middleware.AuthMiddleware
}

// NewCQRSPrintHandler 创建基于 CQRS 的打印处理器
func NewCQRSPrintHandler(
	commandBus *buscommand.Bus, 
	queryBus *busquery.Bus,
	authMiddleware *middleware.AuthMiddleware,
) *CQRSPrintHandler {
	return &CQRSPrintHandler{
		commandBus:     commandBus,
		queryBus:       queryBus,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSPrintHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /print/bill", h.authMiddleware.RequireAuth(h.PrintBill))
	mux.HandleFunc("POST /print/lease", h.authMiddleware.RequireAuth(h.PrintLease))
	mux.HandleFunc("POST /print/invoice", h.authMiddleware.RequireAuth(h.PrintInvoice))
	mux.HandleFunc("GET /print/content/{billId}", h.authMiddleware.RequireAuth(h.GetPrintContent))
	mux.HandleFunc("GET /print/jobs", h.authMiddleware.RequireAuth(h.ListPrintJobs))
}

// PrintBill 打印账单
func (h *CQRSPrintHandler) PrintBill(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BillID string `json:"bill_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := print.PrintBillCommand{
		BillID: req.BillID,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobID := result.(string)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"jobId": jobID})
}

// PrintLease 打印租约
func (h *CQRSPrintHandler) PrintLease(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LeaseID string `json:"lease_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := print.PrintLeaseCommand{
		LeaseID: req.LeaseID,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobID := result.(string)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"jobId": jobID})
}

// PrintInvoice 打印发票
func (h *CQRSPrintHandler) PrintInvoice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BillID string `json:"bill_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := print.PrintInvoiceCommand{
		BillID: req.BillID,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobID := result.(string)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"jobId": jobID})
}

// GetPrintContent 获取打印内容
func (h *CQRSPrintHandler) GetPrintContent(w http.ResponseWriter, r *http.Request) {
	billID := r.PathValue("billId")
	q := print.GetPrintContentQuery{BillID: billID}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := result.([]byte)
	w.Header().Set("Content-Type", "application/rtf")
	_, _ = w.Write(content)
}

// ListPrintJobs 列出打印作业
func (h *CQRSPrintHandler) ListPrintJobs(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	queryParams := r.URL.Query()

	var status string
	var startDateStr, endDateStr string
	var offset, limit int

	if v := queryParams.Get("status"); v != "" {
		status = v
	}
	if v := queryParams.Get("start_date"); v != "" {
		startDateStr = v
	}
	if v := queryParams.Get("end_date"); v != "" {
		endDateStr = v
	}
	if v := queryParams.Get("offset"); v != "" {
		fmt.Sscanf(v, "%d", &offset)
	}
	if v := queryParams.Get("limit"); v != "" {
		fmt.Sscanf(v, "%d", &limit)
	}

	// 构建查询对象
	q := print.ListPrintJobsQuery{
		Status: status,
		Offset: offset,
		Limit:  limit,
	}

	// 解析日期
	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			q.StartDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			q.EndDate = &t
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
