## ADDED Requirements

### Requirement: 命令定义接口
系统 SHALL 定义 Command 接口，所有命令必须实现该接口。

#### Scenario: 命令实现
- **WHEN** 创建一个新命令
- **THEN** 命令必须实现 Command 接口的方法

### Requirement: 命令处理器接口
系统 SHALL 定义 CommandHandler 接口，用于处理特定类型的命令。

#### Scenario: 注册命令处理器
- **WHEN** 向命令总线注册命令处理器
- **THEN** 命令总线将该处理器与对应命令类型关联

### Requirement: 命令分发
系统 SHALL 支持将命令分发到对应的处理器执行。

#### Scenario: 发送命令
- **WHEN** 向命令总线发送一个命令
- **THEN** 对应的命令处理器被调用
- **AND** 处理器返回的结果或错误被返回给调用者

### Requirement: 中间件支持
系统 SHALL 支持命令中间件，用于在命令处理前后执行横切关注点。

#### Scenario: 使用中间件
- **WHEN** 命令总线配置了中间件
- **THEN** 中间件在命令处理器执行前后被调用
