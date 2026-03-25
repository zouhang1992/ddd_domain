## Context

当前位置管理和房间管理使用传统的服务层架构，而打印服务也是独立的组件。之前我们已经实现了完整的总线架构（command bus、query bus、event bus），并重构了房东、租约、账单管理。为了保持系统架构的一致性，我们需要将位置管理、房间管理和打印服务调整为相同的模式。

## Goals / Non-Goals

**Goals:**

1. 将位置管理重构为使用总线架构
2. 将房间管理重构为使用总线架构
3. 将打印服务重构为使用总线架构
4. 将 LocationHandler 和 RoomHandler 改为实现接口，而非函数传参
5. 保持与之前实现的架构风格一致

**Non-Goals:**

1. 不修改 Location 和 Room 的领域模型核心逻辑
2. 不改变打印服务的功能（如 RTFS 生成）
3. 不引入新的外部依赖

## Decisions

### 决定 1: 使用与完整总线架构相同的模式

**理由**：保持架构一致性是首要目标。之前的实现已经验证了这种模式的有效性，包括：

- Command bus 用于处理写操作
- Query bus 用于处理读操作
- Event bus 用于事件发布和订阅

### 决定 2: 为每个 handler 创建接口

**理由**：使用接口而非函数传参可以提供更好的类型安全性和代码组织。例如：

```go
// LocationCommandHandler 位置命令处理器接口
type LocationCommandHandler interface {
    HandleCreateLocation(cmd command.CreateLocationCommand) (*model.Location, error)
    HandleUpdateLocation(cmd command.UpdateLocationCommand) (*model.Location, error)
    HandleDeleteLocation(cmd command.DeleteLocationCommand) error
}

// LocationQueryHandler 位置查询处理器接口
type LocationQueryHandler interface {
    HandleGetLocation(cmd query.GetLocationQuery) (*model.Location, error)
    HandleListLocations(cmd query.ListLocationsQuery) ([]*model.Location, error)
}
```

### 决定 3: 位置和房间管理的事件建模

**理由**：虽然当前需求未要求事件驱动，但我们应该为未来的功能做好准备。可能的事件包括：

- LocationCreatedEvent
- LocationUpdatedEvent
- LocationDeletedEvent
- RoomCreatedEvent
- RoomUpdatedEvent
- RoomDeletedEvent
- PrintJobCreatedEvent

### 决定 4: 打印服务的命令和查询

**理由**：打印服务可以视为具有以下操作：

```go
// PrintBillCommand 打印账单命令
type PrintBillCommand struct {
    BillID string
}

// PrintBillReceiptQuery 查询并打印账单收据
type PrintBillReceiptQuery struct {
    BillID string
}
```

### 决定 5: 保持与 main.go 中现有模式一致

**理由**：在 main.go 中，我们已经有了初始化总线和注册处理器的模式，我们应该遵循：

```go
// 初始化事件总线
eventBus := event.NewBus()

// 初始化处理器
locationCmdHandler := handler.NewLocationCommandHandler(locationRepo, eventBus)
locationQueryHandler := handler.NewLocationQueryHandler(locationRepo)

// 注册到总线
commandBus.Register("create_location", command.HandlerFunc(locationCmdHandler.HandleCreateLocation))
queryBus.Register("get_location", query.HandlerFunc(locationQueryHandler.HandleGetLocation))
```

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| 与现有代码冲突 | 在重构前进行充分的测试 |
| handler 接口变更导致编译错误 | 一次性完成所有相关变更并运行构建 |
| 性能下降 | 使用与之前相同的总线实现，已优化过性能 |
| 迁移期间的功能不可用 | 确保新旧架构可以共存一段时间，或使用功能标志 |
