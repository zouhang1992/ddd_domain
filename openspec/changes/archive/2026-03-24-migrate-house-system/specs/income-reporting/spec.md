# income-reporting Specification

## Purpose
定义收入汇总查询的功能需求，支持按月份查询收入统计。

## Requirements

### Requirement: 按月收入汇总查询
系统 SHALL 支持按月份查询收入汇总。

#### Scenario: 成功查询当月收入汇总
- **WHEN** 用户请求收入汇总，指定月份参数（YYYY-MM 格式）
- **THEN** 系统返回该月份的收入汇总信息，包括：
  - 租金总收入
  - 水电费总收入
  - 押金收入
  - 其他收入
  - 总计收入
- **AND** 按到账时间月份统计，而非账单创建月份
