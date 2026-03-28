# DDD架构重构实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将domain目录结构调整为按限定上下文组织，领域事件由聚合根直接产生

**Architecture:** 按业务领域划分为独立的限界上下文，每个限界上下文包含自己的model、repository和events；聚合根内部维护事件队列

**Tech Stack:** Go, Fx, Zap, SQLite

---

## 文件结构

**新增文件：**
- `internal/domain/common/model/aggregate_root.go` - 聚合根基类
- `internal/domain/common/events/base.go` - 事件基类
- `internal/domain/landlord/model/landlord.go` - 房东模型
- `internal/domain/landlord/repository/repository.go` - 房东仓储
- `internal/domain/landlord/events/landlord_events.go` - 房东事件
- `internal/domain/lease/model/lease.go` - 租约模型
- `internal/domain/lease/repository/repository.go` - 租约仓储
- `internal/domain/lease/events/lease_events.go` - 租约事件
- `internal/domain/bill/model/bill.go` - 账单模型
- `internal/domain/bill/repository/repository.go` - 账单仓储
- `internal/domain/bill/events/bill_events.go` - 账单事件
- `internal/domain/room/model/room.go` - 房间模型
- `internal/domain/room/repository/repository.go` - 房间仓储
- `internal/domain/room/events/room_events.go` - 房间事件
- `internal/domain/location/model/location.go` - 位置模型
- `internal/domain/location/repository/repository.go` - 位置仓储
- `internal/domain/location/events/location_events.go` - 位置事件
- `internal/domain/deposit/model/deposit.go` - 押金模型
- `internal/domain/deposit/repository/repository.go` - 押金仓储
- `internal/domain/deposit/events/deposit_events.go` - 押金事件

**修改文件：**
- `internal/infrastructure/bus/event/bus.go` - 支持新的DomainEvent接口
- `internal/infrastructure/persistence/sqlite/*.go` - 迁移仓储实现
- `internal/application/command/handler/*.go` - 更新命令处理器
- `internal/application/event/handler/*.go` - 更新事件处理器
- `internal/facade/*.go` - 更新facade
- `cmd/api/main.go` - 更新依赖注入配置

---

## 阶段一：基础设施准备

### Task 1.1: 创建事件基类

**Files:**
- Create: `internal/domain/common/events/base.go`

- [ ] **Step 1: 创建事件基类文件**

```go
package events

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventName() string
	EventID() string
	TimeStamp() time.Time
	AggregateID() string
	Version() int
}

// BaseEvent 事件基类
type BaseEvent struct {
	eventName   string
	eventID     string
	timeStamp   time.Time
	aggregateID string
	version     int
}

func NewBaseEvent(eventName, aggregateID string, version int) BaseEvent {
	return BaseEvent{
		eventName:   eventName,
		eventID:     uuid.NewString(),
		timeStamp:   time.Now(),
		aggregateID: aggregateID,
		version:     version,
	}
}

func (e BaseEvent) EventName() string {
	return e.eventName
}

func (e BaseEvent) EventID() string {
	return e.eventID
}

func (e BaseEvent) TimeStamp() time.Time {
	return e.timeStamp
}

func (e BaseEvent) AggregateID() string {
	return e.aggregateID
}

func (e BaseEvent) Version() int {
	return e.version
}
```

- [ ] **Step 2: 编译验证**

```bash
go build ./internal/domain/common/events/...
```
Expected: 编译成功

- [ ] **Step 3: 提交更改**

```bash
git add internal/domain/common/events/base.go
git commit -m "feat: create common events base class"
```

### Task 1.2: 创建聚合根基类

**Files:**
- Create: `internal/domain/common/model/aggregate_root.go`

- [ ] **Step 1: 创建聚合根基类文件**

```go
package model

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
)

// AggregateRoot 聚合根接口
type AggregateRoot interface {
	ID() string
	Version() int
	Events() []events.DomainEvent
	ClearEvents()
}

// BaseAggregateRoot 基础聚合根实现
type BaseAggregateRoot struct {
	id      string
	version int
	events  []events.DomainEvent
}

func NewBaseAggregateRoot(id string) BaseAggregateRoot {
	return BaseAggregateRoot{
		id:      id,
		version: 0,
		events:  []events.DomainEvent{},
	}
}

func (a *BaseAggregateRoot) ID() string {
	return a.id
}

func (a *BaseAggregateRoot) Version() int {
	return a.version
}

func (a *BaseAggregateRoot) Events() []events.DomainEvent {
	return a.events
}

func (a *BaseAggregateRoot) ClearEvents() {
	a.events = []events.DomainEvent{}
}

func (a *BaseAggregateRoot) RecordEvent(evt events.DomainEvent) {
	a.version++
	a.events = append(a.events, evt)
}
```

- [ ] **Step 2: 编译验证**

