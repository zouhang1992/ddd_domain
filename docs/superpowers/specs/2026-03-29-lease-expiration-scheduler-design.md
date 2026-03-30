---
name: Lease Expiration Scheduler Design
description: 新增定时任务检查租约到期时间，当租约到期后调整租约状态
type: design
---

# 租约到期定时检查设计

## 概述

新增定时任务定期检查租约到期时间，当租约到期后自动调整租约状态为 expired，并通过领域事件驱动房间状态更新。

## 背景

当前系统中租约到期需要手动处理，缺少自动过期机制。需要：
1. 定时检查 active 状态的租约是否已到期
2. 自动将到期租约状态更新为 expired
3. 通过领域事件触发房间状态更新为 available
4. 不影响现有业务流程

## 设计决策

### 1. 架构风格：集成在 LeaseService + Application 层

**决策：** 采用方案 A，在 LeaseService 中添加业务逻辑，在 Application 层实现调度器。

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

## 数据模型

无需新增数据模型，复用现有的 Lease 模型。

## 组件设计

### 1. LeaseService 扩展

**文件位置：** `internal/domain/lease/service/service.go`

**新增方法：**
```go
// CheckAndExpireLeases 检查并处理到期租约
// 返回：处理的租约数量，错误
func (s *LeaseService) CheckAndExpireLeases() (int, error)
```

**职责：**
- 查找所有状态为 active 且已到期的租约
- 对每个租约调用 lease.Expire()
- 保存租约
- 返回处理的租约数量

### 2. LeaseExpirationScheduler

**文件位置：** `internal/application/lease/scheduler.go`

```go
type LeaseExpirationScheduler struct {
    leaseService *leaseservice.LeaseService
    eventBus     *event.Bus
    logger       *zap.Logger
    cron         *cron.Cron
    running      bool
}

func NewLeaseExpirationScheduler(...) *LeaseExpirationScheduler
func (s *LeaseExpirationScheduler) Start() error
func (s *LeaseExpirationScheduler) Stop()
```

**职责：**
- 启动/停止定时任务
- 启动时立即执行一次检查
- 每小时执行一次检查
- 记录执行日志

### 3. Lease module.go 更新

**文件位置：** `internal/application/lease/module.go`

新增提供：
- `fx.Provide(NewLeaseExpirationScheduler)`

### 4. main.go 集成

在 main.go 中：
- 添加 `LeaseExpirationScheduler` 依赖
- 通过 `fx.Invoke(startLeaseExpirationScheduler)` 启动调度器

## 数据流程

```
应用启动
  ↓
fx.Invoke(startLeaseExpirationScheduler)
  ↓
LeaseExpirationScheduler.Start()
  ├─→ 立即执行一次 CheckAndExpireLeases()
  └─→ 设置 cron 定时任务（"0 * * * *" - 每小时整点执行）
  ↓
定时触发 / 立即执行
  ↓
LeaseService.CheckAndExpireLeases()
  ├─→ leaseRepo.FindActiveLeasesExpiringBefore(now)
  ├─→ 遍历到期租约:
  │     ├─→ lease.Expire() [记录 lease.expired 事件]
  │     └─→ leaseRepo.Save(lease)
  └─→ 返回处理数量
  ↓
发布领域事件（异步）
  ↓
OperationLogEventHandler 记录日志
  ↓
LeaseRoomEventHandler 更新房间状态为 available
```

## 错误处理

- **调度器启动失败**：记录错误日志，不影响主应用启动
- **租约查询失败**：记录错误，跳过本次执行，等待下次调度
- **单个租约过期失败**：记录错误，继续处理下一个租约
- **幂等性保证**：`lease.Expire()` 方法检查当前状态，避免重复处理

## 实现任务清单

- [ ] 添加 robfig/cron 依赖到 go.mod
- [ ] 在 LeaseService 中添加 CheckAndExpireLeases 方法
- [ ] 创建 LeaseExpirationScheduler
- [ ] 更新 lease/module.go 提供调度器
- [ ] 在 main.go 中集成并启动调度器
- [ ] 编译验证
