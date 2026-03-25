## 位置管理（Location Management）

### 1. 概述

位置管理是系统中的基础模块，用于管理房屋租赁业务中的位置信息。该模块将采用总线架构，包括命令处理器（Command Handler）、查询处理器（Query Handler）和事件处理器（Event Handler）。

### 2. 命令（Commands）

#### 2.1 创建位置命令

**Command Type**: CreateLocationCommand

**Properties**:
- ShortName: 位置简称
- Detail: 详细地址

**验证规则**:
- ShortName 不能为空
- Detail 可以为空

**执行逻辑**:
1. 验证命令
2. 创建 Location 领域模型
3. 保存到位置仓库
4. 发布 LocationCreated 事件

#### 2.2 更新位置命令

**Command Type**: UpdateLocationCommand

**Properties**:
- ID: 位置唯一标识符
- ShortName: 位置简称
- Detail: 详细地址

**验证规则**:
- ID 不能为空
- 至少需要修改一个字段

**执行逻辑**:
1. 验证命令
2. 从仓库中查找位置
3. 更新位置信息
4. 保存到仓库
5. 发布 LocationUpdated 事件

#### 2.3 删除位置命令

**Command Type**: DeleteLocationCommand

**Properties**:
- ID: 位置唯一标识符

**验证规则**:
- ID 不能为空
- 该位置不能有相关联的房间

**执行逻辑**:
1. 验证命令
2. 检查该位置是否有相关联的房间
3. 从仓库中删除位置
4. 发布 LocationDeleted 事件

### 3. 查询（Queries）

#### 3.1 获取位置查询

**Query Type**: GetLocationQuery

**Properties**:
- ID: 位置唯一标识符

**验证规则**:
- ID 不能为空

**执行逻辑**:
1. 验证查询
2. 从仓库中查找位置
3. 返回位置信息

#### 3.2 列出所有位置查询

**Query Type**: ListLocationsQuery

**Properties**:
- 无特定属性

**执行逻辑**:
1. 查询所有位置
2. 返回位置列表

### 4. 事件（Events）

#### 4.1 位置创建事件

**Event Type**: LocationCreatedEvent

**Properties**:
- ID: 位置唯一标识符
- ShortName: 位置简称
- Detail: 详细地址
- CreatedAt: 创建时间

#### 4.2 位置更新事件

**Event Type**: LocationUpdatedEvent

**Properties**:
- ID: 位置唯一标识符
- ShortName: 位置简称
- Detail: 详细地址
- UpdatedAt: 更新时间

#### 4.3 位置删除事件

**Event Type**: LocationDeletedEvent

**Properties**:
- ID: 位置唯一标识符
- DeletedAt: 删除时间

### 5. 处理器接口

#### 5.1 位置命令处理器

```go
// LocationCommandHandler 位置命令处理器接口
type LocationCommandHandler interface {
    HandleCreateLocation(cmd command.CreateLocationCommand) (*model.Location, error)
    HandleUpdateLocation(cmd command.UpdateLocationCommand) (*model.Location, error)
    HandleDeleteLocation(cmd command.DeleteLocationCommand) error
}
```

#### 5.2 位置查询处理器

```go
// LocationQueryHandler 位置查询处理器接口
type LocationQueryHandler interface {
    HandleGetLocation(cmd query.GetLocationQuery) (*model.Location, error)
    HandleListLocations(cmd query.ListLocationsQuery) ([]*model.Location, error)
}
```

#### 5.3 位置事件处理器

```go
// LocationEventHandler 位置事件处理器接口
type LocationEventHandler interface {
    HandleLocationCreated(event event.LocationCreatedEvent) error
    HandleLocationUpdated(event event.LocationUpdatedEvent) error
    HandleLocationDeleted(event event.LocationDeletedEvent) error
}
```
