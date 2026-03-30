## Why

当前系统中租约到期需要手动处理，缺少自动过期机制。需要定时检查 active 状态的租约是否已到期，自动将到期租约状态更新为 expired，并通过领域事件触发房间状态更新。

## What Changes

- 新增 robfig/cron 依赖
- 在 LeaseService 中添加 CheckAndExpireLeases 方法
- 创建 LeaseExpirationScheduler 定时调度器
- 在应用启动时自动启动调度器
- 不修改现有业务代码，通过事件驱动方式更新房间状态

## Capabilities

### New Capabilities
- `lease-expiration-scheduler`: 定时检查租约到期时间，自动将到期租约状态更新为 expired

### Modified Capabilities

## Impact

- 新增依赖：github.com/robfig/cron/v3
- 新增文件：internal/application/lease/scheduler.go
- 修改文件：internal/domain/lease/service/service.go
- 修改文件：internal/application/lease/module.go
- 修改文件：cmd/api/main.go
