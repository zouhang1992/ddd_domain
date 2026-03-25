## ADDED Requirements

### Requirement: Lease list shows deposit information
系统 SHALL 在租约列表中显示押金信息。

#### Scenario: Deposit amount displayed in lease list
- **WHEN** 用户查看租约列表
- **THEN** 表格显示"押金"列，显示押金金额和状态

#### Scenario: Deposit status displayed
- **WHEN** 租约有押金绑定
- **THEN** 使用不同颜色的标签显示押金状态（如：已收取、已退还、已扣除）

### Requirement: Print lease contract from list
系统 SHALL 允许用户从租约列表直接打印合同。

#### Scenario: Print button available in lease list
- **WHEN** 用户查看租约列表
- **THEN** 操作列显示"打印合同"按钮

#### Scenario: Print contract initiated
- **WHEN** 用户点击"打印合同"按钮
- **THEN** 系统下载 RTF 格式的合同文件

#### Scenario: Print button disabled for invalid leases
- **WHEN** 租约处于无效状态
- **THEN** 打印合同按钮禁用

### Requirement: Lease status management
系统 SHALL 显示租约的完整状态。

#### Scenario: Lease status displayed in list
- **WHEN** 用户查看租约列表
- **THEN** 使用颜色编码的标签显示租约状态（待生效、生效中、已过期、已退租）
