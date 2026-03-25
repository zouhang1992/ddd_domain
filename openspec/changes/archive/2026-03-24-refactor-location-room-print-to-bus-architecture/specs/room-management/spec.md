## 房间管理（Room Management）

### 1. 概述

房间管理是系统中的核心模块，用于管理房屋租赁业务中的房间信息。该模块将采用总线架构，包括命令处理器（Command Handler）、查询处理器（Query Handler）和事件处理器（Event Handler）。

### 2. 命令（Commands）

#### 2.1 创建房间命令

**Command Type**: CreateRoomCommand

**Properties**:
- LocationID: 位置标识符
- RoomNumber: 房间号
- Tags: 标签数组

**验证规则**:
- LocationID 不能为空
- RoomNumber 不能为空
- Tags 可以为空

**执行逻辑**:
1. 验证命令
2. 创建 Room 领域模型
3. 保存到房间仓库
4. 发布 RoomCreated 事件

#### 2.2 更新房间命令

**Command Type**: UpdateRoomCommand

**Properties**:
- ID: 房间唯一标识符
- LocationID: 位置标识符
- RoomNumber: 房间号
- Tags: 标签数组

**验证规则**:
- ID 不能为空
- 至少需要修改一个字段

**执行逻辑**:
1. 验证命令
2. 从仓库中查找房间
3. 更新房间信息
4. 保存到仓库
5. 发布 RoomUpdated 事件

#### 2.3 删除房间命令

**Command Type**: DeleteRoomCommand

**Properties**:
- ID: 房间唯一标识符

**验证规则**:
- ID 不能为空
- 该房间不能有相关联的租约

**执行逻辑**:
1. 验证命令
2. 从仓库中删除房间
3. 发布 RoomDeleted 事件

### 3. 查询（Queries）

#### 3.1 获取房间查询

**Query Type**: GetRoomQuery

**Properties**:
- ID: 房间唯一标识符

**验证规则**:
- ID 不能为空

**执行逻辑**:
1. 验证查询
2. 从仓库中查找房间
3. 返回房间信息

#### 3.2 列出所有房间查询

**Query Type**: ListRoomsQuery

**Properties**:
- LocationID (可选): 位置标识符，用于过滤
- Tags (可选): 标签数组，用于过滤

**执行逻辑**:
1. 查询所有符合条件的房间
2. 返回房间列表

#### 3.3 按位置列出房间查询

**Query Type**: ListRoomsByLocationQuery

**Properties**:
- LocationID: 位置标识符

**验证规则**:
- LocationID 不能为空

**执行逻辑**:
1. 查询指定位置的所有房间
2. 返回房间列表

### 4. 事件（Events）

#### 4.1 房间创建事件

**Event Type**: RoomCreatedEvent

**Properties**:
- ID: 房间唯一标识符
- LocationID: 位置标识符
- RoomNumber: 房间号
- Tags: 标签数组
- CreatedAt: 创建时间

#### 4.2 房间更新事件

**Event Type**: RoomUpdatedEvent

**Properties**:
- ID: 房间唯一标识符
- LocationID: 位置标识符
- RoomNumber: 房间号
- Tags: 标签数组
- UpdatedAt: 更新时间

#### 4.3 房间删除事件

**Event Type**: RoomDeletedEvent

**Properties**:
- ID: 房间唯一标识符
- DeletedAt: 删除时间

### 5. 处理器接口

#### 5.1 房间命令处理器

```go
// RoomCommandHandler 房间命令处理器接口
type RoomCommandHandler interface {
    HandleCreateRoom(cmd command.CreateRoomCommand) (*model.Room, error)
    HandleUpdateRoom(cmd command.UpdateRoomCommand) (*model.Room, error)
    HandleDeleteRoom(cmd command.DeleteRoomCommand) error
}
```

#### 5.2 房间查询处理器

```go
// RoomQueryHandler 房间查询处理器接口
type RoomQueryHandler interface {
    HandleGetRoom(cmd query.GetRoomQuery) (*model.Room, error)
    HandleListRooms(cmd query.ListRoomsQuery) ([]*model.Room, error)
    HandleListRoomsByLocation(cmd query.ListRoomsByLocationQuery) ([]*model.Room, error)
}
```

#### 5.3 房间事件处理器

```go
// RoomEventHandler 房间事件处理器接口
type RoomEventHandler interface {
    HandleRoomCreated(event event.RoomCreatedEvent) error
    HandleRoomUpdated(event event.RoomUpdatedEvent) error
    HandleRoomDeleted(event event.RoomDeletedEvent) error
}
```
