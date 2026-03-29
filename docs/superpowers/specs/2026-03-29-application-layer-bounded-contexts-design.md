---
name: 应用层按领域划分重构设计
description: 将应用层从按类型组织重构为按限界上下文组织
---

# 应用层按领域划分重构设计

## 概述

本设计旨在将应用层从当前的按类型（command、query、event、service）组织重构为按限界上下文组织，与domain层保持一致的架构风格。

## 目标

1. **架构一致性** - 应用层与domain层采用相同的限界上下文划分
2. **高内聚低耦合** - 相关的command、query、handler放在同一个限界上下文内
3. **更好的可维护性** - 修改一个业务领域时只需关注一个目录
4. **清晰的模块边界** - 每个限界上下文有自己的module，便于依赖注入

## 当前架构

### 当前应用层结构
```
internal/application/
├── command/
│   ├── base.go
│   ├── landlord.go
│   ├── lease.go
│   ├── bill.go
│   ├── room.go
│   ├── location.go
│   ├── print.go
│   └── handler/
│       ├── module.go
│       ├── landlord_command_handler.go
│       ├── lease_command_handler.go
│       ├── bill_command_handler.go
│       ├── room_command_handler.go
│       ├── location_command_handler.go
│       └── print_command_handler.go
├── query/
│   ├── base.go
│   ├── landlord.go
│   ├── lease.go
│   ├── bill.go
│   ├── room.go
│   ├── location.go
│   ├── print.go
│   ├── operation_log.go
│   └── handler/
│       ├── module.go
│       ├── landlord_query_handler.go
│       ├── lease_query_handler.go
│       ├── bill_query_handler.go
│       ├── room_query_handler.go
│       ├── location_query_handler.go
│       ├── print_query_handler.go
│       └── operation_log_query_handler.go
├── event/handler/
│   ├── module.go
│   ├── lease_room_event_handler.go
│   └── operation_log_handler.go
├── service/
│   ├── auth.go
│   └── print.go
└── (无总module)
```

### 当前Domain层结构（已重构）
```
internal/domain/
├── landlord/
│   ├── model/
│   ├── repository/
│   └── events/
├── lease/
│   ├── model/
│   ├── repository/
│   └── events/
├── bill/
│   ├── model/
│   ├── repository/
│   └── events/
├── room/
│   ├── model/
│   ├── repository/
│   └── events/
├── location/
│   ├── model/
│   ├── repository/
│   └── events/
├── deposit/
│   ├── model/
│   ├── repository/
│   └── events/
└── common/
    ├── model/
    ├── events/
    └── errors/
```

## 新的架构设计

### 目标结构
```
internal/application/
├── landlord/              # 房东限界上下文
│   ├── command.go         # 命令定义
│   ├── query.go           # 查询定义
│   ├── handler.go         # CommandHandler + QueryHandler
│   └── module.go          # Fx module
├── lease/                 # 租约限界上下文
│   ├── command.go
│   ├── query.go
│   ├── handler.go
│   └── module.go
├── bill/                  # 账单限界上下文
│   ├── command.go
│   ├── query.go
│   ├── handler.go
│   └── module.go
├── room/                  # 房间限界上下文
│   ├── command.go
│   ├── query.go
│   ├── handler.go
│   └── module.go
├── location/              # 位置限界上下文
│   ├── command.go
│   ├── query.go
│   ├── handler.go
│   └── module.go
├── deposit/               # 押金限界上下文
│   ├── command.go
│   ├── query.go
│   ├── handler.go
│   └── module.go
├── common/                # 共享组件
│   ├── service/
│   │   ├── print.go
│   │   └── auth.go
│   ├── event/
│   │   └── handler/
│   │       ├── lease_room_event_handler.go
│   │       └── module.go
│   └── module.go
├── command/               # (临时，逐步迁移)
├── query/                 # (临时，逐步迁移)
├── event/                 # (临时，逐步迁移)
├── service/               # (临时，逐步迁移)
└── module.go              # 总的应用层Module
```

