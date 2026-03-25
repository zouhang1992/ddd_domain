package handler

import (
	"github.com/google/uuid"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// Handler 操作日志事件处理器
type Handler struct {
	repo repository.OperationLogRepository
}

// NewHandler 创建操作日志事件处理器
func NewHandler(repo repository.OperationLogRepository) *Handler {
	return &Handler{repo: repo}
}

// Handle 处理领域事件并记录操作日志
func (h *Handler) Handle(evt event.DomainEvent) error {
	data := event.ExtractEventData(evt)
	data.ID = uuid.New().String()

	log := model.NewOperationLog(
		data.ID,
		data.Timestamp,
		data.EventName,
		data.DomainType,
		data.AggregateID,
		data.OperatorID,
		data.Action,
		data.Details,
		data.Metadata,
	)

	return h.repo.Save(log)
}

// SubscribeToAllEvents 订阅所有可能的事件
func (h *Handler) SubscribeToAllEvents(bus *event.Bus) {
	for _, eventName := range event.GetAllEventNames() {
		bus.Subscribe(eventName, h)
	}
}

// HandleFunc 返回包装后的事件处理函数
func (h *Handler) HandleFunc() event.HandlerFunc {
	return event.HandlerFunc(h.Handle)
}
