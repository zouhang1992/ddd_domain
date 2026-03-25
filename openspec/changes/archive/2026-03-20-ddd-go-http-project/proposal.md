## Why

创建一个符合领域驱动设计 (DDD) 的 Go HTTP 项目基础架构，提供完整的分层架构和基础设施支持。

## What Changes

- 初始化 Go 模块和项目目录结构
- 实现 DDD 分层：门面层、应用服务层、领域层、基础设施层
- 实现事件总线 (Event Bus)
- 实现命令总线 (Command Bus)
- 实现 Saga 模式支持分布式事务一致性

## Capabilities

### New Capabilities
- `ddd-architecture`: DDD 分层架构和目录结构
- `command-bus`: 命令总线实现，支持命令分发和处理
- `event-bus`: 事件总线实现，支持领域事件发布和订阅
- `saga-pattern`: Saga 模式实现，处理分布式事务一致性

### Modified Capabilities
(无)

## Impact

- 新创建项目的所有核心代码和目录结构
- 引入 Go 语言标准库及必要的第三方依赖
