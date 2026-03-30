package model

import "time"

// OperationLog 操作日志领域模型
type OperationLog struct {
	ID          string
	Timestamp   time.Time
	EventName   string
	DomainType  string
	AggregateID string
	OperatorID  string
	Action      string // created, updated, deleted, activated, checked-out, etc.
	Details     map[string]interface{}
	Metadata    map[string]interface{}
	CreatedAt   time.Time
}

// NewOperationLog 创建新的操作日志
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
	now := time.Now()
	return &OperationLog{
		ID:          id,
		Timestamp:   timestamp,
		EventName:   eventName,
		DomainType:  domainType,
		AggregateID: aggregateID,
		OperatorID:  operatorID,
		Action:      action,
		Details:     details,
		Metadata:    metadata,
		CreatedAt:   now,
	}
}
