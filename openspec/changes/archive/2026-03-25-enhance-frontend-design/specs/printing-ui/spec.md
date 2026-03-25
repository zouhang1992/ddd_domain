## ADDED Requirements

### Requirement: Print receipt
系统 SHALL 允许用户打印收账收据。

#### Scenario: Print receipt successfully
- **WHEN** 用户点击打印按钮
- **THEN** 系统生成 RTF 格式的收据并打印

#### Scenario: Receipt preview
- **WHEN** 用户查看收据详情
- **THEN** 系统显示收据预览

### Requirement: Print lease contract
系统 SHALL 允许用户打印租房合同。

#### Scenario: Print contract successfully
- **WHEN** 用户点击打印合同按钮
- **THEN** 系统生成 RTF 格式的合同并打印

#### Scenario: Contract preview
- **WHEN** 用户查看合同详情
- **THEN** 系统显示合同预览

### Requirement: Print format selection
系统 SHALL 支持选择打印格式。

#### Scenario: Select print format
- **WHEN** 用户选择打印格式
- **THEN** 系统按照所选格式生成打印内容