## 限界上下文内部结构

### 1. Command定义
每个限界上下文的command.go包含该领域的所有命令定义。

示例（landlord/command.go）：
```go
package landlord

import "errors"

// CreateLandlordCommand 创建房东命令
type CreateLandlordCommand struct {
    Name  string
    Phone string
    Note  string
}

// CommandName 实现 Command 接口
func (c CreateLandlordCommand) CommandName() string {
    return "create_landlord"
}

// Validate 验证命令
func (c CreateLandlordCommand) Validate() error {
    if c.Name == "" {
        return errors.New("name is required")
    }
    return nil
}

// UpdateLandlordCommand 更新房东命令
type UpdateLandlordCommand struct {
    ID    string
    Name  string
    Phone string
    Note  string
}

// CommandName 实现 Command 接口
func (c UpdateLandlordCommand) CommandName() string {
    return "update_landlord"
}

// Validate 验证命令
func (c UpdateLandlordCommand) Validate() error {
    if c.ID == "" {
        return errors.New("id is required")
    }
    if c.Name == "" {
        return errors.New("name is required")
    }
    return nil
}

// DeleteLandlordCommand 删除房东命令
type DeleteLandlordCommand struct {
    ID string
}

// CommandName 实现 Command 接口
func (c DeleteLandlordCommand) CommandName() string {
    return "delete_landlord"
}

// Validate 验证命令
func (c DeleteLandlordCommand) Validate() error {
    if c.ID == "" {
        return errors.New("id is required")
    }
    return nil
}
```

### 2. Query定义
每个限界上下文的query.go包含该领域的所有查询定义。

示例（landlord/query.go）：
```go
package landlord

// GetLandlordQuery 获取房东查询
type GetLandlordQuery struct {
    ID string
}

// QueryName 实现 Query 接口
func (q GetLandlordQuery) QueryName() string {
    return "get_landlord"
}

// ListLandlordsQuery 查询房东列表
type ListLandlordsQuery struct {
    Name     string
    Phone    string
    Offset   int
    Limit    int
}

// QueryName 实现 Query 接口
func (q ListLandlordsQuery) QueryName() string {
    return "list_landlords"
}

// LandlordsQueryResult 房东列表查询结果
type LandlordsQueryResult struct {
    Items []interface{} `json:"items"`
    Total int           `json:"total"`
    Page  int           `json:"page"`
    Limit int           `json:"limit"`
}
```

### 3. Handler
每个限界上下文的handler.go包含CommandHandler和QueryHandler。

