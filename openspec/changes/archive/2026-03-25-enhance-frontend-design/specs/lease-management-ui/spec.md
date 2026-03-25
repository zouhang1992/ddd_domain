## ADDED Requirements

### Requirement: Lease renewal
系统 SHALL 允许用户处理租约续租。

#### Scenario: Successful lease renewal
- **WHEN** 用户点击续租按钮
- **THEN** 系统生成新的租约记录并绑定押金

#### Scenario: Lease renewal with changes
- **WHEN** 用户在续租时修改租金或期限
- **THEN** 系统按照新信息创建租约

### Requirement: Lease checkout
系统 SHALL 允许用户处理租约退租结算。

#### Scenario: Successful lease checkout
- **WHEN** 用户点击退租结算按钮
- **THEN** 系统记录水电费用、押金扣除并退还剩余押金

#### Scenario: Checkout with deductions
- **WHEN** 退租时有押金扣除
- **THEN** 系统显示扣除明细并计算退款金额

### Requirement: Lease receipt generation
系统 SHALL 在租约操作后生成相应的收据。

#### Scenario: Receipt after charge
- **WHEN** 用户完成收账操作
- **THEN** 系统生成收账收据

#### Scenario: Receipt after checkout
- **WHEN** 用户完成退租结算
- **THEN** 系统生成结算收据
