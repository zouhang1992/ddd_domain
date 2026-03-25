## ADDED Requirements

### Requirement: 完整的总线架构集成
系统 SHALL 在应用层使用 commandbus、querybus 和 eventbus 管理所有逻辑。

#### Scenario: 应用服务作为总线分发器
- **WHEN** 调用应用服务的方法
- **THEN** 应用服务将请求转发到相应的总线（命令总线或查询总线）
- **AND** 命令执行后发布领域事件

### Requirement: HTTP 处理器简化
系统 SHALL 简化 HTTP 处理器，使其直接调用总线而不需要通过传统应用服务。

#### Scenario: HTTP 处理器与命令总线交互
- **WHEN** HTTP 处理器收到请求
- **THEN** 处理器解析请求参数，创建命令
- **AND** 将命令发送到命令总线执行
- **AND** 将执行结果返回给客户端

### Requirement: 查询总线集成
系统 SHALL 将所有查询操作通过查询总线统一管理。

#### Scenario: 查询总线使用
- **WHEN** 执行查询操作
- **THEN** 查询被发送到查询总线
- **AND** 查询总线找到对应的查询处理器
- **AND** 返回处理结果

### Requirement: 事件总线集成
系统 SHALL 在命令执行后发布领域事件，并通过事件总线处理。

#### Scenario: 事件发布
- **WHEN** 命令执行成功
- **THEN** 发布对应的领域事件
- **AND** 事件总线将事件分发给所有订阅者

### Requirement: 移除未使用的服务
系统 SHALL 移除未使用的传统应用服务实现。

#### Scenario: 清理未使用的服务
- **WHEN** 重构完成
- **THEN** 系统只保留认证和打印等工具类服务
- **AND** 移除所有传统的应用服务实现