示例（landlord/handler.go）：
```go
package landlord

import (
    "github.com/google/uuid"
    "github.com/zouhang1992/ddd_domain/internal/application/command"
    "github.com/zouhang1992/ddd_domain/internal/application/query"
    landlordmodel "github.com/zouhang1992/ddd_domain/internal/domain/landlord/model"
    landlordrepo "github.com/zouhang1992/ddd_domain/internal/domain/landlord/repository"
    domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
    "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/event"
)

// CommandHandler 房东命令处理器
type CommandHandler struct {
    repo     landlordrepo.LandlordRepository
    eventBus *event.Bus
}

// NewCommandHandler 创建房东命令处理器
func NewCommandHandler(repo landlordrepo.LandlordRepository, eventBus *event.Bus) *CommandHandler {
    return &CommandHandler{repo: repo, eventBus: eventBus}
}

// HandleCreateLandlord 处理创建房东命令
func (h *CommandHandler) HandleCreateLandlord(cmd command.Command) (any, error) {
    createCmd, ok := cmd.(CreateLandlordCommand)
    if !ok {
        return nil, domerrors.ErrInvalidCommand
    }

    if err := createCmd.Validate(); err != nil {
        return nil, err
    }

    id := uuid.NewString()
    landlord := landlordmodel.NewLandlord(id, createCmd.Name, createCmd.Phone, createCmd.Note)
    if err := h.repo.Save(landlord); err != nil {
        return nil, err
    }

    // Publish events from aggregate
    if h.eventBus != nil {
        for _, evt := range landlord.Events() {
            h.eventBus.PublishAsync(evt)
        }
        landlord.ClearEvents()
    }

    return landlord, nil
}

// HandleUpdateLandlord 处理更新房东命令
func (h *CommandHandler) HandleUpdateLandlord(cmd command.Command) (any, error) {
    updateCmd, ok := cmd.(UpdateLandlordCommand)
    if !ok {
        return nil, domerrors.ErrInvalidCommand
    }

    if err := updateCmd.Validate(); err != nil {
        return nil, err
    }

    landlord, err := h.repo.FindByID(updateCmd.ID)
    if err != nil {
        return nil, err
    }
    if landlord == nil {
        return nil, domerrors.ErrNotFound
    }

    landlord.Update(updateCmd.Name, updateCmd.Phone, updateCmd.Note)
    if err := h.repo.Save(landlord); err != nil {
        return nil, err
    }

    // Publish events from aggregate
    if h.eventBus != nil {
        for _, evt := range landlord.Events() {
            h.eventBus.PublishAsync(evt)
        }
        landlord.ClearEvents()
    }

    return landlord, nil
}

// HandleDeleteLandlord 处理删除房东命令
func (h *CommandHandler) HandleDeleteLandlord(cmd command.Command) (any, error) {
    deleteCmd, ok := cmd.(DeleteLandlordCommand)
    if !ok {
        return nil, domerrors.ErrInvalidCommand
    }

    if err := deleteCmd.Validate(); err != nil {
        return nil, err
    }

    if err := h.repo.Delete(deleteCmd.ID); err != nil {
        return nil, err
    }

    return nil, nil
}

// QueryHandler 房东查询处理器
type QueryHandler struct {
    repo landlordrepo.LandlordRepository
}

// NewQueryHandler 创建房东查询处理器
func NewQueryHandler(repo landlordrepo.LandlordRepository) *QueryHandler {
    return &QueryHandler{repo: repo}
}

// HandleGetLandlord 处理获取房东查询
func (h *QueryHandler) HandleGetLandlord(q query.Query) (any, error) {
    getQuery, ok := q.(GetLandlordQuery)
    if !ok {
        return nil, domerrors.ErrInvalidCommand
    }

    landlord, err := h.repo.FindByID(getQuery.ID)
    if err != nil {
        return nil, err
    }
    if landlord == nil {
        return nil, domerrors.ErrNotFound
    }

    return landlord, nil
}

// HandleListLandlords 处理查询房东列表
func (h *QueryHandler) HandleListLandlords(q query.Query) (any, error) {
    listQuery, ok := q.(ListLandlordsQuery)
    if !ok {
        return nil, domerrors.ErrInvalidCommand
    }

    landlords, err := h.repo.FindAll()
    if err != nil {
        return nil, err
    }

    // 简单过滤
    var items []interface{}
    for _, l := range landlords {
        items = append(items, l)
    }

    limit := listQuery.Limit
    if limit <= 0 {
        limit = 10
    }

    offset := listQuery.Offset
    if offset < 0 {
        offset = 0
    }

    // 简单分页
    if offset > len(items) {
        items = []interface{}{}
    } else if offset+limit > len(items) {
        items = items[offset:]
    } else {
        items = items[offset : offset+limit]
    }

    page := 1
    if offset > 0 && limit > 0 {
        page = (offset / limit) + 1
    }

    return &LandlordsQueryResult{
        Items: items,
        Total: len(landlords),
        Page:  page,
        Limit: limit,
    }, nil
}
```

### 4. Module
每个限界上下文的module.go定义该上下文的Fx依赖注入模块。

