## Context

The orchestrator-decision spike (`.docs/orchestrator-decision.md`) committed us to **adopt Hatchet (via `hatchet-lite`) when the first concrete AI workflow lands**, and to keep River permanently for one-shot stateless work. This change is "the first concrete AI workflow landing" — a repo-summarization feature that exercises every primitive the doc names (durable resume, per-step retry, fan-out, idempotent retries on paid LLM calls). It is also the Phase 1 PoC from the doc, executed end-to-end.

The current state has no AI features. The four background jobs (`send_magic_link`, `send_verification_email`, `send_2fa_otp`, `send_login_notification`) and the data-export job all run on River. The bounded contexts (`stats`, `auth`, `notifications`, `exports`) have a clean DDD shape established by the `refactor/backend-clean-architecture` change. SSE progress is already wired (`internal/platform/sse/broker.go`, the export job publishes progress via `exports/infrastructure/jobs/worker.go`).

Constraints from the existing project state:

- **Single-replica Kamal deploy** is the production target. The PoC stays dev-only; no Kamal changes here.
- **DDD bounded contexts are the architectural template.** `internal/stats/` is the canonical small example. The new context follows it verbatim where possible.
- **SSE broker is platform-cross-cutting.** New workflow steps publish via the existing broker; we do not introduce a new transport.
- **No OTel for Hatchet Go SDK** (vendor gap, May 2026). Logs go stdout → Promtail → Loki.
- **Postgres-only mode** for Hatchet (no RabbitMQ — supported up to ~100 events/sec, well above our scale).

## Goals / Non-Goals

**Goals:**

- Ship a working `hatchet-lite` + `ollama` Docker Compose dev stack that another contributor can spin up with `just dev` + a one-time `ollama pull gemma4:e4b`.
- Produce a `internal/aiworkflows/` bounded context that is structurally identical in shape to `internal/stats/` and `internal/exports/` (same four-layer DDD layout, same `Entities()` registry pattern, same SSE-via-publisher pattern).
- Implement an end-to-end repo summarization workflow that proves all five workflow primitives we expect future AI features to need: durable resume mid-step, fan-out over LLM calls, per-step retry policies, paid-token idempotency, SSE-progress.
- Document the Hatchet operational reality (write amplification, no OTel, broadcast-address fragility) in code comments and `infra/README.md` so future maintainers do not re-learn it.

**Non-Goals:**

- Production Kamal deploy of `hatchet-lite` or `ollama`. That is Phase 2 of the orchestrator plan, in a follow-up change.
- Multi-user isolation of workflow runs beyond the basic "users only see their own runs" check. No tenant boundaries, no rate-limit tier per workflow.
- LLM prompt engineering. Use simple "summarize this file in 3 sentences" / "given these per-file summaries, give a 5-sentence repo overview" prompts. v2 territory.
- Caching of identical repo summaries. Every request is fresh.
- Cancellation UI. Workflows auto-cancel on context propagation; no explicit cancel button.
- Cost tracking. With local Ollama there is no per-token bill; the abstraction can be added when we add paid LLMs.
- Migration of any existing River job to Hatchet. River remains the queue for one-shot work, permanently.

## Decisions

### Decision 1: Use `hatchet-lite` (not full distributed Hatchet) for the dev stack

**Why:** `hatchet-lite` bundles engine + API + dashboard + migrator into one Docker image. The Hatchet team's own framing: "designed for development and low-volume use-cases." Matches our scale (one VPS in production, dev laptops in development). Distributed mode adds RabbitMQ + four separate containers and is overkill until we cross ~100 events/sec.

**Alternative considered:** Full distributed Hatchet (engine, API, frontend, queue separately). Rejected for now; revisit only if we ever cross the throughput threshold or need multi-region.

### Decision 2: Dedicated `hatchet` Postgres database on the existing dev Postgres container

**Why:** Keeps `pg_dump` boundaries clean; isolates Hatchet's table partitioning and autovacuum behaviour from the app schema. Matches the orchestrator-decision doc Phase 1 recommendation. No new Postgres container required — same instance, two databases (`nextgopg` and `hatchet`).

**Alternative considered:** Share the app database. Rejected because Hatchet's write amplification (~5 transactions per task) and 30-day-default retention partitioning would clutter `pg_dump` and complicate disaster recovery. Revisit only after Phase 2 production operations show the dedicated DB is unnecessary overhead.

### Decision 3: Ollama in dev, no cloud-LLM fallback

