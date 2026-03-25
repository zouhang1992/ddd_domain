## 1. 文件删除

- [x] 1.1 删除 `internal/application/service/location.go`
- [x] 1.2 删除 `internal/application/service/room.go`
- [x] 1.3 删除 `internal/facade/location_handler.go`
- [x] 1.4 删除 `internal/facade/room_handler.go`

## 2. 代码清理

- [x] 2.1 更新 `cmd/api/main.go` 以移除对传统服务的引用
- [x] 2.2 删除 `internal/application/service` 包中不再需要的导入
- [x] 2.3 验证所有代码编译成功

## 3. 测试

- [x] 3.1 运行所有测试以确保没有破坏功能
- [x] 3.2 手动测试核心功能（位置管理、房间管理）是否正常工作
