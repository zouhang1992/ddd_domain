## ADDED Requirements

### Requirement: 领域事件监听器
系统 SHALL 实现 `OperationLogHandler` 事件处理器，监听所有领域事件。

#### Scenario: 成功监听事件
- **WHEN** 系统启动
- **THEN** `OperationLogHandler` 自动注册到事件总线
- **AND** 监听所有类型的领域事件

#### Scenario: 事件到操作日志的转换
- **WHEN** 领域事件被发布到事件总线
- **THEN** `OperationLogHandler` 接收并处理事件
- **AND** 将事件内容转换为 `OperationLog` 对象

#### Scenario: 事件属性提取
- **WHEN** 处理领域事件时
- **THEN** 系统从事件中提取以下信息：
  - 事件名称 → `event_name` 字段
  - 事件发生时间 → `timestamp` 字段
  - 领域类型（根据事件名称推断） → `domain_type` 字段
  - 聚合根ID（从事件数据中提取） → `aggregate_id` 字段
  - 完整事件数据 → `details` 字段（JSON格式）

### Requirement: 操作类型识别
系统 SHALL 根据事件名称识别操作类型（创建、更新、删除）。

#### Scenario: 识别创建操作
- **WHEN** 事件名称以 "Created" 结尾（如 "RoomCreated"）
- **THEN** 系统识别为创建操作（`action = "created"`）

#### Scenario: 识别更新操作
- **WHEN** 事件名称以 "Updated" 结尾（如 "LeaseUpdated"）
- **THEN** 系统识别为更新操作（`action = "updated"`）

#### Scenario: 识别删除操作
- **WHEN** 事件名称以 "Deleted" 结尾（如 "LandlordDeleted"）
- **THEN** 系统识别为删除操作（`action = "deleted"`）

#### Scenario: 其他操作类型
- **WHEN** 事件名称不符合标准模式
- **THEN** 系统使用默认操作类型（`action = "unknown"`）

### Requirement: 异步事件处理
系统 SHALL 使用异步方式处理事件，以避免阻塞主业务流程。

#### Scenario: 异步记录日志
- **WHEN** 领域事件发生
- **THEN** 事件总线使用异步方式调用 `OperationLogHandler`
- **AND** 主业务流程不等待日志记录完成

#### Scenario: 日志记录失败不影响业务
- **WHEN** 日志记录过程中发生错误（如数据库连接失败）
- **THEN** 系统记录错误日志
- **AND** 不影响主业务操作的完成

### Requirement: 操作人信息记录
系统 SHALL 记录操作人的ID。

#### Scenario: 从认证上下文获取操作人
- **WHEN** 事件处理时认证上下文可用
- **THEN** 系统从认证上下文中获取用户ID作为 `operator_id`
- **AND** 存储到操作日志中

#### Scenario: 无操作人信息
- **WHEN** 事件处理时没有认证上下文（如系统自动操作）
- **THEN** `operator_id` 字段为空（NULL）
