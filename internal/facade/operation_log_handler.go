package facade

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/zouhang1992/ddd_domain/internal/infrastructure/middleware"
	operationlogmodel "github.com/zouhang1992/ddd_domain/internal/domain/operationlog/model"
	operationlogrepo "github.com/zouhang1992/ddd_domain/internal/domain/operationlog/repository"
)

// OperationLogHandler 操作日志 HTTP 处理器
type OperationLogHandler struct {
	repo            operationlogrepo.OperationLogRepository
	authMiddleware  *middleware.AuthMiddleware
}

// NewOperationLogHandler 创建操作日志处理器
func NewOperationLogHandler(
	repo operationlogrepo.OperationLogRepository,
	authMiddleware *middleware.AuthMiddleware,
) *OperationLogHandler {
	return &OperationLogHandler{
		repo:           repo,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册路由
func (h *OperationLogHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /operation-logs", h.authMiddleware.RequireAuth(h.ListOperationLogs))
	mux.HandleFunc("GET /operation-logs/{id}", h.authMiddleware.RequireAuth(h.GetOperationLog))
}

// OperationLogResponse 操作日志响应
type OperationLogResponse struct {
	ID          string         `json:"id"`
	Timestamp   string         `json:"timestamp"`
	EventName   string         `json:"event_name"`
	DomainType  string         `json:"domain_type"`
	AggregateID string         `json:"aggregate_id"`
	OperatorID  *string        `json:"operator_id,omitempty"`
	Action      *string        `json:"action,omitempty"`
	Details     map[string]any `json:"details,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   string         `json:"created_at"`
}

// OperationLogsQueryResponse 操作日志查询响应
type OperationLogsQueryResponse struct {
	Items []OperationLogResponse `json:"items"`
	Total int                    `json:"total"`
	Page  int                    `json:"page"`
	Limit int                    `json:"limit"`
}

// toResponse 将领域模型转换为响应
func toResponse(log *operationlogmodel.OperationLog) OperationLogResponse {
	resp := OperationLogResponse{
		ID:          log.ID,
		Timestamp:   log.Timestamp.Format(time.RFC3339),
		EventName:   log.EventName,
		DomainType:  log.DomainType,
		AggregateID: log.AggregateID,
		Details:     log.Details,
		Metadata:    log.Metadata,
		CreatedAt:   log.CreatedAt.Format(time.RFC3339),
	}

	if log.OperatorID != "" {
		resp.OperatorID = &log.OperatorID
	}
	if log.Action != "" {
		resp.Action = &log.Action
	}

	return resp
}

// ListOperationLogs 列出操作日志
func (h *OperationLogHandler) ListOperationLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// 解析分页参数
	offsetStr := query.Get("offset")
	limitStr := query.Get("limit")
	pageStr := query.Get("page")

	offset := 0
	limit := 20

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

	var logs []*operationlogmodel.OperationLog
	var total int
	var err error

	// 检查查询参数（同时支持驼峰和下划线命名）
	domainType := query.Get("domain_type")
	if domainType == "" {
		domainType = query.Get("domainType")
	}
	aggregateID := query.Get("aggregate_id")
	if aggregateID == "" {
		aggregateID = query.Get("aggregateId")
	}
	startTimeStr := query.Get("start_time")
	if startTimeStr == "" {
		startTimeStr = query.Get("startTime")
	}
	endTimeStr := query.Get("end_time")
	if endTimeStr == "" {
		endTimeStr = query.Get("endTime")
	}

	switch {
	case domainType != "" && aggregateID != "":
		logs, total, err = h.repo.FindByDomainTypeAndAggregateID(domainType, aggregateID, offset, limit)
	case domainType != "":
		logs, total, err = h.repo.FindByDomainType(domainType, offset, limit)
	case aggregateID != "":
		logs, err = h.repo.FindByAggregateID(aggregateID)
		total = len(logs)
	case startTimeStr != "" && endTimeStr != "":
		var startTime, endTime time.Time
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, "invalid start_time format", http.StatusBadRequest)
			return
		}
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			http.Error(w, "invalid end_time format", http.StatusBadRequest)
			return
		}
		logs, total, err = h.repo.FindByTimeRange(startTime, endTime, offset, limit)
	default:
		// 默认查询所有，使用时间范围（最近30天）
		endTime := time.Now()
		startTime := endTime.AddDate(0, 0, -30)
		logs, total, err = h.repo.FindByTimeRange(startTime, endTime, offset, limit)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 转换为响应
	items := make([]OperationLogResponse, len(logs))
	for i, log := range logs {
		items[i] = toResponse(log)
	}

	page := (offset / limit) + 1
	resp := OperationLogsQueryResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// GetOperationLog 获取单个操作日志
func (h *OperationLogHandler) GetOperationLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	log, err := h.repo.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if log == nil {
		http.Error(w, "operation log not found", http.StatusNotFound)
		return
	}

	resp := toResponse(log)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
