---
name: Print Domain Service Refactor Design
description: 将 print 领域的业务逻辑收敛到领域服务中，应用服务主要串联流程
type: design
---

# Print 领域服务重构设计

## 概述

将 print 领域的打印内容生成和打印作业创建逻辑从应用层（CommandHandler）收敛到领域层（PrintService），使应用层仅负责流程编排。

## 背景

当前 print 领域的 CommandHandler 中包含打印内容生成（RTF 格式）和打印作业创建的逻辑，职责不够清晰。需要按照 DDD 原则重构：
- **领域服务** - 包含打印内容生成、打印作业创建等业务逻辑
- **应用服务** - 仅负责流程编排、命令转换、返回结果

## 设计决策

### 1. 架构风格："胖"领域服务

**决策：** 采用"胖"领域服务模式，PrintService 可以直接依赖 repository。

**理由：**
- 业务逻辑集中管理，易于测试和理解
- 与 lease 领域保持一致的架构风格
- 避免过度分层带来的复杂性

### 2. 职责划分

**CommandHandler（应用层）职责：**
- 命令类型转换
- 命令基础验证（Command.Validate()）
- 调用领域服务
- 返回打印作业 ID

**PrintService（领域层）职责：**
- 打印账单内容生成（RTF 格式）
- 打印租约内容生成（RTF 格式）
- 打印发票内容生成（RTF 格式）
- 打印作业创建（加载账单、加载租约）

### 3. PrintService 方法设计

```go
type PrintService struct {
    billRepo  billrepo.BillRepository
    leaseRepo leaserepo.LeaseRepository
}

// CreateBillPrintJob 创建账单打印作业
func (s *PrintService) CreateBillPrintJob(billID string) (string, error)

// CreateLeasePrintJob 创建租约打印作业
func (s *PrintService) CreateLeasePrintJob(leaseID string) (string, error)

// CreateInvoicePrintJob 创建发票打印作业
func (s *PrintService) CreateInvoicePrintJob(billID string) (string, error)

// GenerateBillContent 生成账单打印内容
func (s *PrintService) GenerateBillContent(bill *billmodel.Bill, lease *leasemodel.Lease) []byte

// GenerateLeaseContent 生成租约打印内容
func (s *PrintService) GenerateLeaseContent(lease *leasemodel.Lease) []byte

// GenerateInvoiceContent 生成发票打印内容
func (s *PrintService) GenerateInvoiceContent(bill *billmodel.Bill, lease *leasemodel.Lease) []byte
```

## 数据流程

### 打印账单流程

```
HandlePrintBill(cmd)
  ↓
命令类型转换 + cmd.Validate()
  ↓
printService.CreateBillPrintJob(billID)
  ├─→ billRepo.FindByID(billID)
  ├─→ leaseRepo.FindByID(bill.LeaseID)
  └─→ 返回 jobID (uuid)
  ↓
返回 jobID
```

### 打印租约流程

```
HandlePrintLease(cmd)
  ↓
命令类型转换 + cmd.Validate()
  ↓
printService.CreateLeasePrintJob(leaseID)
  ├─→ leaseRepo.FindByID(leaseID)
  └─→ 返回 jobID (uuid)
  ↓
返回 jobID
```

### 打印发票流程

```
HandlePrintInvoice(cmd)
  ↓
命令类型转换 + cmd.Validate()
  ↓
printService.CreateInvoicePrintJob(billID)
  ├─→ billRepo.FindByID(billID)
  ├─→ leaseRepo.FindByID(bill.LeaseID)
  └─→ 返回 jobID (uuid)
  ↓
返回 jobID
```

### 获取打印内容流程

```
HandleGetPrintContent(cmd)
  ↓
命令类型转换 + cmd.Validate()
  ↓
billRepo.FindByID(billID)
  ↓
leaseRepo.FindByID(bill.LeaseID)
  ↓
printService.GenerateInvoiceContent(bill, lease)
  ↓
返回内容
```

## 实现任务清单

- [ ] 创建 PrintService（internal/domain/print/service/）
- [ ] 实现 CreateBillPrintJob 方法
- [ ] 实现 CreateLeasePrintJob 方法
- [ ] 实现 CreateInvoicePrintJob 方法
- [ ] 实现 GenerateBillContent 方法
- [ ] 实现 GenerateLeaseContent 方法
- [ ] 实现 GenerateInvoiceContent 方法
- [ ] 更新 print/module.go 提供 PrintService
- [ ] 简化 CommandHandler，调用 PrintService
- [ ] 编译验证
