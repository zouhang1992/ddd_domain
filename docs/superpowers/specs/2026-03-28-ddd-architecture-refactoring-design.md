---
name: DDD架构重构设计
description: 将domain目录结构调整为按限定上下文组织，领域事件由聚合根直接产生
---

# DDD架构重构设计

## 概述

本设计旨在重构现有架构，使其更加符合领域驱动设计（DDD）的原则。主要改进包括：
1. 领域事件由聚合根直接产生，通过内部事件队列管理
2. domain目录结构调整为按业务领域限定上下文组织
3. 每个限定上下文包含自己的model、repository和events

## 当前架构问题

### 问题分析

1. **领域事件产生方式不符合DDD原则**：
   - 当前通过独立的 `NewXxxEvent` 函数创建事件
   - 事件产生与领域对象状态变化分离
   - 事件创建逻辑分散

2. **目录结构不合理**：
   - 所有领域模型集中在 `model/` 目录
   - 所有仓储接口集中在 `repository/` 目录
   - 没有按限界上下文组织
   - 领域边界不清晰

3. **模块耦合度高**：
   - 不同业务领域的模型混在一起
   - 难以独立开发和部署

## 新的架构设计

### 1. 目录结构

按业务领域划分为独立的限界上下文：

```
internal/domain/
├── landlord/           # 房东管理限界上下文
│   ├── model/
│   │   └── landlord.go
│   ├── repository/
│   │   └── repository.go
│   └── events/
│       └── landlord_events.go
├── lease/             # 租约管理限界上下文
│   ├── model/
│   │   └── lease.go
│   ├── repository/
│   │   └── repository.go
│   └── events/
│       └── lease_events.go
├── bill/              # 账单管理限界上下文
│   ├── model/
│   │   └── bill.go
│   ├── repository/
│   │   └── repository.go
│   └── events/
│       └── bill_events.go
├── room/              # 房间管理限界上下文
│   ├── model/
│   │   └── room.go
│   ├── repository/
│   │   └── repository.go
│   └── events/
│       └── room_events.go
├── location/          # 位置管理限界上下文
│   ├── model/
│   │   └── location.go
│   ├── repository/
│   │   └── repository.go
│   └── events/
│       └── location_events.go
├── deposit/           # 押金管理限界上下文
│   ├── model/
│   │   └── deposit.go
│   ├── repository/
│   │   └── repository.go
│   └── events/
│       └── deposit_events.go
└── common/            # 共享基础设施
    ├── model/
    │   └── aggregate_root.go  # 聚合根基类
    └── events/
        └── base.go             # 事件基类
```

### 2. 聚合根设计

#### 聚合根基类

所有聚合根都继承自 `BaseAggregateRoot`：

```go
// internal/domain/common/model/aggregate_root.go
package model

import (
    "time"
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

#### Lease 聚合根示例

```go
// internal/domain/lease/model/lease.go
package model

import (
    "time"
    "github.com/zouhang1992/ddd_domain/internal/domain/common/model"
    "github.com/zouhang1992/ddd_domain/internal/domain/lease/events"
)

// LeaseStatus 租约状态
type LeaseStatus string

const (
    LeaseStatusPending  LeaseStatus = "pending"
    LeaseStatusActive   LeaseStatus = "active"
    LeaseStatusExpired  LeaseStatus = "expired"
    LeaseStatusCheckout LeaseStatus = "checkout"
)

