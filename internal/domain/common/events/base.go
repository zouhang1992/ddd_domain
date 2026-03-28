package events

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventName() string
	EventID() string
	TimeStamp() time.Time
	AggregateID() string
	Version() int
}

// BaseEvent 事件基类
type BaseEvent struct {
	eventName   string
	eventID     string
	timeStamp   time.Time
	aggregateID string
	version     int
}

func NewBaseEvent(eventName, aggregateID string, version int) BaseEvent {
	return BaseEvent{
		eventName:   eventName,
		eventID:     uuid.NewString(),
		timeStamp:   time.Now(),
		aggregateID: aggregateID,
		version:     version,
	}
}

func (e BaseEvent) EventName() string {
	return e.eventName
}

func (e BaseEvent) EventID() string {
	return e.eventID
}

func (e BaseEvent) TimeStamp() time.Time {
	return e.timeStamp
}

func (e BaseEvent) AggregateID() string {
	return e.aggregateID
}

func (e BaseEvent) Version() int {
	return e.version
}
