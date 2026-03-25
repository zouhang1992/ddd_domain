## ADDED Requirements

### Requirement: Saga 定义接口
系统 SHALL 定义 Saga 接口，用于编排一系列分布式事务步骤。

#### Scenario: Saga 定义
- **WHEN** 创建一个新的 Saga
- **THEN** Saga 必须定义执行步骤和对应的补偿步骤

### Requirement: Saga 步骤执行
系统 SHALL 按顺序执行 Saga 定义的步骤。

#### Scenario: 成功执行 Saga
- **WHEN** Saga 被触发
- **THEN** 所有步骤按顺序成功执行
- **AND** Saga 状态标记为已完成

### Requirement: 补偿事务
系统 SHALL 在 Saga 步骤失败时执行已完成步骤的补偿操作。

#### Scenario: Saga 失败并回滚
- **WHEN** Saga 某个步骤执行失败
- **THEN** 已执行的步骤按逆序执行补偿操作
- **AND** Saga 状态标记为已回滚

### Requirement: Saga 状态持久化
系统 SHALL 支持将 Saga 状态持久化到 SQLite 数据库。

#### Scenario: 状态保存与恢复
- **WHEN** Saga 执行过程中
- **THEN** Saga 状态被保存到数据库
- **AND** 系统重启后可以恢复 Saga 状态
