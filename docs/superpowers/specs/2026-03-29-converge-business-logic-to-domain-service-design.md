---
name: Converge Business Logic to Domain Service Design
description: 将检验和业务规则相关逻辑收敛到领域服务中，应用服务主要串联流程
type: design
---

# 业务逻辑收敛到领域服务设计

## 概述

将租约相关的业务校验和规则逻辑从应用层（CommandHandler）收敛到领域层（LeaseService），使应用层仅负责流程编排。

## 背景

当前架构中，CommandHandler 既负责流程编排，又包含业务校验逻辑，职责不够清晰。需要按照 DDD 原则重构：
- **领域服务** - 包含所有业务规则和校验逻辑
- **应用服务** - 仅负责流程编排、协调 repository、发布事件

## 设计决策

### 1. 架构风格："胖"领域服务

**决策：** 采用"胖"领域服务模式，LeaseService 可以直接依赖多个 repository。

**理由：**
- 业务逻辑集中管理，易于测试和理解
- 与当前代码库风格一致
- 避免过度分层带来的复杂性

**替代方案：**
- "瘦"领域服务 - 领域服务不依赖 repository，应用层负责加载数据（职责分散）
- 按用例拆分多个服务 - 服务数量过多，过度设计

### 2. 职责划分

**CommandHandler（应用层）职责：**
- 命令类型转换
- 命令基础验证（Command.Validate()）
- 调用领域服务
- 保存聚合根
- 发布领域事件

**LeaseService（领域层）职责：**
- 房间可租性校验
- 删除前置条件校验（账单、押金检查）
- 激活前置条件校验（状态、日期检查）
- 租约创建（含押金创建）协调

### 3. LeaseService 方法设计

```go
type LeaseService struct {
    leaseRepo   leaserepo.LeaseRepository
    depositRepo depositrepo.DepositRepository
    roomRepo    roomrepo.RoomRepository
}

// ValidateRoomForLease 校验房间是否可租
func (s *LeaseService) ValidateRoomForLease(room *roommodel.Room) error

// CreateLease 创建租约（含押金）
func (s *LeaseService) CreateLease(cmd CreateLeaseCommand) (*leasemodel.Lease, *depositmodel.Deposit, error)

// ValidateDelete 校验租约是否可删除
func (s *LeaseService) ValidateDelete(leaseID string) error

// ValidateActivate 校验租约是否可激活
func (s *LeaseService) ValidateActivate(lease *leasemodel.Lease, room *roommodel.Room) error
```

## 数据流程

### 创建租约流程

```
HandleCreateLease(cmd)
  ↓
命令类型转换 + cmd.Validate()
  ↓
leaseService.CreateLease(cmd)
  ├─→ roomRepo.FindByID(roomID)
  ├─→ ValidateRoomForLease(room)
  ├─→ NewLease(...)
  └─→ NewDeposit(...) [如果需要]
  ↓
leaseRepo.Save(lease)
  ↓
depositRepo.Save(deposit) [如果有]
  ↓
发布领域事件
  ↓
返回结果
```

### 删除租约流程

```
HandleDeleteLease(cmd)
  ↓
命令类型转换 + cmd.Validate()
  ↓
leaseService.ValidateDelete(leaseID)
  ├─→ leaseRepo.HasBills(leaseID)
  └─→ leaseRepo.HasDeposit(leaseID)
  ↓
leaseRepo.Delete(leaseID)
  ↓
返回结果
```

### 激活租约流程

```
HandleActivateLease(cmd)
  ↓
命令类型转换 + cmd.Validate()
  ↓
leaseRepo.FindByID(leaseID)
  ↓
roomRepo.FindByID(lease.RoomID)
  ↓
leaseService.ValidateActivate(lease, room)
  ├─→ 检查 lease.Status == Pending
  ├─→ 检查 lease.StartDate <= Now
  └─→ ValidateRoomForLease(room)
  ↓
lease.Activate()
  ↓
leaseRepo.Save(lease)
  ↓
发布领域事件
  ↓
返回结果
```

## 实现任务清单

- [ ] 重构 LeaseService，添加 repository 依赖
- [ ] 实现 CreateLease 方法
- [ ] 实现 ValidateDelete 方法
- [ ] 实现 ValidateActivate 方法
- [ ] 简化 CommandHandler，移除业务逻辑
- [ ] 编译验证
