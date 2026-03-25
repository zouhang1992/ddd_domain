## Design

### Overview

This change will refactor the application's dependency management to use UberFx, a lightweight Go dependency injection library that simplifies component wiring and lifecycle management.

### Architecture

The new architecture will:
1. Add Fx modules for different layers (persistence, application, facade)
2. Provide constructors as Fx providers using `fx.Provide`
3. Automatically inject dependencies using type matching in constructor parameters
4. Centralize configuration and start the server in an `fx.Invoke` block

### Component Organization

1. **Persistence Module**: Provides SQLite connection and all repositories
2. **Bus Module**: Provides CommandBus, QueryBus, and EventBus
3. **Command Handler Module**: Provides all command handlers
4. **Query Handler Module**: Provides all query handlers
5. **Facade Module**: Provides all HTTP handlers

### Changes to Each Package

- **cmd/api/main.go**: Rewritten using fx.New().Provide(...).Run(...)
- **internal/infrastructure/persistence/sqlite**: Add `fx.Provide` functions for all repositories
- **internal/application/command/handler**: Add `fx.Provide` functions for command handlers
- **internal/application/query/handler**: Add `fx.Provide` functions for query handlers
- **internal/facade**: Add `fx.Provide` functions for all HTTP handlers
- **go.mod**: Add `go.uber.org/fx` dependency

### Key Benefits

1. **Simplified initialization**: No need to manually construct dependencies
2. **Type safety**: Fx performs compile-time type checking of dependencies
3. **Better testing**: Easier to mock and replace components
4. **Cleaner code**: Reduces boilerplate in main.go
5. **Lifecycle management**: Fx manages start/stop hooks automatically

### Example Code Snippet

```go
// cmd/api/main.go
package main

import (
	"go.uber.org/fx"
	"net/http"
)

func main() {
	fx.New(
		persistence.Module,
		bus.Module,
		command.Module,
		query.Module,
		facade.Module,
		fx.Invoke(setupServer),
	).Run()
}

func setupServer(
	locationHandler *facade.CQRSLocationHandler,
	roomHandler *facade.CQRSRoomHandler,
	// ... more handlers
) {
	mux := http.NewServeMux()
	locationHandler.RegisterRoutes(mux)
	roomHandler.RegisterRoutes(mux)
	// ... register other routes
	log.Println("Server starting on :8080")
	_ = http.ListenAndServe(":8080", mux)
}

// internal/infrastructure/persistence/sqlite/module.go
package sqlite

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewConnection),
	fx.Provide(NewLandlordRepository),
	fx.Provide(NewLeaseRepository),
	// ... more repositories
)

// internal/application/command/handler/module.go
package handler

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewLandlordCommandHandler),
	fx.Provide(NewLeaseCommandHandler),
	// ... more command handlers
)
```
