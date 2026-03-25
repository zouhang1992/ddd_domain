## 1. 依赖准备

- [x] 1.1 添加 uber/fx 依赖到 go.mod
- [x] 1.2 运行 go mod tidy 下载依赖

## 2. 模块化重构

- [x] 2.1 创建 internal/infrastructure/persistence/sqlite/module.go（persistence 模块）
- [x] 2.2 创建 internal/infrastructure/bus/module.go（bus 模块）
- [x] 2.3 创建 internal/application/command/handler/module.go（command 模块）
- [x] 2.4 创建 internal/application/query/handler/module.go（query 模块）
- [x] 2.5 创建 internal/facade/module.go（facade 模块）

## 3. 重构 main.go

- [x] 3.1 重构 cmd/api/main.go 使用 fx.New 启动应用
- [x] 3.2 引入各个模块到 fx.New
- [x] 3.3 使用 fx.Invoke 注册路由和启动服务器
- [x] 3.4 删除所有手动依赖构造的代码

## 4. 测试与验证

- [x] 4.1 运行 go build ./... 确保代码编译通过
- [x] 4.2 运行 go test ./... 确保所有测试通过
- [ ] 4.3 手动测试核心功能是否正常工作