**Why:** The PoC's entire point is to validate the workflow infrastructure. Cloud-LLM integration is orthogonal — we can swap the `OllamaClient` adapter for an `OpenAIClient` (or anything else) implementing the same `application.LLMClient` port later. Local Ollama eliminates per-token cost during dev and removes a network dependency for CI/integration tests.

**Alternative considered:** Direct OpenAI API client. Rejected for dev; the port stays open for production substitution.

### Decision 4: `gemma4:e4b` as the default model

**Why:** Research (May 2026) shows Gemma 4 family is current state-of-the-art for small local LLMs. E4B (4B params, ~5 GB at Q4) runs without GPU on a typical dev laptop, supports 256K context (helpful for whole-file summarization), and outperforms Qwen 3.5/3.6 at the same size on coding benchmarks (LiveCodeBench 80% vs ~43%). Configurable via `OLLAMA_MODEL` env so a contributor with more RAM can swap to `gemma4:26b-a4b` (the best local coding model per real-world test, ~14-18 GB at Q4).

**Alternative considered:** Qwen 3.5 27B (similar size class, weaker coding benchmarks). Llama 3.2 (older, behind on coding). Phi-4-mini (smaller, less capable for code). All rejected in favour of Gemma 4 E4B as default with 26B-A4B as the heavier opt-in.

### Decision 5: New bounded context `internal/aiworkflows/`, copy structure from `internal/stats/`

**Why:** The DDD layout is established and documented (`.claude/CLAUDE.md`). `internal/stats/` is the canonical small example. Copying its shape minimizes review surface area — the differences will be in domain types and infrastructure adapters, not in directory layout.

**Alternative considered:** Reuse `internal/exports/` as the template (closer feature shape: SSE progress, fan-out). Rejected because `exports/` has more legacy and an in-memory store; `stats/` has the cleaner aggregate + persistence pattern.

### Decision 6: Hatchet `Enqueuer` as an `application` port, adapter in `infrastructure/workflows/`

**Why:** Mirrors the `notifications/application.JobEnqueuer` port that fronts River. The application layer never imports `sdks/go`; the infrastructure adapter does. Lets us swap Hatchet for a different workflow engine in the future without touching application code (this is the explicit migration-path-from-anything-else option in the orchestrator-decision doc).

**Alternative considered:** Directly import `sdks/go` in the HTTP handler. Rejected — violates the inward-only dependency rule and makes the workflow engine a load-bearing choice in the application layer.

### Decision 7: Workflow steps are first-class typed structs in `infrastructure/workflows/steps.go`

**Why:** Hatchet's Go SDK uses typed step input/output. Defining `CloneInput`, `CloneOutput`, `SummarizeFileInput`, etc. as concrete types in `infrastructure/` makes the workflow DAG legible at a glance and serializable in the dashboard. Keep them out of `application/` so the port surface stays workflow-engine-agnostic.

**Alternative considered:** Anonymous map-based step inputs. Rejected — loses type safety and dashboard display fidelity.

### Decision 8: SSE progress events via the existing platform broker, new event type `ai-progress`

**Why:** The `internal/platform/sse/broker.go` is already cross-cutting and used by `stats` + `exports`. Adding `ai-progress` as a new event type keeps the same transport pattern. Frontend `useSSE()` hook extends for the new event with no transport changes.

**Alternative considered:** A workflow-engine-native progress mechanism (Hatchet has its own dashboard/event stream). Rejected because the user-facing UI lives in our frontend; piping progress through our own SSE broker keeps that flow simple and decouples our UI from Hatchet's eventing.

### Decision 9: `RepoCloner` uses `go-git/go-git` (pure Go), not shell-out to `git`

**Why:** `go-git/go-git` is pure Go, no binary dependency on the runner. Container image stays smaller. Easier to test. Adequate for our use (shallow clone, single ref).

**Alternative considered:** `os/exec git clone --depth=1`. Faster runtime in some cases, but adds a binary dependency on the image. Rejected unless we hit performance issues we cannot tune.

### Decision 10: Workflow run state lives in the aggregate; Hatchet is the executor not the source of truth

**Why:** The `RepoSummary` aggregate holds canonical state (status, file summaries, final summary). Hatchet's run history is for debugging and replay; queries from the frontend hit our `Store` port, not Hatchet's API. This keeps Hatchet swappable.

**Alternative considered:** Frontend queries Hatchet's REST API directly. Rejected — couples frontend to Hatchet, breaks the engine-swappability commitment.

