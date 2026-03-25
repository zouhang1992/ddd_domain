## Context

当前 ddd_domain 项目已建立基于领域驱动设计（DDD）的架构，包含位置（Location）和房间（Room）管理。原 house 系统是一个单文件实现的租房管理系统，业务逻辑包括房东、租约、账单、押金管理以及收据打印。为提升代码可维护性和用户体验，需要将 house 系统的后端业务移入 ddd_domain，并使用 React+TS 重写前端。

**架构约束：**
- 保持 DDD 架构风格（domain → application → infrastructure → facade）
- 使用 SQLite 作为持久化存储
- RESTful API 接口设计
- 前端使用 React 18 + TypeScript + Vite

## Goals / Non-Goals

**Goals:**
1. 完整迁移 house 系统的后端业务逻辑到 ddd_domain
2. 保持与现有 Location 和 Room 管理模块的兼容性
3. 实现完整的用户认证和授权功能
4. 开发 React+TS 前端应用，提供现代化用户界面
5. 支持收据和合同打印功能（RTF 格式）
6. 优化代码结构，提高可测试性和可维护性

**Non-Goals:**
1. 不改变现有的 Location 和 Room 领域模型的核心结构
2. 不引入新的数据库技术（保持 SQLite）
3. 不实现复杂的权限管理（仅支持单用户登录）
4. 不实现多租户支持

## Decisions

### 1. 领域模型设计

**决定：** 遵循现有架构风格，在 `internal/domain/model/` 中创建独立的领域模型类

**理由：** 原 house 系统的核心业务实体（Landlord、Lease、Bill、Deposit）具有明确的业务含义和生命周期，适合作为独立的领域模型。

**实现：**
- `Landlord`: 房东信息（姓名、联系方式等）
- `Lease`: 租约信息（租客、房间、起止日期、状态等）
- `Bill`: 账单信息（租约关联、金额、类型、状态等）
- `Deposit`: 押金信息（租约关联、金额、状态等）

### 2. 认证与授权

**决定：** 使用 JWT（JSON Web Token）实现用户认证

**理由：** JWT 是轻量级的认证方案，适合前后端分离架构，便于 React 前端实现。

**实现：**
- 在 `internal/application/service/` 中创建 AuthService
- 在 `internal/facade/` 中创建 AuthHandler
- 密码使用 bcrypt 加密存储
- Token 过期时间设置为 7 天

### 3. 前端架构

**决定：** 使用 React 18 + TypeScript + Vite + Ant Design

**理由：**
- React 是主流的前端框架，社区支持强大
- TypeScript 提供类型安全，提高代码质量
- Vite 提供快速的开发体验
- Ant Design 提供丰富的 UI 组件

**实现：**
- 单页应用（SPA）架构
- 组件化设计，每个业务模块对应一个页面
- 使用 React Router 6 进行路由管理
- 使用 Axios 与后端 API 通信

### 4. 收据/合同打印

**决定：** 保持 RTF 格式，使用 Go 模板生成

**理由：** RTF 格式可直接在 Word 中打开，用户体验良好，且实现简单。

**实现：**
- 在 `internal/application/service/` 中创建 PrintService
- 使用 Go 文本模板生成 RTF 内容
- 提供 HTTP 接口返回 RTF 内容，浏览器可直接下载

## Risks / Trade-offs

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 数据迁移风险 | 原 house 系统的数据需要迁移到新系统 | 编写数据迁移脚本，支持从原数据库导入数据 |
| 前端性能风险 | React 应用可能在低性能设备上运行缓慢 | 优化组件渲染，使用虚拟列表，代码分割 |
| API 兼容性风险 | 原 house 系统的 API 接口与新系统不同 | 提供旧 API 路由的兼容层 |
| 测试覆盖风险 | 新代码可能缺乏足够测试 | 为每个领域模型编写单元测试，为应用服务编写集成测试 |
