---
name: add-operation-logging
change: add-operation-logging
---

## Context

### 当前状态

系统已实现基于领域事件的操作日志系统，主要组件包括：
- **领域模型**: `internal/domain/model/OperationLog.go` - 操作日志聚合根
- **事件处理器**: `internal/application/event/handler/operation_log_handler.go` - 监听领域事件并创建操作日志
- **HTTP API**: `internal/facade/cqrs_operation_log_handler.go` - 提供操作日志查询接口
- **存储**: `internal/infrastructure/persistence/sqlite/operation_log_repo.go` - SQLite存储实现
- **查询模型**: `internal/application/query/operation_log.go` - 支持分页和筛选查询

### 现有限制

尽管系统已有操作日志功能，但仍存在以下需要完善的地方：

1. **前端集成不足**：操作日志仅在独立页面展示，未与各业务模块深度集成
2. **打印操作日志缺失**：打印任务操作尚未完整记录在操作日志中
3. **操作日志可见性差**：用户在执行业务操作时无法直接查看相关操作历史
4. **查询功能优化**：需要进一步优化操作日志的查询和展示体验

## Goals / Non-Goals

### Goals

1. **前端集成**：为所有列表页面添加操作日志查看功能，使用模态框展示
2. **打印操作日志**：补充打印任务相关的操作日志记录
3. **用户体验优化**：提升操作日志的可见性和易用性
4. **API优化**：确保操作日志查询API与现有系统风格一致

### Non-Goals

1. **重新设计操作日志架构**：当前架构已成熟，无需重写
2. **操作日志分析功能**：本次不包含日志分析、统计或报表功能
3. **多租户支持**：当前系统为单租户设计，本次不涉及多租户改造

## Decisions

### Decision 1: 操作日志展示方式

**方案**：为每个列表页面添加"操作日志"按钮，点击打开模态框展示该实体的操作历史

**理由**：
- 保持现有页面布局不变，不破坏用户习惯
- 提供快速访问入口，提升操作日志的可见性
- 模态框形式不占用额外页面空间

**技术实现**：使用Ant Design的Modal组件，内置Table组件展示操作日志

### Decision 2: 打印操作日志实现

**方案**：补充打印事件处理，确保所有打印操作都有对应的操作日志记录

**理由**：
- 打印操作是重要的业务场景，需要完整记录
- 利用现有的事件驱动架构，无需新增基础设施

**实现细节**：
- 在 `internal/application/event/handler/operation_log_handler.go` 中添加打印事件处理
- 支持的打印事件：BillPrinted、LeasePrinted、InvoicePrinted、PrintJobFailed
- 为打印任务添加操作日志记录

### Decision 3: 操作日志组件设计

**方案**：创建通用的OperationLogModal组件，支持在所有页面复用

**理由**：
- 避免代码重复
- 统一操作日志的展示风格
- 简化维护和扩展

**组件接口**：
```typescript
interface OperationLogModalProps {
  visible: boolean;
  onCancel: () => void;
  domainType: string;
  aggregateID?: string;
}
```

### Decision 4: 查询参数优化

**方案**：为操作日志查询添加 `domainType` 和 `aggregateID` 筛选条件

**理由**：
- 实现按实体类型和具体实体筛选操作日志
- 提升查询效率，避免返回过多无关记录

**API变更**：
```go
type ListOperationLogsQuery struct {
    DomainType  string     // 领域类型（如 landlord, lease, bill 等）
    AggregateID string     // 聚合根ID（实体ID）
    // 现有参数...
    Offset      int
    Limit       int
}
```

## Risks / Trade-offs

### 风险1：性能影响

**风险**：操作日志查询可能对数据库性能造成影响

**缓解**：
- 确保查询使用适当的索引
- 限制单次查询的返回条数（默认20条）
- 优化SQL查询，避免全表扫描

### 风险2：前端代码重复

**风险**：在多个页面添加操作日志功能可能导致代码重复

**缓解**：
- 创建通用的OperationLogModal组件
- 使用高阶组件或Hook封装查询逻辑
- 统一API调用方式

### 风险3：操作日志完整性

**风险**：打印操作日志可能存在遗漏

**缓解**：
- 全面梳理所有打印相关事件
- 添加单元测试确保事件处理正确性
- 在开发和测试阶段充分验证

## Open Questions

1. **操作员ID获取**：当前系统未实现用户认证，操作员ID为空。需要确定如何获取和展示操作员信息。
2. **操作详情优化**：操作日志详情目前为JSON格式，是否需要优化为更友好的展示方式？

## Migration Plan

### 部署步骤

1. 部署后端代码变更，包括：
   - 打印事件处理逻辑
   - 操作日志查询API优化

2. 部署前端代码变更，包括：
   - OperationLogModal组件
   - 各页面操作日志按钮和模态框集成

3. 数据库迁移：无需迁移，使用现有 `operation_logs` 表结构

### 回滚策略

1. 如有问题，可回退到之前的代码版本
2. 操作日志数据不会受影响，因为仅添加查询和展示功能
