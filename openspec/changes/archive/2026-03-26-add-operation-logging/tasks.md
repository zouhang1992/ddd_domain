---
name: tasks
change: add-operation-logging
description: 操作日志功能实现任务清单
---

## 1. 前端实现

### 1.1 OperationLogModal 通用组件
- [x] 1.1.1 创建 OperationLogModal.tsx 组件
- [x] 1.1.2 实现组件基础结构（Modal + Table）
- [x] 1.1.3 添加操作日志API调用方法
- [x] 1.1.4 实现查询参数处理（domainType, aggregateID）
- [x] 1.1.5 实现分页功能
- [x] 1.1.6 添加组件样式和交互优化

### 1.2 页面集成
- [x] 1.2.1 在 Landlords.tsx 页面添加操作日志按钮和Modal集成
- [x] 1.2.2 在 Leases.tsx 页面添加操作日志按钮和Modal集成
- [x] 1.2.3 在 Bills.tsx 页面添加操作日志按钮和Modal集成
- [x] 1.2.4 在 Locations.tsx 页面添加操作日志按钮和Modal集成
- [x] 1.2.5 在 Rooms.tsx 页面添加操作日志按钮和Modal集成
- [x] 1.2.6 在 Print.tsx 页面添加操作日志按钮和Modal集成

## 2. 后端实现

### 2.1 打印操作日志事件处理
- [x] 2.1.1 在 operation_log_handler.go 中添加打印事件处理逻辑
- [x] 2.1.2 处理 BillPrinted 事件
- [x] 2.1.3 处理 LeasePrinted 事件
- [x] 2.1.4 处理 InvoicePrinted 事件
- [x] 2.1.5 处理 PrintJobFailed 事件

### 2.2 打印查询API完善
- [x] 2.2.1 在 internal/application/query/print.go 中补充打印操作查询
- [x] 2.2.2 在 internal/domain/repository/print.go 中补充仓库接口
- [x] 2.2.3 在 internal/infrastructure/persistence/sqlite/print_repo.go 中补充SQLite实现

### 2.3 操作日志查询API优化
- [x] 2.3.1 确保操作日志查询API与现有系统风格一致
- [x] 2.3.2 验证查询参数处理的正确性
- [x] 2.3.3 优化操作日志查询的性能

## 3. 验证测试

### 3.1 功能测试
- [x] 3.1.1 测试操作日志模态框显示和关闭
- [x] 3.1.2 测试操作日志查询功能
- [x] 3.1.3 测试操作日志分页功能
- [x] 3.1.4 测试打印操作日志记录
- [x] 3.1.5 测试各页面操作日志集成

### 3.2 兼容性测试
- [x] 3.2.1 验证所有现有功能正常工作
- [x] 3.2.2 验证操作日志功能与现有系统兼容性
