package event

import "sync"

// Bus 事件总线
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]EventHandler
}

// NewBus 创建事件总线
func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (b *Bus) Subscribe(eventName string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventName] = append(b.subscribers[eventName], handler)
}

// Publish 发布事件（同步）
func (b *Bus) Publish(event DomainEvent) error {
	b.mu.RLock()
	handlers, ok := b.subscribers[event.EventName()]
	b.mu.RUnlock()

	if !ok {
		return nil
	}

	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			return err
		}
	}
	return nil
}

// PublishAsync 发布事件（异步）
func (b *Bus) PublishAsync(event DomainEvent) {
	b.mu.RLock()
	handlers, ok := b.subscribers[event.EventName()]
	handlersCopy := make([]EventHandler, len(handlers))
	if ok {
		copy(handlersCopy, handlers)
	}
	b.mu.RUnlock()

	if !ok {
		return
	}

	go func() {
		for _, handler := range handlersCopy {
			_ = handler.Handle(event)
		}
	}()
}
