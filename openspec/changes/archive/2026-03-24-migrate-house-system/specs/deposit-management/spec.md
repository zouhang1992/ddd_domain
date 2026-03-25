# deposit-management Specification

## Purpose
定义押金管理的功能需求，包括押金的收取、退还、扣除管理。

## Requirements

### Requirement: 押金收取
系统 SHALL 支持在创建租约时收取押金。

#### Scenario: 成功收取押金
- **WHEN** 用户创建租约并指定押金金额
- **THEN** 系统创建押金记录，状态为"已收取"
- **AND** 押金记录与租约关联

### Requirement: 押金退还
系统 SHALL 支持在退租时退还押金。

#### Scenario: 成功退还押金
- **WHEN** 用户退租结算并选择退还押金
- **AND** 提供退还金额
- **THEN** 系统更新押金状态为"已退还"
- **AND** 记录退还金额和时间

### Requirement: 押金扣除
系统 SHALL 支持在退租时扣除部分或全部押金。

#### Scenario: 成功扣除押金
- **WHEN** 用户退租结算并选择扣除押金
- **AND** 提供扣除金额和原因
- **THEN** 系统更新押金状态为"已扣除"
- **AND** 记录扣除金额、原因和时间

### Requirement: 押金状态查询
系统 SHALL 支持查询押金状态和历史。

#### Scenario: 查询押金详情
- **WHEN** 用户请求指定押金的详情
- **THEN** 系统返回押金的完整信息，包括状态变更历史

### Requirement: 续租时押金自动绑定
系统 SHALL 在续租时将原租约的押金自动绑定到新租约。

#### Scenario: 续租时押金自动绑定
- **WHEN** 用户对有押金的租约进行续租
- **THEN** 系统将原租约的押金记录关联到新租约
- **AND** 原租约不再关联该押金
