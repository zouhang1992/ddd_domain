## Why

在之前的变更中，我们已经将多个核心模块（房东管理、租约管理、账单管理、位置管理、房间管理和打印服务）从传统的应用服务架构重构为总线架构（Command/Query/Event Bus）。这些模块现在已经完全通过总线管理，之前的应用服务层代码（如 `LocationService`、`RoomService` 等）已经不再被使用。

为了保持代码库的整洁和一致性，我们需要删除这些已经不再被使用的传统应用服务层代码。

## What Changes

**BREAKING**

### 删除的文件：
- `internal/application/service/location.go` - 位置管理传统应用服务（已由总线架构替代）
- `internal/application/service/room.go` - 房间管理传统应用服务（已由总线架构替代）
- `internal/facade/location_handler.go` - 位置管理传统 HTTP 处理器（已由 CQRSLocationHandler 替代）
- `internal/facade/room_handler.go` - 房间管理传统 HTTP 处理器（已由 CQRSRoomHandler 替代）

**修改的文件：**
- `cmd/api/main.go` - 删除对传统应用服务的引用

## Impact

### Affected Code
- `internal/application/service/*` - 删除位置和房间管理的应用服务
- `internal/facade/*` - 删除传统的位置和房间管理 HTTP 处理器
- `cmd/api/main.go` - 清理传统服务的初始化代码

### Capabilities Modified
- `bus-integration` - 更新为只包含总线架构的代码
- `application-services` - 移除非总线管理的应用服务

## Migration

1. 不再需要使用传统的应用服务
2. 所有请求现在都应该通过 HTTP Handler -> Command/Query Bus -> Handler 的路径处理
3. 对于位置和房间管理，现在应该使用 `NewCQRSLocationHandler` 和 `NewCQRSRoomHandler` 而非传统的 `NewLocationHandler` 和 `NewRoomHandler`
