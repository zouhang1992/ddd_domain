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
	IDField      string               `json:"id"`
	VersionField int                  `json:"-"`
	EventsField  []events.DomainEvent `json:"-"`
}

func NewBaseAggregateRoot(id string) BaseAggregateRoot {
	return BaseAggregateRoot{
		IDField:      id,
		VersionField: 0,
		EventsField:  []events.DomainEvent{},
	}
}

func (a *BaseAggregateRoot) ID() string {
	return a.IDField
}

func (a *BaseAggregateRoot) Version() int {
	return a.VersionField
}

func (a *BaseAggregateRoot) Events() []events.DomainEvent {
	return a.EventsField
}

func (a *BaseAggregateRoot) ClearEvents() {
	a.EventsField = []events.DomainEvent{}
}

func (a *BaseAggregateRoot) RecordEvent(evt events.DomainEvent) {
	a.VersionField++
	a.EventsField = append(a.EventsField, evt)
}
