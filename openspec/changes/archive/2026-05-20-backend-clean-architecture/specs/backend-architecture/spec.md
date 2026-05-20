## ADDED Requirements

### Requirement: Layered package structure

The backend SHALL organize its Go packages into four horizontal layers under `backend/internal/`: `domain`, `application`, `infrastructure`, and `composition`. Each layer SHALL be a sibling directory.

#### Scenario: Required layer directories exist

- **WHEN** a contributor inspects `backend/internal/`
- **THEN** the directories `domain/`, `application/`, `infrastructure/persistence/`, and `composition/` exist
- **AND** each contains at least one Go file owned by that layer

#### Scenario: HTTP, jobs, and middleware retain their own top-level packages

- **WHEN** a contributor inspects `backend/internal/`
- **THEN** the packages `handler/`, `middleware/`, and `jobs/` continue to exist alongside the four core layers
- **AND** they depend on `application/` and `domain/` for business types, not on `infrastructure/...` directly

#### Scenario: Adapter implementations live under infrastructure

- **WHEN** a contributor inspects `backend/internal/infrastructure/`
- **THEN** it contains `persistence/` (GORM-backed adapters), `sse/` (the SSE broker that satisfies `application.EventBroadcaster`), and `email/` (the gomail-backed adapter that satisfies `application.EmailSender`)
- **AND** any new concrete implementation of an `application/` port lands here, not at the top level of `internal/`

### Requirement: Pure domain layer

The `backend/internal/domain` package SHALL contain only pure Go types and functions. It MUST NOT import any persistence, HTTP, ORM, or framework package.

#### Scenario: Domain package has no forbidden imports

- **WHEN** the import graph of `backend/internal/domain` is inspected
- **THEN** it does not include `gorm.io/...`, `database/sql`, `net/http`, `github.com/riverqueue/...`, or any package under `backend/internal/infrastructure`, `backend/internal/handler`, or `backend/internal/application`

#### Scenario: Domain types carry no GORM tags

- **WHEN** any struct in `backend/internal/domain` is inspected
- **THEN** no field uses a `gorm:"..."` struct tag

#### Scenario: Domain compiles standalone

- **WHEN** `go build ./internal/domain/...` is run from `backend/`
- **THEN** the build succeeds without requiring any non-stdlib persistence dependency

### Requirement: Value objects for invariant-bearing identifiers and enums

The domain SHALL expose value-object types for identifiers and enumerated fields whose invariants are enforced at the type boundary. At minimum, `UserID` and `StatField` SHALL exist as domain types.

#### Scenario: UserID rejects empty input

- **WHEN** `domain.NewUserID("")` is called
- **THEN** it returns a non-nil error
- **AND** no `UserID` value is produced

#### Scenario: StatField is a closed set

- **WHEN** a caller attempts to construct a `StatField` from an unknown string
- **THEN** the constructor returns an error naming the allowed values

### Requirement: Application layer defines ports and use cases

The `backend/internal/application` package SHALL declare the repository interfaces (ports) consumed by use cases, and SHALL expose each use case as a struct with an `Execute(ctx context.Context, ...)` method.

#### Scenario: Repository interfaces live in application/ports.go

- **WHEN** a contributor inspects `backend/internal/application/ports.go`
- **THEN** the file declares `StatsRepository` (and any further repository interfaces introduced by later features)
- **AND** the interfaces use only `context.Context` and domain types in their signatures — no GORM or HTTP types

#### Scenario: Use cases are structs with explicit dependencies

- **WHEN** a use case type in `backend/internal/application` is inspected
- **THEN** its dependencies are declared as interface-typed fields on the struct
- **AND** its public surface includes a single `Execute(ctx, …)` method
- **AND** all such methods take `context.Context` as their first parameter

#### Scenario: Business logic does not live in repositories

- **WHEN** any file under `backend/internal/infrastructure/persistence` is inspected
- **THEN** it contains only persistence translation, CRUD, and query code
- **AND** invariant enforcement (e.g. non-negative counters) lives on the relevant domain type, invoked by a use case

### Requirement: Persistence adapters in infrastructure layer

Concrete repository implementations SHALL live under `backend/internal/infrastructure/persistence/`. Each implementation SHALL satisfy a port declared in `backend/internal/application` and SHALL translate between persistence models and domain types via explicit mapper functions.

#### Scenario: Repository implementation satisfies an application port

- **WHEN** a repository implementation type in `infrastructure/persistence/` is inspected
- **THEN** it is assignable to the corresponding interface in `application/ports.go`
- **AND** a compile-time assertion `var _ application.<Port> = (*<Impl>)(nil)` exists in the same file

