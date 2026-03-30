## Why

当前系统在创建租约时没有校验房间是否存在，也没有校验房间是否处于可租状态。这可能导致为不存在的房间创建租约，或为已出租/维修中的房间创建租约。需要在租约相关操作中添加房间校验逻辑。

## What Changes

- 新增房间相关错误类型（房间不存在、房间不可租）
- 在 lease 限界上下文中创建 LeaseService 租约领域服务
- 在 LeaseService 中实现房间校验逻辑
- 在租约创建、更新、激活操作中添加房间校验
- 更新依赖注入配置

## Capabilities

### New Capabilities

### Modified Capabilities
- `lease-management`: 添加房间存在性和状态校验需求

## Impact

- 影响代码：`internal/application/lease/`、`internal/domain/common/errors/`
- 无 API 变更
- 无外部依赖变更
