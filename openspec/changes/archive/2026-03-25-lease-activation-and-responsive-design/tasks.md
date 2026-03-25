## 1. 响应式设计优化

- [x] 1.1 优化 Dashboard.tsx 的响应式布局，调整统计卡片在不同屏幕下的显示
- [x] 1.2 优化 Leases.tsx 表格的响应式设计，添加横向滚动和列宽调整
- [x] 1.3 优化 Bills.tsx 表格的响应式设计，确保在小屏幕下的可用性
- [x] 1.4 优化其他页面（Locations、Rooms、Landlords、Print、Income）的响应式设计
- [x] 1.5 优化 Layout.tsx 的侧边栏在小屏幕下的显示和折叠

## 2. 租约列表增强

- [x] 2.1 在 Leases.tsx 表格中添加押金状态列
- [x] 2.2 在 Leases.tsx 表格中添加押金金额列并格式化显示
- [x] 2.3 在 Leases.tsx 操作列添加"打印合同"按钮
- [x] 2.4 在 leaseApi.ts 中添加下载合同的 API 方法
- [x] 2.5 优化租约状态标签的颜色和显示

## 3. 租约生效功能

- [x] 3.1 在 domain/model/lease.go 中添加租约生效的领域事件
- [x] 3.2 在 application/command/lease.go 中添加 ActivateLease 命令
- [x] 3.3 在 application/command/handler/lease_command_handler.go 中添加激活租约的处理逻辑
- [x] 3.4 在 internal/facade/cqrs_lease_handler.go 中添加激活租约的 API 端点
- [x] 3.5 在 web/src/api/lease.ts 中添加激活租约的 API 方法
- [x] 3.6 在 Leases.tsx 中添加"生效"按钮和相关逻辑

## 4. 押金收入/支出追踪

- [x] 4.1 在 internal/facade/income_handler.go 中更新收入汇总逻辑，包含押金的收取和退还
- [x] 4.2 在 Income.tsx 中添加押金收入/支出的分类显示
- [x] 4.3 在 Income.tsx 中更新收入统计，区分租金、水电、押金等不同类型
- [x] 4.4 确保押金收取显示为正收入，押金退还显示为负收入
- [x] 4.5 优化 Income.tsx 的响应式设计和数据显示

## 5. 测试和验证

- [x] 5.1 测试响应式设计在不同浏览器宽度下的效果
- [x] 5.2 测试租约列表押金信息的显示和打印功能
- [x] 5.3 测试租约生效功能的业务逻辑
- [x] 5.4 测试押金收入/支出追踪功能
- [x] 5.5 进行全面的回归测试确保没有破坏现有功能
