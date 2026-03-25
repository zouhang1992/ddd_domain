## Design

### Overview

This change focuses on removing obsolete application service layer code that has already been replaced by the bus architecture (Command/Query/Event Bus) in previous implementations.

### Decisions

1. **Remove unused service implementations**:
   - `LocationService` and `RoomService` are no longer used since their functionality is now handled by `LocationCommandHandler`, `LocationQueryHandler`, `RoomCommandHandler`, and `RoomQueryHandler`.

2. **Remove old HTTP handlers**:
   - The traditional `LocationHandler` and `RoomHandler` that directly used the application service layer have been replaced by `CQRSLocationHandler` and `CQRSRoomHandler` that interact directly with the command and query buses.

3. **Cleanup main.go initialization**:
   - Remove the initialization of traditional services from `main.go` to avoid unused variables and potential confusion.

### Implementation Approach

1. **Files to Delete**:
   - `internal/application/service/location.go`
   - `internal/application/service/room.go`
   - `internal/facade/location_handler.go`
   - `internal/facade/room_handler.go`

2. **Main.go Changes**:
   - Remove imports for `service` package
   - Remove unused variable declarations for `locationService` and `roomService`

### Code Changes Summary

```diff
// cmd/api/main.go changes
- import "github.com/zouhang1992/ddd_domain/internal/application/service"
...
- locationService := service.NewLocationService(locationRepo)
- roomService := service.NewRoomService(roomRepo)
```

All other services that are still using the traditional architecture will remain intact for now.