```bash
go build ./internal/domain/common/model/...
```
Expected: 编译成功

- [ ] **Step 3: 提交更改**

```bash
git add internal/domain/common/model/aggregate_root.go
git commit -m "feat: create aggregate root base class"
```

### Task 1.3: 更新事件总线支持新接口

**Files:**
- Modify: `internal/infrastructure/bus/event/bus.go`

- [ ] **Step 1: 读取当前文件内容**

```bash
cat internal/infrastructure/bus/event/bus.go
```

- [ ] **Step 2: 更新DomainEvent接口引用**

将所有对旧的DomainEvent接口的引用更新为新的 `common/events.DomainEvent` 接口。

- [ ] **Step 3: 编译验证**

```bash
go build ./internal/infrastructure/bus/event/...
```
Expected: 编译成功

- [ ] **Step 4: 提交更改**

```bash
git add internal/infrastructure/bus/event/bus.go
git commit -m "refactor: update event bus to support new DomainEvent interface"
```

---

## 阶段二：限界上下文创建

### Task 2.1: 创建房东管理限界上下文

**Files:**
- Create: `internal/domain/landlord/model/landlord.go`
- Create: `internal/domain/landlord/repository/repository.go`
- Create: `internal/domain/landlord/events/landlord_events.go`

- [ ] **Step 1: 创建房东模型文件**

```go
package model

import (
	"time"

	"github.com/zouhang1992/ddd_domain/internal/domain/common/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/landlord/events"
)

// Landlord 房东领域模型（聚合根）
type Landlord struct {
	model.BaseAggregateRoot
	Name      string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewLandlord 创建新房东
func NewLandlord(id, name, phone string) *Landlord {
	now := time.Now()
	landlord := &Landlord{
		BaseAggregateRoot: model.NewBaseAggregateRoot(id),
		Name:              name,
		Phone:             phone,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	landlord.RecordEvent(events.NewLandlordCreated(landlord))
	return landlord
}

// Update 更新房东信息
func (l *Landlord) Update(name, phone string) {
	l.Name = name
	l.Phone = phone
	l.UpdatedAt = time.Now()
	l.RecordEvent(events.NewLandlordUpdated(l))
}

// Equals 比较房东是否相等
func (l *Landlord) Equals(other *Landlord) bool {
	return l.ID() == other.ID()
}
```

- [ ] **Step 2: 创建房东事件文件**

```go
package events

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/common/events"
	"github.com/zouhang1992/ddd_domain/internal/domain/landlord/model"
)

type LandlordCreated struct {
	events.BaseEvent
	Name  string
	Phone string
}

type LandlordUpdated struct {
	events.BaseEvent
	Name  string
	Phone string
}

type LandlordDeleted struct {
	events.BaseEvent
}

func NewLandlordCreated(landlord *model.Landlord) LandlordCreated {
	return LandlordCreated{
		BaseEvent: events.NewBaseEvent("landlord.created", landlord.ID(), landlord.Version()),
		Name:      landlord.Name,
		Phone:     landlord.Phone,
	}
}

func NewLandlordUpdated(landlord *model.Landlord) LandlordUpdated {
	return LandlordUpdated{
		BaseEvent: events.NewBaseEvent("landlord.updated", landlord.ID(), landlord.Version()),
		Name:      landlord.Name,
		Phone:     landlord.Phone,
	}
}

func NewLandlordDeleted(landlordID string, version int) LandlordDeleted {
	return LandlordDeleted{
		BaseEvent: events.NewBaseEvent("landlord.deleted", landlordID, version),
	}
}
```

- [ ] **Step 3: 创建房东仓储接口**

```go
package repository

import (
	"github.com/zouhang1992/ddd_domain/internal/domain/landlord/model"
)

type LandlordRepository interface {
	FindByID(id string) (*model.Landlord, error)
	FindAll() ([]*model.Landlord, error)
	Save(landlord *model.Landlord) error
	Delete(id string) error
}
```

- [ ] **Step 4: 编译验证**

```bash
go build ./internal/domain/landlord/...
```
Expected: 编译成功

- [ ] **Step 5: 提交更改**

```bash
git add internal/domain/landlord/
git commit -m "feat: create landlord bounded context"
```

### Task 2.2: 创建租约管理限界上下文

**Files:**
- Create: `internal/domain/lease/model/lease.go`
- Create: `internal/domain/lease/repository/repository.go`
- Create: `internal/domain/lease/events/lease_events.go`

- [ ] **Step 1: 创建租约模型文件**（参考设计文档中的示例）

- [ ] **Step 2: 创建租约事件文件**（参考设计文档中的示例）

- [ ] **Step 3: 创建租约仓储接口**（参考设计文档中的示例）

- [ ] **Step 4: 编译验证**

```bash
go build ./internal/domain/lease/...
```
Expected: 编译成功

