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

- [x] 4.1 Created `composition/composition.go` exporting `Build(ctx, Inputs) (*App, error)` and `(*App).Shutdown(ctx)`. Constructs DB (with retry), AutoMigrate, persistence repos, use cases, SSE broker, River client, handlers, middleware, full HTTP router, and `http.Server`. Health endpoints live as a private `healthEndpoints` type inside composition.
- [x] 4.2 All wiring (DB connect, repos, use cases, broker, River, handlers, router) moved out of `cmd/server/main.go` into composition.
- [x] 4.3 `cmd/server/main.go` is **91 lines** (target was <150). Contents: swagger annotations, build vars, config load + validate, logger init, metrics init, `composition.Build`, goroutine to run server, signal handling, `app.Shutdown`.
- [x] 4.4 Workers wired through the application port: `WorkerDeps.Events` is `application.EventBroadcaster`, `WorkerDeps.StatsRepo` is `application.StatsRepository`. `DataExportWorker` no longer takes `*sse.Broker`. No new `JobEnqueuer` port needed — the existing `jobs.JobEnqueuer` is sufficient and webhook uses it directly.
- [x] 4.5 Only one `*sse.Broker` reference remains outside the SSE package itself: `composition.routerDeps.sseBroker`, used to register the broker AS the `/events` route handler (the broker implements `http.Handler`). Every other consumer (use case, worker) sees only `application.EventBroadcaster`. Verified via `grep -rn "sse.Broker" backend --include="*.go" | grep -v _test.go`.
- [x] 4.6 `go build ./... && go vet ./... && go test ./...` — green (189 tests, 17 packages).
- [ ] 4.7 _Manual smoke test pending — runs after Phase 5._ Confirms: `just dev`, login flow, dashboard stats display, stats SSE update, magic-link delivery via River.

_Cleanup landed in Phase 4:_ Removed the unused `sseBroker` field from `APIHandler` (broadcasting now lives in the `IncrementStatField` use case). Constructor signature: `NewAPIHandler(frontendURL, *GetUserStats, *IncrementStatField)`.

## 5. Phase 5 — Tooling, docs, cleanup

- [x] 5.1 `.goca.yaml` updated — `architecture.layers.usecase.directory` → `internal/application`, `architecture.layers.repository.directory` → `internal/infrastructure/persistence`. Domain and handler unchanged.
- [x] 5.2 CLAUDE.md rewritten: new four-layer tree, explicit inward-dependency arrows, replaced "`internal/usecase/` doesn't exist yet" note with a layer-ownership summary that points at the refactor.
- [x] 5.3 Updated the Goca "Adding a New Feature" walkthrough — six explicit steps including manual GORM-tag relocation, mapper authoring, registry update, port declaration with compile-time assertion, use-case wiring, and composition-root touch.
- [x] 5.4 `just swagger` produces `backend/docs/swagger.json` **byte-identical** to the pre-refactor version (`diff` on 26315-byte file: zero output).
- [x] 5.5 Frontend: `bun run lint` clean (Biome + Steiger, 75 files), `bunx tsc --noEmit` exits 0.
- [x] 5.6 Final pass: `go build && go vet && go test ./...` green (189 tests, 17 packages). Frontend lint + typecheck green. Manual smoke test of `just dev` to be run by maintainer before merge.

## 6. Cleanup & dead-code audit (added per user request)

- [x] 6.1 `deadcode ./...` filtered to refactored areas: **zero** new unreachable functions from this refactor. `NewUserID` is reachable via handlers + tests. Killed two leftovers found mid-audit: `APIHandler.GetAuthMiddleware()` + `authMiddleware` field (composition wires `combinedAuth` separately) and the unused `betterAuthURL` constructor parameter.
- [x] 6.2 No orphaned dirs: `internal/repository/` gone, no remaining imports of `"backend/internal/repository"` (`grep -rn` returns empty).
- [x] 6.3 No duplicate types: only intentional twins are `domain.UserStats` ↔ `persistence.gormUserStats` and `domain.User` (projection used by the `UserDirectory` port). All connected via mappers / `betterAuthUserRow`.
- [ ] 6.4 Follow-up: add `knip` (frontend dead-export detection) and wire `deadcode`/`staticcheck` into a `just deadcode` recipe. Captured here as a known follow-up — out of scope for this refactor since it's tooling addition, not architecture work.

## 7. Verification & archive

- [ ] 7.1 Run `/opsx:verify` and resolve any flagged gaps
- [ ] 7.2 Run `openspec validate backend-clean-architecture --strict`
- [ ] 7.3 Run `/opsx:archive` once verification passes and the change has shipped
