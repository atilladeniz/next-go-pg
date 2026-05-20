## 1. Phase A — Compose + smoke test

- [ ] 1.1 Add `hatchet-lite` service to `infra/compose/docker-compose.dev.yml` (pin image SHA, gRPC `7077` on internal docker network only, dashboard `8888` exposed to localhost, env `SERVER_URL=postgres://postgres:postgres@db:5432/hatchet?sslmode=disable`, depend on `db` service)
- [ ] 1.2 Add `CREATE DATABASE hatchet;` step to the dev Postgres init (either an `init-db.sql` mounted into `db` or a one-shot init container) so the `hatchet` database exists before `hatchet-lite` migrates
- [ ] 1.3 Add `ollama` service to `infra/compose/docker-compose.dev.yml` (`ollama/ollama` image pinned to a specific tag, named volume `ollama_data:/root/.ollama`, port `11434` exposed on the docker network only)
- [ ] 1.4 Add `COMPOSE_PROFILES` support so contributors can opt out of `ollama` + `hatchet-lite` (profile name `ai`) if they only need the core dev stack
- [ ] 1.5 Add `just ai-up` and `just ai-down` recipes wrapping `docker compose --profile ai up/down`
- [ ] 1.6 Verify: `just dev` brings up core services; `just ai-up` brings up Hatchet + Ollama; Hatchet dashboard reachable at `http://localhost:8888`; `curl http://localhost:11434/api/tags` returns 200
- [ ] 1.7 Update `infra/README.md` with: AI dev stack setup, model pull instructions, expected resource usage, known gotchas (broadcast address fragility, first-pull time)

## 2. Phase B — Backend skeleton: domain layer

- [ ] 2.1 Create directory `backend/internal/aiworkflows/{domain,application,infrastructure,interfaces}/` (mirror `internal/stats/` layout)
- [ ] 2.2 `domain/repo_url.go`: `RepoURL` value object with `NewRepoURL(string)` constructor that validates URL shape (http/https only, ends with `.git` optional, no shell metachars)
- [ ] 2.3 `domain/status.go`: `Status` value object with constants `StatusPending`, `StatusRunning`, `StatusCompleted`, `StatusFailed`, `StatusCancelled`; constructor rejects unknown values
- [ ] 2.4 `domain/file_summary.go`: `FileSummary` value object with `NewFileSummary(filename, summary string)` constructor; rejects empty filename
- [ ] 2.5 `domain/repo_summary.go`: `RepoSummary` aggregate embedding `shared.AggregateBase`; fields `ID`, `UserID`, `RepoURL`, `Status`, `Files []FileSummary`, `RepoSummary string`, `StartedAt`, `CompletedAt`; methods `MarkStarted`, `AppendFileSummary`, `MarkCompleted(summary string)`, `MarkFailed(reason string)`, `MarkCancelled`
- [ ] 2.6 `domain/events.go`: `SummaryStarted`, `FileSummarized`, `SummaryCompleted`, `SummaryFailed`, `SummaryCancelled` domain events; all implement `EventName() string`
- [ ] 2.7 Domain unit tests in `domain/repo_summary_test.go`: aggregate transitions raise correct events, value-object invariants hold, illegal state transitions return errors
- [ ] 2.8 Verify: `go build ./internal/aiworkflows/domain/...` clean; `go test ./internal/aiworkflows/domain/...` green

## 3. Phase B continued — Application layer

- [ ] 3.1 `application/ports.go`: declare ports `HatchetEnqueuer`, `LLMClient`, `RepoCloner`, `Store`, `ProgressPublisher` (domain-event-driven SSE adapter)
- [ ] 3.2 `application/usecases.go`: `SummarizeRepo` use case (validates input, enqueues workflow, returns run ID) and `GetRepoSummary` use case (loads aggregate, enforces ownership)
- [ ] 3.3 `application/usecases_test.go`: table-driven tests with mocked ports — happy path, ownership violation returns not-found, invalid URL returns error
- [ ] 3.4 Verify: `go test ./internal/aiworkflows/application/...` green

## 4. Phase C — Workflow definition

