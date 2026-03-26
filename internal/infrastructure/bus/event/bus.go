package event

import (
	"sync"

	"go.uber.org/zap"
)

// Bus 事件总线
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]EventHandler
	log         *zap.Logger
}

// NewBus 创建事件总线
func NewBus(logger *zap.Logger) *Bus {
	return &Bus{
		subscribers: make(map[string][]EventHandler),
		log:         logger,
	}
}

// Subscribe 订阅事件
func (b *Bus) Subscribe(eventName string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventName] = append(b.subscribers[eventName], handler)
	b.log.Info("Event subscriber added", zap.String("event", eventName))
}

// Publish 发布事件（同步）
func (b *Bus) Publish(event DomainEvent) error {
	eventName := event.EventName()
	b.log.Info("Publishing event", zap.String("event", eventName))

	b.mu.RLock()
	handlers, ok := b.subscribers[eventName]
	b.mu.RUnlock()

	if !ok {
		b.log.Debug("No subscribers for event", zap.String("event", eventName))
		return nil
	}

	b.log.Debug("Event handlers found",
		zap.String("event", eventName),
		zap.Int("handler_count", len(handlers)))

	for _, handler := range handlers {
		b.log.Debug("Calling event handler", zap.String("event", eventName))
		if err := handler.Handle(event); err != nil {
			b.log.Error("Event handler failed",
				zap.String("event", eventName),
				zap.Error(err))
			return err
		}
	}

	b.log.Info("Event published successfully", zap.String("event", eventName))
	return nil
}

// PublishAsync 发布事件（异步）
func (b *Bus) PublishAsync(event DomainEvent) {
	eventName := event.EventName()
	b.log.Info("Publishing event asynchronously", zap.String("event", eventName))

	b.mu.RLock()
	handlers, ok := b.subscribers[eventName]
	handlersCopy := make([]EventHandler, len(handlers))
	if ok {
		copy(handlersCopy, handlers)
	}
	b.mu.RUnlock()

	if !ok {
		b.log.Debug("No subscribers for event", zap.String("event", eventName))
		return
	}

	b.log.Debug("Event handlers found for async publish",
		zap.String("event", eventName),
		zap.Int("handler_count", len(handlersCopy)))

	go func() {
		for _, handler := range handlersCopy {
			b.log.Debug("Calling event handler asynchronously", zap.String("event", eventName))
			if err := handler.Handle(event); err != nil {
				b.log.Error("Async event handler failed",
					zap.String("event", eventName),
					zap.Error(err))
			}
		}
		b.log.Info("Async event published successfully", zap.String("event", eventName))
	}()
}
