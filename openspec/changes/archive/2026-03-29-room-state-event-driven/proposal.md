## Why

当前系统中，租约状态的变化不会自动影响房间状态，需要手动维护，这可能导致数据不一致和用户体验问题。我们需要将这种耦合关系调整为事件驱动的方式，当租约状态发生变化时，自动更新对应的房间状态。

## What Changes

- 在房间模型中添加状态变更事件的记录功能
- 创建事件处理器来监听租约状态变化事件
- 当租约激活时，自动将房间状态设置为已出租
- 当租约退租或过期时，自动将房间状态设置为可出租
- 更新相关的依赖注入配置

## Capabilities

### New Capabilities

- `room-state-event-driven`: 实现房间状态的事件驱动管理

### Modified Capabilities

- `lease-management`: 租约状态变化时会发布领域事件
- `room-management`: 房间状态会根据租约事件自动更新

## Impact

- `internal/domain/room/model/room.go`: 添加状态变更事件记录功能
- `internal/application/event/handler/`: 创建新的事件处理器
- `cmd/api/main.go`: 更新依赖注入配置
- `internal/domain/lease/model/lease.go`: 确保租约状态变化时发布事件