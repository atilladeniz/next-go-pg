## Context

Today the Go backend (`backend/internal/`) has three layers: `domain/` (GORM structs), `repository/` (concrete CRUD types), and `handler/` (HTTP). Handlers depend on concrete repository structs, repositories carry business rules (negative-count clamping, conditional field updates), and `cmd/server/main.go` is a 498-line file that wires everything inline. CLAUDE.md acknowledges the gap ("`internal/usecase/` reserved … doesn't exist yet because no feature has needed a dedicated usecase layer").

Our reference implementation is the sibling project `quiz-engine/engine/`:
- `internal/domain/` is pure — no DB imports, value objects (`QuestionType`, `Option`, `CorrectAnswer`), aggregate methods (`Quiz.IsMock()`, `CorrectAnswer.IsCorrect()`).
- `internal/application/ports.go` defines hexagonal interfaces (`SessionRepository`, `CourseRepository`).
- `internal/application/*.go` holds use-case structs (`LoadQuiz`, `SaveSession`) with `Execute(ctx, …)` methods that take ports via DI.
- `internal/infrastructure/persistence/` holds the actual SQLite/Postgres adapters.
- `internal/composition/composition.go` wires the graph; main.go is a thin entry point.

The surface area in `next-go-pg/backend` is still small — one domain entity (`UserStats`), one repository, a handful of handlers, plus River workers and SSE plumbing. This is the right moment to refactor: cheap now, expensive after several more features land.

**Stakeholders:** sole maintainer (project owner). No external API consumers depend on internal Go package paths. Frontend is consumed via Swagger/Orval, which is unaffected.

## Goals / Non-Goals

**Goals:**
- Establish four-layer architecture: `domain → application → infrastructure → composition`, matching quiz-engine.
- Domain layer is pure: zero imports of `gorm.io/*`, `database/sql`, `net/http`, or any framework. Compilable in isolation.
- All repository access goes through interfaces declared in `application/`; concrete impls live in `infrastructure/persistence/`.
- Each handler depends on use-case interfaces (or directly on use-case structs), never on repository structs.
- Composition root wires the dependency graph; `cmd/server/main.go` shrinks to a slim bootstrap (target: <150 lines).
- Existing HTTP API surface, Swagger output, and frontend integration are unchanged.
- Goca workflow updated so future features generate into the new layout.

**Non-Goals:**
- No new product features.
- No database schema changes.
- No event sourcing, CQRS, or other DDD tactical patterns beyond what's justified for the current entities. Aggregates stay simple.
- No replacement of GORM. Persistence stays on GORM; only the package boundary moves.
- No rewrite of River workers' internals. They get the same DI treatment (interfaces in their construction signature) but otherwise stay put.
- No new test framework. Existing test setup stays.

## Decisions

### D1: Layer layout matches quiz-engine

```
backend/internal/
├── domain/                   # Pure entities, value objects, aggregate methods
│   ├── user_stats.go
│   ├── value_objects.go      # UserID, count types
│   └── registry.go           # entity list (kept here for AutoMigrate target list)
├── application/              # Use cases + ports
│   ├── ports.go              # All repository interfaces in one file
│   ├── stats_usecases.go     # GetUserStats, IncrementStatField, …
│   └── dto.go                # Application-layer DTOs (if any)
├── infrastructure/
│   └── persistence/
│       ├── gorm_models.go    # GORM-tagged structs
│       ├── user_stats_mapper.go  # domain ↔ gorm translation
│       └── user_stats_repo.go    # implements application.StatsRepository
├── handler/                  # HTTP — depends on application use cases
├── middleware/               # unchanged
├── jobs/                     # River workers — also wired via interfaces
├── sse/                      # unchanged (already infrastructure-ish)
└── composition/
    └── composition.go        # DI wiring
```

