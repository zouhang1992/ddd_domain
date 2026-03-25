## Context

本项目是一个全新的 Go HTTP 项目，采用领域驱动设计 (DDD) 架构模式。项目需要建立完整的分层架构和基础设施支持，包括命令总线、事件总线和 Saga 模式。使用 SQLite 作为数据库。

## Goals / Non-Goals

**Goals:**
- 建立符合 DDD 的项目目录结构
- 实现命令总线 (Command Bus) 支持命令分发
- 实现事件总线 (Event Bus) 支持领域事件
- 实现 Saga 模式处理分布式事务一致性
- 使用 SQLite 作为数据库，提供 repository 实现
- 提供清晰的分层：门面层、应用服务层、领域层、基础设施层

**Non-Goals:**
- 不实现复杂的业务领域逻辑（仅提供基础示例）
- 不实现分布式部署支持
- 不实现外部服务集成

## Decisions

### 1. 项目目录结构
**决策:** 采用标准 DDD 分层目录结构
```
ddd_domain/
├── cmd/                    # 应用入口
│   └── api/
│       └── main.go
├── internal/
│   ├── facade/             # 门面层
│   ├── application/        # 应用服务层
│   │   ├── service/
│   │   ├── command/
│   │   └── query/
│   ├── domain/             # 领域层
│   │   ├── model/
│   │   ├── service/
│   │   ├── event/
│   │   └── repository/
│   └── infrastructure/     # 基础设施层
│       ├── bus/
│       │   ├── command/
│       │   └── event/
│       ├── saga/
│       └── persistence/
│           └── sqlite/
├── go.mod
└── go.sum
```

### 2. 数据库选择
**决策:** 使用 SQLite 作为数据库
- 使用 `modernc.org/sqlite` 纯 Go 驱动（无 CGO）
- 数据库文件存储在项目根目录 `data/` 文件夹
- 提供基础的 repository 实现和事务支持

### 3. 命令总线实现
**决策:** 使用内存同步命令总线，支持中间件扩展
- 支持命令注册和分发
- 支持命令验证中间件
- 支持事务中间件
- 接口定义清晰，易于替换为分布式实现

### 4. 事件总线实现
**决策:** 使用内存事件总线，支持同步和异步处理
- 支持事件订阅
- 支持事件发布
- 支持多个订阅者
- 易于扩展为消息队列实现 (Kafka, RabbitMQ)

### 5. Saga 模式
**决策:** 实现 Saga 编排模式 (Choreography + Orchestration 混合)
- 提供 Saga 定义接口
- 支持补偿事务
- 支持 Saga 状态持久化到 SQLite
- 提供本地事务管理

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| 内存总线在多实例部署下无法共享 | 接口设计允许未来替换为消息队列 |
| Saga 实现复杂度较高 | 先提供基础编排能力，逐步完善 |
| SQLite 不适合高并发写入 | 本项目定位简单应用，可接受此限制 |
