## Why

原 house 系统（收账系统）是一个功能完整的租房管理系统，但代码全部在单个 main.go 文件中，缺乏模块化设计。当前 ddd_domain 项目已经建立了良好的领域驱动设计架构，并且已有位置和房间管理的基础能力。将 house 系统的业务逻辑迁移到 ddd_domain，可以利用现有的 DDD 架构提升代码可维护性，同时使用 React+TS 重写前端以提升用户体验。

## What Changes

- **新增**：房东管理（Landlord）领域模型和业务逻辑
- **新增**：租约管理（Lease）领域模型和业务逻辑，包括创建、查询、修改、续租、退租结算
- **新增**：账单管理（Bill）领域模型和业务逻辑，包括收账记录、收据打印
- **新增**：押金管理（Deposit）领域模型和业务逻辑，包括收取、退还、扣除
- **新增**：用户认证（Auth）功能，支持登录/登出
- **新增**：收入汇总查询功能
- **新增**：React+TypeScript 前端应用，替换原有 HTML 模板
- **迁移**：将 house 系统的数据库 schema 适配到当前 SQLite 持久化层

## Capabilities

### New Capabilities
- `landlord-management`: 房东信息的创建、查询、列表管理
- `lease-management`: 租约的创建、查询、修改、续租、退租结算等完整生命周期管理
- `bill-management`: 账单记录的创建、查询、修改、删除，以及收据打印
- `deposit-management`: 押金的收取、退还、扣除管理，续租时自动绑定
- `auth-management`: 用户登录、登出、鉴权功能
- `income-reporting`: 按月收入汇总查询
- `react-frontend`: React+TypeScript 单页应用，提供完整的用户界面

### Modified Capabilities
- `location-management`: 扩展位置管理，与房东关联
- `room-management`: 扩展房间管理，支持租约状态、押金状态展示

## Impact

- **后端代码**：在 `internal/domain/model/` 新增多个领域模型，在 `internal/domain/repository/` 新增仓储接口，在 `internal/application/service/` 新增应用服务，在 `internal/facade/` 新增 HTTP 处理器
- **API 接口**：新增大量 REST API 端点，包括 `/landlords`、`/leases`、`/bills`、`/login`、`/income` 等
- **数据库**：扩展 SQLite schema，新增多个表（landlords、leases、bills、deposits 等）
- **前端**：新增 `web/` 目录，包含完整的 React+TS 应用
- **依赖**：后端可能需要新增 JWT 相关依赖，前端需要 React、TypeScript、Vite 等依赖
