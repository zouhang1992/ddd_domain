## Context

这是租房收账系统的第一个业务模块，基于已有的 DDD 框架和 SQLite 数据库。

## Goals / Non-Goals

**Goals:**
- 实现位置、房间的 CRUD 功能
- 标签作为房间的字段存储（逗号分隔字符串）
- 实现删除位置时的关联检查
- 使用已有的 DDD 架构和基础设施

**Non-Goals:**
- 不实现租户和收账功能
- 不实现权限控制
- 不实现单独的标签管理

## Decisions

### 1. 数据模型
- **locations**: id, short_name, detail, created_at, updated_at
- **rooms**: id, location_id, room_number, tags (逗号分隔), created_at, updated_at

### 2. 删除策略
删除位置时，先查询是否有关联房间，有关联则返回错误，不允许删除。

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| 标签作为字符串不利于查询 | 按当前需求，此方案可接受 |
| 并发删除和更新可能导致数据不一致 | 使用数据库事务 |
