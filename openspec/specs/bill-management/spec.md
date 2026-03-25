# bill-management Specification

## Purpose
定义账单管理的功能需求，包括收账记录、收据打印、账单修改和删除。

## Requirements

### Requirement: 收账记录
系统 SHALL 支持记录收账，包括租金和水电费。

#### Scenario: 成功记录收账
- **WHEN** 用户提交收账信息（租约、租金覆盖区间、水电单价和用量、到账时间、备注）
- **THEN** 系统创建账单记录，计算总金额
- **AND** 更新租约的最后收账时间

### Requirement: 账单查询
系统 SHALL 支持查询账单列表和单个账单详情。

#### Scenario: 查询账单列表
- **WHEN** 用户请求账单列表（可按房间、租约、月份过滤）
- **THEN** 系统返回符合条件的账单列表

#### Scenario: 查询账单详情
- **WHEN** 用户请求指定 ID 的账单
- **THEN** 系统返回该账单的详细信息

### Requirement: 账单修改
系统 SHALL 支持修改账单信息（到账时间、备注、金额等），不同类型限制不同。

#### Scenario: 成功修改账单
- **WHEN** 用户提交账单更新信息
- **AND** 符合该账单类型的修改限制
- **THEN** 系统更新账单信息并返回更新后的账单

### Requirement: 账单删除
系统 SHALL 支持删除账单，checkout 类型会回滚房间/押金/租约状态。

#### Scenario: 成功删除 charge 类型账单
- **WHEN** 用户删除 charge 类型账单
- **THEN** 系统删除该账单

#### Scenario: 成功删除 checkout 类型账单
- **WHEN** 用户删除 checkout 类型账单
- **THEN** 系统删除该账单
- **AND** 回滚房间状态、押金状态、租约状态

### Requirement: 收据打印
系统 SHALL 支持打印收据，生成 RTF 格式。

#### Scenario: 成功打印收据
- **WHEN** 用户请求打印指定账单的收据
- **THEN** 系统生成 RTF 格式的收据文件
- **AND** 返回给用户下载
