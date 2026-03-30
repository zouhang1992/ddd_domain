## Context

当前系统在创建租约时没有校验房间是否存在，也没有校验房间是否处于可租状态。这可能导致为不存在的房间创建租约，或为已出租/维修中的房间创建租约。

当前架构：
- 应用层按限界上下文组织（lease、room、landlord 等）
- 使用 CQRS 模式分离命令和查询
- 使用 Fx 进行依赖注入
- 房间状态：available（可租）、rented（已租）、maintain（维修中）

## Goals / Non-Goals

**Goals:**
- 在租约创建时校验房间是否存在
- 在租约创建时校验房间是否处于可租状态
- 在租约更新时校验（如变更 roomID）
- 在租约激活时校验房间状态
- 使用专门的错误类型区分不同校验失败场景

**Non-Goals:**
- 不修改房间状态管理逻辑
- 不修改 API 契约
- 不添加新的外部依赖

## Decisions

### 1. 领域服务位置：LeaseService（非 RoomValidationService）
**Decision:** 在 lease 限界上下文中创建 LeaseService
**Rationale:** 房间校验是租约业务逻辑的一部分，应该作为租约领域服务的职责。
**Alternatives:**
- RoomValidationService 在 room 限界上下文 - 会增加跨限界上下文依赖
- 直接在 CommandHandler 中校验 - 代码重复，难以复用

### 2. 错误类型位置：新增专门错误
**Decision:** 在 internal/domain/common/errors/ 下新增房间相关错误
**Rationale:** 与现有错误组织方式一致，便于复用
**Alternatives:**
- 在 lease 包中定义错误 - 其他上下文无法复用

### 3. 校验范围：创建、更新、激活
**Decision:** 在创建、更新（变更 roomID 时、激活时进行校验
**Rationale:** 确保所有可能改变租约与房间关联的操作都经过校验
**Alternatives:**
- 仅创建时校验 - 不够安全

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| 更新租约时如果 roomID 不变，不需要校验 | 在校验逻辑中检查 roomID 是否变更 |
| 房间状态在校验后和租约创建前可能变化 | 这是并发问题，当前 scope 外，后续可考虑分布式锁 |

## Migration Plan

无需迁移，这是新增功能，向后兼容。

## Open Questions

无
