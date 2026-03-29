package lease

import (
	"reflect"

	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
	"github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
	"go.uber.org/zap"
)

// LeaseRoomEventHandler 租约房间事件处理器
type LeaseRoomEventHandler struct {
	roomRepo roomrepo.RoomRepository
	log      *zap.Logger
}

// NewLeaseRoomEventHandler 创建租约房间事件处理器
func NewLeaseRoomEventHandler(roomRepo roomrepo.RoomRepository, logger *zap.Logger) *LeaseRoomEventHandler {
	return &LeaseRoomEventHandler{
		roomRepo: roomRepo,
		log:      logger,
	}
}

// Handle 处理领域事件
func (h *LeaseRoomEventHandler) Handle(evt event.DomainEvent) error {
	h.log.Info("Processing lease event for room state",
		zap.String("event", evt.EventName()))

	switch evt.EventName() {
	case "lease.activated":
		return h.handleLeaseEvent(evt, "rented")
	case "lease.checkout":
		return h.handleLeaseEvent(evt, "available")
	case "lease.expired":
		return h.handleLeaseEvent(evt, "available")
	default:
		h.log.Debug("Event type not handled", zap.String("event", evt.EventName()))
		return nil
	}
}

func (h *LeaseRoomEventHandler) handleLeaseEvent(evt event.DomainEvent, targetState string) error {
	val := reflect.ValueOf(evt)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	roomIDField := val.FieldByName("RoomID")
	if !roomIDField.IsValid() {
		h.log.Error("Event does not have RoomID field", zap.String("event", evt.EventName()))
		return nil
	}

	roomID, ok := roomIDField.Interface().(string)
	if !ok || roomID == "" {
		h.log.Error("Invalid RoomID in event", zap.String("event", evt.EventName()))
		return nil
	}

	h.log.Info("Handling lease event",
		zap.String("event", evt.EventName()),
		zap.String("room_id", roomID),
		zap.String("target_state", targetState))

	room, err := h.roomRepo.FindByID(roomID)
	if err != nil {
		h.log.Error("Failed to find room",
			zap.String("room_id", roomID),
			zap.Error(err))
		return err
	}

	if room == nil {
		h.log.Debug("Room not found for event", zap.String("room_id", roomID))
		return nil
	}

	h.log.Debug("Room found", zap.String("room_id", roomID), zap.String("current_status", string(room.Status)))

	switch targetState {
	case "rented":
		room.MarkRented()
	case "available":
		room.MarkAvailable()
	default:
		h.log.Warn("Unknown target state", zap.String("state", targetState))
		return nil
	}

	if err := h.roomRepo.Save(room); err != nil {
		h.log.Error("Failed to save room state",
			zap.String("room_id", roomID),
			zap.Error(err))
		return err
	}

	h.log.Info("Room state updated successfully",
		zap.String("room_id", roomID),
		zap.String("new_status", string(room.Status)))

	return nil
}
