---
name: Lease Room Validation Design
description: 为租约订单创建时的房间校验功能设计
type: design
---

# 租约房间校验功能设计

## 概述

在创建、更新和激活租约时，添加房间存在性和可租状态校验。

## 背景

当前系统在创建租约时没有校验房间是否存在，也没有校验房间是否处于可租状态。这可能导致：
- 为不存在的房间创建租约
- 为已出租或维修中的房间创建租约

## 需求

### 功能需求

1. **房间存在性校验：校验传入的 roomID 必须对应一个存在的房间
2. **房间状态校验：** 房间必须处于 `available` 状态才能创建/激活租约
3. **校验范围：** 在所有租约相关操作都需要校验（创建、更新、激活）

### 非功能需求

1. 使用专门的错误类型区分不同的校验失败场景
2. 校验逻辑作为租约领域服务的一部分
3. 保持现有代码结构不变

## 设计

### 架构

在 lease 限界上下文中创建 `LeaseService`（租约领域服务），负责处理租约相关的业务逻辑，包括房间校验。

### 组件设计

#### 1. 新增错误类型

**文件：** `internal/domain/common/errors/room_errors.go`

```go
package errors

import "errors"

var (
    // ErrRoomNotFound 房间不存在
    ErrRoomNotFound = errors.New("room not found")
    // ErrRoomNotAvailable 房间不可租
    ErrRoomNotAvailable = errors.New("room not available for lease")
)
```

#### 2. LeaseService 领域服务

**文件：** `internal/application/lease/service.go`

```go
package lease

import (
    roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
    roommodel "github.com/zouhang1992/ddd_domain/internal/domain/room/model"
    domerrors "github.com/zouhang1992/ddd_domain/internal/domain/common/errors"
)

// LeaseService 租约领域服务
type LeaseService struct {
    roomRepo roomrepo.RoomRepository
}

// NewLeaseService 创建租约领域服务
func NewLeaseService(roomRepo roomrepo.RoomRepository) *LeaseService {
    return &LeaseService{roomRepo: roomRepo}
}

// ValidateRoomForLease 校验房间是否可用于租约
func (s *LeaseService) ValidateRoomForLease(roomID string) error {
    // 1. 检查房间是否存在
    room, err := s.roomRepo.FindByID(roomID)
    if err != nil {
        return err
    }
    if room == nil {
        return domerrors.ErrRoomNotFound
    }

    // 2. 检查房间状态是否可租
    if room.Status != roommodel.RoomStatusAvailable {
        return domerrors.ErrRoomNotAvailable
    }

    return nil
}
```

#### 3. 更新 CommandHandler

**文件：** `internal/application/lease/handler.go`

修改 `CommandHandler` 结构体，注入 `LeaseService`：

```go
// CommandHandler 租约命令处理器
type CommandHandler struct {
    repo         leaserepo.LeaseRepository
    depositRepo  depositrepo.DepositRepository
    billRepo     billrepo.BillRepository
    eventBus     *event.Bus
    leaseService *LeaseService
}
```

在以下方法中添加房间校验：

- `HandleCreateLease` - 创建租约时
- `HandleUpdateLease` - 更新租约时（如果变更了 roomID）
- `HandleActivateLease` - 激活租约时

#### 4. 更新 Module

**文件：** `internal/application/lease/module.go`

添加 `LeaseService` 的依赖注入：

```go
// Module provides lease application components
var Module = fx.Options(
    fx.Provide(NewCommandHandler),
    fx.Provide(NewQueryHandler),
    fx.Provide(NewLeaseRoomEventHandler),
    fx.Provide(NewLeaseService),
)
```

### 数据流程

#### 创建租约流程

```
创建租约请求
    ↓
CommandHandler.HandleCreateLease
    ↓
LeaseService.ValidateRoomForLease(roomID)
    ├─→ RoomRepository.FindByID(roomID)
    │   └─→ 不存在 → 返回 ErrRoomNotFound
    └─→ 检查房间状态
        └─→ 不是 available → 返回 ErrRoomNotAvailable
    ↓
校验通过，继续创建租约
    ↓
保存租约
    ↓
返回结果
```

### 需要校验的操作

| 操作 | 校验时机 | 说明 |
|------|----------|------|
| 创建租约 | HandleCreateLease | 必须校验 |
| 更新租约 | HandleUpdateLease | 如果变更了 roomID 则校验 |
| 激活租约 | HandleActivateLease | 必须校验 |
| 续租 | HandleRenewLease | 不需要（不涉及 room 变更） |
| 退租 | HandleCheckoutLease | 不需要 |
| 删除租约 | HandleDeleteLease | 不需要 |

## 实现任务清单

- [ ] 新增房间相关错误类型
- [ ] 创建 LeaseService 领域服务
- [ ] 更新 lease/module.go 添加 LeaseService
- [ ] 更新 CommandHandler 注入 LeaseService
- [ ] 在 HandleCreateLease 中添加校验
- [ ] 在 HandleUpdateLease 中添加校验
- [ ] 在 HandleActivateLease 中添加校验
- [ ] 编译验证
- [ ] 测试验证
