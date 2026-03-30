## Why

当前系统已有完整的领域事件机制，但缺少操作日志审计能力。需要记录所有聚合的操作历史，支持查询操作日志，通过事件监听方式实现，不修改现有业务代码。

## What Changes

- 新增 OperationLog 领域模型和仓储接口
- 新增通用的 OperationLogEventHandler，订阅所有领域事件
- 新增 SQLite operation_logs 表持久化
- 不修改现有业务代码，通过事件监听方式实现

## Capabilities

### New Capabilities
- `operation-logging`: 为所有领域聚合的操作添加操作日志记录功能，通过监听领域事件的方式实现

### Modified Capabilities

## Impact

- 新增文件：internal/domain/operationlog/ 相关文件
- 新增文件：internal/application/operationlog/ 相关文件
- 新增文件：internal/infrastructure/persistence/sqlite/operation_log_repo.go
- 修改文件：cmd/api/main.go（注册事件处理器）
- 新增数据库表：operation_logs
