## Why

The Go backend mixes responsibilities: handlers depend on concrete repository structs, GORM tags live inside `internal/domain/`, and business rules (e.g. negative-count clamping, conditional updates) leak into `repository/user_stats.go`. There is no application/usecase layer despite CLAUDE.md reserving one. Compared to our reference implementation in `quiz-engine/engine/` — which has explicit ports (`application/ports.go`), pure domain types, value objects, and a composition root — `next-go-pg` is one feature away from tangled dependencies. We want to close the gap now, while the surface area is still small (one repository, a handful of handlers), instead of after several more features have been stacked on top.

## What Changes

- Introduce `backend/internal/application/` with hexagonal ports (repository interfaces) and explicit use-case structs (`Execute(ctx, ...)` pattern, matching quiz-engine).
- Move all business logic currently sitting in handlers and repositories into use cases.
- Make `backend/internal/domain/` pure: remove GORM struct tags, introduce minimal value objects (`UserID`, count types with invariants).
- Add `backend/internal/infrastructure/persistence/` with GORM models + mappers (domain ↔ persistence translation).
- Repository implementations consume the domain types via mappers and satisfy the application-layer interfaces.
- Add a composition root at `backend/internal/composition/` so `cmd/server/main.go` becomes a thin entry point — current 498-line `main.go` shrinks substantially.
- **BREAKING (internal-only)**: `APIHandler` and other handler constructors take application-layer interfaces, not concrete `*repository.UserStatsRepository`. No HTTP API surface changes.
- Update Goca workflow guidance in CLAUDE.md to reflect the new layer layout (`goca make usecase` becomes part of the normal flow, not "on demand").

## Capabilities

### New Capabilities

- `backend-architecture`: Defines the layered architecture contract for the Go backend — domain purity, application-layer ports & use cases, persistence mappers, composition root, and the rules that govern dependencies between layers.

### Modified Capabilities

_None — this is the first capability spec in the repo._

## Impact

- **Code touched**: `backend/internal/domain/`, `backend/internal/repository/` (moves to `infrastructure/persistence/`), `backend/internal/handler/` (consumes interfaces), `backend/cmd/server/main.go` (delegates to composition root), `backend/.goca.yaml` (enable usecase layer), CLAUDE.md (layer guide).
- **Public API**: no change. HTTP routes, request/response shapes, and Swagger output stay identical. Frontend untouched.
- **Migrations**: none. GORM entities keep the same table schema; only their package location changes. AutoMigrate registry (`internal/domain/registry.go`) follows the GORM models to their new home.
- **Tests**: enables real unit tests for use cases against in-memory repository fakes (currently impossible because handlers hold concrete types).
- **Risk**: medium. Pure refactor with no behavior change, but it touches the wiring of every handler. Mitigated by phased rollout (see design.md) and the existing test surface.
