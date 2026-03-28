package model

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

// AggregateRoot 聚合根接口
type AggregateRoot interface {
	ID() string
	Version() int
	Events() []events.DomainEvent
	ClearEvents()
}

// BaseAggregateRoot 基础聚合根实现
type BaseAggregateRoot struct {
	id      string
	version int
	events  []events.DomainEvent
}

func NewBaseAggregateRoot(id string) BaseAggregateRoot {
	return BaseAggregateRoot{
		id:      id,
		version: 0,
		events:  []events.DomainEvent{},
	}
}

func (a *BaseAggregateRoot) ID() string {
	return a.id
}

func (a *BaseAggregateRoot) Version() int {
	return a.version
}

func (a *BaseAggregateRoot) Events() []events.DomainEvent {
	return a.events
}

func (a *BaseAggregateRoot) ClearEvents() {
	a.events = []events.DomainEvent{}
}

func (a *BaseAggregateRoot) RecordEvent(evt events.DomainEvent) {
	a.version++
	a.events = append(a.events, evt)
}
