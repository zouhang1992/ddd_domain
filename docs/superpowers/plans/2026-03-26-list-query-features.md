# 各列表页面查询功能实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为所有列表页面（房东、租约、账单、位置、房间、打印）添加查询功能，支持根据查询条件筛选和分页

**Architecture:** 采用与操作日志查询相同的架构模式，使用 CQRS 查询模式，支持后端分页

**Tech Stack:**
- 后端：Go, SQLite, Uber Fx
- 前端：React + TypeScript, Ant Design
- 架构：DDD + CQRS + 事件驱动

---

## 文件结构分析

### 需要创建/修改的文件：

#### 后端查询模型（internal/application/query）
- `internal/application/query/landlord.go` - 更新房东查询结构
- `internal/application/query/lease.go` - 更新租约查询结构
- `internal/application/query/bill.go` - 更新账单查询结构
- `internal/application/query/location.go` - 更新位置查询结构
- `internal/application/query/room.go` - 更新房间查询结构
- `internal/application/query/print.go` - 更新打印查询结构

#### 后端查询处理器（internal/application/query/handler）
- `internal/application/query/handler/landlord_query_handler.go` - 更新房东查询处理器
- `internal/application/query/handler/lease_query_handler.go` - 更新租约查询处理器
- `internal/application/query/handler/bill_query_handler.go` - 更新账单查询处理器
- `internal/application/query/handler/location_query_handler.go` - 更新位置查询处理器
- `internal/application/query/handler/room_query_handler.go` - 更新房间查询处理器
- `internal/application/query/handler/print_query_handler.go` - 更新打印查询处理器

#### 后端仓储接口（internal/domain/repository）
- `internal/domain/repository/landlord.go` - 更新房东仓储接口
- `internal/domain/repository/lease.go` - 更新租约仓储接口
- `internal/domain/repository/bill.go` - 更新账单仓储接口
- `internal/domain/repository/location.go` - 更新位置仓储接口
- `internal/domain/repository/room.go` - 更新房间仓储接口
- `internal/domain/repository/print.go` - 更新打印仓储接口

#### 后端 SQLite 实现（internal/infrastructure/persistence/sqlite）
- `internal/infrastructure/persistence/sqlite/landlord_repo.go` - 更新房东仓储实现
- `internal/infrastructure/persistence/sqlite/lease_repo.go` - 更新租约仓储实现
- `internal/infrastructure/persistence/sqlite/bill_repo.go` - 更新账单仓储实现
- `internal/infrastructure/persistence/sqlite/location_repo.go` - 更新位置仓储实现
- `internal/infrastructure/persistence/sqlite/room_repo.go` - 更新房间仓储实现
- `internal/infrastructure/persistence/sqlite/print_repo.go` - 更新打印仓储实现

#### 后端 HTTP 处理器（internal/facade）
- `internal/facade/cqrs_landlord_handler.go` - 更新房东 HTTP 处理器
- `internal/facade/cqrs_lease_handler.go` - 更新租约 HTTP 处理器
- `internal/facade/cqrs_bill_handler.go` - 更新账单 HTTP 处理器
- `internal/facade/cqrs_location_handler.go` - 更新位置 HTTP 处理器
- `internal/facade/cqrs_room_handler.go` - 更新房间 HTTP 处理器
- `internal/facade/cqrs_print_handler.go` - 更新打印 HTTP 处理器

#### 前端 API 客户端（web/src/api）
- `web/src/api/landlord.ts` - 更新房东 API 客户端
- `web/src/api/lease.ts` - 更新租约 API 客户端
- `web/src/api/bill.ts` - 更新账单 API 客户端
- `web/src/api/location.ts` - 更新位置 API 客户端
- `web/src/api/room.ts` - 更新房间 API 客户端
- `web/src/api/print.ts` - 更新打印 API 客户端

#### 前端页面组件（web/src/pages）
- `web/src/pages/Landlords.tsx` - 更新房东页面，添加查询功能
- `web/src/pages/Leases.tsx` - 更新租约页面，添加查询功能
- `web/src/pages/Bills.tsx` - 更新账单页面，添加查询功能
- `web/src/pages/Locations.tsx` - 更新位置页面，添加查询功能
- `web/src/pages/Rooms.tsx` - 更新房间页面，添加查询功能
- `web/src/pages/Print.tsx` - 更新打印页面，添加查询功能

---

## 任务分解

### 阶段 1：准备工作
- [ ] 确保已创建项目的基础查询架构（已有操作日志作为参考）
- [ ] 统一查询结果格式，支持分页信息