示例（landlord/module.go）：
```go
package landlord

import "go.uber.org/fx"

// Module provides landlord application components
var Module = fx.Options(
    fx.Provide(NewCommandHandler),
    fx.Provide(NewQueryHandler),
)
```

## Common目录结构

### 共享服务
```
internal/application/common/
├── service/
│   ├── print.go       # 打印服务
│   └── auth.go        # 认证服务
├── event/
│   └── handler/
│       ├── lease_room_event_handler.go
│       └── module.go
└── module.go
```

### common/module.go
```go
package common

import "go.uber.org/fx"

// Module provides common application components
var Module = fx.Options(
    eventhandler.Module,
    fx.Provide(NewPrintService),
    fx.Provide(NewAuthService),
)
```

## 总的Application Module

### internal/application/module.go
```go
package application

import "go.uber.org/fx"

// Module provides all application components
var Module = fx.Options(
    landlord.Module,
    lease.Module,
    bill.Module,
    room.Module,
    location.Module,
    deposit.Module,
    common.Module,
)
```

## 重构策略

### 渐进式重构（推荐）

采用一个领域一个领域地重构，降低风险。

#### 阶段1：Landlord
1. 创建 `internal/application/landlord/` 目录
2. 创建 `command.go`、`query.go`、`handler.go`、`module.go`
3. 从旧的 `command/landlord.go` 和 `query/landlord.go` 移动代码
4. 从旧的 `command/handler/landlord_command_handler.go` 和 `query/handler/landlord_query_handler.go` 移动代码
5. 更新所有导入路径
6. 更新旧的 `command/handler/module.go` 和 `query/handler/module.go` 以引用新的handler
7. 运行编译验证
8. 运行测试（如果有）

#### 阶段2：Lease
同阶段1，处理租约限界上下文

#### 阶段3：Bill
同阶段1，处理账单限界上下文

#### 阶段4：Room
同阶段1，处理房间限界上下文

#### 阶段5：Location
同阶段1，处理位置限界上下文

#### 阶段6：Deposit
同阶段1，处理押金限界上下文

#### 阶段7：Common
1. 创建 `internal/application/common/` 目录
2. 移动 `service/print.go` 和 `service/auth.go`
3. 移动 `event/handler/lease_room_event_handler.go`
4. 更新导入路径
5. 创建 `common/module.go`

#### 阶段8：清理
1. 创建 `internal/application/module.go` 总的模块
2. 更新 `cmd/api/main.go` 使用新的 `application.Module`
3. 删除旧的 `command/`、`query/`、`event/`、`service/` 目录
4. 最终编译验证
5. 提交更改

## 与其他层的集成

### Facade层
Facade层保持当前结构不变，只需要更新导入路径以引用新的按领域组织的command和query类型。

### Main.go
更新 `cmd/api/main.go` 以使用新的 `application.Module`：

```go
// 旧的
fx.Options(
    handler.Module,        // command handlers
    eventhandler.Module,   // event handlers
    queryhandler.Module,   // query handlers
    // ...
),

// 新的
fx.Options(
    application.Module,    // 总的应用层模块
    // ...
),
```

## 优势

1. **架构一致性** - 应用层与domain层采用相同的限界上下文划分
2. **高内聚** - 相关的command、query、handler放在同一个目录
3. **低耦合** - 限界上下文之间通过明确的接口交互
4. **更好的可维护性** - 修改一个业务领域时只需关注一个目录
5. **清晰的模块边界** - 每个限界上下文有自己的module
6. **渐进式迁移** - 可以一个领域一个领域地重构，风险低

## 注意事项

1. **向后兼容** - 重构期间保持旧代码可用，逐步迁移
2. **导入路径更新** - 需要更新所有引用旧路径的文件
3. **Facade层不变** - Facade层保持当前结构，只更新导入
4. **测试更新** - 如果有测试，需要同步更新测试的导入路径
