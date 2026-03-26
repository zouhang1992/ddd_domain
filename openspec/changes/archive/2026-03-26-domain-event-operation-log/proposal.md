## Why

为了满足系统审计和问题追踪需求，需要对各领域操作进行完整记录和查询。当前系统缺乏统一的操作日志管理机制，无法追踪领域事件的发生和影响，影响了问题排查和系统透明度。

## What Changes

- 新增操作日志领域模型 `OperationLog`，记录完整的事件信息
- 实现事件监听器 `OperationLogHandler`，自动监听并记录所有领域事件
- 新增 SQLite 数据库表 `operation_logs` 存储操作日志
- 提供查询API支持按领域类型、时间范围、聚合ID等条件查询
- 新增前端操作日志查询页面，支持详细信息展示

## Capabilities

### New Capabilities

- **operation-log-management**: 操作日志的创建、存储和查询管理
- **event-logging**: 领域事件到操作日志的转换和记录机制

### Modified Capabilities

- **event-bus**: 添加操作日志处理器，但不修改现有事件总线接口
- **sqlite-persistence**: 添加操作日志表，但不修改现有表结构

## Impact

- **Infrastructure**: 新增操作日志表 `operation_logs` 和相应的查询API
- **Domain**: 新增 `OperationLog` 聚合根和存储仓库
- **API**: 新增查询接口 `/api/operation-logs`，支持分页和条件查询
- **Frontend**: 新增 `/operation-logs` 页面，提供查询和详情展示功能
- **Dependencies**: 无新增外部依赖，使用项目现有技术栈

## Risks and Mitigations

- **性能影响**: 操作日志记录可能影响事件处理性能 → 采用异步事件处理方式
- **存储增长**: 操作日志可能快速增长 → 考虑添加日志清理策略（可作为后续改进）
