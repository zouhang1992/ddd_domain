## 1. Setup

- [x] 1.1 Add robfig/cron dependency to go.mod
- [x] 1.2 Run go mod tidy

## 2. Domain Layer

- [x] 2.1 Add CheckAndExpireLeases method to LeaseService

## 3. Application Layer

- [x] 3.1 Create LeaseExpirationScheduler at internal/application/lease/scheduler.go
- [x] 3.2 Update lease/module.go to provide scheduler

## 4. Integration

- [x] 4.1 Add scheduler startup in main.go

## 5. Verification

- [x] 5.1 Run go build to verify compilation