- [ ] 4.1 Add `github.com/hatchet-dev/hatchet/sdks/go` to `go.mod` (`go get`), pin to a specific version
- [ ] 4.2 `infrastructure/workflows/types.go`: typed step inputs/outputs (`CloneInput`, `CloneOutput`, `TraverseInput`, `TraverseOutput`, `SummarizeFileInput`, `SummarizeFileOutput`, `AggregateInput`, `AggregateOutput`, `StoreInput`, `StoreOutput`)
- [ ] 4.3 `infrastructure/workflows/steps.go`: each step function signature accepts the typed input + dependencies, returns the typed output; step functions are pure adapters that call ports
- [ ] 4.4 `infrastructure/workflows/workflow.go`: register `SummarizeRepoWorkflow` with five steps in the correct DAG order, with per-step retry policies (Clone: 3× linear; Traverse: no retry; SummarizeFile: 5× exponential with fan-out per file from Traverse output; Aggregate: 3× linear; Store: 3× linear)
- [ ] 4.5 `infrastructure/workflows/enqueuer.go`: implement `application.HatchetEnqueuer` port; assert `var _ aiapp.HatchetEnqueuer = (*Enqueuer)(nil)`
- [ ] 4.6 `infrastructure/workflows/worker.go`: bootstrap a worker that registers the workflow and consumes from Hatchet; called from composition root
- [ ] 4.7 Smoke test: build the binary, start Hatchet + worker locally, register the workflow, verify it appears in the dashboard at `localhost:8888` with all five steps visible (steps can still be no-ops at this point)
- [ ] 4.8 Verify: `go build ./internal/aiworkflows/...` clean; manual dashboard check

## 5. Phase D — Infrastructure adapters

- [ ] 5.1 Add `github.com/go-git/go-git/v5` to `go.mod`
- [ ] 5.2 `infrastructure/git/cloner.go`: `RepoCloner` adapter using `go-git` for shallow clones; enforces size limit (~50 MB unpacked, fail with clear error if exceeded); cleans up `/tmp` directory on context cancellation
- [ ] 5.3 `infrastructure/llm/ollama.go`: `LLMClient` adapter — HTTP client to `http://ollama:11434/api/generate`; configurable via `OLLAMA_URL` and `OLLAMA_MODEL` env; sensible timeouts (60 s default, override via env)
- [ ] 5.4 `infrastructure/persistence/{model,mapper,repository,registry}.go`: GORM `RepoSummaryModel` with `gorm.Model` base + columns for `UserID`, `RepoURL`, `Status`, `Files JSONB`, `RepoSummary`, `StartedAt`, `CompletedAt`; mapper to/from domain; repository implements `application.Store`; `Entities() []any` registry
- [ ] 5.5 `infrastructure/events/publisher.go`: domain-event → SSE adapter publishing `ai-progress` events to the existing platform broker
- [ ] 5.6 Wire workflow steps in `infrastructure/workflows/steps.go` to the real adapters (replace any no-ops from Phase C)
- [ ] 5.7 End-to-end manual run against a small public repo (e.g. `https://github.com/atilladeniz/next-go-pg`): confirm clone, traverse, fan-out summarize, aggregate, store all complete; verify Hatchet dashboard shows step inputs/outputs; verify persisted `RepoSummary` in DB
- [ ] 5.8 Mid-run crash test: trigger a workflow, `docker kill` the backend worker container during SummarizeFile, restart, confirm workflow resumes from the last incomplete file (does not re-run Clone or Traverse)

## 6. Phase E — HTTP interface

