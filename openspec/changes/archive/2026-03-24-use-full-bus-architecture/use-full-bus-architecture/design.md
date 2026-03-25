## Context

当前 ddd_domain 项目已建立基于领域驱动设计（DDD）的架构，并实现了基本的 CQRS 模式和命令总线。然而，应用层仍存在一些问题：
- 存在传统应用服务和 CQRS 应用服务两种实现，造成重复代码
- 没有实现查询总线，查询逻辑分散在 CQRS 应用服务中
- 没有实现事件总线，命令执行后没有发布领域事件
- 存在一些未使用的服务实现

**架构约束：**
- 保持 DDD 架构风格（domain → application → infrastructure → facade）
- 保留现有的命令总线实现
- 使用 SQLite 作为持久化存储
- RESTful API 接口设计

## Goals / Non-Goals

**Goals:**
1. 实现完整的查询总线，统一管理所有查询操作
2. 实现完整的事件总线，支持领域事件发布和订阅
3. 重构 CQRS 应用服务，使其仅作为命令/查询的分发器
4. 移除未使用的传统应用服务实现
5. 简化 HTTP 处理器，使其直接调用总线
6. 实现命令执行后的领域事件发布机制

**Non-Goals:**
1. 不改变现有的领域模型和仓储实现
2. 不引入新的数据库技术
3. 不实现复杂的事件溯源（Event Sourcing）模式
4. 不实现 Saga 模式（在本阶段）

## Decisions

### 1. 查询总线设计

**决定：** 在 `internal/infrastructure/bus/query/` 中实现查询总线，与命令总线保持类似的架构

**理由：**
- 查询操作需要统一管理，便于优化和缓存
- 与命令总线对称的架构设计，便于理解和维护
- 查询可以有自己的中间件链，便于实现日志、性能监控等

**实现：**
```go
// Query 定义查询接口
type Query interface {
    QueryName() string
}

// QueryBus 查询总线
type QueryBus struct {
    handlers map[string]QueryHandler
    middleware []Middleware
}
```

### 2. 事件总线设计

**决定：** 在 `internal/infrastructure/bus/event/` 中实现简单的事件总线

**理由：**
- 命令执行后需要发布领域事件，实现系统解耦
- 简单的同步事件总线足够当前需求
- 便于后续扩展为异步事件处理

**实现：**
```go
// Event 定义事件接口
type Event interface {
    EventName() string
}

// EventBus 事件总线
type EventBus struct {
    subscribers map[string][]EventHandler
}
```

### 3. 移除未使用的服务

**决定：** 移除 `internal/application/service/` 中未被使用的传统应用服务实现

**理由：**
- 这些服务与 CQRS 应用服务功能重复
- 减少维护成本，避免混淆
- CQRS 应用服务已经通过命令总线实现了相同的功能

**需要移除的文件：**
- `internal/application/service/landlord.go`
- `internal/application/service/lease.go`
- `internal/application/service/bill.go`
- `internal/application/service/deposit.go`

**保留的服务：**
- `internal/application/service/auth.go`（身份认证与业务逻辑无关）
- `internal/application/service/print.go`（打印服务是工具类）

### 4. 简化 CQRS 应用服务

**决定：** 重构 CQRS 应用服务，使其仅作为命令/查询的分发器，不包含业务逻辑

**理由：**
- 业务逻辑应该放在命令处理器和查询处理器中
- CQRS 应用服务应该只是一个薄的适配层
- 便于测试和维护

**实现：**
```go
// LandlordApplicationService 房东 CQRS 应用服务
type LandlordApplicationService struct {
    landlordRepo repository.LandlordRepository
    commandBus   *CommandBusService
    queryBus     *QueryBusService
}

// HandleCreateLandlordCommand 处理创建房东命令
func (s *LandlordApplicationService) HandleCreateLandlordCommand(cmd command.CreateLandlordCommand) (*model.Landlord, error) {
    result, err := s.commandBus.SendCommand(cmd)
    if err != nil {
        return nil, err
    }
    return result.(*model.Landlord), nil
}
```

### 5. 简化 HTTP 处理器

**决定：** HTTP 处理器直接与总线交互，不需要通过 CQRS 应用服务

**理由：**
- 减少中间层，提高性能
- 简化代码结构
- HTTP 处理器的职责就是将 HTTP 请求转换为命令/查询

**实现：**
```go
// Create 创建房东
func (h *CQRSLandlordHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name  string `json:"name"`
        Phone string `json:"phone"`
        Note  string `json:"note"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    cmd := command.CreateLandlordCommand{
        Name:  req.Name,
        Phone: req.Phone,
        Note:  req.Note,
    }

    result, err := h.commandBus.SendCommand(cmd)
    // ... 处理响应
}
```

### 6. 事件发布机制

**决定：** 在命令处理器中执行命令后发布领域事件

**理由：**
- 领域事件是领域模型的重要组成部分
- 事件发布是命令执行的副作用
- 便于实现后续的业务逻辑扩展

**实现：**
```go
// HandleCreateLandlord 处理创建房东命令
func (h *LandlordCommandHandler) HandleCreateLandlord(cmd command.Command) (any, error) {
    // ... 创建 landlord
    if err := h.repo.Save(landlord); err != nil {
        return nil, err
    }

    // 发布领域事件
    h.eventBus.Publish(events.NewLandlordCreated(landlord))

    return landlord, nil
}
```

## Risks / Trade-offs

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 重构风险 | 大规模代码改动可能引入新的 bug | 保持接口向后兼容，充分测试 |
| 性能风险 | 增加的抽象层可能影响性能 | 性能关键路径进行基准测试，必要时优化 |
| 回滚风险 | 移除旧服务后难以回滚 | 先标记为 deprecated，保留一段时间后再移除 |
| 测试风险 | 重构可能破坏现有测试 | 更新所有相关测试，保持高测试覆盖率 |
| 学习曲线 | 新的架构模式需要团队适应 | 提供详细的文档和示例代码 |

## Migration Plan

1. **第一阶段：基础设施准备**
   - 实现查询总线
   - 实现事件总线
   - 编写基础设施测试

2. **第二阶段：重构查询处理**
   - 创建查询处理器
   - 重构 CQRS 应用服务使用查询总线
   - 更新 HTTP 处理器

3. **第三阶段：重构命令处理**
   - 增强命令处理器，添加事件发布
   - 确保所有命令处理器已正确实现
   - 更新 HTTP 处理器直接使用命令总线

4. **第四阶段：清理**
   - 标记传统服务为 deprecated
   - 移除未使用的代码
   - 更新文档

## Open Questions

1. **事件是否需要持久化？** 当前设计不要求，后续如果需要可以添加事件存储。
2. **是否需要支持异步事件处理？** 当前使用同步处理，后续可以扩展为异步。
3. **是否需要查询缓存？** 当前不实现，后续可以在查询总线中间件中添加。
