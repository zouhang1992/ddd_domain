## Context

当前项目是一个基于 DDD 架构的租房管理后端系统，使用 Go 语言开发，已经实现了完整的业务逻辑和 REST API 接口。系统包含位置管理、房间管理、房东管理、租约管理、账单管理、打印服务等功能模块。

**约束条件：**
- 后端 API 运行在 `http://localhost:8080`
- 使用标准的 REST 架构
- 需要支持用户认证
- 打印功能需要处理 RTF 文件下载

## Goals / Non-Goals

**Goals:**
- 提供一个现代化、易用的 Web 前端界面
- 使用 React + TypeScript 构建，确保类型安全
- 实现所有后端 API 对应的前端功能
- 支持响应式设计，适配不同屏幕尺寸
- 提供清晰的用户反馈和错误处理

**Non-Goals:**
- 不实现后端业务逻辑（复用现有 API）
- 不实现移动端原生应用
- 不实现复杂的数据可视化（保留简单的报表）
- 不实现实时通信功能

## Decisions

### 1. 前端框架选型
**决策：** 使用 React 18 + TypeScript 5 + Vite 5

**理由：**
- React 生态成熟，组件库丰富
- TypeScript 提供类型安全，减少运行时错误
- Vite 提供快速的开发体验和构建速度
- 团队对 React 技术栈熟悉

**替代方案考虑：**
- Vue 3：生态也成熟，但团队 React 经验更丰富
- Angular：功能完整但学习曲线较陡
- Svelte：轻量但生态较小

### 2. 状态管理
**决策：** 使用 React Context API + useReducer（轻量级场景），不引入额外状态管理库

**理由：**
- 应用规模适中，Context API 足够使用
- 避免引入 Redux 等复杂库带来的额外复杂度
- 减少依赖，降低维护成本

**替代方案考虑：**
- Redux Toolkit：功能强大但对于当前规模过于复杂
- Zustand：轻量但增加额外依赖
- Jotai：原子化状态管理，但同样增加依赖

### 3. 路由管理
**决策：** 使用 React Router v6

**理由：**
- React Router 是 React 生态最流行的路由库
- v6 版本 API 简洁，类型支持好
- 支持嵌套路由、动态路由等高级功能

**替代方案考虑：**
- TanStack Router：类型安全但生态较小
- 无路由库：使用简单的条件渲染不适合多页面应用

### 4. HTTP 客户端
**决策：** 使用 Axios

**理由：**
- API 简洁，支持拦截器（用于统一处理认证、错误）
- 类型支持好
- 社区广泛使用，文档完善

**替代方案考虑：**
- Fetch API：原生但需要更多封装
- TanStack Query：功能强大但主要用于数据缓存，超出当前需求

### 5. UI 组件库
**决策：** 使用 Ant Design（或类似组件库）

**理由：**
- 组件丰富，开箱即用
- TypeScript 支持好
- 设计风格统一，开发效率高

**替代方案考虑：**
- Material-UI：组件丰富但定制化较复杂
- Chakra UI：可访问性好但组件相对较少
- 纯 CSS/SCSS：需要自己实现所有组件，开发效率低

### 6. 项目结构
**决策：** 按功能模块组织代码

```
web/
├── src/
│   ├── api/          # API 客户端
│   ├── components/   # 通用组件
│   ├── pages/        # 页面组件
│   │   ├── locations/
│   │   ├── rooms/
│   │   ├── landlords/
│   │   ├── leases/
│   │   ├── bills/
│   │   ├── print/
│   │   ├── income/
│   │   └── auth/
│   ├── hooks/        # 自定义 Hooks
│   ├── types/        # TypeScript 类型定义
│   ├── utils/        # 工具函数
│   ├── context/      # React Context
│   ├── App.tsx
│   └── main.tsx
├── package.json
├── tsconfig.json
└── vite.config.ts
```

**理由：**
- 按功能组织便于维护和扩展
- 相关代码放在一起，易于查找
- 符合 React 社区最佳实践

## Risks / Trade-offs

### 风险 1：后端 API 变更
**风险：** 后端 API 可能在前端开发过程中发生变更
**缓解措施：**
- 使用 TypeScript 类型定义 API 接口
- API 客户端集中管理，便于统一修改
- 前后端沟通保持同步

### 风险 2：认证状态管理
**风险：** 用户登录状态、Token 刷新等逻辑容易出错
**缓解措施：**
- 使用 Axios 拦截器统一处理认证
- Token 存储在 localStorage 或 sessionStorage
- 提供清晰的登录/登出流程

### 风险 3：打印功能的用户体验
**风险：** RTF 文件下载和打印的用户体验可能不够友好
**缓解措施：**
- 提供明确的下载提示
- 支持在浏览器中预览（如果可能）
- 提供下载和打印操作的反馈

## Open Questions

1. 是否需要支持多语言？（当前假设只需要中文）
2. 是否需要深色模式？（当前假设使用默认浅色主题）
3. 是否需要支持数据导入导出功能？（当前假设只需要前端展示和操作）
