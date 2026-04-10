# DDD House 部署指南

## 概述

本项目包含完整的Docker和Kubernetes部署配置，支持一键部署到本地K8s集群。

## 项目架构

- **域名**: `ddd-house.zouhang.com`
- **后端API路径**: `/api/*`
- **认证服务**: Keycloak (`keycloak.zouhang.com`)
- **前端**: React + TypeScript + Ant Design (Nginx)
- **后端**: Go + SQLite
- **Ingress**: Nginx Ingress (统一处理前后端路由)

## 文件结构

```
deploy/
├── k8s/
│   ├── keycloak/
│   │   ├── helm-values.yaml    # Keycloak Helm配置
│   │   └── keycloak.yaml       # Keycloak自定义部署（已使用）
│   └── app/
│       └── ddd-house.yaml       # 应用K8s部署文件
└── nginx/
    └── nginx.conf               # 前端Nginx配置（仅静态文件）

Dockerfile.backend        # 后端Docker镜像
Dockerfile.frontend       # 前端Docker镜像
.env                      # 环境配置
.env.example             # 环境配置示例
```

## 本地开发

### 后端开发

```bash
# 运行后端（Go）
go run cmd/api/main.go
```

### 前端开发

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

访问地址:
- 前端: http://localhost:5173
- 后端: http://localhost:8080
- API: http://localhost:8080/api/*

## Kubernetes 部署

### 前置条件

1. 本地K8s集群（Docker Desktop）
2. Keycloak已部署（参考 [deploy/k8s/keycloak/](./deploy/k8s/keycloak/)）
3. Harbor镜像仓库（harbor.zouhang.com）

### 1. 配置Hosts

在 `/etc/hosts` 中添加：

```
127.0.0.1 ddd-house.zouhang.com
127.0.0.1 keycloak.zouhang.com
```

### 2. 配置Keycloak

确保Keycloak中已配置：

1. **Realm**: `master`（或创建新的realm）
2. **Client**: `ddd-app`
   - Access Type: `confidential`
   - Valid Redirect URIs: `http://ddd-house.zouhang.com/*`
   - Web Origins: `http://ddd-house.zouhang.com`
3. **Client Secret**: 从Keycloak获取并更新到配置中

### 3. 构建并推送镜像

```bash
# 构建后端镜像
docker build -f Dockerfile.backend -t harbor.zouhang.com/library/ddd-house-backend:latest .

# 构建前端镜像
docker build -f Dockerfile.frontend -t harbor.zouhang.com/library/ddd-house-frontend:latest .

# 登录Harbor
docker login harbor.zouhang.com

# 推送镜像
docker push harbor.zouhang.com/library/ddd-house-backend:latest
docker push harbor.zouhang.com/library/ddd-house-frontend:latest
```

### 4. 部署到K8s

```bash
# 部署应用
kubectl apply -f deploy/k8s/app/ddd-house.yaml

# 查看部署状态
kubectl get pods -n ddd-house
kubectl get svc -n ddd-house
kubectl get ingress -n ddd-house

# 查看日志
kubectl logs -f deployment/ddd-house-backend -n ddd-house
kubectl logs -f deployment/ddd-house-frontend -n ddd-house
```

### 5. 访问应用

打开浏览器访问: http://ddd-house.zouhang.com

## 配置说明

### 后端配置 (.env)

| 配置项 | 说明 | 示例值 |
|--------|------|--------|
| `HTTP_ADDR` | HTTP监听地址 | `:8080` |
| `DATABASE_DSN` | SQLite数据库路径 | `/data/ddd.db` |
| `LOG_ENVIRONMENT` | 日志环境 | `production` |
| `OIDC_DEV_MODE` | 开发模式（跳过OIDC） | `false` |
| `OIDC_ISSUER_URL` | Keycloak Issuer URL | `http://keycloak.keycloak.svc.cluster.local:8080/realms/master` |
| `OIDC_CLIENT_ID` | Keycloak Client ID | `ddd-app` |
| `OIDC_CLIENT_SECRET` | Keycloak Client Secret | `***` |
| `OIDC_REDIRECT_URL` | OIDC回调地址 | `http://ddd-house.zouhang.com/oauth2/callback` |
| `OIDC_FRONTEND_URL` | 前端地址 | `http://ddd-house.zouhang.com` |

### 前端配置

- API_BASE_URL: `/api` (在 `web/src/api/request.ts` 中配置)
- Vite代理: `web/vite.config.ts` (仅用于本地开发)

## 网络路径说明

### 生产环境 (K8s Ingress)

```
用户浏览器
    ↓
ddd-house.zouhang.com
    ↓
K8s Ingress (nginx)
    ├─ / → 前端 (ddd-house-frontend:80)
    ├─ /api/* → 后端 (ddd-house-backend:8080)
    ├─ /oauth2/* → 后端 (ddd-house-backend:8080)
    ├─ /health → 后端 (ddd-house-backend:8080)
    └─ /metrics → 后端 (ddd-house-backend:8080)
```

### 本地开发 (Vite代理)

```
浏览器 (localhost:5173)
    ↓
Vite Dev Server
    ├─ / → 前端
    ├─ /api/* → http://localhost:8080/api/*
    ├─ /oauth2/* → http://localhost:8080/oauth2/*
    └─ ...
```

## 故障排查

### Pod无法启动

```bash
# 查看Pod事件
kubectl describe pod <pod-name> -n ddd-house

# 查看日志
kubectl logs <pod-name> -n ddd-house
```

### 镜像拉取失败

确保:
1. 镜像已推送到Harbor
2. 节点能访问harbor.zouhang.com
3. 镜像名称和标签正确

### OIDC认证失败

检查:
1. Keycloak Pod是否运行正常
2. `OIDC_ISSUER_URL` 配置正确（K8s内部使用service域名）
3. Client配置中的Redirect URI白名单包含 `http://ddd-house.zouhang.com/*`
4. Client Secret正确

### Ingress路由问题

```bash
# 查看Ingress配置
kubectl get ingress ddd-house -n ddd-house -o yaml

# 查看Ingress Controller日志
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
```

## 清理

```bash
# 删除应用部署
kubectl delete -f deploy/k8s/app/ddd-house.yaml

# 删除namespace（会删除所有相关资源）
kubectl delete namespace ddd-house
```
