## ADDED Requirements

### Requirement: Operation log record creation
The system SHALL create an operation log record for every domain event published by aggregates.

#### Scenario: Room created event logged
- **WHEN** a room.created event is published
- **THEN** an operation log record is created with event name "room.created"

#### Scenario: Lease created event logged
- **WHEN** a lease.created event is published
- **THEN** an operation log record is created with event name "lease.created"

#### Scenario: Bill paid event logged
- **WHEN** a bill.paid event is published
- **THEN** an operation log record is created with event name "bill.paid"

### Requirement: Operation log contains complete event information
The system SHALL store complete event information in the operation log record.

#### Scenario: Log contains timestamp
- **WHEN** an operation log is created
- **THEN** the log record includes the timestamp when the event occurred

#### Scenario: Log contains domain type
- **WHEN** an operation log is created
- **THEN** the log record includes the domain type (e.g., "room", "lease", "bill")

#### Scenario: Log contains aggregate ID
- **WHEN** an operation log is created
- **THEN** the log record includes the aggregate ID associated with the event

#### Scenario: Log contains operator ID
- **WHEN** an operation log is created
- **THEN** the log record includes the operator ID from the authentication context

#### Scenario: Log contains action type
- **WHEN** an operation log is created
- **THEN** the log record includes the action type (e.g., "created", "updated", "deleted", "activated", "paid")

#### Scenario: Log contains event details
- **WHEN** an operation log is created
- **THEN** the log record includes the complete event payload as JSON

### Requirement: Operation log query by aggregate ID
The system SHALL allow querying operation logs by aggregate ID.

#### Scenario: Query logs for specific aggregate
- **WHEN** querying operation logs by aggregate ID "abc123"
- **THEN** all operation logs for that aggregate are returned in chronological order

### Requirement: Operation log query by domain type
The system SHALL allow querying operation logs by domain type with pagination.

#### Scenario: Query logs for room domain
- **WHEN** querying operation logs for domain type "room" with offset 0 and limit 10
- **THEN** up to 10 most recent room operation logs are returned

### Requirement: Operation log query by time range
The system SHALL allow querying operation logs by time range with pagination.

#### Scenario: Query logs in time range
- **WHEN** querying operation logs between start time "2024-01-01 and end time "2024-01-31 with offset 0 and limit 20
- **THEN** up to 20 operation logs within that time range are returned

### Requirement: Operation log persistence
The system SHALL persist operation logs to SQLite database.

#### Scenario: Log saved to database
- **WHEN** an operation log is created
- **THEN** it is saved to the operation_logs table in SQLite

### Requirement: Asynchronous logging
The system SHALL record operation logs asynchronously without blocking the main business flow.

#### Scenario: Event published without waiting for log
- **WHEN** a domain event is published
- **THEN** the main business flow continues without waiting for the operation log to be persisted
