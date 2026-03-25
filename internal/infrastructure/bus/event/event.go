package event

import "time"

// DomainEvent 定义领域事件接口
type DomainEvent interface {
	// EventName 返回事件名称
	EventName() string
	// OccurredAt 返回事件发生时间
	OccurredAt() time.Time
}

// EventHandler 定义事件处理器接口
type EventHandler interface {
	// Handle 处理事件
	Handle(event DomainEvent) error
}

// HandlerFunc 函数类型适配器
type HandlerFunc func(event DomainEvent) error

// Handle 实现 EventHandler 接口
func (f HandlerFunc) Handle(event DomainEvent) error {
	return f(event)
}

// BaseEvent 基础事件实现，可嵌入到具体事件中
type BaseEvent struct {
	name       string
	occurredAt time.Time
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(name string) BaseEvent {
	return BaseEvent{
		name:       name,
		occurredAt: time.Now(),
	}
}

// EventName 实现 DomainEvent 接口
func (e BaseEvent) EventName() string {
	return e.name
}

// OccurredAt 实现 DomainEvent 接口
func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}