#### Scenario: GORM models are separate from domain types

- **WHEN** persistence models are inspected
- **THEN** the GORM-tagged structs live in `infrastructure/persistence/gorm_models.go` (or a sibling file in the same package)
- **AND** they are distinct types from the domain entities, connected only via mapper functions

#### Scenario: AutoMigrate operates on persistence models

- **WHEN** `cmd/server` starts and invokes AutoMigrate
- **THEN** the entity list passed to GORM comes from `infrastructure/persistence` and contains the GORM-tagged models, not domain types

### Requirement: Handlers depend on application use cases

HTTP handlers in `backend/internal/handler` SHALL receive their dependencies as application use cases or interfaces declared in `application/`. Handlers MUST NOT import `backend/internal/infrastructure/persistence` or `gorm.io/...`.

#### Scenario: Handler constructor signature uses application types

- **WHEN** a handler constructor in `backend/internal/handler` is inspected
- **THEN** every dependency parameter is either a use-case struct from `application/` or an interface declared in `application/`
- **AND** no parameter is a concrete repository struct from `infrastructure/persistence`

#### Scenario: Handler package has no persistence imports

- **WHEN** the import graph of `backend/internal/handler` is inspected
- **THEN** it does not include `gorm.io/...` or any package under `backend/internal/infrastructure`

### Requirement: Composition root wires the dependency graph

A `backend/internal/composition` package SHALL build the full application dependency graph and expose it to `cmd/server`. `cmd/server/main.go` MUST delegate wiring to this composition root rather than instantiating repositories, use cases, and handlers inline.

#### Scenario: composition.go assembles the graph

- **WHEN** `backend/internal/composition/composition.go` is inspected
- **THEN** it contains a function (e.g. `Build(cfg Config) (*App, error)`) that constructs repositories, use cases, handlers, the SSE broker, and any background workers
- **AND** all constructor calls flow through it

#### Scenario: main.go is a thin entry point

- **WHEN** `backend/cmd/server/main.go` is inspected
- **THEN** it does not call `repository.New*`, `application.<UseCase>{…}`, or `handler.New*` directly
- **AND** it invokes the composition root and then starts the HTTP server

### Requirement: Layer dependency rules

The codebase SHALL enforce a strict inward dependency direction: outer layers may depend on inner layers, never the reverse.

The allowed direction is:
`composition → handler/jobs → application → domain`
and
`composition → infrastructure/{persistence,sse} → application → domain`

#### Scenario: Domain depends on nothing internal

- **WHEN** the import graph of `backend/internal/domain` is inspected
- **THEN** it does not import any other `backend/internal/...` package

#### Scenario: Application depends only on domain

- **WHEN** the import graph of `backend/internal/application` is inspected
- **THEN** the only `backend/internal/...` package it imports is `backend/internal/domain`

#### Scenario: Infrastructure depends only on application and domain

- **WHEN** the import graph of `backend/internal/infrastructure/...` is inspected
- **THEN** it imports only `backend/internal/application` and `backend/internal/domain` among the project's internal packages

### Requirement: Goca and contributor docs reflect the new layout

`backend/.goca.yaml` SHALL be updated so generated code lands in the new layer layout, and `CLAUDE.md` SHALL document the four-layer architecture, the dependency rules, and the new "add a feature" workflow.

#### Scenario: .goca.yaml enables the usecase layer

- **WHEN** `backend/.goca.yaml` is inspected
- **THEN** `architecture.layers` includes `usecase` (or equivalent application layer) alongside the existing layers

#### Scenario: CLAUDE.md documents the architecture

- **WHEN** `CLAUDE.md` is read
- **THEN** it describes the four layers (`domain`, `application`, `infrastructure`, `composition`), the inward-dependency rule, and the updated Goca workflow for adding a new feature
- **AND** outdated guidance stating that `internal/usecase/` does not exist is removed

### Requirement: HTTP and frontend behavior preserved

The refactor MUST NOT change any HTTP route, request shape, response shape, Swagger document, or frontend integration. It is an internal restructuring.

#### Scenario: Swagger output is unchanged

- **WHEN** `just api` is run before and after the refactor
- **THEN** the generated `backend/docs/swagger.json` is structurally identical (same routes, same request/response schemas)

#### Scenario: Frontend continues to work without changes

- **WHEN** the frontend is built and run against the refactored backend
- **THEN** no frontend source file under `frontend/src/` requires modification
- **AND** the dashboard, stats SSE, auth, and data-export flows behave as before
