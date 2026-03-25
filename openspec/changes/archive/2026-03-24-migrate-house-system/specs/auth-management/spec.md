# auth-management Specification

## Purpose
定义用户认证和授权的功能需求，包括登录、登出和鉴权。

## Requirements

### Requirement: 用户登录
系统 SHALL 支持用户登录功能。

#### Scenario: 成功登录
- **WHEN** 用户提交正确的用户名和密码
- **THEN** 系统验证凭证
- **AND** 生成 JWT Token
- **AND** 设置 Cookie 或返回 Token

#### Scenario: 登录失败（凭证错误）
- **WHEN** 用户提交错误的用户名或密码
- **THEN** 系统返回错误，不生成 Token

### Requirement: 用户登出
系统 SHALL 支持用户登出功能。

#### Scenario: 成功登出
- **WHEN** 用户请求登出
- **THEN** 系统清除登录状态（清除 Cookie 或使 Token 失效）

### Requirement: 鉴权
系统 SHALL 支持对受保护的 API 端点进行鉴权。

#### Scenario: 成功通过鉴权
- **WHEN** 用户请求受保护的 API 端点
- **AND** 携带有效的 JWT Token
- **THEN** 系统允许访问

#### Scenario: 鉴权失败（无 Token）
- **WHEN** 用户请求受保护的 API 端点
- **AND** 未携带 Token
- **THEN** 系统返回 401 未授权错误

#### Scenario: 鉴权失败（Token 无效）
- **WHEN** 用户请求受保护的 API 端点
- **AND** 携带无效或过期的 Token
- **THEN** 系统返回 401 未授权错误

### Requirement: 获取当前用户信息
系统 SHALL 支持获取当前登录用户的信息。

#### Scenario: 成功获取用户信息
- **WHEN** 用户请求 /me 端点
- **AND** 携带有效的 Token
- **THEN** 系统返回当前用户的信息
