## Why

当前系统的依赖关系管理较为复杂，主要体现在：

1. 在 `cmd/api/main.go` 中，大量的组件初始化代码分散在各个函数中
2. 依赖关系手动传递，容易出错且难以维护
3. 每个组件的创建都需要显式地构造所有依赖，导致代码冗余
4. 缺乏统一的依赖注入容器，使得测试和代码复用变得困难

为了解决这些问题，我们将引入 UberFx 依赖注入容器。UberFx 是一个轻量级、快速且易用的 Go 依赖注入库，它通过反射和结构体标签自动管理依赖关系，并提供了清晰的生命周期管理。

## What Changes

**BREAKING CHANGE**

### 修改的文件：
- `go.mod` - 添加 uber/fx 依赖
- `cmd/api/main.go` - 重构为使用 UberFx 管理所有组件的依赖关系

### 更新的依赖：
- 添加 `go.uber.org/fx v1.21.0+` 到 go.mod

### Capabilities Modified
- `dependency-injection` - 替换手动依赖管理为 UberFx 自动注入
- `main-entry-point` - 重构 cmd/api/main.go 为 Fx 应用

## Migration

1. 所有组件初始化将通过 Fx 提供和注入，不再需要手动创建和传递
2. 组件间的依赖关系通过构造函数的参数类型自动匹配
3. 可以通过 Fx 的功能实现更好的测试和代码复用

## Impact

### Affected Code
- `cmd/api/main.go` - 重写以使用 UberFx
- `internal/infrastructure/persistence/sqlite/` - 添加 Fx 提供函数
- `internal/application/command/handler/` - 添加 Fx 提供函数
- `internal/application/query/handler/` - 添加 Fx 提供函数
- `internal/facade/` - 添加 Fx 提供函数
