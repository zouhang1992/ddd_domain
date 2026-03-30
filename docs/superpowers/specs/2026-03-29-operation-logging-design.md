---
name: Operation Logging Design
description: 为聚合的每个操作添加操作日志，通过监听领域事件的方式来补充操作日志
type: design
---

# 操作日志功能设计

## 概述

为所有领域聚合的操作添加操作日志记录功能，通过监听领域事件的方式实现，不影响现有业务流程。

## 背景

当前系统已有完整的领域事件机制，但缺少操作日志审计能力。需要：
1. 记录所有聚合的操作历史
2. 支持查询操作日志
3. 通过事件监听方式实现，不修改现有业务代码

## 设计决策

### 1. 架构风格：通用 OperationLogEventHandler

**决策：** 采用方案A，创建通用的 OperationLogEventHandler，订阅所有领域事件。

**理由：**
- 通用性强，新增事件类型无需修改日志处理器
- 与现有事件总线架构完全契合
- 通过反射动态提取事件信息，保持代码简洁

**替代方案：**
- 按事件类型分别处理 - 需要为每个事件创建 handler，过于繁琐
- 事件总线中间件 - 耦合度高，难以定制

### 2. 持久化方式：SQLite

**决策：** 使用现有的 SQLite 数据库，新增 operation_logs 表。

**理由：**
- 无需引入新技术栈
- 与现有持久化架构一致
- 查询性能足够
- 事务支持完善

### 3. 记录方式：异步

**决策：** 使用事件总线的 PublishAsync 方式异步记录日志。

**理由：**
- 不阻塞主业务流程
- 降低对性能的影响
- 即使日志记录失败也不影响主要操作

## 数据模型

### OperationLog 模型

**文件位置：** `internal/domain/operationlog/model/operation_log.go`

```go
type OperationLog struct {
    ID          string
    Timestamp   time.Time
    EventName   string
    DomainType  string
    AggregateID string
    OperatorID  string
    Action      string // created, updated, deleted, activated, checked-out, etc.
    Details     map[string]interface{} // 完整事件内容
    Metadata    map[string]interface{} // 元数据（IP、User-Agent等）
    CreatedAt   time.Time
}
```

### 数据库表结构

**operation_logs 表：**

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | TEXT | 主键 | PK |
| timestamp | DATETIME | 操作发生时间 | INDEX |
| event_name | TEXT | 事件名称 | INDEX |
| domain_type | TEXT | 领域类型 | INDEX |
| aggregate_id | TEXT | 关联的聚合根ID | INDEX |
| operator_id | TEXT | 操作人ID | INDEX |
| action | TEXT | 操作类型 | |
| details | JSON | 详细数据 | |
| metadata | JSON | 元数据 | |
| created_at | DATETIME | 创建时间 | INDEX |

## 组件设计

### 1. OperationLogEventHandler

**文件位置：** `internal/application/operationlog/event_handler.go`

**职责：**
- 订阅所有领域事件
- 反射提取事件信息
- 转换为 OperationLog
- 调用仓储保存

**订阅的事件：**
- room.created, room.updated, room.deleted
- lease.created, lease.activated, lease.checkout, lease.expired, lease.renewed, lease.deleted
- landlord.created, landlord.updated, landlord.deleted
- bill.created, bill.updated, bill.deleted, bill.paid
- location.created, location.updated, location.deleted

### 2. OperationLogRepository

**文件位置：** `internal/domain/operationlog/repository/repository.go`

```go
type OperationLogRepository interface {
    Save(log *model.OperationLog) error
    FindByID(id string) (*model.OperationLog, error)
    FindByAggregateID(aggregateID string) ([]*model.OperationLog, error)
    FindByDomainType(domainType string, offset, limit int) ([]*model.OperationLog, int, error)
    FindByTimeRange(start, end time.Time, offset, limit int) ([]*model.OperationLog, int, error)
}
```

### 3. SqliteOperationLogRepository

**文件位置：** `internal/infrastructure/persistence/sqlite/operation_log_repo.go`

### 4. 查询支持

在 facade 层添加查询 API：
- GET /operation-logs - 分页查询
- GET /operation-logs/:id - 获取单个日志详情
- GET /operation-logs/aggregate/:id - 按聚合查询

## 数据流程

```
领域操作（创建/更新/删除）
  ↓
聚合根.RecordEvent(evt)
  ↓
CommandHandler 保存聚合
  ↓
发布领域事件（PublishAsync）
  ↓
OperationLogEventHandler.Handle(evt)
  ├─→ 反射提取事件信息
  ├─→ 从认证上下文获取操作人
  ├─→ 创建 OperationLog
  └─→ repo.Save(log)
  ↓
存储到 SQLite operation_logs 表
```

## 实现任务清单

- [ ] 创建 OperationLog 领域模型
- [ ] 创建 OperationLogRepository 接口
- [ ] 创建 operation_logs 表的 SQLite 实现
- [ ] 创建 OperationLogEventHandler
- [ ] 在 main.go 中注册事件处理器
- [ ] 添加查询 API（可选，后续实现）
- [ ] 编译验证
