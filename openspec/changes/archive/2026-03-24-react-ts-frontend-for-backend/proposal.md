## Why

后端租房管理系统已有完整的业务逻辑和 API 接口，但缺少对应的前端界面。为了方便用户使用系统功能，需要开发一个基于 React + TypeScript 的前端应用。

## What Changes

- 创建前端项目结构（React + TypeScript + Vite）
- 实现位置管理页面（增删改查）
- 实现房间管理页面（增删改查）
- 实现房东管理页面（增删改查）
- 实现租约管理页面（增删改查、续租、退租）
- 实现账单管理页面（增删改查、打印收据）
- 实现打印功能页面（打印账单、租约、发票）
- 实现收入报表页面
- 实现用户认证页面
- 配置 API 客户端与后端交互

## Capabilities

### New Capabilities
- `frontend-app`: React + TypeScript 前端应用
- `location-management-ui`: 位置管理用户界面
- `room-management-ui`: 房间管理用户界面
- `landlord-management-ui`: 房东管理用户界面
- `lease-management-ui`: 租约管理用户界面
- `bill-management-ui`: 账单管理用户界面
- `print-ui`: 打印功能用户界面
- `income-report-ui`: 收入报表用户界面
- `auth-ui`: 用户认证界面

### Modified Capabilities
无

## Impact

- 新增 `web/` 目录存放前端代码
- 前端通过 REST API 与后端 `:8080` 端口交互
- 依赖 React、TypeScript、React Router、Axios 等前端库
