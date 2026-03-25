package event

import (
	"reflect"
	"strings"
	"time"
)

// OperationLogEventData 操作日志事件数据（用于在不同包之间传递）
type OperationLogEventData struct {
	ID          string
	Timestamp   time.Time
	EventName   string
	DomainType  string
	AggregateID string
	OperatorID  string
	Action      string
	Details     map[string]interface{}
	Metadata    map[string]interface{}
}

// ExtractEventData 从 DomainEvent 中提取数据
func ExtractEventData(event DomainEvent) *OperationLogEventData {
	// 解析事件信息
	timestamp := event.OccurredAt()
	eventName := event.EventName()

	// 从事件名称中获取领域类型（如 "landlord" 从 "landlord.created"）
	domainType, action := parseEventName(eventName)

	// 从事件中提取聚合ID
	aggregateID := extractAggregateID(event)

	// 从事件中提取操作人信息（暂时留空，待认证系统完善）
	operatorID := ""

	// 提取事件详情
	details := extractEventDetails(event)

	// 元数据
	metadata := map[string]interface{}{
		"event_type": reflect.TypeOf(event).Name(),
	}

	return &OperationLogEventData{
		Timestamp:   timestamp,
		EventName:   eventName,
		DomainType:  domainType,
		AggregateID: aggregateID,
		OperatorID:  operatorID,
		Action:      action,
		Details:     details,
		Metadata:    metadata,
	}
}

// parseEventName 解析事件名称，提取领域类型和操作类型
func parseEventName(eventName string) (domainType, action string) {
	parts := strings.SplitN(eventName, ".", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "unknown", "unknown"
}

// extractAggregateID 从事件中提取聚合根ID
func extractAggregateID(event DomainEvent) string {
	// 使用反射从事件结构体中提取ID字段
	val := reflect.ValueOf(event)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i).Interface()

		// 查找包含ID的字段，如 LandlordID, LeaseID, RoomID 等
		if strings.Contains(strings.ToLower(field.Name), "id") && fieldVal != "" {
			idStr, ok := fieldVal.(string)
			if ok && idStr != "" {
				return idStr
			}
		}
	}

	return ""
}

// extractEventDetails 提取事件的详细信息
func extractEventDetails(event DomainEvent) map[string]interface{} {
	// 使用反射将事件转换为map
	val := reflect.ValueOf(event)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	details := make(map[string]interface{})

	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i).Interface()

		// 跳过 BaseEvent 字段，避免重复
		if field.Name == "BaseEvent" {
			continue
		}

		// 处理不同类型字段
		switch v := fieldVal.(type) {
		case string:
			if v != "" {
				details[field.Name] = v
			}
		case int, int8, int16, int32, int64:
			if v != 0 {
				details[field.Name] = v
			}
		case float32, float64:
			if v != 0.0 {
				details[field.Name] = v
			}
		case bool:
			details[field.Name] = v
		case time.Time:
			if !v.IsZero() {
				details[field.Name] = v.Format(time.RFC3339)
			}
		case []byte:
			if len(v) > 0 {
				// 对于字节数组，我们不存储完整内容，避免日志过大
				details[field.Name] = "[binary data]"
			}
		case []string:
			if len(v) > 0 {
				details[field.Name] = v
			}
		}
	}

	return details
}

// GetAllEventNames 获取所有事件名称列表
func GetAllEventNames() []string {
	return []string{
		// 房东相关事件
		"landlord.created",
		"landlord.updated",
		"landlord.deleted",

		// 租约相关事件
		"lease.created",
		"lease.updated",
		"lease.deleted",
		"lease.renewed",
		"lease.checkout",
		"lease.activated",

		// 账单相关事件
		"bill.created",
		"bill.updated",
		"bill.deleted",
		"bill.paid",

		// 位置相关事件
		"location.created",
		"location.updated",
		"location.deleted",

		// 房间相关事件
		"room.created",
		"room.updated",
		"room.deleted",

		// 打印相关事件
		"bill.printed",
		"lease.printed",
		"invoice.printed",
		"print.failed",
	}
}
