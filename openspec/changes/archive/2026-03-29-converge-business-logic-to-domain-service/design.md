## Context

当前架构中，CommandHandler 既负责流程编排，又包含业务校验逻辑，职责不够清晰。LeaseService 目前只有一个简单的 ValidateRoomForLease 方法。

## Goals / Non-Goals

**Goals:**
- 将所有业务校验和规则逻辑收敛到 LeaseService
- CommandHandler 仅负责流程编排
- 保持功能不变，仅重构代码组织

**Non-Goals:**
- 不改变 API 契约
- 不修改外部行为
- 不添加新功能

## Decisions

### 1. 架构风格："胖"领域服务

**Decision:** 采用"胖"领域服务模式，LeaseService 可以直接依赖多个 repository。

**Rationale:**
- 业务逻辑集中管理，易于测试和理解
- 与当前代码库风格一致
- 避免过度分层带来的复杂性

**Alternatives considered:**
- "瘦"领域服务 - 领域服务不依赖 repository，应用层负责加载数据（职责分散）
- 按用例拆分多个服务 - 服务数量过多，过度设计

### 2. 职责划分

**CommandHandler（应用层）职责：**
- 命令类型转换
- 命令基础验证（Command.Validate()）
- 调用领域服务
- 保存聚合根
- 发布领域事件

**LeaseService（领域层）职责：**
- 房间可租性校验
- 删除前置条件校验（账单、押金检查）
- 激活前置条件校验（状态、日期检查）
- 租约创建（含押金创建）协调

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| 领域服务依赖多个 repository | 这是有意为之的设计，便于业务逻辑集中 |
| 重构范围较大 | 分步骤实现，每步编译验证 |

## Migration Plan

增量重构，无需特殊迁移策略。
