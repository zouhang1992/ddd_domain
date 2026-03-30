## ADDED Requirements

### Requirement: Scheduled lease expiration check
The system SHALL automatically check for expired active leases on a scheduled basis.

#### Scenario: Check runs on application startup
- **WHEN** the application starts
- **THEN** the lease expiration check runs immediately

#### Scenario: Check runs hourly
- **WHEN** the application is running
- **THEN** the lease expiration check runs once every hour at minute 0

### Requirement: Expired lease processing
The system SHALL mark active leases as expired when their end date has passed.

#### Scenario: Lease is expired when end date is in past
- **WHEN** an active lease has end_date <= current time
- **THEN** the lease status is updated to "expired"

#### Scenario: Lease expiration event is published
- **WHEN** a lease is marked as expired
- **THEN** a "lease.expired" event is published

#### Scenario: Expired lease is skipped
- **WHEN** a lease is already in "expired", "checkout", or "pending" state
- **THEN** the lease is not processed

### Requirement: Room status update on lease expiration
The system SHALL update the room status to "available" when a lease expires.

#### Scenario: Room becomes available on lease expiration
- **WHEN** a lease expires
- **THEN** the associated room's status is updated to "available"

### Requirement: Error handling for expiration job
The system SHALL handle errors gracefully during lease expiration processing.

#### Scenario: Single lease failure does not block others
- **WHEN** processing multiple expired leases and one fails
- **THEN** the failure is logged and processing continues with other leases

#### Scenario: Job failure does not affect application
- **WHEN** the scheduled lease expiration check fails
- **THEN** the error is logged and the application continues to run

### Requirement: Processing count reporting
The system SHALL report the number of leases processed during each expiration check.

#### Scenario: Count of processed leases is returned
- **WHEN** the expiration check completes
- **THEN** the number of leases processed is returned

#### Scenario: Count is logged
- **WHEN** the expiration check completes
- **THEN** the number of leases processed is logged
