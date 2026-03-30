## Why

后端架构调整后，前端页面没有跟着调整，导致很多前端页面展示有问题。需要修复所有前端页面与后端 API 的数据格式匹配问题。

## What Changes

- 修复所有 CQRS handler 的返回格式，确保与前端期望一致
- 修复前端 API 调用和数据展示
- 确保所有分页查询返回正确的格式
- 确保所有创建/更新操作返回正确的数据结构
- 逐个修复：landlord、lease、bill、room、location、print 页面

## Capabilities

### New Capabilities

### Modified Capabilities

## Impact

- 修改文件：internal/facade/*.go（所有 CQRS handler）
- 修改文件：web/src/api/*.ts（所有前端 API 调用）
- 修改文件：web/src/pages/*.tsx（所有前端页面）
- 不涉及后端业务逻辑变更，仅修复数据格式
