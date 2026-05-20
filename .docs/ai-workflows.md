# AI Workflows (Hatchet + Ollama)

The `aiworkflows` bounded context owns long-running, multi-step AI work
backed by [Hatchet](https://hatchet.run) (durable workflow engine) and a
local [Ollama](https://ollama.com) runtime. Use it when River is too
small a hammer — see `.docs/background-jobs.md` for when each fits.

This doc covers **how to add a new workflow**. For the rationale on
two-queue split see `.docs/orchestrator-decision.md`.

## Layout

```
backend/internal/aiworkflows/
├── domain/                       # Pure types. No SDK, no I/O.
├── application/                  # Ports + use cases.
│   ├── ports.go                  # HatchetEnqueuer, LLMClient,
│   │                             # RepoCloner, Store, ProgressPublisher
│   └── usecases.go               # SummarizeRepo, GetRepoSummary
├── infrastructure/
│   ├── workflows/                # The ONLY place that imports the SDK
│   │   ├── types.go              # Typed step inputs/outputs
│   │   ├── steps.go              # Step functions, closures over Deps
│   │   ├── workflow.go           # DAG + per-step retry policies
│   │   ├── enqueuer.go           # implements aiapp.HatchetEnqueuer
│   │   └── worker.go             # bootstrap + StartBlocking goroutine
│   ├── git/                      # go-git adapter
│   ├── llm/                      # Ollama HTTP adapter
│   ├── persistence/              # GORM model + repo + Entities()
│   └── events/                   # SSE adapter (domain events → broker)
└── interfaces/http/              # HTTP handlers (Swagger-annotated)
```

The composition root (`internal/composition/composition.go`) is the
only place that wires the SDK client, builds the worker, and gates
everything on `HATCHET_CLIENT_TOKEN` — when the token is missing the
package boots in degraded mode (GET endpoints still answer existing
rows, POST returns 503).

## Add a new workflow — checklist

1. **Domain.** New aggregate + value objects under `domain/` (mirror
   `RepoSummary`). State machine lives here; events implement
   `EventName() string` with the `aiworkflows.` prefix so the SSE
   adapter can route them.

2. **Application port.** Add use-case structs in `application/`. They
   depend on the same `HatchetEnqueuer` port — typed payloads go into a
   new `EnqueueXInput` struct, keeping the SDK out of the application
   layer.

3. **Workflow definition.** In `infrastructure/workflows/`:
   - typed `XInput` / `XOutput` per step
   - one step function per node, closure over `Deps`
   - register in a `Build…Workflow` helper using
     `wf.NewTask(name, fn, hatchet.WithParents(prev), hatchet.WithRetries(N), hatchet.WithRetryBackoff(factor, max))`
   - fan-out → register the child as a `StandaloneTask` and call
     `childTask.Run(ctx, in)` inside the parent step, wrapping in
     goroutines + `sync.WaitGroup`

4. **Worker registration.** Extend `NewWorker` so the new workflow +
   child task are passed to `hatchet.WithWorkflows(...)`.

5. **HTTP handler.** Add Swagger-annotated endpoints under
   `interfaces/http/`. Map `aiapp.ErrNotFound` to 404 — never leak
   existence of other users' runs.

6. **Composition wire.** Extend `buildAIWorkflowsHandler` to instantiate
   the new use case and pass it to the handler. Keep the
   `HATCHET_CLIENT_TOKEN` gate.

7. **Frontend.** Add a FSD slice under `frontend/src/features/ai-*/`
   following `features/ai-summarize/`. Subscribe to the existing
   `ai-progress` SSE event type — no new event channels needed.

8. **OpenSpec.** Non-trivial workflows go through `/opsx:propose` first
   (see `.docs/openspec.md`).

## Step retry-policy recipes

| Step shape                                     | Recipe                                                 |
|------------------------------------------------|--------------------------------------------------------|
| Network fetch / clone                          | `WithRetries(3)` linear, fail fast on 4xx              |
| Deterministic computation, no I/O              | No retries — any error here is a bug                   |
| LLM call (Ollama) — idempotent, cost-free      | `WithRetries(5)` + `WithRetryBackoff(2, 60)`           |
| LLM call against paid API (token cost)         | `WithRetries(3)` + `worker.NewNonRetryableError(err)` on 4xx |
| DB write at end of workflow                    | `WithRetries(3)` linear                                |

## Observability

- **Logs:** `logger.WithContext(ctx)` in each step → stdout → Promtail
  → Loki. Hatchet engine logs flow the same way (`service="hatchet"`).
  No OpenTelemetry for the Go SDK as of May 2026.
- **Metrics:** every terminal run increments
  `ai_workflows_completed_total{status="success|failed|cancelled"}` —
  the events publisher owns the counter so it stays in sync with the
  actual emitted domain events.
- **Dashboard:** Hatchet UI at `http://localhost:8888` shows the
  workflow DAG, every step's input/output, retry history, and the
  current queue depth. Indispensable when debugging.

## Gotchas

- **`HATCHET_CLIENT_TOKEN`** must be generated in the Hatchet dashboard
  (default-tenant flow) and added to your env. Without it the AI
  bounded context boots in degraded mode.
- **Broadcast address** (`SERVER_GRPC_BROADCAST_ADDRESS`) must be
  reachable from anywhere the worker dials Hatchet. The dev compose
  uses `127.0.0.1:7077` because the dev backend runs on the host —
  changing it later requires re-registering workers.
- **Goroutine fan-out** is parallel but bounded by `WithSlots(N)` on
  the worker (default 10 in this repo).
- **Ollama first call** is slow (model load into memory). The retry
  config absorbs this on the first per-file summary.

## See also

- `.docs/orchestrator-decision.md` — why Hatchet, not Temporal/DBOS
- `.docs/background-jobs.md` — when to use River instead
- `infra/README.md` — `just ai-up`, broadcast-address pitfalls
