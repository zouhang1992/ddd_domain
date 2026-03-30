## 1. Summary

将租约状态影响房间状态调整为事件驱动的方式，实现租约状态变化时自动更新对应房间状态的功能。

## 2. Requirements

### 2.1 当租约激活时
- **场景**: 租约被成功激活
- **触发条件**: 收到 `LeaseActivated` 事件
- **动作**: 将对应的房间状态设置为 "rented"
- **事件记录**: 房间状态变更时记录 `RoomRented` 事件

### 2.2 当租约退租时
- **场景**: 租约被成功退租
- **触发条件**: 收到 `LeaseCheckout` 事件
- **动作**: 将对应的房间状态设置为 "available"
- **事件记录**: 房间状态变更时记录 `RoomAvailable` 事件

### 2.3 当租约过期时
- **场景**: 租约自然过期
- **触发条件**: 收到 `LeaseExpired` 事件
- **动作**: 将对应的房间状态设置为 "available"
- **事件记录**: 房间状态变更时记录 `RoomAvailable` 事件

### 2.4 事件处理
- **可靠性**: 事件处理过程中遇到错误需要记录详细日志
- **一致性**: 保证租约状态与房间状态的最终一致性

## 3. Domain Events

### 3.1 现有事件（输入）

#### LeaseActivated
```go
type LeaseActivated struct {
    events.BaseEvent
    RoomID string
}
```

#### LeaseCheckout
```go
type LeaseCheckout struct {
    events.BaseEvent
    RoomID string
}
```

#### LeaseExpired
```go
type LeaseExpired struct {
    events.BaseEvent
    RoomID string
}
```

### 3.2 新事件（输出）

#### RoomRented
```go
type RoomRented struct {
    events.BaseEvent
}
```

#### RoomAvailable
```go
type RoomAvailable struct {
    events.BaseEvent
}
```

## 4. Architecture

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Lease Model │────▶│  Event Bus   │────▶│ Event Handler│
└──────────────┘     └──────────────┘     └──────────────┘
                                               ▲
                                               │
                                               ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Room Model  │◀────│  Room Repo   │◀────│ Event Handler│
