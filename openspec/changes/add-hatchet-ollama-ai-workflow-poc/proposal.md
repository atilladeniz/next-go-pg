## Why

The orchestrator-decision spike (`.docs/orchestrator-decision.md`, closes #57) committed us to adopting Hatchet for AI agent workflows when the first concrete proposal lands. This is that proposal — a "summarize a Git repository" feature that exercises every primitive we expect future AI workflows to need (durable multi-step pipeline, fan-out over LLM calls, per-step retry policies, SSE-progress, idempotent retries on paid-by-the-token operations).

Doing this end-to-end now, even before any AI feature ships to users, gives us:

1. **A working `hatchet-lite` deployment on Docker Compose** that we can later promote to Kamal accessory once we are ready to ship.
2. **A working local `ollama` runtime** with a state-of-the-art small model (`gemma4:e4b`), so future AI features can be prototyped without paying for cloud LLMs during dev.
3. **A `internal/aiworkflows/` bounded context** that becomes the template for every future AI feature (DeepWiki-style indexer, agentic researcher, etc.) — the same way `internal/stats/` is the template for new aggregates.
4. **Confidence that the Phase 1 PoC of the orchestrator plan is buildable on our actual stack** before we commit it to production infra.

## What Changes

- **Add `hatchet-lite` to `infra/compose/docker-compose.dev.yml`** with a dedicated `hatchet` Postgres database, Postgres-only mode (no RabbitMQ), gRPC on 7077 (internal docker network only), dashboard on 8888 (exposed to localhost).
- **Add `ollama` to `infra/compose/docker-compose.dev.yml`** with a persisted `ollama_data` volume, default model pulled on first start (`OLLAMA_MODEL=gemma4:e4b`, configurable via env).
- **New bounded context `backend/internal/aiworkflows/`** with full four-layer DDD layout:
  - `domain/`: `RepoSummary` aggregate, value objects (`RepoURL`, `Status`, `FileSummary`), domain events (`SummaryStarted`, `FileSummarized`, `SummaryCompleted`).
  - `application/`: ports (`HatchetEnqueuer`, `OllamaClient`, `RepoCloner`, `Store`), use cases (`SummarizeRepo`, `GetRepoSummary`).
  - `infrastructure/`: GORM persistence (`RepoSummaryModel` + mapper), Hatchet adapter (`sdks/go`), Ollama HTTP adapter, git cloner.
  - `interfaces/http/`: `POST /api/v1/ai/summarize-repo`, `GET /api/v1/ai/summaries/{id}` with Swagger annotations.
- **Hatchet workflow definition (`SummarizeRepoWorkflow`)** with five durable steps and per-step retry policies:
  - `Clone` (3× linear, fail-fast on auth/network),
  - `Traverse` (deterministic, no retry),
  - `SummarizeFile` (fan-out per file, 5× exponential backoff),
  - `Aggregate` (3× linear),
  - `Store` (3× linear).
- **Composition root wiring** for the new context, including registering `aipersist.Entities()` for AutoMigrate, building the Hatchet worker, and the existing SSE broker as the progress publisher.
- **Frontend feature slice `frontend/src/features/ai-summarize/`** with progress UI (German), Orval-generated hooks, new page under `app/(protected)/ai/summarize/`.
- **SSE event type `ai-progress`** keyed by run ID; reuses the existing platform SSE broker (`internal/platform/sse/`).
- **New Prometheus counter** `ai_workflows_completed_total{status="success|failed|cancelled"}`.
- **Documentation update** to `.docs/orchestrator-decision.md` marking Phase 1 (PoC alongside River) as in-progress / completed once this lands.

## Capabilities

### New Capabilities

- `ai-workflows-platform`: Hatchet engine integration + Ollama runtime + the cross-cutting workflow worker registration pattern. This is the platform-level capability — the new bounded context's infrastructure surface plus the Compose services.
- `repo-summarization`: The user-facing feature itself. Enqueue a summarization run for a Git repo URL, follow progress via SSE, retrieve the final summary.

### Modified Capabilities

None. No existing canonical specs exist in `openspec/specs/`, and no current capability has spec-level requirements that change.

## Impact

- **Code:** New backend bounded context `backend/internal/aiworkflows/` (~6 files across four layers, similar shape to `internal/stats/`). New frontend feature `frontend/src/features/ai-summarize/`. New page `app/(protected)/ai/summarize/page.tsx`. Composition root updates.
- **APIs:** Two new HTTP endpoints (`POST /api/v1/ai/summarize-repo`, `GET /api/v1/ai/summaries/{id}`). One new SSE event type (`ai-progress`). Regenerated Orval client + frontend hooks.
- **Database:** One new GORM-managed table (`repo_summaries`) plus a dedicated `hatchet` Postgres database for Hatchet's own schema. Both created via AutoMigrate / Hatchet migrator on container start.
- **Dependencies (Go):** `github.com/hatchet-dev/hatchet/sdks/go` (workflow engine SDK), `github.com/go-git/go-git/v5` or shell-out to `git` for `RepoCloner` (TBD in design.md).
- **Dependencies (infra):** `hatchet-lite` Docker image, `ollama/ollama` Docker image. Both pinned to a specific tag/SHA in `infra/compose/docker-compose.dev.yml`.
- **Disk:** Ollama model storage (~5 GB for `gemma4:e4b`) in a named docker volume.
- **Observability:** Hatchet engine stdout logs flow through existing Promtail → Loki pipeline. No OTel for the Go SDK (vendor gap, Python-only in May 2026 — documented in orchestrator-decision.md).
- **Production deploy:** **Not in scope here.** `infra/kamal/deploy.yml` is unchanged. This PoC is dev-only. Promoting `hatchet-lite` + `ollama` (or a hosted equivalent) to production is Phase 2 of the orchestrator plan, in a follow-up change.
- **Operational risk on dev:** Two new containers + ~5 GB model download. First `docker compose up` takes longer; documented in `infra/README.md`.