---

### 阶段 2：房东管理页面查询功能（Landlords）

**任务 2.1：更新房东查询模型**
- [ ] 修改 `internal/application/query/landlord.go`，添加查询结构和结果类型
- [ ] 更新 `ListLandlordsQuery`，添加查询参数
- [ ] 添加 `LandlordsQueryResult` 结果类型，支持分页

**任务 2.2：更新房东查询处理器**
- [ ] 修改 `internal/application/query/handler/landlord_query_handler.go`
- [ ] 添加 `HandleListLandlords` 方法的查询逻辑
- [ ] 支持按姓名和电话模糊搜索

**任务 2.3：更新房东仓储接口**
- [ ] 修改 `internal/domain/repository/landlord.go`
- [ ] 添加查询接口方法 `FindByCriteria` 和 `CountByCriteria`

**任务 2.4：更新房东 SQLite 仓储实现**
- [ ] 修改 `internal/infrastructure/persistence/sqlite/landlord_repo.go`
- [ ] 实现查询方法，支持按姓名和电话模糊搜索
- [ ] 实现计数方法

**任务 2.5：更新房东 HTTP 处理器**
- [ ] 修改 `internal/facade/cqrs_landlord_handler.go`
- [ ] 更新 `List` 方法，支持查询参数和分页
- [ ] 解析 URL 查询参数

**任务 2.6：更新房东前端 API**
- [ ] 修改 `web/src/api/landlord.ts`
- [ ] 支持传递查询参数和分页参数

**任务 2.7：更新房东前端页面**
- [ ] 修改 `web/src/pages/Landlords.tsx`
- [ ] 添加查询表单，支持姓名和电话搜索
- [ ] 实现查询和重置功能
- [ ] 支持分页显示

---

### 阶段 3：租约管理页面查询功能（Leases）

**任务 3.1：更新租约查询模型**
- [ ] 修改 `internal/application/query/lease.go`，添加查询结构和结果类型
- [ ] 更新 `ListLeasesQuery`，添加查询参数（租户姓名、租户电话、状态、位置、房间、时间范围）

**任务 3.2：更新租约查询处理器**
- [ ] 修改 `internal/application/query/handler/lease_query_handler.go`
- [ ] 添加查询逻辑，支持多条件筛选

**任务 3.3：更新租约仓储接口**
- [ ] 修改 `internal/domain/repository/lease.go`
- [ ] 添加查询接口方法

**任务 3.4：更新租约 SQLite 仓储实现**
- [ ] 修改 `internal/infrastructure/persistence/sqlite/lease_repo.go`
- [ ] 实现查询方法，支持时间范围查询

**任务 3.5：更新租约 HTTP 处理器**
- [ ] 修改 `internal/facade/cqrs_lease_handler.go`
- [ ] 支持查询参数和分页

**任务 3.6：更新租约前端 API**
- [ ] 修改 `web/src/api/lease.ts`

**任务 3.7：更新租约前端页面**
- [ ] 修改 `web/src/pages/Leases.tsx`
- [ ] 添加查询表单，支持复杂筛选条件
- [ ] 实现查询和重置功能

---

### 阶段 4：账单管理页面查询功能（Bills）

**任务 4.1：更新账单查询模型**
- [ ] 修改 `internal/application/query/bill.go`，添加查询结构
- [ ] 更新 `ListBillsQuery`，添加查询参数（类型、状态、租约、金额范围、时间范围）

**任务 4.2：更新账单查询处理器**
- [ ] 修改 `internal/application/query/handler/bill_query_handler.go`
- [ ] 支持金额范围和时间范围查询

**任务 4.3：更新账单仓储接口**
- [ ] 修改 `internal/domain/repository/bill.go`

**任务 4.4：更新账单 SQLite 仓储实现**
- [ ] 修改 `internal/infrastructure/persistence/sqlite/bill_repo.go`

**任务 4.5：更新账单 HTTP 处理器**
- [ ] 修改 `internal/facade/cqrs_bill_handler.go`

**任务 4.6：更新账单前端 API**
- [ ] 修改 `web/src/api/bill.ts`

**任务 4.7：更新账单前端页面**
- [ ] 修改 `web/src/pages/Bills.tsx`

---

### 阶段 5：位置管理页面查询功能（Locations）

**任务 5.1：更新位置查询模型**
- [ ] 修改 `internal/application/query/location.go`
- [ ] 更新 `ListLocationsQuery`，添加查询参数（简称和详情模糊搜索）

