## MODIFIED Requirements

### Requirement: 按月收入汇总查询
系统 SHALL 支持按月份查询收入汇总，包含押金的收入和支出。

#### Scenario: 成功查询当月收入汇总
- **WHEN** 用户请求收入汇总，指定月份参数（YYYY-MM 格式）
- **THEN** 系统返回该月份的收入汇总信息，包括：
  - 租金总收入
  - 水电费总收入
  - 押金收入（押金收取时）
  - 押金支出（押金退还时，显示为负收入）
  - 其他收入
  - 总计收入
- **AND** 按到账时间月份统计，而非账单创建月份

#### Scenario: 押金收入显示
- **WHEN** 月份中有押金收取的账单（bill_type = 'deposit'）
- **THEN** 押金金额显示为正收入

#### Scenario: 押金支出显示
- **WHEN** 月份中有押金退还的账单（bill_type = 'checkout' 且包含 deposit_refund）
- **THEN** 押金退还金额显示为负收入
