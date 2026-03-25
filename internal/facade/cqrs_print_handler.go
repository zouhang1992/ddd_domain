package facade

import (
	"encoding/json"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"net/http"
)

// CQRSPrintHandler 基于 CQRS 的打印 HTTP 处理器
type CQRSPrintHandler struct {
	commandBus *buscommand.Bus
	queryBus   *busquery.Bus
}

// NewCQRSPrintHandler 创建基于 CQRS 的打印处理器
func NewCQRSPrintHandler(commandBus *buscommand.Bus, queryBus *busquery.Bus) *CQRSPrintHandler {
	return &CQRSPrintHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSPrintHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /print/bill", h.PrintBill)
	mux.HandleFunc("POST /print/lease", h.PrintLease)
	mux.HandleFunc("POST /print/invoice", h.PrintInvoice)
	mux.HandleFunc("GET /print/content/{billId}", h.GetPrintContent)
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

	cmd := command.PrintBillCommand{
		BillID: req.BillID,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobID := result.(string)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
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

	cmd := command.PrintLeaseCommand{
		LeaseID: req.LeaseID,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobID := result.(string)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
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

	cmd := command.PrintInvoiceCommand{
		BillID: req.BillID,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobID := result.(string)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}

// GetPrintContent 获取打印内容
func (h *CQRSPrintHandler) GetPrintContent(w http.ResponseWriter, r *http.Request) {
	billID := r.PathValue("billId")
	q := query.GetPrintContentQuery{BillID: billID}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := result.([]byte)
	w.Header().Set("Content-Type", "application/rtf")
	_, _ = w.Write(content)
}
