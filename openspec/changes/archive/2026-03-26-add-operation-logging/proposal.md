---
name: add-operation-logging
change: add-operation-logging
description: 补充系统操作日志功能，确保所有业务操作都有完整的操作记录
labels: ["feature", "logging"]
---

## Why

当前系统已具备基础的操作日志功能，但需要补充以下关键点：
1. 操作日志在前端页面的集成展示
2. 打印操作日志的补充
3. 操作日志与各业务模块的关联优化
4. 确保所有新增业务功能都有对应的操作日志记录

操作日志是系统审计、问题排查和用户行为分析的重要依据，完善这一功能可以提升系统的可维护性和用户体验。

## What Changes

1. **前端集成**：为所有列表页面添加操作日志查看功能
2. **打印操作日志**：补充打印任务相关的操作日志记录
3. **操作日志优化**：完善操作日志的详情展示和筛选功能
4. **业务关联**：确保所有新增的查询、筛选、分页功能都有对应的操作记录
5. **API完善**：优化操作日志的查询和展示API

## Capabilities

### New Capabilities

- `operation-logging-integration`: 操作日志与业务模块的集成
- `print-job-logging`: 打印任务操作日志
- `frontend-logging-ui`: 前端操作日志展示界面

### Modified Capabilities

- `landlord-management`: 新增操作日志查看功能
- `lease-management`: 新增操作日志查看功能
- `bill-management`: 新增操作日志查看功能
- `location-management`: 新增操作日志查看功能
- `room-management`: 新增操作日志查看功能
- `print-management`: 新增操作日志记录和查看功能

## Impact

**后端代码受影响：**
- `/internal/application/query/print.go` - 补充打印操作查询
- `/internal/domain/repository/print.go` - 补充打印仓库接口
- `/internal/infrastructure/persistence/sqlite/print_repo.go` - 补充SQLite实现
- `/internal/application/event/handler/operation_log_handler.go` - 补充打印事件处理

**前端代码受影响：**
- `/web/src/pages/*.tsx` - 所有列表页面新增操作日志查看功能
- `/web/src/components/OperationLogModal.tsx` - 新增操作日志模态框组件
- `/web/src/api/print.ts` - 补充打印操作日志API

**数据库：**
- `operation_logs` 表将新增打印操作相关的记录
