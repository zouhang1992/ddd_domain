## Context

当前系统已有完整的租约领域模型和事件机制。租约有 `active` 状态，`end_date` 字段表示到期时间。当租约到期后，需要手动处理，缺少自动过期机制。

**现有架构：**
- Lease 聚合根已有 `Expire()` 方法，会发布 `lease.expired` 事件
- LeaseRoomEventHandler 已订阅 `lease.expired` 事件，会将房间状态更新为 `available`
- LeaseRepository 已有 `FindActiveLeasesExpiringBefore()` 方法

## Goals / Non-Goals

**Goals:**
- 定时检查 active 状态的租约是否已到期
- 自动将到期租约状态更新为 expired
- 通过领域事件触发房间状态更新为 available
- 应用启动后立即执行一次检查
- 每小时执行一次检查

**Non-Goals:**
- 手动触发检查（后续可添加 API）
- 租约到期通知（后续可扩展）
- 批量操作的 UI

## Decisions

### 1. 架构风格：集成在 LeaseService + Application 层

**决策：** 在 LeaseService 中添加业务逻辑，在 Application 层实现调度器。

**理由：**
- 业务逻辑在领域服务中，符合 DDD 原则
- 调度逻辑清晰，职责分明
- 与现有架构风格一致
- 易于测试

**替代方案：**
- 独立的 Scheduler 服务 - 过度设计
- 事件驱动方式 - 不如 cron 灵活

### 2. 调度库：robfig/cron

**决策：** 使用 `github.com/robfig/cron/v3` 库实现定时任务。

**理由：**
- 成熟稳定，广泛使用
- 支持标准 cron 表达式
- 支持启动时立即执行
- 与 Go 生态集成良好

### 3. 检查频率：每小时

**决策：** 每小时执行一次检查（cron 表达式：`0 * * * *`）。

**理由：**
- 平衡实时性和资源消耗
- 对大多数场景足够
- 可通过修改 cron 表达式调整

### 4. 启动行为：立即执行一次

**决策：** 应用启动后立即执行一次检查，然后按计划执行。

**理由：**
- 及时处理已到期的租约
- 避免启动后等待第一个调度点
- 确保系统启动后数据一致性

## Risks / Trade-offs

**[风险] 调度器启动失败** → 记录错误日志，不影响主应用启动

**[风险] 租约查询失败** → 记录错误，跳过本次执行，等待下次调度

**[风险] 单个租约过期失败** → 记录错误，继续处理下一个租约

**[风险] 重复处理** → `lease.Expire()` 方法检查当前状态，避免重复处理
