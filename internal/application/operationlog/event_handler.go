package operationlog

import (
	"reflect"
	"strings"

	"github.com/google/uuid"

	operationlogmodel "github.com/zouhang1992/ddd_domain/internal/domain/operationlog/model"
	operationlogrepo "github.com/zouhang1992/ddd_domain/internal/domain/operationlog/repository"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
	"go.uber.org/zap"
)

// OperationLogEventHandler 操作日志事件处理器
type OperationLogEventHandler struct {
	repo operationlogrepo.OperationLogRepository
	log  *zap.Logger
}

// NewOperationLogEventHandler 创建操作日志事件处理器
func NewOperationLogEventHandler(repo operationlogrepo.OperationLogRepository, logger *zap.Logger) *OperationLogEventHandler {
	return &OperationLogEventHandler{
		repo: repo,
		log:  logger,
	}
}

// Handle 处理领域事件
func (h *OperationLogEventHandler) Handle(evt event.DomainEvent) error {
	h.log.Info("Processing event for operation log",
		zap.String("event", evt.EventName()))

	// 通过反射提取事件信息
	val := reflect.ValueOf(evt)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	// 从事件名称中提取 domain type 和 action
	eventName := evt.EventName()
	domainType, action := parseEventName(eventName)

	// 尝试通过接口获取 AggregateID
	aggregateID := ""
	if aggIDGetter, ok := evt.(interface{ AggregateID() string }); ok {
		aggregateID = aggIDGetter.AggregateID()
	}

	// 如果接口方式不行，尝试反射提取
	if aggregateID == "" {
		aggregateID = extractFieldAsString(val, "ID")
	}
	if aggregateID == "" {
		aggregateID = extractFieldAsString(val, "AggregateID")
	}

	// 构建事件详情
	details := make(map[string]any)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.IsExported() {
			details[field.Name] = val.Field(i).Interface()
		}
	}

	// 创建操作日志
	log := operationlogmodel.NewOperationLog(
		uuid.NewString(),
		evt.OccurredAt(),
		eventName,
		domainType,
		aggregateID,
		"", // OperatorID - 从认证上下文获取，暂留空
		action,
		details,
		make(map[string]any), // Metadata - 暂留空
	)

	if err := h.repo.Save(log); err != nil {
		h.log.Error("Failed to save operation log",
			zap.String("event", eventName),
			zap.Error(err))
		return err
	}

	h.log.Info("Operation log saved",
		zap.String("event", eventName),
		zap.String("log_id", log.ID))

	return nil
}

// parseEventName 解析事件名称，提取 domain type 和 action
// 格式: "domain.action" → domainType: "domain", action: "action"
func parseEventName(eventName string) (domainType, action string) {
	parts := strings.SplitN(eventName, ".", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "unknown", eventName
}

// extractFieldAsString 通过反射提取字符串字段
func extractFieldAsString(val reflect.Value, fieldName string) string {
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return ""
	}
	if field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}
