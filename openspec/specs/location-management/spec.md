# location-management Specification

## Purpose
TBD - created by archiving change rental-room-management. Update Purpose after archive.
## Requirements
### Requirement: 位置创建
系统 SHALL 支持创建新位置，包含位置简称和详细信息。

#### Scenario: 成功创建位置
- **WHEN** 用户提交位置简称和详细信息
- **THEN** 系统保存位置信息并返回创建的位置

### Requirement: 位置查询
系统 SHALL 支持查询单个位置和所有位置列表。

#### Scenario: 查询位置详情
- **WHEN** 用户请求指定 ID 的位置
- **THEN** 系统返回该位置的详细信息

#### Scenario: 查询位置列表
- **WHEN** 用户请求所有位置
- **THEN** 系统返回所有位置的列表

### Requirement: 位置更新
系统 SHALL 支持更新位置信息。

#### Scenario: 成功更新位置
- **WHEN** 用户提交位置更新信息
- **THEN** 系统更新位置信息并返回更新后的位置

### Requirement: 位置删除
系统 SHALL 支持删除位置，但删除前需检查是否有关联房间，有关联则不允许删除。

#### Scenario: 成功删除无关联房间的位置
- **WHEN** 用户删除无关联房间的位置
- **THEN** 系统删除该位置

#### Scenario: 拒绝删除有关联房间的位置
- **WHEN** 用户删除有关联房间的位置
- **THEN** 系统返回错误，不执行删除操作

