## ADDED Requirements

### Requirement: Bill charge recording
系统 SHALL 允许用户记录收账信息，包括租金和水电费。

#### Scenario: Successful charge recording
- **WHEN** 用户点击收账按钮并填写信息
- **THEN** 系统保存收账记录并更新租约状态

#### Scenario: Charge with utilities
- **WHEN** 用户记录包含水电费的收账
- **THEN** 系统分别计算租金和水电费金额

### Requirement: Bill editing
系统 SHALL 允许用户修改已有的收账记录。

#### Scenario: Successful bill editing
- **WHEN** 用户修改收账记录
- **THEN** 系统更新记录并重新计算金额

#### Scenario: Edit with validation
- **WHEN** 用户尝试保存无效数据
- **THEN** 系统显示验证错误并拒绝保存

### Requirement: Bill deletion
系统 SHALL 允许用户删除收账记录。

#### Scenario: Successful bill deletion
- **WHEN** 用户删除收账记录
- **THEN** 系统从数据库中移除记录

#### Scenario: Delete with constraints
- **WHEN** 用户尝试删除有约束的记录
- **THEN** 系统显示错误信息并阻止删除

### Requirement: Bill filtering and searching
系统 SHALL 允许用户按条件筛选和搜索账单。

#### Scenario: Filter bills by room
- **WHEN** 用户按房间筛选账单
- **THEN** 系统只显示该房间的账单

#### Scenario: Search bills by keyword
- **WHEN** 用户搜索账单
- **THEN** 系统显示匹配的账单记录