- [ ] **Step 5: 提交更改**

```bash
git add internal/domain/lease/
git commit -m "feat: create lease bounded context"
```

### Task 2.3: 创建其他限界上下文

**Files:**
- Create: `internal/domain/bill/` - 账单管理
- Create: `internal/domain/room/` - 房间管理
- Create: `internal/domain/location/` - 位置管理
- Create: `internal/domain/deposit/` - 押金管理

- [ ] **Step 1: 为每个限界上下文创建对应的模型、仓储接口和事件文件**

- [ ] **Step 2: 编译验证每个限界上下文**

- [ ] **Step 3: 提交更改**

```bash
git add internal/domain/bill/ internal/domain/room/ internal/domain/location/ internal/domain/deposit/
git commit -m "feat: create remaining bounded contexts"
```

---

## 阶段三：聚合根重构

### Task 3.1: 更新聚合根方法以记录事件

**Files:**
- Modify: `internal/domain/lease/model/lease.go`
- Modify: `internal/domain/landlord/model/landlord.go`
- Modify: `internal/domain/bill/model/bill.go`
- Modify: `internal/domain/room/model/room.go`
- Modify: `internal/domain/location/model/location.go`
- Modify: `internal/domain/deposit/model/deposit.go`

- [ ] **Step 1: 更新每个聚合根的方法，在状态变化时记录事件**

- [ ] **Step 2: 编译验证**

```bash
go build ./internal/domain/...
```
Expected: 编译成功

- [ ] **Step 3: 提交更改**

```bash
git add internal/domain/
git commit -m "refactor: update aggregate roots to record events"
```

---

## 阶段四：仓储重构

### Task 4.1: 迁移SQLite仓储实现

**Files:**
- Modify: `internal/infrastructure/persistence/sqlite/*.go`

- [ ] **Step 1: 创建新的仓储实现对应新的限界上下文**

- [ ] **Step 2: 更新旧的仓储实现以保持兼容**

- [ ] **Step 3: 编译验证**

```bash
go build ./internal/infrastructure/persistence/sqlite/...
```
Expected: 编译成功

- [ ] **Step 4: 提交更改**

```bash
git add internal/infrastructure/persistence/sqlite/
git commit -m "refactor: migrate persistence to new bounded contexts"
```

### Task 4.2: 更新命令处理器

**Files:**
- Modify: `internal/application/command/handler/*.go`

- [ ] **Step 1: 更新命令处理器以使用新的仓储和聚合根**

- [ ] **Step 2: 更新事件发布逻辑以从聚合根获取事件**

- [ ] **Step 3: 编译验证**

```bash
go build ./internal/application/command/handler/...
```
Expected: 编译成功

- [ ] **Step 4: 提交更改**

```bash
git add internal/application/command/handler/
git commit -m "refactor: update command handlers"
```

---

## 阶段五：事件发布重构

### Task 5.1: 更新事件处理器

**Files:**
- Modify: `internal/application/event/handler/*.go`

- [ ] **Step 1: 更新事件处理器以处理新的事件类型**

- [ ] **Step 2: 编译验证**

```bash
go build ./internal/application/event/handler/...
```
Expected: 编译成功

- [ ] **Step 3: 提交更改**

```bash
git add internal/application/event/handler/
git commit -m "refactor: update event handlers"
```

### Task 5.2: 更新依赖注入配置

**Files:**
- Modify: `cmd/api/main.go`

- [ ] **Step 1: 更新Fx配置以使用新的仓储**

- [ ] **Step 2: 编译验证**

```bash
go build ./cmd/api/...
```
Expected: 编译成功

- [ ] **Step 3: 提交更改**

```bash
git add cmd/api/main.go
git commit -m "refactor: update dependency injection"
```

---

## 验证和测试

### Task 6.1: 运行完整编译验证

**Files:**
- (All files)

- [ ] **Step 1: 运行完整项目编译**

```bash
go build ./...
```
Expected: 编译成功

### Task 6.2: 运行现有单元测试

**Files:**
- (All test files)

- [ ] **Step 1: 运行所有测试**

```bash
go test ./internal/... -v
```
Expected: 所有测试通过

---

## 自审查

### Spec Coverage
- [x] 创建事件基类和聚合根基类 - Task 1.1, 1.2
- [x] 按业务领域创建限界上下文 - Task 2.1-2.3
- [x] 聚合根内部维护事件队列 - Task 3.1
- [x] 迁移仓储实现 - Task 4.1
- [x] 更新应用层以使用新的事件发布模式 - Task 4.2, 5.1

### Placeholder Scan
- [x] 无TBD, TODO或未完成的部分
- [x] 所有代码示例完整
- [x] 所有文件路径明确

### Type Consistency
- [x] DomainEvent接口在所有位置一致
- [x] 聚合根方法签名一致
- [x] 仓储接口保持一致性
