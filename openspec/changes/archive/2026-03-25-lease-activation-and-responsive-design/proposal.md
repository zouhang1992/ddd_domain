## Why

当前租约管理系统存在以下问题：
1. 前端页面在不同浏览器宽度下布局不佳，缺乏响应式设计
2. 租约列表缺少押金信息展示和直接打印功能
3. 租约无法生效，需要添加租约生效的业务能力
4. 收入汇总中缺少押金的收入（收取时）和支出（退租时）信息

## What Changes

- **新增响应式设计**：前端页面调整为根据浏览器宽度自适应布局，优化表格、表单和卡片在不同屏幕尺寸下的显示
- **租约列表增强**：在租约列表页面添加押金信息展示，并提供直接打印租约合同的功能
- **租约生效能力**：添加租约生效的业务逻辑，支持将待生效租约转换为生效状态
- **收入汇总增强**：在收入汇总中补充押金的收入（收取时）和支出（退租时）信息显示

## Capabilities

### New Capabilities

- `responsive-design`: 响应式设计能力，优化前端页面在不同屏幕尺寸下的布局
- `lease-activation`: 租约生效业务能力，支持激活待生效的租约
- `lease-list-enhancements`: 租约列表功能增强，添加押金信息和打印能力

### Modified Capabilities

- `lease-management`: 现有的租约管理能力，将进行功能增强
- `react-frontend`: 现有的 React 前端应用，将进行响应式设计优化
- `income-reporting`: 现有的收入报告能力，将增强以包含押金的收入和支出追踪

## Impact

- **前端代码**：在 `web/src/pages/Leases.tsx`、`web/src/pages/Dashboard.tsx`、`web/src/pages/Bills.tsx`、`web/src/pages/Income.tsx` 等页面中添加响应式设计优化和押金收入/支出显示
- **后端 API**：新增租约生效相关的 API 接口
- **业务逻辑**：在 `internal/domain/model/lease.go`、`internal/application/command/handler/lease_command_handler.go` 中添加租约生效和押金追踪的业务逻辑
- **数据模型**：可能需要更新租约的状态字段和相关表结构
- **收入查询**：在 `internal/facade/income_handler.go` 中更新收入汇总逻辑，包含押金的收取和退还