// Lease 租约领域模型（聚合根）
type Lease struct {
    model.BaseAggregateRoot
    RoomID         string
    LandlordID     string
    TenantName     string
    TenantPhone    string
    StartDate      time.Time
    EndDate        time.Time
    RentAmount     int64
    DepositAmount  int64
    Status         LeaseStatus
    Note           string
    LastChargeAt   *time.Time
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

// NewLease 创建新租约
func NewLease(id, roomID, landlordID, tenantName, tenantPhone string,
    startDate, endDate time.Time, rentAmount, depositAmount int64, note string) *Lease {
    now := time.Now()
    lease := &Lease{
        BaseAggregateRoot: model.NewBaseAggregateRoot(id),
        RoomID:         roomID,
        LandlordID:     landlordID,
        TenantName:     tenantName,
        TenantPhone:    tenantPhone,
        StartDate:      startDate,
        EndDate:        endDate,
        RentAmount:     rentAmount,
        DepositAmount:  depositAmount,
        Status:         LeaseStatusPending,
        Note:           note,
        CreatedAt:      now,
        UpdatedAt:      now,
    }

    // 记录创建事件
    lease.RecordEvent(events.NewLeaseCreated(lease))
    return lease
}

// Activate 激活租约
func (l *Lease) Activate() {
    l.Status = LeaseStatusActive
    l.UpdatedAt = time.Now()
    l.RecordEvent(events.NewLeaseActivated(l))
}

// Checkout 退租
func (l *Lease) Checkout() {
    l.Status = LeaseStatusCheckout
    l.UpdatedAt = time.Now()
    l.RecordEvent(events.NewLeaseCheckout(l))
}

// Expire 标记租约为过期状态
func (l *Lease) Expire() {
    l.Status = LeaseStatusExpired
    l.UpdatedAt = time.Now()
    l.RecordEvent(events.NewLeaseExpired(l))
}
```

### 3. 领域事件设计

#### 事件基类

```go
// internal/domain/common/events/base.go
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

#### Lease 领域事件

```go
// internal/domain/lease/events/lease_events.go
package events

import (
    "time"
    "github.com/zouhang1992/ddd_domain/internal/domain/common/events"
    "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
)

type LeaseCreated struct {
    events.BaseEvent
    RoomID        string
    LandlordID    string
    TenantName    string
}

type LeaseActivated struct {
    events.BaseEvent
    RoomID        string
}

type LeaseCheckout struct {
    events.BaseEvent
    RoomID        string
}

type LeaseExpired struct {
    events.BaseEvent
    RoomID        string
}

type LeaseRenewed struct {
    events.BaseEvent
    NewEndDate    string
}

type LeaseDeleted struct {
    events.BaseEvent
}

func NewLeaseCreated(lease *model.Lease) LeaseCreated {
    return LeaseCreated{
        BaseEvent: events.NewBaseEvent("lease.created", lease.ID(), lease.Version()),
        RoomID:     lease.RoomID,
        LandlordID: lease.LandlordID,
        TenantName: lease.TenantName,
    }
}

func NewLeaseActivated(lease *model.Lease) LeaseActivated {
    return LeaseActivated{
        BaseEvent: events.NewBaseEvent("lease.activated", lease.ID(), lease.Version()),
        RoomID:     lease.RoomID,
    }
}

func NewLeaseCheckout(lease *model.Lease) LeaseCheckout {
    return LeaseCheckout{
        BaseEvent: events.NewBaseEvent("lease.checkout", lease.ID(), lease.Version()),
        RoomID:     lease.RoomID,
    }
}

func NewLeaseExpired(lease *model.Lease) LeaseExpired {
    return LeaseExpired{
        BaseEvent: events.NewBaseEvent("lease.expired", lease.ID(), lease.Version()),
        RoomID:     lease.RoomID,
    }
}

func NewLeaseRenewed(lease *model.Lease) LeaseRenewed {
    return LeaseRenewed{
        BaseEvent: events.NewBaseEvent("lease.renewed", lease.ID(), lease.Version()),
        NewEndDate: lease.EndDate.Format("2006-01-02"),
    }
}

func NewLeaseDeleted(leaseID string, version int) LeaseDeleted {
    return LeaseDeleted{
        BaseEvent: events.NewBaseEvent("lease.deleted", leaseID, version),
    }
}
```

### 4. 事件发布模式

#### 应用层服务示例

```go
func (s *LeaseService) ActivateLease(leaseID string) error {
    // 从仓储获取聚合根
    lease, err := s.leaseRepo.FindByID(leaseID)
    if err != nil {
        return err
    }

    // 执行领域操作（会自动记录事件）
    lease.Activate()

    // 保存聚合根到仓储
    if err := s.leaseRepo.Save(lease); err != nil {
        return err
    }

    // 获取并发布领域事件
    for _, evt := range lease.Events() {
        s.eventBus.PublishAsync(evt)
    }

    // 清除已发布的事件
    lease.ClearEvents()

    return nil
}
```

### 5. 仓储接口重构

每个限界上下文有自己的仓储接口：

```go
// internal/domain/lease/repository/repository.go
package repository

import (
    "time"
    "github.com/zouhang1992/ddd_domain/internal/domain/lease/model"
)

type LeaseRepository interface {
    FindByID(id string) (*model.Lease, error)
    FindAll() ([]*model.Lease, error)
    Save(lease *model.Lease) error
    Delete(id string) error
    FindActiveLeasesExpiringBefore(time time.Time) ([]*model.Lease, error)
    HasBills(leaseID string) (bool, error)
    HasDeposit(leaseID string) (bool, error)
}
```

## 迁移策略

### 阶段一：基础设施准备
- 创建 `common/` 包和基类
- 重构事件总线以支持新的 DomainEvent 接口

### 阶段二：限界上下文创建
- 按业务领域创建新的目录结构
- 迁移模型和事件到对应限界上下文

### 阶段三：聚合根重构
- 为每个聚合根集成事件队列
- 修改领域方法以记录事件

### 阶段四：仓储重构
- 迁移仓储接口和实现
- 调整应用层服务以使用新的仓储

### 阶段五：事件发布重构
- 调整应用层服务以使用新的事件发布模式
- 更新事件处理器以处理新的事件类型

## 优势

### 1. 符合 DDD 原则
- 领域事件由聚合根产生，符合领域事件的定义
- 清晰的限界上下文边界

### 2. 更好的代码组织
- 按业务领域组织，更容易理解和维护
- 高内聚，低耦合

### 3. 更好的可扩展性
- 新的业务领域可以独立添加
- 便于团队独立开发

### 4. 更好的可测试性
- 聚合根测试时可以验证事件产生
- 事件处理可以独立测试
