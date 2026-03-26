---
name: add-execution-logging
change: add-execution-logging
---

## Context

### 当前状态

系统目前使用 Go 标准库的 `log` 包进行简单的日志记录，主要存在以下问题：

1. **基础的日志记录**：仅使用 `log.Println` 和 `log.Fatal` 进行输出
2. **缺少上下文信息**：每条日志仅包含时间戳和消息，缺少上下文数据
3. **无结构化支持**：不支持 JSON 格式等结构化输出
4. **日志级别限制**：没有日志级别控制，所有消息都被输出
5. **难以调试**：生产环境和开发环境日志格式相同，调试困难

### 技术约束

系统使用 Uber FX 依赖注入框架，这为我们提供了良好的组件管理和生命周期管理能力。我们希望新的日志系统能够：
- 与 FX 框架无缝集成
- 提供灵活的配置选项
- 支持结构化日志输出
- 保持高性能和低开销

## Goals / Non-Goals

### Goals

1. **替换标准库 log**：使用 Uber Zap 作为日志库，提供结构化日志功能
2. **统一日志接口**：为整个系统提供一致的日志接口
3. **分级日志**：支持 Debug/Info/Warn/Error 四个日志级别
4. **环境配置**：开发环境使用控制台输出，生产环境使用 JSON 格式输出
5. **FX 集成**：将日志系统集成到 Uber FX 依赖注入系统中
6. **性能优化**：确保日志系统具有低开销和高性能

### Non-Goals

1. **复杂的日志路由**：本次不实现日志转发、分割或归档功能
2. **分布式追踪**：不包含与 OpenTelemetry 等分布式追踪系统的集成
3. **前端日志**：不涉及前端应用的日志记录
4. **操作审计**：这不是操作日志功能（该功能已在 `add-operation-logging` 中实现）

## Decisions

### Decision 1: 选择 Uber Zap 作为日志库

**方案**：使用 Uber Zap 作为主要的日志库

**理由**：
- 高性能：Zap 被设计为非常快速的结构化日志库
- 低分配：使用对象池和避免反射来减少内存分配
- 结构化输出：原生支持 JSON 格式
- 灵活配置：支持开发和生产环境的不同配置
- 已经是间接依赖：从 go.sum 可以看到 zap 已经存在

**代码示例**：
```go
// 开发环境配置
logger, err := zap.NewDevelopment()
if err != nil {
    log.Fatal(err)
}
defer logger.Sync()

// 生产环境配置
logger, err := zap.NewProduction()
if err != nil {
    log.Fatal(err)
}
defer logger.Sync()
```

### Decision 2: 日志配置方式

**方案**：创建统一的日志配置模块，支持开发和生产环境

**配置结构**：
```go
type Config struct {
    Environment string `json:"environment"` // "development" or "production"
    Level       string `json:"level"`       // "debug", "info", "warn", "error"
    OutputPath  string `json:"outputPath"`  // 日志输出路径，默认为 stdout
}
```

**理由**：
- 简单且灵活的配置方式
- 支持环境变量和配置文件
- 符合系统现有的配置风格

### Decision 3: FX 集成方案

**方案**：创建一个日志模块，使用 FX 提供和注入 logger 实例

**代码结构**：
```go
// internal/infrastructure/logging/
package logging

import (
    "go.uber.org/fx"
    "go.uber.org/zap"
)

func Module() fx.Option {
    return fx.Options(
        fx.Provide(NewLogger),
    )
}

func NewLogger() (*zap.Logger, error) {
    // 配置日志器
    cfg := zap.NewProductionConfig()
    cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

    return cfg.Build()
}
```

**理由**：
- 与系统的依赖注入风格一致
- 提供单一的日志器实例
- 方便单元测试时进行 mock

### Decision 4: 日志上下文和字段

**方案**：支持日志字段和上下文信息的添加

**代码示例**：
```go
// 基本使用
logger.Info("User logged in",
    zap.String("user_id", userID),
    zap.String("method", "POST"),
)

// 有错误时
logger.Error("Failed to save user",
    zap.String("user_id", userID),
    zap.Error(err),
)

// 创建带上下文的 logger
reqLogger := logger.With(zap.String("request_id", reqID))
reqLogger.Debug("Processing request",
    zap.String("path", r.URL.Path),
)
```

**理由**：
- 提供丰富的上下文信息
- 便于问题排查和分析
- 支持结构化的字段查询

## Risks / Trade-offs

### Risk 1: 代码迁移复杂度

**风险**：需要修改所有使用标准库 log 的地方

**缓解措施**：
- 逐步迁移，使用适配器模式暂时兼容
- 提供详细的迁移指南
- 在编译时检查未修改的代码

### Risk 2: 学习曲线

**风险**：开发人员需要学习 Zap 的 API

**缓解措施**：
- 提供统一的日志接口
- 创建示例代码和最佳实践文档
- 代码审查时进行指导

### Risk 3: 依赖冲突

**风险**：可能与其他依赖产生版本冲突

**缓解措施**：
- 固定 Zap 的版本
- 确保与现有的 Uber 库版本兼容
- 测试依赖的兼容性

## Migration Plan

### 部署步骤

1. 添加 Zap 依赖到 go.mod
2. 创建日志配置和模块
3. 更新 FX 引导过程
4. 迁移各个组件的日志代码
5. 测试并验证日志功能

### 回滚策略

如果新系统有问题，可通过以下步骤回滚：

1. 移除新的日志模块
2. 恢复对标准库 log 的使用
3. 回滚 go.mod/go.sum 的变更
4. 重新部署

## Open Questions

1. **日志输出位置**：是否需要支持文件输出还是仅保持标准输出？
2. **日志格式**：是否需要支持文本和 JSON 格式的灵活切换？
3. **日志级别配置**：是否需要支持运行时动态调整日志级别？
