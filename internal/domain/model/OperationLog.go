package model

import (
	"encoding/json"
	"time"
)

// OperationLog 操作日志聚合根
type OperationLog struct {
	id          string
	timestamp   time.Time
	eventName   string
	domainType  string
	aggregateID string
	operatorID  string
	action      string
	details     map[string]interface{}
	metadata    map[string]interface{}
	createdAt   time.Time
}

// NewOperationLog 创建操作日志
func NewOperationLog(
	id string,
	timestamp time.Time,
	eventName string,
	domainType string,
	aggregateID string,
	operatorID string,
	action string,
	details map[string]interface{},
	metadata map[string]interface{},
) *OperationLog {
	return &OperationLog{
		id:          id,
		timestamp:   timestamp,
		eventName:   eventName,
		domainType:  domainType,
		aggregateID: aggregateID,
		operatorID:  operatorID,
		action:      action,
		details:     details,
		metadata:    metadata,
		createdAt:   time.Now(),
	}
}

// ID 获取ID
func (l OperationLog) ID() string {
	return l.id
}

// Timestamp 获取操作发生时间
func (l OperationLog) Timestamp() time.Time {
	return l.timestamp
}

// EventName 获取事件名称
func (l OperationLog) EventName() string {
	return l.eventName
}

// DomainType 获取领域类型
func (l OperationLog) DomainType() string {
	return l.domainType
}

// AggregateID 获取关联的聚合根ID
func (l OperationLog) AggregateID() string {
	return l.aggregateID
}

// OperatorID 获取操作人ID
func (l OperationLog) OperatorID() string {
	return l.operatorID
}

// Action 获取操作类型
func (l OperationLog) Action() string {
	return l.action
}

// Details 获取详细数据
func (l OperationLog) Details() map[string]interface{} {
	return l.details
}

// Metadata 获取元数据
func (l OperationLog) Metadata() map[string]interface{} {
	return l.metadata
}

// CreatedAt 获取创建时间
func (l OperationLog) CreatedAt() time.Time {
	return l.createdAt
}

// SetDetails 设置详细数据
func (l *OperationLog) SetDetails(details map[string]interface{}) {
	l.details = details
}

// SetMetadata 设置元数据
func (l *OperationLog) SetMetadata(metadata map[string]interface{}) {
	l.metadata = metadata
}

// MarshalDetails 将详细数据转换为JSON字符串
func (l OperationLog) MarshalDetails() (string, error) {
	data, err := json.Marshal(l.details)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalDetails 从JSON字符串解析详细数据
func (l *OperationLog) UnmarshalDetails(data string) error {
	var details map[string]interface{}
	if err := json.Unmarshal([]byte(data), &details); err != nil {
		return err
	}
	l.details = details
	return nil
}

// MarshalMetadata 将元数据转换为JSON字符串
func (l OperationLog) MarshalMetadata() (string, error) {
	data, err := json.Marshal(l.metadata)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalMetadata 从JSON字符串解析元数据
func (l *OperationLog) UnmarshalMetadata(data string) error {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		return err
	}
	l.metadata = metadata
	return nil
}
