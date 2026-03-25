## ADDED Requirements

### Requirement: Dashboard displays key metrics
系统 SHALL 在仪表盘页面显示位置、房间、房东、租约等关键指标。

#### Scenario: Dashboard loads successfully
- **WHEN** 用户访问仪表盘页面
- **THEN** 系统显示位置、房间、房东、租约的数量统计

#### Scenario: Metrics update in real-time
- **WHEN** 数据发生变化
- **THEN** 仪表盘指标自动更新

### Requirement: Dashboard provides quick access
系统 SHALL 在仪表盘提供快速访问其他功能的入口。

#### Scenario: Quick navigation to sections
- **WHEN** 用户点击仪表盘的快速访问入口
- **THEN** 系统导航到对应的功能页面
