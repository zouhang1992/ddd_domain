## 1. 数据库设计与迁移

- [x] 1.1 在 `internal/infrastructure/persistence/sqlite/` 目录下创建操作日志表的迁移脚本
- [x] 1.2 添加 `operation_logs` 表的 SQLite 建表语句（包含所有字段和索引）

## 2. 领域模型与存储

- [x] 2.1 在 `internal/domain/model/` 目录下创建 `OperationLog.go`，定义操作日志聚合根
- [x] 2.2 在 `internal/domain/repository/` 目录下创建 `operation_log.go`，定义存储仓库接口
- [x] 2.3 在 `internal/infrastructure/persistence/sqlite/` 目录下创建 `operation_log_repo.go`，实现操作日志存储仓库

## 3. 事件处理器

- [x] 3.1 在 `internal/infrastructure/bus/event/` 目录下创建事件数据提取工具函数
- [x] 3.2 在 `internal/application/event/handler/` 目录下创建 `operation_log_handler.go`，实现 `EventHandler` 接口
- [x] 3.3 实现事件到 `OperationLog` 的转换逻辑
- [x] 3.4 在事件总线初始化时注册 `OperationLogHandler`

## 4. 查询API实现

- [x] 4.1 在 `internal/application/query/` 目录下创建查询操作日志的查询结构和处理逻辑
- [x] 4.2 在 `internal/facade/handler/` 目录下创建HTTP API处理器
- [x] 4.3 在 `api/` 目录下创建路由和请求/响应结构

## 5. 前端查询界面

- [x] 5.1 在 `web/src/pages/` 目录下创建 `OperationLogs.tsx` 页面组件
- [x] 5.2 实现查询表单组件，支持领域类型、时间范围、聚合ID等条件
- [x] 5.3 实现日志列表组件，支持分页显示
- [x] 5.4 实现日志详情弹窗组件，显示完整事件信息
- [x] 5.5 在 `web/src/components/Layout.tsx` 中添加导航菜单入口

## 6. 测试与验证

- [x] 6.1 为 OperationLog 领域模型创建单元测试
- [x] 6.2 为 OperationLogHandler 事件处理器创建集成测试
- [x] 6.3 为查询API创建端到端测试
- [x] 6.4 手动测试：启动应用程序，执行操作，验证日志记录和查询功能
