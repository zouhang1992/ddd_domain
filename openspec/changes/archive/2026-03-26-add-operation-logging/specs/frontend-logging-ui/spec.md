# frontend-logging-ui Specification

## Purpose
定义前端操作日志展示界面的功能需求，包括通用组件设计和交互体验。

## ADDED Requirements

### Requirement: OperationLogModal 通用组件
系统 SHALL 提供通用的 OperationLogModal 组件，支持在所有页面复用。

#### Scenario: 组件接收 domainType 参数
- **WHEN** OperationLogModal 组件被传入 domainType 参数
- **THEN** 组件根据 domainType 筛选并展示对应类型的操作日志

#### Scenario: 组件接收 aggregateID 参数
- **WHEN** OperationLogModal 组件被传入 aggregateID 参数
- **THEN** 组件筛选并展示特定实体的操作日志

#### Scenario: 组件控制可见性
- **WHEN** OperationLogModal 组件的 visible 属性为 true
- **THEN** 组件显示在页面上

#### Scenario: 组件关闭回调
- **WHEN** 用户点击关闭按钮或模态框外部
- **THEN** 组件触发 onCancel 回调

### Requirement: 操作日志表格展示
系统 SHALL 在 OperationLogModal 组件中使用表格展示操作日志。

#### Scenario: 表格列定义
- **WHEN** 操作日志数据加载完成
- **THEN** 表格展示以下列：操作时间、事件名称、操作类型、详情

#### Scenario: 详情展示
- **WHEN** 操作日志详情字段较长
- **THEN** 系统支持折叠/展开展示，或使用Tooltip显示完整内容

### Requirement: 操作日志分页
系统 SHALL 支持操作日志分页查看，避免一次加载过多数据。

#### Scenario: 默认分页大小
- **WHEN** 组件首次加载操作日志
- **THEN** 系统默认加载20条记录

#### Scenario: 翻页操作
- **WHEN** 用户点击分页控件的翻页按钮
- **THEN** 系统加载对应页的操作日志记录

#### Scenario: 显示总条数
- **WHEN** 操作日志数据加载完成
- **THEN** 分页控件显示总记录数
