## ADDED Requirements

### Requirement: Income query by month
系统 SHALL 允许用户按月份查询收入汇总。

#### Scenario: Successful income query
- **WHEN** 用户选择月份并点击查询
- **THEN** 系统显示该月份的收入汇总数据

#### Scenario: No income data
- **WHEN** 所选月份无收入数据
- **THEN** 系统显示无数据提示

### Requirement: Income report export
系统 SHALL 允许用户导出收入报表。

#### Scenario: Export income report
- **WHEN** 用户点击导出按钮
- **THEN** 系统下载收入报表文件

### Requirement: Income trend display
系统 SHALL 显示收入趋势图。

#### Scenario: Income trend visualization
- **WHEN** 用户查看收入查询页面
- **THEN** 系统显示收入趋势图表

### Requirement: Income breakdown
系统 SHALL 提供收入分类明细。

#### Scenario: Income breakdown by type
- **WHEN** 用户查看收入明细
- **THEN** 系统显示租金、水电费等收入分类
