## 打印服务（Print Service）

### 1. 概述

打印服务是系统中的辅助模块，用于生成和打印各种文档。该模块将采用总线架构，包括命令处理器（Command Handler）、查询处理器（Query Handler）和事件处理器（Event Handler）。

### 2. 命令（Commands）

#### 2.1 打印账单命令

**Command Type**: PrintBillCommand

**Properties**:
- BillID: 账单唯一标识符

**验证规则**:
- BillID 不能为空

**执行逻辑**:
1. 验证命令
2. 查询账单和相关租约信息
3. 生成 RTF 格式的账单内容
4. 发布 BillPrinted 事件

#### 2.2 打印租约命令

**Command Type**: PrintLeaseCommand

**Properties**:
- LeaseID: 租约唯一标识符

**验证规则**:
- LeaseID 不能为空

**执行逻辑**:
1. 验证命令
2. 查询租约信息
3. 生成 RTF 格式的租约内容
4. 发布 LeasePrinted 事件

#### 2.3 打印发票命令

**Command Type**: PrintInvoiceCommand

**Properties**:
- BillID: 账单唯一标识符

**验证规则**:
- BillID 不能为空

**执行逻辑**:
1. 验证命令
2. 查询账单和相关租约信息
3. 生成 JSON 格式的发票内容
4. 发布 InvoicePrinted 事件

### 3. 查询（Queries）

#### 3.1 获取打印作业查询

**Query Type**: GetPrintJobQuery

**Properties**:
- JobID: 打印作业唯一标识符

**验证规则**:
- JobID 不能为空

**执行逻辑**:
1. 验证查询
2. 从打印作业仓库中查找作业
3. 返回打印作业信息

#### 3.2 列出打印作业查询

**Query Type**: ListPrintJobsQuery

**Properties**:
- Status (可选): 打印作业状态，用于过滤
- StartDate (可选): 开始日期，用于过滤
- EndDate (可选): 结束日期，用于过滤

**执行逻辑**:
1. 查询符合条件的打印作业
2. 返回打印作业列表

#### 3.3 获取打印内容查询

**Query Type**: GetPrintContentQuery

**Properties**:
- BillID: 账单唯一标识符

**验证规则**:
- BillID 不能为空

**执行逻辑**:
1. 验证查询
2. 查询账单和相关租约信息
3. 生成 RTF 格式的内容
4. 返回内容

### 4. 事件（Events）

#### 4.1 账单打印事件

**Event Type**: BillPrintedEvent

**Properties**:
- JobID: 打印作业唯一标识符
- BillID: 账单唯一标识符
- PrintedAt: 打印时间
- Content: 打印内容（可选）

#### 4.2 租约打印事件

**Event Type**: LeasePrintedEvent

**Properties**:
- JobID: 打印作业唯一标识符
- LeaseID: 租约唯一标识符
- PrintedAt: 打印时间
- Content: 打印内容（可选）

#### 4.3 发票打印事件

**Event Type**: InvoicePrintedEvent

**Properties**:
- JobID: 打印作业唯一标识符
- BillID: 账单唯一标识符
- PrintedAt: 打印时间
- Content: 打印内容（可选）

#### 4.4 打印作业失败事件

**Event Type**: PrintJobFailedEvent

**Properties**:
- JobID: 打印作业唯一标识符
- BillID (可选): 账单唯一标识符
- LeaseID (可选): 租约唯一标识符
- FailedAt: 失败时间
- Error: 错误信息

### 5. 处理器接口

#### 5.1 打印命令处理器

```go
// PrintCommandHandler 打印命令处理器接口
type PrintCommandHandler interface {
    HandlePrintBill(cmd command.PrintBillCommand) (*model.PrintJob, error)
    HandlePrintLease(cmd command.PrintLeaseCommand) (*model.PrintJob, error)
    HandlePrintInvoice(cmd command.PrintInvoiceCommand) (*model.PrintJob, error)
}
```

#### 5.2 打印查询处理器

```go
// PrintQueryHandler 打印查询处理器接口
type PrintQueryHandler interface {
    HandleGetPrintJob(cmd query.GetPrintJobQuery) (*model.PrintJob, error)
    HandleListPrintJobs(cmd query.ListPrintJobsQuery) ([]*model.PrintJob, error)
    HandleGetPrintContent(cmd query.GetPrintContentQuery) ([]byte, error)
}
```

#### 5.3 打印事件处理器

```go
// PrintEventHandler 打印事件处理器接口
type PrintEventHandler interface {
    HandleBillPrinted(event event.BillPrintedEvent) error
    HandleLeasePrinted(event event.LeasePrintedEvent) error
    HandleInvoicePrinted(event event.InvoicePrintedEvent) error
    HandlePrintJobFailed(event event.PrintJobFailedEvent) error
}
```
