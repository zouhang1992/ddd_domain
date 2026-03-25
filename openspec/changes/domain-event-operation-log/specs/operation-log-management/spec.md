## ADDED Requirements

### Requirement: 操作日志存储
系统 SHALL 将操作日志存储在 SQLite 数据库的 `operation_logs` 表中。

#### Scenario: 成功保存操作日志
- **WHEN** 领域事件发生并被处理
- **THEN** 系统在 `operation_logs` 表中创建一条新记录
- **AND** 记录包含完整的事件信息（事件名称、时间、领域类型、聚合ID、详情等）

#### Scenario: 查询操作日志
- **WHEN** 用户查询操作日志
- **THEN** 系统返回符合查询条件的日志记录
- **AND** 支持按领域类型、时间范围、聚合ID、操作人等条件筛选

#### Scenario: 分页查询操作日志
- **WHEN** 查询结果数量超过页面大小
- **THEN** 系统返回分页结果
- **AND** 支持指定页码和每页大小

### Requirement: 操作日志查询API
系统 SHALL 提供 RESTful API 用于查询操作日志。

#### Scenario: 获取日志列表
- **WHEN** 发送 GET 请求到 `/api/operation-logs`
- **THEN** 系统返回操作日志列表
- **AND** 支持查询参数：`domainType`, `eventName`, `startTime`, `endTime`, `aggregateId`, `operatorId`

#### Scenario: 获取单条日志详情
- **WHEN** 发送 GET 请求到 `/api/operation-logs/:id`
- **THEN** 系统返回该日志的完整信息
- **AND** `details` 字段包含完整的事件数据（JSON格式）

### Requirement: 前端操作日志查询页面
系统 SHALL 提供前端界面用于查询和展示操作日志。

#### Scenario: 显示查询表单
- **WHEN** 用户访问 `/operation-logs` 页面
- **THEN** 页面显示查询表单
- **AND** 支持选择领域类型、时间范围、聚合ID等查询条件

#### Scenario: 显示日志列表
- **WHEN** 用户点击查询按钮
- **THEN** 页面显示符合条件的日志列表
- **AND** 每条日志显示时间、事件名、领域类型、操作人
- **AND** 支持点击查看详情

#### Scenario: 显示日志详情
- **WHEN** 用户点击日志的查看详情按钮
- **THEN** 系统显示包含完整信息的详情弹窗
- **AND** 显示详细的事件数据（JSON格式）

#### Scenario: 分页显示
- **WHEN** 查询结果有多个页面
- **THEN** 页面显示分页控件
- **AND** 支持页码跳转和每页大小调整
