package facade

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/zouhang1992/ddd_domain/internal/application/query"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
)

// CQRSOperationLogHandler 基于 CQRS 的操作日志 HTTP 处理器
type CQRSOperationLogHandler struct {
	queryBus *busquery.Bus
}

// NewCQRSOperationLogHandler 创建基于 CQRS 的操作日志处理器
func NewCQRSOperationLogHandler(queryBus *busquery.Bus) *CQRSOperationLogHandler {
	return &CQRSOperationLogHandler{queryBus: queryBus}
}

// RegisterRoutes 注册路由
func (h *CQRSOperationLogHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /operation-logs", h.List)
	mux.HandleFunc("GET /operation-logs/{id}", h.Get)
}

// List 列出操作日志
func (h *CQRSOperationLogHandler) List(w http.ResponseWriter, r *http.Request) {
	var (
		domainType  = r.URL.Query().Get("domainType")
		eventName   = r.URL.Query().Get("eventName")
		aggregateID = r.URL.Query().Get("aggregateId")
		operatorID  = r.URL.Query().Get("operatorId")
	)

	// 解析时间范围
	var (
		startTime *time.Time
		endTime   *time.Time
		err       error
	)

	if startStr := r.URL.Query().Get("startTime"); startStr != "" {
		t, parseErr := time.Parse(time.RFC3339, startStr)
		if parseErr == nil {
			startTime = &t
		}
	}

	if endStr := r.URL.Query().Get("endTime"); endStr != "" {
		t, parseErr := time.Parse(time.RFC3339, endStr)
		if parseErr == nil {
			endTime = &t
		}
	}

	// 解析分页参数
	var (
		offset = 0
		limit  = 20
	)

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if v, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && v >= 0 {
			offset = int(v)
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if v, err := strconv.ParseInt(limitStr, 10, 32); err == nil && v > 0 && v <= 100 {
			limit = int(v)
		}
	}

	// 构建查询
	q := query.ListOperationLogsQuery{
		DomainType:  domainType,
		EventName:   eventName,
		AggregateID: aggregateID,
		OperatorID:  operatorID,
		StartTime:   startTime,
		EndTime:     endTime,
		Offset:      offset,
		Limit:       limit,
	}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult, ok := result.(*query.OperationLogsQueryResult)
	if !ok {
		http.Error(w, "invalid query result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult)
}

// Get 获取单条操作日志
func (h *CQRSOperationLogHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	q := query.GetOperationLogQuery{ID: id}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult, ok := result.(*query.OperationLogQueryResult)
	if !ok {
		http.Error(w, "invalid query result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.OperationLog)
}
