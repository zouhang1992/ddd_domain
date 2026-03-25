## ADDED Requirements

### Requirement: 房间创建
系统 SHALL 支持创建新房间，包含位置、房号和标签信息（标签为逗号分隔字符串）。

#### Scenario: 成功创建房间
- **WHEN** 用户提交房间信息（位置、房号、标签）
- **THEN** 系统保存房间信息并返回创建的房间

### Requirement: 房间查询
系统 SHALL 支持查询单个房间和所有房间列表。

#### Scenario: 查询房间详情
- **WHEN** 用户请求指定 ID 的房间
- **THEN** 系统返回该房间的详细信息，包括位置和标签

#### Scenario: 查询房间列表
- **WHEN** 用户请求所有房间
- **THEN** 系统返回所有房间的列表

### Requirement: 房间更新
系统 SHALL 支持更新房间信息，包括位置、房号和标签。

#### Scenario: 成功更新房间
- **WHEN** 用户提交房间更新信息
- **THEN** 系统更新房间信息并返回更新后的房间

### Requirement: 房间删除
系统 SHALL 支持删除房间。

#### Scenario: 成功删除房间
- **WHEN** 用户删除房间
- **THEN** 系统删除该房间
