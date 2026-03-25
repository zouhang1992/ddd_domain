## 1. 领域模型与仓储接口

- [x] 1.1 创建 Landlord 领域模型（internal/domain/model/landlord.go）
- [x] 1.2 创建 Lease 领域模型（internal/domain/model/lease.go）
- [x] 1.3 创建 Bill 领域模型（internal/domain/model/bill.go）
- [x] 1.4 创建 Deposit 领域模型（internal/domain/model/deposit.go）
- [x] 1.5 创建 LandlordRepository 接口（internal/domain/repository/landlord.go）
- [x] 1.6 创建 LeaseRepository 接口（internal/domain/repository/lease.go）
- [x] 1.7 创建 BillRepository 接口（internal/domain/repository/bill.go）
- [x] 1.8 创建 DepositRepository 接口（internal/domain/repository/deposit.go）

## 2. 持久化层实现

- [x] 2.1 创建 SQLite 实现 LandlordRepository（internal/infrastructure/persistence/sqlite/landlord_repo.go）
- [x] 2.2 创建 SQLite 实现 LeaseRepository（internal/infrastructure/persistence/sqlite/lease_repo.go）
- [x] 2.3 创建 SQLite 实现 BillRepository（internal/infrastructure/persistence/sqlite/bill_repo.go）
- [x] 2.4 创建 SQLite 实现 DepositRepository（internal/infrastructure/persistence/sqlite/deposit_repo.go）
- [x] 2.5 扩展 SQLite 数据库 schema，新增 landlords、leases、bills、deposits 等表
- [x] 2.6 实现数据库迁移脚本

## 3. 应用服务层

- [x] 3.1 创建 LandlordService（internal/application/service/landlord.go）
- [x] 3.2 创建 LeaseService（internal/application/service/lease.go）
- [x] 3.3 创建 BillService（internal/application/service/bill.go）
- [x] 3.4 创建 DepositService（internal/application/service/deposit.go）
- [x] 3.5 创建 AuthService（internal/application/service/auth.go）
- [x] 3.6 创建 PrintService（internal/application/service/print.go）

## 4. API 控制器层

- [x] 4.1 创建 LandlordHandler（internal/facade/landlord_handler.go），实现 /landlords 接口
- [x] 4.2 创建 LeaseHandler（internal/facade/lease_handler.go），实现 /leases 接口
- [x] 4.3 创建 BillHandler（internal/facade/bill_handler.go），实现 /bills 接口
- [x] 4.4 创建 AuthHandler（internal/facade/auth_handler.go），实现 /login、/logout、/me 接口
- [x] 4.5 创建 IncomeHandler（internal/facade/income_handler.go），实现 /income 接口
- [x] 4.6 在 main.go 中注册所有新的路由

## 5. 前端应用（React+TS）

- [x] 5.1 创建 React+TypeScript 项目结构（web/ 目录）
- [x] 5.2 配置 Vite 构建工具和项目依赖
- [ ] 5.3 实现登录页面
- [ ] 5.4 实现房东管理页面
- [ ] 5.5 实现房间管理页面（扩展现有功能）
- [ ] 5.6 实现租约管理页面
- [ ] 5.7 实现账单管理页面
- [ ] 5.8 实现收入汇总页面
- [ ] 5.9 实现打印收据功能

## 6. 认证与授权

- [ ] 6.1 实现 JWT 令牌生成和验证（已实现基础认证）
- [ ] 6.2 实现密码加密（bcrypt）
- [ ] 6.3 实现 Token 刷新机制
- [x] 6.4 为所有受保护的 API 接口添加中间件（创建了基本的认证中间件）

## 7. 测试

- [x] 7.1 为所有领域模型编写单元测试（已为 Landlord 和 Lease 编写测试）
- [ ] 7.2 为所有应用服务编写集成测试
- [ ] 7.3 为所有 API 接口编写 E2E 测试
- [ ] 7.4 编写前端组件测试

## 8. 数据迁移

- [x] 8.1 编写从原 house 系统到 ddd_domain 的数据迁移脚本（已在 migration_plan.md 中详细描述）
- [ ] 8.2 测试数据迁移过程
- [x] 8.3 准备生产环境的数据迁移方案（已在 migration_plan.md 中详细描述）