## Risks / Trade-offs

- **[Risk] First-`docker compose up` is slow** (Ollama pulls ~5 GB for `gemma4:e4b`, Hatchet image + migrator startup is ~30s) → **Mitigation:** Document explicitly in `infra/README.md`; add a `just ai-up` recipe that surfaces progress. Model is cached after first pull.
- **[Risk] Hatchet write amplification on Postgres** (5 tx per task × N fan-out files) → **Mitigation:** Dedicated `hatchet` database isolates impact; document the autovacuum tuning gate in `infra/README.md`. Within our scale this is non-issue, but the comment prevents footguns when someone later increases fan-out limits.
- **[Risk] OTel gap for Go SDK** → **Mitigation:** Stdout → Promtail → Loki via existing infra; dashboard at `:8888` covers per-step debugging. Document the gap in the workflow worker bootstrap comments so nobody wastes time looking for OTel config that isn't there.
- **[Risk] GRPC broadcast address is sticky** (changing it later invalidates the token) → **Mitigation:** Pin the broadcast address in Compose env from day one. Use the docker network DNS name (`hatchet`) so it stays stable across host environments.
- **[Risk] `go-git` may be slow on very large repos** → **Mitigation:** PoC scope is small repos (~50 files). Add a size limit at clone time (fail fast if repo exceeds ~50 MB). Production scale-up is a follow-up.
- **[Risk] Workflow engine drift between dev (Compose) and production (Kamal)** → **Mitigation:** Production-deploy story is in scope of a future change; for now, do not let the dev Compose shape suggest production design (e.g. no hard-coded `host.docker.internal` references in workflow code).
- **[Trade-off] Local Ollama eats RAM during dev** (5-18 GB depending on model) → developers can override `OLLAMA_MODEL` to a smaller variant, or skip the AI service entirely with a compose profile (`COMPOSE_PROFILES=core` excludes Ollama).
- **[Trade-off] DDD layout adds ~6 files for what could be ~2** → consistent with `stats/` and `exports/`; predictability beats compactness. Pays back the second time we add an AI feature.

## Migration Plan

This is a **new feature, not a migration**, so there is no rollback of data. Rollback if needed = revert the merge commit; the new Hatchet database + Ollama volume can be left behind or removed via `docker volume rm`.

Phased rollout within the change itself:

1. **Phase A — Compose + smoke test.** Add `hatchet-lite` + `ollama` to `infra/compose/docker-compose.dev.yml`. Verify Hatchet engine boots, dashboard reachable, Ollama responds. Backend wires to nothing yet.
2. **Phase B — Backend skeleton.** Create `internal/aiworkflows/` directory with `domain/` types compiled but no integration. Confirm `go build ./...` passes.
3. **Phase C — Workflow definition + worker.** Register `SummarizeRepoWorkflow` with Hatchet; run a no-op smoke test (each step returns immediately). Confirm Hatchet dashboard shows the workflow definition.
4. **Phase D — Adapters.** Implement `RepoCloner`, `OllamaClient`, `Store`. Wire the workflow steps to real adapters. Run end-to-end against a small public repo.
5. **Phase E — HTTP + SSE.** Add the two HTTP endpoints, wire SSE progress publisher, update Swagger via `just api`.
6. **Phase F — Frontend.** New feature slice, new page, Orval-generated hooks, German UI text.
7. **Phase G — Tests + docs.** Unit tests, application-layer tests, integration test (testcontainers, skippable in CI). Update `infra/README.md`, `.docs/orchestrator-decision.md` (Phase 1 marked complete).

Each phase is independently committable; later phases can be pushed as follow-ups if scope balloons.

## Open Questions

1. **CI integration test strategy.** Do we run the Ollama integration test in CI (download the model, run on GitHub Actions runners) or mark it `// +build integration` and skip in CI? Recommendation: skip in CI by default, run nightly on a self-hosted runner if we have one — TBD during Phase G.
2. **Hatchet dashboard auth in dev.** Default `hatchet-lite` has no auth on the dashboard. For dev that is acceptable (localhost-only). Document the gap in `infra/README.md`; production Phase 2 will need real auth.
3. **`go-git` vs shell `git`.** Default decision is `go-git`. Reassess at Phase D if performance / feature parity issues surface.
4. **Should `RepoSummary.Status` be a value object or a typed enum const?** Value object adds construction invariants; enum const is simpler. Decision: value object with `NewStatus(string)` constructor, consistent with `internal/exports/domain.Status`.
