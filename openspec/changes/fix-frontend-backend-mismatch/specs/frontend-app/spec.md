## MODIFIED Requirements

### Requirement: Frontend Application Bootstrap
系统 SHALL 提供一个基于 React + TypeScript 的前端应用，能够通过 HTTP 访问。前端 SHALL 在 API 层处理 snake_case（后端）与 camelCase（前端）之间的字段名转换。

#### Scenario: Application starts
- **WHEN** user navigates to the application URL
- **THEN** application renders the login page

### Requirement: Location Management UI
系统 SHALL 提供位置管理功能的用户界面，支持增删改查操作。位置列表查询 SHALL 正确处理分页结果格式（items, total, page, limit）。

#### Scenario: List locations
- **WHEN** user navigates to the location management page
- **THEN** system displays a list of all locations with correct pagination

#### Scenario: Create location
- **WHEN** user fills in location details and clicks create
- **THEN** system sends a POST request to /locations and updates the list

#### Scenario: Update location
- **WHEN** user edits location details and saves
- **THEN** system sends a PUT request to /locations/:id and updates the list

#### Scenario: Delete location
- **WHEN** user clicks delete on a location
- **THEN** system sends a DELETE request to /locations/:id and updates the list

### Requirement: Room Management UI
系统 SHALL 提供房间管理功能的用户界面，支持增删改查操作。房间列表查询 SHALL 正确处理分页结果格式（items, total, page, limit）。

#### Scenario: List rooms
- **WHEN** user navigates to the room management page
- **THEN** system displays a list of all rooms with correct pagination

#### Scenario: Create room
- **WHEN** user fills in room details and clicks create
- **THEN** system sends a POST request to /rooms and updates the list

#### Scenario: Update room
- **WHEN** user edits room details and saves
- **THEN** system sends a PUT request to /rooms/:id and updates the list

#### Scenario: Delete room
- **WHEN** user clicks delete on a room
- **THEN** system sends a DELETE request to /rooms/:id and updates the list

### Requirement: Landlord Management UI
系统 SHALL 提供房东管理功能的用户界面，支持增删改查操作。房东列表查询 SHALL 正确处理分页结果格式（items, total, page, limit）。

#### Scenario: List landlords
- **WHEN** user navigates to the landlord management page
- **THEN** system displays a list of all landlords with correct pagination

#### Scenario: Create landlord
- **WHEN** user fills in landlord details and clicks create
- **THEN** system sends a POST request to /landlords and updates the list

#### Scenario: Update landlord
- **WHEN** user edits landlord details and saves
- **THEN** system sends a PUT request to /landlords/:id and updates the list

#### Scenario: Delete landlord
- **WHEN** user clicks delete on a landlord
- **THEN** system sends a DELETE request to /landlords/:id and updates the list

### Requirement: Lease Management UI
系统 SHALL 提供租约管理功能的用户界面，支持增删改查、续租、退租操作。租约列表查询 SHALL 正确处理分页结果格式（items, total, page, limit）。

#### Scenario: List leases
- **WHEN** user navigates to the lease management page
- **THEN** system displays a list of all leases with correct pagination

#### Scenario: Create lease
- **WHEN** user fills in lease details and clicks create
- **THEN** system sends a POST request to /leases and updates the list

#### Scenario: Update lease
- **WHEN** user edits lease details and saves
- **THEN** system sends a PUT request to /leases/:id and updates the list

#### Scenario: Delete lease
- **WHEN** user clicks delete on a lease
- **THEN** system sends a DELETE request to /leases/:id and updates the list

#### Scenario: Renew lease
- **WHEN** user clicks renew on a lease
- **THEN** system sends a POST request to /leases/:id/renew and updates the list

#### Scenario: Checkout lease
- **WHEN** user clicks checkout on a lease
- **THEN** system sends a POST request to /leases/:id/checkout and updates the list

### Requirement: Bill Management UI
系统 SHALL 提供账单管理功能的用户界面，支持增删改查和打印收据操作。账单列表查询 SHALL 正确处理分页结果格式（items, total, page, limit）。

#### Scenario: List bills
- **WHEN** user navigates to the bill management page
- **THEN** system displays a list of all bills with correct pagination

#### Scenario: Create bill
- **WHEN** user fills in bill details and clicks create
- **THEN** system sends a POST request to /bills and updates the list

#### Scenario: Update bill
- **WHEN** user edits bill details and saves
- **THEN** system sends a PUT request to /bills/:id and updates the list

#### Scenario: Delete bill
- **WHEN** user clicks delete on a bill
- **THEN** system sends a DELETE request to /bills/:id and updates the list

#### Scenario: Print receipt
- **WHEN** user clicks print receipt on a bill
- **THEN** system downloads an RTF file of the receipt

### Requirement: Print UI
系统 SHALL 提供打印功能的用户界面，支持打印账单、租约、发票。打印任务列表查询 SHALL 正确处理分页结果格式（items, total, page, limit）。

#### Scenario: Print bill
- **WHEN** user fills in bill ID and clicks print
- **THEN** system sends a POST request to /print/bill and returns a job ID

#### Scenario: Print lease
- **WHEN** user fills in lease ID and clicks print
- **THEN** system sends a POST request to /print/lease and returns a job ID

#### Scenario: Print invoice
- **WHEN** user fills in bill ID and clicks print
- **THEN** system sends a POST request to /print/invoice and returns a job ID

### Requirement: Income Report UI
系统 SHALL 提供收入报表功能的用户界面。

#### Scenario: View income report
- **WHEN** user navigates to the income report page
- **THEN** system displays an income report based on bill data

### Requirement: Authentication UI
系统 SHALL 提供用户认证功能的用户界面。

#### Scenario: User login
- **WHEN** user enters valid credentials and clicks login
- **THEN** system stores a token and redirects to the dashboard

#### Scenario: User logout
- **WHEN** user clicks logout
- **THEN** system clears the token and redirects to the login page
