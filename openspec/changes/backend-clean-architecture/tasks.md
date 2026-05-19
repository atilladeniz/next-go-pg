## 1. Phase 1 — Scaffold new layers (no behavior change)

- [x] 1.1 Create empty packages: `backend/internal/application/`, `backend/internal/infrastructure/persistence/`, `backend/internal/composition/` (each with a single placeholder `.go` file so `go build ./...` succeeds)
- [x] 1.2 Add `backend/internal/application/ports.go` with the `StatsRepository` and `EventBroadcaster` interfaces. _(Reshuffle: full method signatures landed in Phase 1 instead of Phase 3 — empty interfaces tripped gopls's `any`-suggestion. Phase 3 task 3.1 is now a confirm-only step.)_
- [x] 1.3 `go build ./...` and `go vet ./...` green. (`just lint` is frontend-only — Biome/Steiger; backend uses `go vet`/gopls in this repo.)

## 2. Phase 2 — Domain split for UserStats

- [x] 2.1 Added `backend/internal/domain/value_objects.go` with `UserID` (newtype, empty rejected) and `StatField` (int enum, parsed via `NewStatField`).
- [x] 2.2 Rewrote `backend/internal/domain/user_stats.go`: GORM tags stripped, `UserID` field typed as `domain.UserID`, `IncrementField(StatField, int)` method clamps at zero.
- [x] 2.3 Created `backend/internal/infrastructure/persistence/gorm_models.go` with unexported `gormUserStats` (all original GORM tags preserved) + `TableName()`.
- [x] 2.4 Created `backend/internal/infrastructure/persistence/user_stats_mapper.go` (`userStatsToDomain` / `userStatsFromDomain`).
- [x] 2.5 Created `backend/internal/infrastructure/persistence/registry.go` returning `*gormUserStats`. Deleted `backend/internal/domain/registry.go`. Wired `cmd/server/main.go` to `persistence.AllEntities()`.
- [x] 2.6 `go list -deps ./internal/domain | grep -E "gorm|database/sql|net/http|backend/internal"` returns only the package itself — zero forbidden dependencies.
- [x] 2.7 `go build ./... && go vet ./... && go test ./...` — green (185 tests pass).

_Reshuffle notes for Phase 2:_
- Pulled task 3.2 forward: `persistence.UserStatsRepository` is implemented now (satisfies `application.StatsRepository`, has compile-time assertion). Necessary because the OLD `internal/repository/user_stats.go` had to keep delegating somewhere once domain went pure.
- Old `internal/repository/user_stats.go` is now a thin facade (delegates to persistence repo, keeps `IncrementField` because handler still calls it). Phase 3 deletes the facade.
- `IncrementField` clamping moved from facade into domain method (`UserStats.IncrementField`); the wire-level field name is parsed via `domain.NewStatField`.
- Handler's `*UserStatsResponse.UserID` is `string`, domain's is `domain.UserID` — added `string(stats.UserID)` conversions at the two boundary points in `handler/api.go`.

## 3. Phase 3 — Application layer & handler rewire

- [x] 3.1 Updated `application/ports.go` — `StatsRepository.GetOrCreate` now takes `domain.UserID`. Also added `UserDirectory` port (`UserByID`, `HasKnownDevice`) for webhook handler's Better Auth lookups.
- [x] 3.2 _(landed in Phase 2)_ `backend/internal/infrastructure/persistence/user_stats_repo.go` implements `application.StatsRepository` with compile-time assertion.
- [x] 3.3 Created `application/stats_usecases.go` with `GetUserStats` and `IncrementStatField`, both as structs with `Execute(ctx, ...)` and interface-typed dependencies.
- [x] 3.4 Rewrote `handler/api.go`: constructor now takes `*GetUserStats` + `*IncrementStatField`; handlers call `Execute()` and validate inputs via `domain.NewUserID` / `domain.NewStatField` at the HTTP boundary.
- [x] 3.5 `go list -deps ./internal/handler | grep -E "gorm.io|backend/internal/infrastructure"` returns empty — handler has zero gorm and zero persistence imports.
- [x] 3.6 `infrastructure/persistence/user_stats_repo.go` contains only persistence translation (`GetOrCreate`, `Save` via mapper + GORM). All clamping and field parsing live in domain.
- [x] 3.7 `internal/repository/user_stats.go` and the empty `internal/repository/` directory deleted.
- [x] 3.8 Added `application/stats_usecases_test.go` with `fakeStatsRepo` + `fakeBroadcaster`. Four tests cover happy paths, clamping invariant, and Save-failure / no-broadcast semantics. All pass.
- [x] 3.9 `go build && go vet && go test ./...` — 189 tests pass across 17 packages. `deadcode` reports only pre-existing unreachables (`webhook.ComputeHMAC`, `APIHandler.GetAuthMiddleware` — both unrelated to this refactor).

_Scope expansion in Phase 3:_ `webhook.go` was also refactored to consume `application.UserDirectory` instead of `*gorm.DB`, with concrete `persistence.UserDirectoryRepository` adapter. Added `domain.User` for the projection. This was the only way to satisfy the spec's "no `gorm.io/...` imports in `handler/`" scenario without scoping it down. Better Auth tables stay where they are — only the access path is wrapped.

## 4. Phase 4 — Composition root & River workers

- [ ] 4.1 Create `backend/internal/composition/composition.go` exporting `Build(cfg Config) (*App, error)` (or similar) that constructs: DB connection, persistence registry, repositories, use cases, SSE broker, River client, handlers
- [ ] 4.2 Move all repository / use-case / handler constructor calls out of `backend/cmd/server/main.go` into the composition root
- [ ] 4.3 Confirm `backend/cmd/server/main.go` contains: env loading, calling `composition.Build`, registering routes (or letting composition return a `http.Handler`), `http.ListenAndServe`, and graceful shutdown — and nothing else. Target: <150 lines.
- [ ] 4.4 Refactor River worker registration (`backend/internal/jobs/registry.go` and call sites in `cmd/server/main.go`) so workers receive application interfaces instead of `*gorm.DB`. Add a small `application.JobEnqueuer` port if any use case enqueues jobs.
- [ ] 4.5 Confirm SSE broker is consumed via the `application.EventBroadcaster` interface, not as a concrete type, in every use case that broadcasts
- [ ] 4.6 `cd backend && go build ./... && go test ./... && just lint` — green
- [ ] 4.7 Manual smoke test: `just dev`, log in, exercise dashboard (stats), trigger SSE update (e.g. via webhook or repeated request), trigger one River job (e.g. magic-link request) — all behave as before

## 5. Phase 5 — Tooling, docs, cleanup

- [ ] 5.1 Update `backend/.goca.yaml`: point `architecture.layers.usecase.directory` at `internal/application`, point `architecture.layers.repository.directory` at `internal/infrastructure/persistence`. Domain and handler directories stay.
- [ ] 5.2 Update `CLAUDE.md`: replace the "internal/usecase/ doesn't exist yet" note with the new four-layer description (domain / application / infrastructure / composition), the inward-dependency rule, and the updated "add a new feature" workflow (Goca + manual mapper)
- [ ] 5.3 Update `backend/`'s Goca example in CLAUDE.md so new entities land in `domain/` (pure) and a mapper stub goes into `infrastructure/persistence/`. Document the manual post-`goca make entity` steps: (a) strip GORM tags from the generated domain entity, (b) add a `gormFoo` to `infrastructure/persistence/gorm_models.go`, (c) write the mapper, (d) append the model to `persistence.AllEntities()`.
- [ ] 5.4 Confirm the Swagger output is unchanged: `just api` produces a `backend/docs/swagger.json` with the same routes and schemas as before this change (diff against the pre-refactor file)
- [ ] 5.5 Confirm frontend builds without modification: `cd frontend && bun run lint && bunx tsc --noEmit`
- [ ] 5.6 Final pass: `just lint`, `just typecheck`, `cd backend && go test ./...`, manual smoke test of `just dev` — everything green

## 6. Cleanup & dead-code audit (added per user request)

- [ ] 6.1 Run `deadcode ./...` from `backend/` after Phase 4 — confirm no new unreachable functions were introduced by this refactor. (`NewUserID` becomes reachable via use cases in Phase 3.) Pre-existing dead code in logger/middleware/config is out of scope for this PR.
- [ ] 6.2 Verify no orphaned directories: after Phase 3 deletes `internal/repository/`, ensure no stale references remain (grep for `"backend/internal/repository"` imports).
- [ ] 6.3 Verify no duplicate types: `domain.UserStats` ↔ `persistence.gormUserStats` are intentional twins via mapper. No other duplicates expected.
- [ ] 6.4 Open follow-up issue (or note in CLAUDE.md `## Follow-ups`) for: install `knip` for frontend dead-export detection (Steiger covers FSD rules but not unused exports), and wire `deadcode`/`staticcheck` into the project's quality recipes (e.g. `just deadcode`).

## 7. Verification & archive

- [ ] 7.1 Run `/opsx:verify` and resolve any flagged gaps
- [ ] 7.2 Run `openspec validate backend-clean-architecture --strict`
- [ ] 7.3 Run `/opsx:archive` once verification passes and the change has shipped
