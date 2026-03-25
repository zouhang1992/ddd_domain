## ADDED Requirements

### Requirement: Lease activation capability
系统 SHALL 允许用户将待生效租约转换为生效状态。

#### Scenario: Activate button available for pending leases
- **WHEN** 用户查看待生效的租约
- **THEN** 操作列显示"生效"按钮

#### Scenario: Lease activation successful
- **WHEN** 用户点击"生效"按钮且租约符合生效条件
- **THEN** 租约状态变为"生效中"，系统记录生效时间

#### Scenario: Lease activation fails with invalid status
- **WHEN** 用户尝试激活非待生效状态的租约
- **THEN** 系统显示错误信息，租约状态保持不变

### Requirement: Activation validation
系统 SHALL 验证租约生效的条件。

#### Scenario: Activation requires valid start date
- **WHEN** 用户尝试激活开始日期未到的租约
- **THEN** 系统显示错误信息，租约无法激活

#### Scenario: Activation updates lease status
- **WHEN** 租约成功激活
- **THEN** 租约状态从"待生效"变为"生效中"，激活时间记录到系统

### Requirement: API endpoint for lease activation
系统 SHALL 提供租约生效的 API 接口。

#### Scenario: Activate lease via API
- **WHEN** 客户端发送激活租约的请求
- **THEN** 系统验证条件并更新租约状态

#### Scenario: API responds with appropriate status
- **WHEN** 租约激活成功
- **THEN** API 响应 200 状态码和成功信息

#### Scenario: API responds with validation errors
- **WHEN** 租约激活条件不满足
- **THEN** API 响应 400/409 状态码和错误信息
