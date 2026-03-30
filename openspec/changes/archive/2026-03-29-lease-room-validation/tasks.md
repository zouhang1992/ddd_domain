## 1. 新增错误类型

- [x] 1.1 在 `internal/domain/common/errors/` 下新增 `room_errors.go`
- [x] 1.2 定义 `ErrRoomNotFound` 错误
- [x] 1.3 定义 `ErrRoomNotAvailable` 错误

## 2. 创建 LeaseService

- [x] 2.1 在 `internal/application/lease/` 下新增 `service.go`
- [x] 2.2 定义 `LeaseService` 结构体，依赖 `RoomRepository`
- [x] 2.3 实现 `NewLeaseService` 构造函数
- [x] 2.4 实现 `ValidateRoomForLease` 方法

## 3. 更新依赖注入配置

- [x] 3.1 更新 `internal/application/lease/module.go`，添加 `fx.Provide(NewLeaseService)`

## 4. 更新 CommandHandler

- [x] 4.1 在 `CommandHandler` 结构体中添加 `leaseService *LeaseService` 字段
- [x] 4.2 更新 `NewCommandHandler` 构造函数，注入 `LeaseService`

## 5. 在租约操作中添加校验

- [x] 5.1 在 `HandleCreateLease` 中添加房间校验
- [x] 5.2 在 `HandleUpdateLease` 中添加房间校验（变更 roomID 时）- 注：UpdateLeaseCommand 不支持变更 roomID
- [x] 5.3 在 `HandleActivateLease` 中添加房间校验

## 6. 验证

- [x] 6.1 编译验证通过
- [x] 6.2 运行现有测试验证 - 无现有测试
