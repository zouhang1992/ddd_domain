# OIDC/Keycloak 集成设计文档

**日期**: 2026-04-01  
**状态**: 待审核  
**版本**: 1.0

## 概述

将现有的简单认证系统替换为基于 OIDC (OpenID Connect) 的认证系统，使用 Keycloak 作为授权服务器。实现完整的 RBAC (基于角色的访问控制) 功能。

## 需求回顾

| 需求项 | 选择 |
|--------|------|
| 认证流程 | 完全替换当前认证 |
| Keycloak 部署 | 外部服务 |
| 授权范围 | 完整 RBAC |
| 登录体验 | 标准 OIDC 流程（重定向） |
| Token 处理 | 后端代理模式 |

## 架构设计

### 整体流程图

```
┌─────────────┐                    ┌─────────────┐                    ┌─────────────┐
│   Frontend  │                    │   Backend   │                    │   Keycloak  │
│   (React)   │                    │    (Go)     │                    │             │
└──────┬──────┘                    └──────┬──────┘                    └──────┬──────┘
       │                                    │                                    │
       │  1. 用户点击登录                    │                                    │
       ├───────────────────────────────────>│                                    │
       │                                    │                                    │
       │  2. 重定向到 Keycloak              │                                    │
       │<───────────────────────────────────┤                                    │
       │                                    │                                    │
       │  3. 用户在 Keycloak 登录           │                                    │
       ├──────────────────────────────────────────────────────────────────────>│
       │                                    │                                    │
       │  4. 返回授权码 (code)              │                                    │
       │<──────────────────────────────────────────────────────────────────────┤
       │                                    │                                    │
       │  5. 发送 code 到后端               │                                    │
       ├───────────────────────────────────>│                                    │
       │                                    │                                    │
       │                                    │  6. 用 code 换 tokens             │
       │                                    ├───────────────────────────────────>│
       │                                    │                                    │
       │                                    │  7. 返回 access_token             │
       │                                    │<───────────────────────────────────┤
       │                                    │                                    │
       │                                    │  8. 验证 token，创建 session        │
       │                                    │                                    │
       │  9. 返回 session cookie           │                                    │
       │<───────────────────────────────────┤                                    │
       │                                    │                                    │
       │  10. 后续请求带 cookie             │                                    │
       ├───────────────────────────────────>│                                    │
       │                                    │                                    │
       │  11. 验证 session，检查 RBAC       │                                    │
       │<───────────────────────────────────┤                                    │
```

### 设计原则

1. **后端代理模式**: 前端不直接接触 OAuth2 tokens，由后端完全管理
2. **Session 管理**: 使用 HttpOnly, Secure Cookie 存储 Session ID
3. **完整 RBAC**: 基于 Keycloak 的 realm roles 和 resource roles
4. **标准 OIDC**: 使用 Authorization Code Flow 确保安全性
5. **无第三方依赖**: 仅使用 Go 标准库 `golang.org/x/oauth2`

## 后端组件设计

### 目录结构变更

```
internal/
├── application/
│   ├── auth/              # 新增：OIDC 认证服务
│   │   ├── oidc_service.go
│   │   ├── rbac_service.go
│   │   └── module.go
│   └── common/
│       └── service/
│           └── auth.go    # 删除：旧的简单认证服务
├── infrastructure/
│   ├── middleware/        # 新增：认证中间件
│   │   ├── auth.go
│   │   └── rbac.go
│   └── persistence/
│       └── sqlite/
│           ├── session_repo.go  # 新增：Session 仓储
│           └── module.go
└── facade/
    ├── auth_handler.go    # 删除：旧的认证处理器
    └── oidc_handler.go    # 新增：OIDC HTTP 处理器
```

### 核心组件详解

#### 1. OIDC Service (`internal/application/auth/oidc_service.go`)

**职责**:
- 实现 OAuth2 Authorization Code Flow
- 管理 OIDC Discovery 配置缓存
- 管理 JWKS (JSON Web Key Set) 缓存
- 验证 ID Token 和 Access Token
- 刷新 Access Token (使用 Refresh Token)

**接口定义**:
```go
type OIDCService interface {
    GetAuthURL(state string) string
    ExchangeCode(code string) (*TokenSet, error)
    VerifyToken(token string) (*UserClaims, error)
    RefreshToken(refreshToken string) (*TokenSet, error)
}
```

