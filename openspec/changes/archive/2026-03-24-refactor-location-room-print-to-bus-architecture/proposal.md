## Why

当前位置管理和房间管理仍然使用传统的服务层架构，而打印服务也是独立的组件。为了保持整个系统架构的一致性，我们需要将这些组件调整为与之前实现的完整总线架构（command bus、query bus、event bus）相同的模式。

这样可以带来以下好处：
- 统一的架构风格，降低维护成本
- 更好的解耦和可扩展性
- 支持事件驱动的架构
- 便于添加新功能和修改现有功能

## What Changes

### 位置管理（Location Management）

**BREAKING**
- 将 LocationService 重构为使用总线架构
- 创建位置管理的 Command 和 Query
- 实现对应的 CommandHandler 和 QueryHandler
- 重构 LocationHandler 直接使用总线而非传统服务
- 将 handler 改成实现对应的接口，不再使用函数方式传参

### 房间管理（Room Management）

**BREAKING**
- 将 RoomService 重构为使用总线架构
- 创建房间管理的 Command 和 Query
- 实现对应的 CommandHandler 和 QueryHandler
- 重构 RoomHandler 直接使用总线而非传统服务
- 将 handler 改成实现对应的接口，不再使用函数方式传参

### 打印服务（Print Service）

**BREAKING**
- 将 PrintService 重构为使用总线架构
- 可能调整为事件驱动或其他方式
- 实现打印命令和查询的处理

## Capabilities

### New Capabilities

- `location-management`: 位置管理的总线架构实现
- `room-management`: 房间管理的总线架构实现
- `print-service`: 打印服务的总线架构实现

### Modified Capabilities

- `bus-integration`: 需要更新总线集成规范，添加对位置管理和房间管理的支持

## Impact

### Affected Code

- `internal/application/command/*` - 新增位置和房间管理的命令
- `internal/application/query/*` - 新增位置和房间管理的查询
- `internal/application/command/handler/*` - 新增命令处理器
- `internal/application/query/handler/*` - 新增查询处理器
- `internal/application/service/location.go` - 废弃，不再使用
- `internal/application/service/room.go` - 废弃，不再使用
- `internal/facade/location_handler.go` - 重构
- `internal/facade/room_handler.go` - 重构
- `internal/domain/model/*` - 可能需要调整位置和房间模型以支持事件
- `internal/domain/repository/*` - 更新接口以支持查询和命令
- `internal/infrastructure/bus/*` - 可能需要更新总线接口
- `cmd/api/main.go` - 重构初始化流程
