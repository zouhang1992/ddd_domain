## ADDED Requirements

### Requirement: Frontend pages are responsive to browser width
系统 SHALL 使前端页面根据浏览器宽度自适应布局。

#### Scenario: Dashboard responds to small screens
- **WHEN** 用户在宽度小于 992px 的浏览器中查看仪表盘
- **THEN** 统计卡片布局调整为单列或两列显示，表格允许横向滚动

#### Scenario: Tables respond to small screens
- **WHEN** 用户在宽度小于 1200px 的浏览器中查看表格
- **THEN** 表格显示横向滚动条，关键列优先显示

#### Scenario: Forms respond to small screens
- **WHEN** 用户在宽度小于 768px 的浏览器中查看表单
- **THEN** 表单字段调整为单列布局

### Requirement: Navigation is responsive
系统 SHALL 使导航栏在小屏幕上自适应。

#### Scenario: Menu collapses on small screens
- **WHEN** 用户在宽度小于 768px 的浏览器中查看侧边栏
- **THEN** 侧边栏自动折叠，菜单图标显示

### Requirement: Grid layout adapts to screen sizes
系统 SHALL 使用响应式网格系统。

#### Scenario: Grid responds to different breakpoints
- **WHEN** 用户在不同尺寸的浏览器中查看页面
- **THEN** 网格布局根据 xs/sm/md/lg/xl/xxl 断点自动调整列数
