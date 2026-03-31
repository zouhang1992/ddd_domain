# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Backend (Go)

**Build the project:**
```bash
go build ./...
```

**Run tests:**
```bash
go test ./...
```

**Run a specific test:**
```bash
go test ./path/to/package -run TestName -v
```

**Run the application:**
```bash
go run cmd/api/main.go
```

### Frontend (React + TypeScript)

**Install dependencies:**
```bash
cd web && npm install
```

**Start development server:**
```bash
cd web && npm run dev
```

**Build for production:**
```bash
cd web && npm run build
```

**Lint code:**
```bash
cd web && npm run lint
```

## Architecture Overview

This is a Domain-Driven Design (DDD) rental property management system with CQRS and Event Sourcing patterns.

### Layers

1. **Domain Layer** (`internal/domain/`)
   - Aggregate roots: Lease, Room, Bill, Deposit, Landlord, Location
   - Domain events and repositories
   - Business logic and invariants

2. **Application Layer** (`internal/application/`)
   - Command handlers: Write operations
   - Query handlers: Read operations
   - Services: Domain service orchestration

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - SQLite persistence
   - Command/Query/Event buses
   - Logging (Zap)
   - SAGA pattern implementation

4. **Facade Layer** (`internal/facade/`)
   - HTTP handlers
   - CQRS facade for external API

5. **Web Layer** (`web/`)
   - React + TypeScript + Ant Design frontend
   - Vite as build tool

### Key Patterns

- **Dependency Injection**: Using Uber Fx for DI
- **Event-Driven**: Domain events drive side effects (e.g., lease status changes affect room status)
- **CQRS**: Separate command and query models
- **Aggregate Roots**: Each aggregate manages its own consistency boundary

### Important Domain Relationships

- Location → Room (one-to-many)
- Room → Lease (one-to-many)
- Lease → Bill (one-to-many)
- Lease → Deposit (one-to-one)
- Landlord → Lease (one-to-many)

### Event Bus

Events are published to the event bus and handled by subscribers:
- `lease.activated`, `lease.checkout`, `lease.expired`: Update room status
- All domain events: Record operation logs
