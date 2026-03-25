# landlord-management Specification

## Purpose
定义房东信息管理的功能需求，包括房东信息的创建、查询、更新和删除。

## Requirements

### Requirement: 房东创建
系统 SHALL 支持创建新房东，包含姓名、联系方式等信息。

#### Scenario: 成功创建房东
- **WHEN** 用户提交房东信息（姓名、电话、备注）
- **THEN** 系统保存房东信息并返回创建的房东

### Requirement: 房东查询
系统 SHALL 支持查询单个房东和所有房东列表。

#### Scenario: 查询房东详情
- **WHEN** 用户请求指定 ID 的房东
- **THEN** 系统返回该房东的详细信息

#### Scenario: 查询房东列表
- **WHEN** 用户请求所有房东
- **THEN** 系统返回所有房东的列表

### Requirement: 房东更新
系统 SHALL 支持更新房东信息。

#### Scenario: 成功更新房东
- **WHEN** 用户提交房东更新信息
- **THEN** 系统更新房东信息并返回更新后的房东

### Requirement: 房东删除
系统 SHALL 支持删除房东，但删除前需检查是否有关联租约，有关联则不允许删除。

#### Scenario: 成功删除无关联租约的房东
- **WHEN** 用户删除无关联租约的房东
- **THEN** 系统删除该房东

#### Scenario: 拒绝删除有关联租约的房东
- **WHEN** 用户删除有关联租约的房东
- **THEN** 系统返回错误，不执行删除操作