**Why this layout over alternatives:**
- *Vertical (per-feature) slicing*: Considered (`internal/userstats/{domain,app,infra}`). Rejected because we don't yet have enough features to justify the boilerplate, and quiz-engine — our reference — uses horizontal layering. Symmetry beats abstract purity here.
- *Keeping repository/ as-is and just adding usecase/*: Rejected. It leaves domain GORM-tainted and skips the ports/adapters split that makes testing easy.

### D2: Pure domain, mapper-based persistence

Domain types are POGOs (plain old Go objects) with no struct tags. A separate GORM model in `infrastructure/persistence/gorm_models.go` carries the `gorm:""` tags. A mapper translates both directions:

```go
// infrastructure/persistence/user_stats_mapper.go
func toDomain(m gormUserStats) domain.UserStats { … }
func fromDomain(s domain.UserStats) gormUserStats { … }
```

**Why:** lets `domain` compile without a database driver, enables in-memory fakes, and matches quiz-engine's `infrastructure/persistence/sqlite_session_repo.go` pattern.

**Alternative considered:** mark GORM tags optional via build tags. Rejected — fragile, hides the dependency from `go mod graph`.

### D3: AutoMigrate target list moves to persistence layer

`internal/domain/registry.go` currently returns `[]interface{}{ &UserStats{} }` for GORM AutoMigrate. After the refactor, AutoMigrate must point at the GORM models, not the domain types. The registry moves to `infrastructure/persistence/registry.go` and returns the GORM models. `cmd/server/main.go` (via composition root) calls it.

### D4: Use cases are structs with `Execute(ctx, …)`

Matching quiz-engine, each use case is a struct holding its dependencies via interface fields, with one entry point:

```go
type IncrementStatField struct {
    Repo StatsRepository
    SSE  EventBroadcaster
}
func (uc IncrementStatField) Execute(ctx context.Context, userID domain.UserID, field domain.StatField, delta int) error
```

**Why over thin functions:** explicit dependencies show up in the type signature; easy to swap fakes in tests; matches quiz-engine; lets us evolve a use case (e.g. add caching, retries) without changing the call site.

### D5: Value objects, but scoped

Introduce **only** the value objects whose invariants are actively violated today:
- `UserID` (string newtype, blocks empty values at construction)
- `StatField` (enum of valid increment targets: `"projects"`, `"tasks"`, `"completed"`, …) — replaces the magic-string switch in `repository/user_stats.go`

**Not** introducing wrapper types for every count field (`ProjectCount`, `TaskCount`, etc.) — overkill at this scale. Negative-count clamping moves into a domain method `UserStats.IncrementField(field StatField, delta int)`.

### D6: Composition root, not Wire/Fx

Hand-written DI in `internal/composition/composition.go` mirroring quiz-engine. No `google/wire`, no `uber-go/fx`. Hand-written is fine at this scale; both alternatives add tooling cost without payoff for <30 components.

### D7: Repository interfaces live in `application/`, not `domain/`

Quiz-engine puts ports in `application/ports.go`. Some Clean Architecture purists put them in `domain/`. We follow quiz-engine to keep symmetry — the domain stays free of even abstract dependencies on persistence.

### D8: SSE broker stays where it is

The SSE broker (`internal/sse/broker.go`) is infrastructure but doesn't need a `persistence/` sibling treatment. Use cases that broadcast events depend on a small `EventBroadcaster` interface declared in `application/ports.go`; the broker satisfies it.

## Risks / Trade-offs

- **Risk: GORM hooks on entities silently break.** GORM hooks (`BeforeCreate`, etc.) attached to domain types would not run on the new GORM model types. → **Mitigation:** grep for `func (.*) BeforeCreate|BeforeSave|AfterFind` in the codebase before starting; current `UserStats` has none, but verify before each subsequent entity migrates.
- **Risk: mapper duplication.** Two structs per entity (domain + GORM model) plus two mapper functions. → **Trade-off accepted.** This is the cost of domain purity. Mappers are mechanical and trivial to test.
- **Risk: Goca regenerates into the old layout.** `goca make entity` emits files into `internal/domain/` with GORM tags by default. → **Mitigation:** update `backend/.goca.yaml` (`architecture.layers`) and document the new manual steps in CLAUDE.md. Long-term: contribute a config flag upstream or wrap Goca with a `just` recipe.
- **Risk: `cmd/server/main.go` history gets noisy.** A 498→<150 line shrink is a big diff. → **Mitigation:** split into two commits (move composition code out → simplify main).
- **Risk: River workers also need wiring.** They currently get a concrete `*gorm.DB` and call repositories directly. → Handled in tasks.md phase 4. Without this, the new layer would be bypassed by jobs.

## Migration Plan

Phased rollout, one PR per phase, verify (`just lint && just typecheck && just dev-backend`) between phases:

1. **Phase 1 — Scaffolding (no behavior change):**
   - Create `application/`, `infrastructure/persistence/`, `composition/` packages.
   - Define `application/ports.go` interfaces (initially empty bodies).
   - No imports flipped yet. Existing code keeps working.

2. **Phase 2 — Domain split for UserStats:**
   - Introduce `domain/value_objects.go` (`UserID`, `StatField`).
   - Strip GORM tags from `domain/user_stats.go`; move tagged copy to `infrastructure/persistence/gorm_models.go`.
   - Write mapper + repo impl satisfying `application.StatsRepository`.
   - Move `AllEntities()` to persistence package.

3. **Phase 3 — Use cases + handler rewire:**
   - Write `application/stats_usecases.go` with `Execute(ctx, …)` methods.
   - Repoint handlers to consume use cases.
   - Drop business logic from repository (negative clamping, field switch) — moves into domain method called by use case.

4. **Phase 4 — Composition root + River:**
   - Move DI out of `cmd/server/main.go` into `composition/composition.go`.
   - Rewire River workers to consume application interfaces, not `*gorm.DB`.
   - Slim main.go to flag parsing + composition + `http.ListenAndServe`.

5. **Phase 5 — Cleanup:**
   - Delete now-empty `internal/repository/`.
   - Update CLAUDE.md and `.goca.yaml`.
   - Run `just lint && just typecheck`, manual smoke test of dashboard + stats SSE.

**Rollback:** each phase is a separate PR. Revert the PR — earlier phases keep compiling because they only add new packages without removing old ones until Phase 5.

## Open Questions

- **Q1:** Should `internal/domain/registry.go` (entity list) stay in `domain/` (purely an enumeration, no persistence types) or move with the GORM models? → *Tentative answer: move to `infrastructure/persistence/`, because after the refactor the list contains GORM models, not domain types.*
- **Q2:** Do River job args (`MagicLinkArgs`, `DataExportArgs` in `internal/jobs/`) count as domain or application DTOs? → *Tentative answer: leave them in `jobs/`. They're River-framework-shaped; treating them as domain would force a second mapper. Revisit if a use case needs to enqueue jobs through an `application/JobEnqueuer` port.*
- **Q3:** Add basic use-case unit tests in Phase 3, or follow up separately? → *Recommend: at least one happy-path test per use case in Phase 3, to lock in the new interfaces. Comprehensive coverage can follow.*
