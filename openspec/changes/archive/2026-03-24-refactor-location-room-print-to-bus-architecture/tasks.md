## 1. 位置管理 - Command/Query/Event

- [x] 1.1 创建位置管理的 Command（CreateLocationCommand, UpdateLocationCommand, DeleteLocationCommand）
- [x] 1.2 创建位置管理的 Query（GetLocationQuery, ListLocationsQuery）
- [x] 1.3 创建位置管理的 Event（LocationCreatedEvent, LocationUpdatedEvent, LocationDeletedEvent）

## 2. 位置管理 - Handler 实现

- [x] 2.1 实现 LocationCommandHandler 接口
- [x] 2.2 实现 LocationQueryHandler 接口
- [ ] 2.3 定义 LocationEventHandler 接口（可选实现）

## 3. 房间管理 - Command/Query/Event

- [x] 3.1 创建房间管理的 Command（CreateRoomCommand, UpdateRoomCommand, DeleteRoomCommand）
- [x] 3.2 创建房间管理的 Query（GetRoomQuery, ListRoomsQuery, ListRoomsByLocationQuery）
- [x] 3.3 创建房间管理的 Event（RoomCreatedEvent, RoomUpdatedEvent, RoomDeletedEvent）

## 4. 房间管理 - Handler 实现

- [x] 4.1 实现 RoomCommandHandler 接口
- [x] 4.2 实现 RoomQueryHandler 接口
- [ ] 4.3 定义 RoomEventHandler 接口（可选实现）

## 5. 打印服务 - Command/Query/Event

- [x] 5.1 创建打印服务的 Command（PrintBillCommand, PrintLeaseCommand, PrintInvoiceCommand）
- [x] 5.2 创建打印服务的 Query（GetPrintJobQuery, ListPrintJobsQuery, GetPrintContentQuery）
- [x] 5.3 创建打印服务的 Event（BillPrintedEvent, LeasePrintedEvent, InvoicePrintedEvent, PrintJobFailedEvent）

## 6. 打印服务 - Handler 实现

- [x] 6.1 实现 PrintCommandHandler 接口
- [x] 6.2 实现 PrintQueryHandler 接口
- [ ] 6.3 定义 PrintEventHandler 接口（可选实现）

## 7. HTTP Handler 重构

- [x] 7.1 重构 LocationHandler 为接口实现，直接使用总线
- [x] 7.2 重构 RoomHandler 为接口实现，直接使用总线
- [x] 7.3 创建 PrintHandler 实现，使用总线

## 8. 集成和清理

- [x] 8.1 更新 main.go 初始化总线和处理器
- [x] 8.2 注册所有 command/query handler 到总线
- [x] 8.3 删除旧的 LocationService 和 RoomService
- [x] 8.4 删除旧的传统 handler（如果存在）
- [x] 8.5 运行构建确保编译通过
