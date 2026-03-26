---
name: tasks
change: add-execution-logging
description: 执行日志功能实现任务清单
---

## 1. 项目配置和依赖

- [x] 1.1 在 go.mod 中添加 uber-go/zap 直接依赖
- [x] 1.2 运行 go mod tidy 确保依赖正确

## 2. 日志基础设施

- [x] 2.1 创建 internal/infrastructure/logging 目录
- [x] 2.2 实现 Config 配置结构体
- [x] 2.3 实现 NewLogger 函数，支持开发/生产环境配置
- [x] 2.4 实现 Module() 函数，提供 FX 模块
- [x] 2.5 添加日志单元测试

## 3. 应用启动集成

- [x] 3.1 在 cmd/api/main.go 中导入 logging 模块
- [x] 3.2 将 logging.Module 添加到 FX 应用配置中
- [x] 3.3 修改 startServer 函数，注入 logger
- [x] 3.4 替换 log.Println 为 logger.Info
- [x] 3.5 替换 log.Fatalf 为 logger.Fatal

## 4. 数据库连接模块集成

- [x] 4.1 修改 internal/infrastructure/persistence/sqlite/connection.go
- [x] 4.2 为 Connection 添加 logger 字段
- [x] 4.3 记录数据库连接成功的日志
- [x] 4.4 记录数据库初始化操作的日志
- [x] 4.5 记录错误和警告信息

## 5. 命令总线模块集成

- [x] 5.1 修改 internal/infrastructure/bus/command 模块
- [x] 5.2 为命令总线添加 logger 字段
- [x] 5.3 记录命令处理开始和完成的日志
- [x] 5.4 记录命令处理错误的日志
- [x] 5.5 包含命令名称和相关上下文

## 6. 查询总线模块集成

- [x] 6.1 修改 internal/infrastructure/bus/query 模块
- [x] 6.2 为查询总线添加 logger 字段
- [x] 6.3 记录查询处理开始和完成的日志
- [x] 6.4 记录查询处理错误的日志
- [x] 6.5 包含查询名称和相关上下文

## 7. 事件总线模块集成

- [x] 7.1 修改 internal/infrastructure/bus/event 模块
- [x] 7.2 为事件总线添加 logger 字段
- [x] 7.3 记录事件发布和处理的日志
- [x] 7.4 记录事件处理错误的日志
- [x] 7.5 包含事件名称和相关上下文

## 8. 模块配置更新

- [x] 8.1 更新 sqlite 模块配置，注入 logger
- [x] 8.2 更新 bus 模块配置，注入 logger

## 9. 应用层和 Facade 层集成

- [ ] 9.1 修改 internal/application/command/handler 模块（部分完成）
- [ ] 9.2 为命令处理器添加 logger 字段（部分完成）
- [ ] 9.3 记录命令执行的关键步骤（部分完成）
- [x] 9.4 修改 internal/application/event/handler 模块
- [x] 9.5 为事件处理器添加 logger 字段
- [x] 9.6 记录事件处理的关键步骤
- [ ] 9.7 修改 internal/facade 中的 HTTP 处理器
- [ ] 9.8 为 HTTP 处理器添加 logger 字段
- [ ] 9.9 记录请求到达和完成的日志
- [ ] 9.10 包含请求路径、方法和状态码
- [ ] 9.11 记录请求处理错误的日志

## 10. 验证和测试

- [x] 10.1 编译并运行应用，验证日志输出
- [x] 10.2 测试开发环境下的控制台格式
- [x] 10.3 测试生产环境下的 JSON 格式
- [x] 10.4 验证日志级别控制
- [x] 10.5 运行所有现有测试，确保没有破坏现有功能
