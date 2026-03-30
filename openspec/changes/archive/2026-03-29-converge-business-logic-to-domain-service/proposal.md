## Why

当前架构中，CommandHandler 既负责流程编排，又包含业务校验逻辑，职责不够清晰。需要按照 DDD 原则重构，将业务规则和校验逻辑收敛到领域服务中。

## What Changes

- 重构 LeaseService，添加多个 repository 依赖
- 将创建租约逻辑（含押金创建）移到 LeaseService
- 将删除租约校验逻辑（账单、押金检查）移到 LeaseService
- 将激活租约校验逻辑（状态、日期、房间检查）移到 LeaseService
- 简化 CommandHandler，仅保留流程编排职责

## Capabilities

### New Capabilities

### Modified Capabilities

## Impact

- 影响代码：`internal/domain/lease/service/`、`internal/application/lease/`
- 无 API 变更
- 无外部依赖变更
