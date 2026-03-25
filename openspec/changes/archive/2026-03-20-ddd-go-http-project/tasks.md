## 1. 项目初始化

- [x] 1.1 初始化 Go 模块 (go mod init)
- [x] 1.2 创建完整的 DDD 目录结构
- [x] 1.3 添加 SQLite 驱动依赖 (modernc.org/sqlite)

## 2. 基础设施层 - 总线

- [x] 2.1 定义 Command 和 CommandHandler 接口
- [x] 2.2 实现内存命令总线，支持命令分发
- [x] 2.3 实现命令中间件支持
- [x] 2.4 定义 DomainEvent 和 EventHandler 接口
- [x] 2.5 实现内存事件总线，支持事件发布和订阅

## 3. 基础设施层 - Saga 模式

- [x] 3.1 定义 Saga 接口和步骤定义
- [x] 3.2 实现 Saga 编排器
- [x] 3.3 实现 Saga 状态管理
- [x] 3.4 实现 Saga 状态持久化到 SQLite

## 4. 基础设施层 - 持久化

- [x] 4.1 实现 SQLite 连接管理
- [x] 4.2 定义基础 Repository 接口
- [x] 4.3 实现事务支持

## 5. 应用服务层

- [x] 5.1 创建应用服务目录结构
- [x] 5.2 定义基础 Command 和 Query 结构

## 6. 领域层

- [x] 6.1 创建领域层目录结构
- [x] 6.2 定义领域实体和值对象基础
- [x] 6.3 定义 Repository 接口

## 7. 门面层和入口

- [x] 7.1 创建门面层目录结构
- [x] 7.2 实现 HTTP 服务入口 (main.go)
- [x] 7.3 配置 HTTP 服务器基础路由