- [ ] 6.1 `interfaces/http/handler.go`: `POST /api/v1/ai/summarize-repo` (validates body, calls `SummarizeRepo` use case, returns 202 with run ID); `GET /api/v1/ai/summaries/{id}` (calls `GetRepoSummary` use case, returns 404 on missing or other-user's run, 200 with payload otherwise)
- [ ] 6.2 Add Swagger annotations (`@Summary`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`) to both handlers
- [ ] 6.3 Register routes in composition root behind the `combinedAuth` middleware
- [ ] 6.4 `interfaces/http/handler_test.go`: table-driven HTTP tests with mocked use cases — happy paths, auth failure, validation failure, ownership violation
- [ ] 6.5 Run `just api` to regenerate Swagger + Orval client
- [ ] 6.6 Verify: `go test ./internal/aiworkflows/...` all green; new endpoints visible in Swagger UI

## 7. Phase E continued — Composition root + Prometheus

- [ ] 7.1 Compose adapters in `internal/composition/composition.go`: instantiate `Cloner`, `OllamaClient`, `Repository` (when DB available), `Enqueuer` (when Hatchet client available), `ProgressPublisher` (from existing SSE broker); build `SummarizeRepo` + `GetRepoSummary` use cases; build the HTTP handler
- [ ] 7.2 Register `aipersist.Entities()` in `runAutoMigrations`
- [ ] 7.3 Start the Hatchet worker as a goroutine in `composition.Build`, wire to `App.Shutdown` for clean exit
- [ ] 7.4 Add `metrics.AIWorkflowsCompleted = promauto.NewCounterVec(...)` with `status` label in `backend/pkg/metrics/metrics.go`
- [ ] 7.5 Increment the counter in the `Store`/`Failed`/`Cancelled` use-case branches
- [ ] 7.6 Add a feature-flag-style gate: if `HATCHET_CLIENT_TOKEN` is unset, log a warning and continue running without the AI feature wired (so the dev stack still boots without the AI profile)
- [ ] 7.7 Verify: `go build ./...` clean; `go test ./... -race` green for all packages

## 8. Phase F — Frontend

- [ ] 8.1 Create `frontend/src/features/ai-summarize/` (FSD slice) with `ui/`, `model/`, `index.ts`
- [ ] 8.2 `model/use-summarize.ts`: hook wrapping the Orval-generated `usePostAiSummarizeRepo` mutation
- [ ] 8.3 `model/use-summary.ts`: hook wrapping `useGetAiSummariesById`
- [ ] 8.4 `model/use-ai-sse.ts` or extension of existing `features/stats/model/use-sse.ts`: listen for `ai-progress` events and invalidate the summary query on transitions
- [ ] 8.5 `ui/summarize-form.tsx`: input + submit button, German labels ("Repository-URL", "Zusammenfassen", validation messages)
- [ ] 8.6 `ui/summary-progress.tsx`: progress bar showing current step in German ("Repository klonen…", "Dateien analysieren…", "Zusammenfassen…", "Aggregieren…", "Speichern…"), plus per-file counter for SummarizeFile
- [ ] 8.7 `ui/summary-view.tsx`: render finished summary (repo-level text + collapsible per-file table)
- [ ] 8.8 Create `frontend/src/app/(protected)/ai/summarize/page.tsx` — Server Component using `getSession` for auth gate, `HydrationBoundary` if prefetch makes sense (probably not, since a new run starts empty)
- [ ] 8.9 Add a navigation entry to the Header or Dashboard so the page is reachable
- [ ] 8.10 Verify: `bun run typecheck`, `bun run lint`, manual end-to-end through the browser

## 9. Phase G — Tests + observability + docs

- [ ] 9.1 Integration test in `internal/aiworkflows/infrastructure/llm/ollama_integration_test.go` using testcontainers (Ollama container, model pulled in `TestMain`); marked with build tag `integration` and skipped if `TESTCONTAINERS_RYUK_DISABLED` is set or if the host has no Docker
- [ ] 9.2 Update `infra/README.md`: AI dev stack section, model variants and trade-offs, troubleshooting (token regeneration, model not found, port collisions)
- [ ] 9.3 Update `.docs/orchestrator-decision.md`: mark Phase 1 (PoC alongside River) as **complete**, link to this change in the archive once shipped, list what was learned (footprint observed, gotchas encountered, OTel-stdout decision validated)
- [ ] 9.4 Update `.docs/background-jobs.md` to clarify the split: River for one-shot jobs, Hatchet for multi-step AI workflows
- [ ] 9.5 Add a short `.docs/ai-workflows.md` explaining how to define new workflows, the port pattern, and where new bounded contexts go
- [ ] 9.6 Update `.claude/CLAUDE.md`: new context in the project structure tree, new `.docs/ai-workflows.md` index entry, brief Hatchet+Ollama section under tech stack
- [ ] 9.7 Run full check suite: `just lint`, `just typecheck`, `go test ./... -race`, `bun run test:run`, `just security-scan`
- [ ] 9.8 Manual smoke: full flow through the browser against a small public repo; mid-run crash + restart confirms durable resume; Hatchet dashboard at `localhost:8888` shows all steps with inputs/outputs

## 10. Wrap-up

- [ ] 10.1 PR description summarizing what shipped, what was deliberately out of scope, links to relevant `.docs/` sections
- [ ] 10.2 `/opsx:verify` against the artifacts before archive
- [ ] 10.3 `/opsx:archive` once verified green
- [ ] 10.4 Close issue #57's epic (if still open) referencing this change as Phase 1 completion
- [ ] 10.5 Open follow-up issues for Phase 2 (production Kamal deploy of `hatchet-lite` + cloud LLM substitution), captured as bullet items in the orchestrator-decision doc
