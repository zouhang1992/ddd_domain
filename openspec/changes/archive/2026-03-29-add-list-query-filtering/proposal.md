## Why

项目当前的列表页面（房东、租约、账单、位置、房间、打印）缺乏查询筛选和分页功能，导致在处理大量数据时用户体验不佳。用户需要能够根据特定条件快速定位所需信息，并支持分页显示以提高加载速度和浏览效率。

## What Changes

- 为所有列表页面添加查询功能，支持根据查询条件筛选数据
- 实现后端分页功能，提高大数据量场景下的性能
- 更新前端查询表单和表格组件，提供良好的用户交互体验
- 统一查询接口和数据格式，保持代码一致性

## Capabilities

### New Capabilities
- `list-query-filtering`: 为所有列表页面提供查询和分页功能

### Modified Capabilities
无

## Impact

**受影响的代码文件：**

### 后端
- `internal/application/query/landlord.go` - 更新房东查询模型
- `internal/application/query/lease.go` - 更新租约查询模型
- `internal/application/query/bill.go` - 更新账单查询模型
- `internal/application/query/location.go` - 更新位置查询模型
- `internal/application/query/room.go` - 更新房间查询模型
- `internal/application/query/print.go` - 更新打印查询模型
- `internal/application/query/handler/landlord_query_handler.go` - 更新房东查询处理器
- `internal/application/query/handler/lease_query_handler.go` - 更新租约查询处理器
- `internal/application/query/handler/bill_query_handler.go` - 更新账单查询处理器
- `internal/application/query/handler/location_query_handler.go` - 更新位置查询处理器
- `internal/application/query/handler/room_query_handler.go` - 更新房间查询处理器
- `internal/application/query/handler/print_query_handler.go` - 更新打印查询处理器
- `internal/domain/repository/landlord.go` - 更新房东仓储接口
- `internal/domain/repository/lease.go` - 更新租约仓储接口
- `internal/domain/repository/bill.go` - 更新账单仓储接口
- `internal/domain/repository/location.go` - 更新位置仓储接口
- `internal/domain/repository/room.go` - 更新房间仓储接口
- `internal/domain/repository/print.go` - 更新打印仓储接口
- `internal/infrastructure/persistence/sqlite/landlord_repo.go` - 更新房东仓储实现
- `internal/infrastructure/persistence/sqlite/lease_repo.go` - 更新租约仓储实现
- `internal/infrastructure/persistence/sqlite/bill_repo.go` - 更新账单仓储实现
- `internal/infrastructure/persistence/sqlite/location_repo.go` - 更新位置仓储实现
- `internal/infrastructure/persistence/sqlite/room_repo.go` - 更新房间仓储实现
- `internal/infrastructure/persistence/sqlite/print_repo.go` - 更新打印仓储实现
- `internal/facade/cqrs_landlord_handler.go` - 更新房东 HTTP 处理器
- `internal/facade/cqrs_lease_handler.go` - 更新租约 HTTP 处理器
- `internal/facade/cqrs_bill_handler.go` - 更新账单 HTTP 处理器
- `internal/facade/cqrs_location_handler.go` - 更新位置 HTTP 处理器
- `internal/facade/cqrs_room_handler.go` - 更新房间 HTTP 处理器
- `internal/facade/cqrs_print_handler.go` - 更新打印 HTTP 处理器

### 前端
- `web/src/api/landlord.ts` - 更新房东 API 客户端
- `web/src/api/lease.ts` - 更新租约 API 客户端
- `web/src/api/bill.ts` - 更新账单 API 客户端
- `web/src/api/location.ts` - 更新位置 API 客户端
- `web/src/api/room.ts` - 更新房间 API 客户端
- `web/src/api/print.ts` - 更新打印 API 客户端
- `web/src/pages/Landlords.tsx` - 更新房东页面，添加查询功能
- `web/src/pages/Leases.tsx` - 更新租约页面，添加查询功能
- `web/src/pages/Bills.tsx` - 更新账单页面，添加查询功能
- `web/src/pages/Locations.tsx` - 更新位置页面，添加查询功能
- `web/src/pages/Rooms.tsx` - 更新房间页面，添加查询功能
- `web/src/pages/Print.tsx` - 更新打印页面，添加查询功能

**架构模式：**
- 采用与操作日志查询相同的 CQRS 查询模式
- 支持后端分页，使用 OFFSET 和 LIMIT 实现
- 查询条件支持模糊搜索（LIKE '%value%'）和时间范围查询（BETWEEN）

**技术栈：**
- 后端：Go, SQLite, Uber Fx
- 前端：React + TypeScript, Ant Design
- 架构：DDD + CQRS + 事件驱动
