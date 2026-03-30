## 1. 重构 LeaseService

- [x] 1.1 在 LeaseService 中添加 LeaseRepository 依赖
- [x] 1.2 在 LeaseService 中添加 DepositRepository 依赖
- [x] 1.3 在 LeaseService 中添加 RoomRepository 依赖
- [x] 1.4 更新 NewLeaseService 构造函数

## 2. 实现 CreateLease 方法

- [x] 2.1 实现 CreateLease 方法（加载房间、校验、创建租约、创建押金）

## 3. 实现 ValidateDelete 方法

- [x] 3.1 实现 ValidateDelete 方法（检查账单、押金）

## 4. 实现 ValidateActivate 方法

- [x] 4.1 实现 ValidateActivate 方法（检查状态、日期、房间）

## 5. 简化 CommandHandler

- [x] 5.1 简化 HandleCreateLease，调用 leaseService.CreateLease
- [x] 5.2 简化 HandleDeleteLease，调用 leaseService.ValidateDelete
- [x] 5.3 简化 HandleActivateLease，调用 leaseService.ValidateActivate

## 6. 验证

- [x] 6.1 编译验证通过