#### 2. RBAC Service (`internal/application/auth/rbac_service.go`)

**职责**:
- 从 token claims 中提取用户角色和权限
- 检查用户是否拥有指定角色
- 检查用户是否拥有指定权限
- 缓存角色-权限映射

**接口定义**:
```go
type RBACService interface {
    HasRole(claims *UserClaims, role string) bool
    HasPermission(claims *UserClaims, permission string) bool
    GetRoles(claims *UserClaims) []string
    GetPermissions(claims *UserClaims) []string
}
```

**Claims 结构**:
```go
type UserClaims struct {
    Sub          string   `json:"sub"`           // 用户 ID
    Email        string   `json:"email"`         // 邮箱
    Name         string   `json:"name"`          // 姓名
    RealmRoles   []string `json:"realm_roles"`   // Realm 级别角色
    ResourceRoles map[string][]string `json:"resource_roles"` // 资源级别角色
    Permissions  []string `json:"permissions"`   // 权限列表
    Exp          int64    `json:"exp"`           // 过期时间
}
```

#### 3. Session Repository (`internal/infrastructure/persistence/sqlite/session_repo.go`)

**职责**:
- 创建、读取、删除 Session
- 关联 Session 与 Keycloak tokens
- Session TTL 管理

**数据模型**:
```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    id_token TEXT NOT NULL,
    claims JSON NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

#### 4. Auth Middleware (`internal/infrastructure/middleware/auth.go`)

**职责**:
- 从 Cookie 中读取 Session ID
- 验证 Session 是否有效
- 将用户信息注入请求上下文
- 处理未认证请求

**使用方式**:
```go
mux.Handle("GET /protected", authMiddleware.RequireAuth(protectedHandler))
```

#### 5. RBAC Middleware (`internal/infrastructure/middleware/rbac.go`)

**职责**:
- 检查用户是否拥有所需角色/权限
- 拒绝未授权访问
- 提供灵活的权限检查方式

**使用方式**:
```go
mux.Handle("GET /admin", 
    authMiddleware.RequireAuth(
        rbacMiddleware.RequireRole("admin", adminHandler)))
```

#### 6. OIDC Handler (`internal/facade/oidc_handler.go`)

**HTTP 路由**:
- `GET /oauth2/login` - 启动 OIDC 认证流程
- `GET /oauth2/callback` - 处理 Keycloak 回调
- `POST /oauth2/logout` - 登出
- `GET /oauth2/userinfo` - 获取当前用户信息

### 配置项

新增环境变量:

```bash
# Keycloak OIDC 配置
KEYCLOAK_ISSUER_URL=https://keycloak.example.com/realms/your-realm
KEYCLOAK_CLIENT_ID=ddd-domain-app
KEYCLOAK_CLIENT_SECRET=your-client-secret-here
KEYCLOAK_REDIRECT_URL=http://localhost:8080/oauth2/callback

# Session 配置
SESSION_SECRET=your-session-secret-key-here-min-32-chars
SESSION_TTL=86400  # Session 有效期（秒），默认 24 小时

