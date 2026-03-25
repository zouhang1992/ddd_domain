## Why

租房收账系统需要房间和位置管理功能，这是系统的基础模块。

## What Changes

- 新增位置管理，包含增删查改，删除时检查关联房间
- 新增房间管理，包含增删查改，关联位置，标签作为字段存储

## Capabilities

### New Capabilities
- `location-management`: 位置管理，包含简称和详细信息
- `room-management`: 房间管理，关联位置，标签作为字段

### Modified Capabilities
(无)

## Impact

- 新增领域模型（Location、Room）
- 新增 Repository 接口和实现
- 新增应用服务和命令
- 新增门面层 HTTP 接口
- 数据库新增两张表
