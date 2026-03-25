## ADDED Requirements

### Requirement: DDD 分层目录结构
项目 SHALL 建立符合领域驱动设计的分层目录结构，包括：门面层、应用服务层、领域层、基础设施层。

#### Scenario: 项目初始化
- **WHEN** 项目初始化完成
- **THEN** 目录结构包含 `cmd/`, `internal/facade/`, `internal/application/`, `internal/domain/`, `internal/infrastructure/` 目录

### Requirement: Go 模块初始化
项目 SHALL 初始化 Go 模块，包含必要的依赖声明。

#### Scenario: 模块初始化
- **WHEN** 项目初始化完成
- **THEN** 根目录存在 `go.mod` 文件，包含模块名称和 Go 版本
- **AND** 包含 SQLite 驱动依赖 `modernc.org/sqlite`

### Requirement: 入口应用
项目 SHALL 提供 HTTP 服务的入口应用程序。

#### Scenario: 启动 HTTP 服务
- **WHEN** 运行 `cmd/api/main.go`
- **THEN** HTTP 服务器在指定端口启动并监听请求