**任务 5.2：更新位置查询处理器**
- [ ] 修改 `internal/application/query/handler/location_query_handler.go`

**任务 5.3：更新位置仓储接口**
- [ ] 修改 `internal/domain/repository/location.go`

**任务 5.4：更新位置 SQLite 仓储实现**
- [ ] 修改 `internal/infrastructure/persistence/sqlite/location_repo.go`

**任务 5.5：更新位置 HTTP 处理器**
- [ ] 修改 `internal/facade/cqrs_location_handler.go`

**任务 5.6：更新位置前端 API**
- [ ] 修改 `web/src/api/location.ts`

**任务 5.7：更新位置前端页面**
- [ ] 修改 `web/src/pages/Locations.tsx`

---

### 阶段 6：房间管理页面查询功能（Rooms）

**任务 6.1：更新房间查询模型**
- [ ] 修改 `internal/application/query/room.go`
- [ ] 更新 `ListRoomsQuery`，添加查询参数（位置、房间号、标签、时间范围）

**任务 6.2：更新房间查询处理器**
- [ ] 修改 `internal/application/query/handler/room_query_handler.go`

**任务 6.3：更新房间仓储接口**
- [ ] 修改 `internal/domain/repository/room.go`

**任务 6.4：更新房间 SQLite 仓储实现**
- [ ] 修改 `internal/infrastructure/persistence/sqlite/room_repo.go`

**任务 6.5：更新房间 HTTP 处理器**
- [ ] 修改 `internal/facade/cqrs_room_handler.go`

**任务 6.6：更新房间前端 API**
- [ ] 修改 `web/src/api/room.ts`

**任务 6.7：更新房间前端页面**
- [ ] 修改 `web/src/pages/Rooms.tsx`

---

### 阶段 7：打印管理页面查询功能（Print）

**任务 7.1：更新打印查询模型**
- [ ] 修改 `internal/application/query/print.go`
- [ ] 更新 `ListPrintJobsQuery`，添加查询参数

**任务 7.2：更新打印查询处理器**
- [ ] 修改 `internal/application/query/handler/print_query_handler.go`

**任务 7.3：更新打印仓储接口**
- [ ] 修改 `internal/domain/repository/print.go`

**任务 7.4：更新打印 SQLite 仓储实现**
- [ ] 修改 `internal/infrastructure/persistence/sqlite/print_repo.go`

**任务 7.5：更新打印 HTTP 处理器**
- [ ] 修改 `internal/facade/cqrs_print_handler.go`

**任务 7.6：更新打印前端 API**
- [ ] 修改 `web/src/api/print.ts`

**任务 7.7：更新打印前端页面**
- [ ] 修改 `web/src/pages/Print.tsx`

---

### 阶段 8：测试和验证

**任务 8.1：单元测试**
- [ ] 为所有新添加的查询功能编写单元测试
- [ ] 测试查询处理器逻辑
- [ ] 测试仓储查询实现

**任务 8.2：集成测试**
- [ ] 测试完整查询流程（从 API 到数据库）
- [ ] 测试分页功能

**任务 8.3：前端测试**
- [ ] 测试查询表单交互
- [ ] 测试查询和重置功能
- [ ] 测试分页显示

**任务 8.4：手动测试**
- [ ] 启动开发服务器
- [ ] 访问各列表页面
- [ ] 测试查询功能
- [ ] 验证分页功能

---

## 执行策略

### 开发顺序建议：
1. 先完成房东管理页面（最简单，字段最少）
2. 然后完成位置管理页面（相对简单）
3. 接下来完成租约管理页面（中等复杂）
4. 再完成账单管理页面（复杂，涉及金额范围）
5. 最后完成打印和房间管理页面

### 技术注意事项：
- 使用操作日志查询功能作为参考实现
- 保持查询结构和命名一致
- 支持模糊搜索使用 `LIKE '%value%'`
- 时间范围查询使用 `BETWEEN`
- 分页使用 `OFFSET` 和 `LIMIT`
- 错误处理与现有模式保持一致

### 前端设计原则：
- 查询表单采用内联布局
- 查询按钮使用主色调
- 重置按钮使用默认样式
- 查询条件和结果分离
- 分页控件使用 Ant Design 组件

---

## 相关文档

- 操作日志查询功能：可作为参考实现
- API 设计规范：遵循项目现有模式
- 数据库查询规范：SQLite 查询最佳实践
- React 查询组件：可重用查询表单组件