└──────────────┘     └──────────────┘     └──────────────┘
```

### 4.1 组件说明

1. **租约模型**: 负责租约状态管理并产生领域事件
2. **事件总线**: 负责事件的分发和传递
3. **事件处理器**: 监听租约事件并协调房间状态更新
4. **房间仓储**: 负责房间数据的持久化
5. **房间模型**: 负责房间状态管理和事件记录

## 5. Implementation

### 5.1 房间模型增强
修改 `internal/domain/room/model/room.go`：

```go
// Room 房间领域模型（聚合根）
type Room struct {
    model.BaseAggregateRoot
    LocationID string
    RoomNumber string
    Status     RoomStatus
    Tags       []string
    Note       string
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// 房间事件
type roomRented struct {
    events.BaseEvent
}

type roomAvailable struct {
    events.BaseEvent
}

// MarkRented 标记房间为已出租（修改版本）
func (r *Room) MarkRented() {
    r.Status = RoomStatusRented
    r.UpdatedAt = time.Now()
    evt := roomRented{
        BaseEvent: events.NewBaseEvent("room.rented", r.ID(), r.Version()),
    }
    r.RecordEvent(evt)
}

// MarkAvailable 标记房间为可出租（修改版本）
func (r *Room) MarkAvailable() {
    r.Status = RoomStatusAvailable
    r.UpdatedAt = time.Now()
    evt := roomAvailable{
        BaseEvent: events.NewBaseEvent("room.available", r.ID(), r.Version()),
    }
    r.RecordEvent(evt)
}
```

### 5.2 事件处理器实现
创建 `internal/application/event/handler/lease_room_event_handler.go`：

```go
package handler

import (
    "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
    "github.com/zouhang1992/ddd_domain/internal/domain/room/model"
    "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
    "go.uber.org/zap"
)

// LeaseRoomEventHandler 租约房间事件处理器
type LeaseRoomEventHandler struct {
    roomRepo repository.RoomRepository
    log      *zap.Logger
}

// NewLeaseRoomEventHandler 创建租约房间事件处理器
func NewLeaseRoomEventHandler(roomRepo repository.RoomRepository, logger *zap.Logger) *LeaseRoomEventHandler {
    return &LeaseRoomEventHandler{
        roomRepo: roomRepo,
        log:      logger,
    }
}

// Handle 处理领域事件
func (h *LeaseRoomEventHandler) Handle(evt event.DomainEvent) error {
    h.log.Info("Processing lease event for room state",
        zap.String("event", evt.EventName()),
        zap.String("aggregate", evt.AggregateID()))

    switch e := evt.(type) {
    case event.LeaseActivated:
        return h.handleLeaseActivated(e)
    case event.LeaseCheckout:
        return h.handleLeaseCheckout(e)
    case event.LeaseExpired:
        return h.handleLeaseExpired(e)
    default:
        h.log.Debug("Event type not handled", zap.String("event_type", evt.EventName()))
        return nil
    }
}

func (h *LeaseRoomEventHandler) handleLeaseActivated(evt event.LeaseActivated) error {
    room, err := h.roomRepo.FindByID(evt.RoomID)
    if err != nil {
        h.log.Error("Failed to find room",
            zap.String("room_id", evt.RoomID),
            zap.Error(err))
        return err
    }

    if room != nil && room.Status != model.RoomStatusRented {
        room.MarkRented()
        if err := h.roomRepo.Save(room); err != nil {
            h.log.Error("Failed to update room state",
                zap.String("room_id", evt.RoomID),
                zap.Error(err))
            return err
        }

        h.log.Info("Room state updated to rented",
            zap.String("room_id", evt.RoomID))
    }

    return nil
}

func (h *LeaseRoomEventHandler) handleLeaseCheckout(evt event.LeaseCheckout) error {
    room, err := h.roomRepo.FindByID(evt.RoomID)
    if err != nil {
        h.log.Error("Failed to find room",
            zap.String("room_id", evt.RoomID),
            zap.Error(err))
        return err
    }

    if room != nil && room.Status != model.RoomStatusAvailable {
        room.MarkAvailable()
        if err := h.roomRepo.Save(room); err != nil {
            h.log.Error("Failed to update room state",
                zap.String("room_id", evt.RoomID),
                zap.Error(err))
            return err
        }

        h.log.Info("Room state updated to available",
            zap.String("room_id", evt.RoomID))
    }

    return nil
}

func (h *LeaseRoomEventHandler) handleLeaseExpired(evt event.LeaseExpired) error {
    room, err := h.roomRepo.FindByID(evt.RoomID)
    if err != nil {
        h.log.Error("Failed to find room",
            zap.String("room_id", evt.RoomID),
            zap.Error(err))
        return err
    }

    if room != nil && room.Status != model.RoomStatusAvailable {
        room.MarkAvailable()
        if err := h.roomRepo.Save(room); err != nil {
            h.log.Error("Failed to update room state",
                zap.String("room_id", evt.RoomID),
                zap.Error(err))
            return err
        }

        h.log.Info("Room state updated to available",
            zap.String("room_id", evt.RoomID))
    }

    return nil
}
```

### 5.3 依赖注入配置
更新 `cmd/api/main.go`：

```go
// 添加房间仓储到依赖图
fx.Provide(
    sqlite.NewRoomRepository,
    // ... 其他依赖
),

// 添加事件处理器到依赖图
fx.Provide(
    handler.NewLeaseRoomEventHandler,
    // ... 其他事件处理器
),
```

## 6. Verification

### 6.1 测试场景

1. **租约激活时房间状态变为已出租**
   - 创建一个房间
   - 创建一个租约关联到该房间
   - 激活租约
   - 验证房间状态变为 "rented"

2. **租约退租时房间状态变为可出租**
   - 创建一个已出租的房间
   - 退租关联的租约
   - 验证房间状态变为 "available"

3. **租约过期时房间状态变为可出租**
   - 创建一个已出租的房间
   - 将租约状态设置为过期
   - 验证房间状态变为 "available"

### 6.2 预期行为

- 事件处理器成功处理事件时，房间状态会正确更新
- 事件处理失败时会记录详细错误日志
- 所有状态变化都会记录相应的领域事件