# OIDC 可选配置
KEYCLOAK_SCOPES=openid,profile,email,roles  # 请求的 scopes
KEYCLOAK_ACR_VALUES=                      # ACR 值
```

## 前端改动

### 移除文件
- `web/src/pages/Login.tsx` - 旧的登录页面

### 修改文件

#### 1. `web/src/context/AuthContext.tsx`

**变更**:
- 移除 token 管理逻辑
- 改为检查 `/oauth2/userinfo` 端点来判断认证状态
- 登录按钮直接跳转到 `/oauth2/login`

#### 2. `web/src/App.tsx`

**变更**:
- 移除 Login 路由
- 调整受保护路由的认证检查逻辑

#### 3. `web/src/api/auth.ts`

**变更**:
- 移除 login/logout 调用
- 简化为检查认证状态的 API

### 新增文件
- 无

### 登录流程简化

**旧流程**:
1. 用户在应用内输入用户名密码
2. 前端调用 `/login` API
3. 后端验证硬编码凭证，返回 token
4. 前端存储 token 到 localStorage

**新流程**:
1. 用户点击"登录"按钮
2. 前端跳转到 `/oauth2/login`
3. 后端生成 state，重定向到 Keycloak
4. 用户在 Keycloak 登录
5. Keycloak 回调 `/oauth2/callback?code=...&state=...`
6. 后端验证 state，用 code 换取 tokens
7. 后端验证 tokens，创建 Session
8. 后端设置 HttpOnly Cookie，重定向到首页

## 数据库变更

### 新增表

```sql
-- Session 表
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    id_token TEXT NOT NULL,
    claims TEXT NOT NULL,  -- JSON 格式存储 claims
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
```

### 删除旧认证相关代码
- 旧的 AuthService 中硬编码的 admin/admin
- 简单 token 生成/验证逻辑

## Keycloak 配置要求

### Realm 配置
- 创建或使用现有的 Realm
- 配置用户和角色

### Client 配置
1. **Client Type**: Confidential
2. **Access Type**: Confidential
3. **Standard Flow Enabled**: On
4. **Direct Access Grants Enabled**: Off
5. **Authorization Enabled**: Off
6. **Valid Redirect URIs**: `http://localhost:8080/oauth2/callback`（开发环境）
7. **Web Origins**: `http://localhost:8080`（开发环境）
8. **Client Authentication**: On
9. **Service Accounts Roles**: On

### Client Scopes
添加以下 scopes:
- `openid` (必需)
- `profile` (获取用户基本信息)
- `email` (获取用户邮箱)
- `roles` (获取用户角色)

### Roles 配置
建议配置以下 Realm Roles:
- `admin` - 管理员权限
- `user` - 普通用户权限

### Mappers 配置
配置以下 Protocol Mappers:
1. **User Realm Role**: 将 Realm Roles 添加到 token
2. **Audience**: 添加 client ID 到 audience
3. 确保 `email` 和 `name` claims 包含在 ID Token 中

## 安全考虑

### 1. CSRF 防护
- OIDC state 参数防止 CSRF
- Session Cookie 设置 SameSite=Lax

### 2. XSS 防护
- Token 不暴露给前端
- Session Cookie 设置 HttpOnly 标志
- Session Cookie 设置 Secure 标志（生产环境）

### 3. Token 安全
- Access Token 不发送给前端
- Refresh Token 安全存储在后端
- Token 验证使用 JWKS

### 4. Session 安全
- Session ID 使用加密随机数生成
- Session 有明确的过期时间
- 登出时立即销毁 Session

## 迁移策略

### 步骤 1: 准备 Keycloak
- 配置 Keycloak Realm 和 Client
- 创建初始用户和角色

### 步骤 2: 后端改造
- 添加新的 OIDC 相关组件
- 保持旧认证代码不变，使用 feature flag 切换

### 步骤 3: 前端改造
- 修改认证流程
- 移除旧的登录页面

### 步骤 4: 测试
- 完整测试 OIDC 认证流程
- 测试 RBAC 权限控制
- 测试 Session 管理

### 步骤 5: 上线
- 配置生产环境 Keycloak
- 部署更新后的应用
- 监控认证相关指标

## 回滚计划

如果遇到问题，可以快速回滚：
1. 使用环境变量切换回旧认证方式
2. 或者部署上一个版本

## 测试计划

### 单元测试
- OIDC Service 测试（使用 mock Keycloak）
- RBAC Service 测试
- Middleware 测试

### 集成测试
- 完整的 OIDC 流程测试
- Session 管理测试
- RBAC 权限测试

### E2E 测试
- 登录/登出流程
- 受保护页面访问
- 权限检查

## 后续优化方向

1. **单点登出 (SLO)**: 实现 OIDC Single Logout
2. **权限缓存**: 实现分布式权限缓存
3. **多租户**: 支持多 Realm 或多 Client
4. **Social Login**: 添加 GitHub/Google 等社交登录
5. **MFA**: 支持多因素认证

## 总结

本设计文档描述了将现有简单认证系统替换为基于 Keycloak 的 OIDC 认证系统的完整方案。采用纯 Go 标准库实现，确保安全性和可维护性，同时实现完整的 RBAC 功能。

---

**审批状态**: 待审批  
**下一步**: 用户审阅后，调用 writing-plans skill 创建实现计划
