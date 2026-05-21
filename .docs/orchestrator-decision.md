# Workflow Orchestrator Decision (Spike for #57)

**Date:** 2026-05-20
**Status:** Decision recorded — implementation gated on the first
concrete AI workflow proposal.
**Decision:** **Adopt Hatchet (via `hatchet-lite`)** as the workflow
engine for AI agent features. **Keep River for one-shot stateless work
(emails, exports) — permanently, not transitionally.** Two queues,
one Postgres, separate concerns by workload shape. Do **not** adopt
DBOS Transact in production yet — re-evaluate Nov 2026.

This document is the deliverable of issue #57. It is a *decision-support
artifact* — the implementation will land in follow-up issues, gated on
the first concrete AI workflow proposal (e.g. a DeepWiki-style code
indexer, an agentic researcher, a multi-step LLM pipeline).

**Reference research (May 18 2026):** A deep dive on production-state of
DBOS Transact Go, Hatchet operational reality on small infra, and
head-to-head migration from River informed this version. Cited
empirical claims below come from that research; vendor-marketing
claims are explicitly excluded.

---

## Why this decision changed

The first version of this spike recommended *defer Hatchet until a
concrete AI workflow lands*. That recommendation was wrong. Reviewing
the actual shape of AI agent workflows we are likely to build — using
[Cognition's DeepWiki](https://docs.devin.ai/work-with-devin/deepwiki)
as a representative example — makes clear that River is structurally
the wrong tool for this class of feature, and the right move is to
plan the migration now rather than retrofitting later.

### What a DeepWiki-style feature looks like as a workflow

A repo-indexer pipeline has at least five discrete stages per repo:

1. **Clone** — network I/O, can fail on auth / network errors.
2. **File traversal + filter** — fast, deterministic.
3. **Embedding generation** — LLM API calls per file, expensive
   ($$$), can fail mid-stream, partial results have value.
4. **RAG indexing** — CPU/memory intensive, idempotent.
5. **Documentation generation** — LLM calls for wikis, diagrams,
   Q&A index, expensive, can fail.

The properties of this workflow are exactly what makes it a
*workflow*, not a job:

- **Crash mid-stream is expected, recovery must not redo expensive
  steps.** A crash during stage 5 should not redo stage 3 (paid
  embedding tokens). Each step must have a durable checkpoint.
- **Per-step retry policy is different.** Clone: 3× linear backoff,
  fail fast. LLM embedding: 5× exponential backoff (rate-limit
  aware). RAG indexing: 2× retry. These are not one global retry
  policy.
- **Fan-out / fan-in over LLM calls.** Embedding 500 files in
  parallel, gathering results before stage 4. Job queues serialize
  per worker; workflow engines model parallel sub-tasks natively.
- **Cancellation must be cheap and propagating.** Stopping an
  in-progress $10 indexing run should kill every spawned LLM call,
  not let them complete and waste tokens.
- **Replay UI matters.** Debugging an agent that did 12 LLM calls
  over 8 minutes requires step-by-step inspection. Job-log lines do
  not give you that.

### What River gives you, and where it falls short

River ([`.docs/river.md`](./river.md)) is excellent at what it is: a
Postgres-native single-step job queue. Enqueue → work → done, with
retries. For stateless work — sending a magic-link email, generating
a CSV export — it is the right tool and we should keep using it.

What it does not give you, that a workflow engine does:

| Need | River | Workflow engine (Hatchet, Temporal) |
|------|-------|--------------------------------------|
| Multi-step durable resume | hand-rolled: chain jobs by enqueueing the next at the end of the previous, hold state in your own DB rows | native: each step has a durable checkpoint, replay from step N with same inputs |
| Per-step retry policy | one retry policy per job kind | per-step retry policies, configurable per step within a workflow |
| Fan-out / fan-in | enqueue N children, poll for completion in a follow-up job, manage aggregation yourself | native primitive: `Parallel(...)`, `WaitForAll(...)`, etc. |
| Human-in-the-loop pause | not supported; you would build a `pending_approval` table + cron sweeper | native: durable sleep, event waits, signal-and-resume |
| Replay / step-level observability | structured logs at job granularity | UI showing every step, input/output, retries, durations |

Trying to build a DeepWiki-style pipeline on River means recreating
all five of these primitives on top of one-shot jobs. Doable, slow,
fragile, and you end up with a worse version of Hatchet inside your
own codebase.

---

## Why Hatchet, specifically

- **Go-native, Go SDK is first-class.** The
  [Hatchet authors' own argument for Go agents](https://hatchet.run/blog/go-agents)
  highlights goroutine concurrency (~2 kB per goroutine), uniform
  `context.Context` cancellation (cancel a $10 run, propagate to
  every sub-call), and built-in `runtime/pprof` profiling. These
  matter for long-running agents — Python and Node alternatives
  carry real overhead at agent scale.
- **Postgres durability.** Same datastore we already operate; we
  already know how to back it up
  ([`.docs/disaster-recovery.md`](./disaster-recovery.md)), monitor
  it, restore it. No Redis, no Cassandra, no separate state store.
- **Self-host story is single-binary + Postgres.** The `hatchet-lite`
  image bundles engine + API + dashboard + migrator into one Docker
  image. Drops cleanly into `infra/kamal/deploy.yml` as an accessory,
  like Loki and Grafana already are. Ports: 7077 (gRPC, internal
  only), 8888 (dashboard).
- **Replay UI is built in.** Dashboard on `:8888` shows every step's
  input/output, retry counts, durations, with re-run-from-UI. This
  is the killer feature for debugging multi-step LLM pipelines that
  River fundamentally cannot provide.
- **No per-workflow spawn limit** (Temporal caps at 51,200 spawned
  tasks per workflow — fine for now, but a known ceiling).
- **Postgres-only mode is supported.** For sub-100 events/sec
  (well above our scale), Hatchet can skip RabbitMQ entirely and use
  Postgres as the message queue — one less moving part.

### Caveats (operational reality, not marketing)

These are real findings from production reports, not vendor pages:

- **Write amplification.** Per Hatchet's founder on HN, every task is
  *at minimum 5 Postgres transactions*. At DeepWiki-pipeline scale
  with 500 fan-out embed-file steps, that's ~2,500 transactions just
  for the embedding stage. Fine on our scale; verify Postgres
  `max_connections` accommodates Hatchet's default
  `DATABASE_MAX_CONNS=50` plus our app's pool.
- **OTel for Go is not yet shipped.** Per Hatchet docs (May 2026):
  "OpenTelemetry support is currently only available for the Python
  SDK." For Go workers, log shipping is **stdout → Promtail → Loki**,
  same path we already use for the app. No auto-instrumented trace
  export. Worker spans are missing; engine logs and the dashboard
  cover the observability gap in practice.
- **`hatchet-lite` is officially "for development and low-volume
  use-cases".** That's our scale exactly, but it means the moment
  we go horizontal we'd switch to the full distributed deployment
  (separate engine / API / queue containers, RabbitMQ back in the
  loop). Not today's problem.
- **GRPC broadcast address is sticky.** Per docs: "modifying the GRPC
  broadcast address or server URL will require re-issuing an API
  token." Pin the address in Kamal config from day one.
- **Idle footprint is not published.** No vendor-published "MB RAM
  idle" figure exists for `hatchet-lite` [unverified]. Validation
  step: run a 24-hour measurement on the smallest VPS before
  committing to a server SKU.

### Why not Temporal

Temporal is excellent and battle-tested at scale ([raised $300M at
$5B valuation Feb 2026](https://hatchet.run/versus/hatchet-vs-temporal)).
But:

- The cluster (history + matching + frontend services, plus Postgres
  or Cassandra) is a real jump from where we are. We are
  single-replica Kamal with one Postgres.
- Temporal Cloud is the easy-button, but it is paid infra at a SaaS
  per-execution price. We have no need for that yet.
- The Go SDK is great, but the platform itself feels over-spec for
  our scale.

Worth revisiting *if* we outgrow Hatchet. Today, Hatchet is the right
fit.

### Why not Inngest

Inngest has best-in-class AI agent SDK
([AgentKit](https://agentkit.inngest.com/overview)), but it is
TypeScript-first. Our backend is Go. The Go SDK is second-class and
the platform's happy path is SaaS. Wrong language profile for our
stack.

### Why not Restate

Modern durable-promise model, Rust core, interesting design. Too new
(small ecosystem, fewer production deployments) to bet on for core
infrastructure right now. Watch for the future.

### Why not DBOS Transact (Go) — yet

DBOS Transact's library-only model is genuinely elegant: import a Go
library, point at Postgres, annotate workflow functions. Zero new
infrastructure containers, single binary, same datastore as the app.
For a Postgres-first stack like ours that is structurally attractive.
On the technical primitives it covers what we need — per-step retry
policies (`WithStepMaxRetries`, `WithBaseInterval`,
`WithBackoffFactor`, `WithMaxInterval`), durable fan-out via
`dbos.Go(...)` and `dbos.Select(...)`, queues for stateless work,
cron schedules, and workflow patching for breaking changes.

What rules it out for May 2026 adoption:

1. **Workflow UI is paywalled.** Per `docs.dbos.dev`: self-hosted
   Conductor is "released under a proprietary license. Commercial or
   production use requires a paid license key." The free tier is
   capped at one executor. The OSS escape hatch is "build a Grafana
   dashboard on the `dbos.*` Postgres tables" — workable, but tax we
   did not plan for. Hatchet ships its dashboard for free.
2. **Maturity gap.** `dbos-inc/dbos-transact-golang` was at v0.11.0
   (Feb 10 2026) at research time — still pre-1.0 with ~3 months
   between tags. 682 stars, 58 forks, 6 importers on `pkg.go.dev`,
   no publicly named Go production users. Go-launch Show HN drew
   5 points and 1 comment (from the DBOS CEO). The Python and TS
   SDKs are older and have real adopters; the Go SDK is on its own.
3. **Workflow versioning is operationally heavy for single-replica
   deploys.** DBOS recommends blue-green (keep old-version processes
   alive until in-flight workflows drain). That contradicts our
   single-replica Kamal posture. The escape hatch — `Patch` /
   `DeprecatePatch` conditionals at breakage points — adds permanent
   `if dbos.Patch(...)` noise to workflow code.

**Revisit date: November 2026.** Flip to DBOS for new workflows if
*all four* are true: (a) DBOS Go is ≥ v1.0, (b) an OSS workflow UI
ships or Conductor's free tier expands beyond one executor, (c) at
least two named Go production case studies exist, (d) Hatchet's
operational overhead (token rotation, dashboard auth resets,
broadcast-address fragility) costs more than ~1 hour/month sustained.

---

## Migration plan

The migration runs in four phases. Each phase is independently
shippable; we can pause between phases if a phase reveals something
unexpected.

### Phase 1 — PoC alongside River (week of the first AI workflow proposal)

**Status (2026-05-20):** **In progress on `feat/hatchet-ollama-ai-workflow-poc`.**
The OpenSpec change `add-hatchet-ollama-ai-workflow-poc` implements
this phase end-to-end: dedicated dev Compose services, a new
`internal/aiworkflows/` bounded context, a 5-step DAG (Clone →
Traverse → SummarizeFile fan-out → Aggregate → Store), and a
`/ai/summarize` page that exercises the SSE-progress integration.
Runtime verification (image-tag boot, broadcast-address sanity check,
mid-run crash test) is the remaining gate before this section flips
to **complete**.

Trigger: an AI feature lands in the roadmap (e.g. issue
*"Add code-indexer workflow"*).

- Add **`hatchet-lite`** as a Kamal accessory in
  `infra/kamal/deploy.yml`. Pin to a specific image SHA, not a
  floating tag. Expose ports 7077 (gRPC, *internal network only* —
  Tailscale or Docker network) and 8888 (dashboard, behind auth).
- Provision a **dedicated `hatchet` Postgres database** on the
  existing instance (separate database, same cluster). Keeps
  `pg_dump` boundaries clean and isolates Hatchet's table
  partitioning + autovacuum behaviour from the app schema. The
  research recommends *against* sharing the app DB until we have
  Hatchet operational experience.
- Use **Postgres-only mode** (skip RabbitMQ). Supported up to
  ~100 events/sec, well above our scale.
- Generate API token via `hatchet-admin token create --tenant-id …`
  and store `HATCHET_CLIENT_TOKEN` in Kamal secrets. Pin the GRPC
  broadcast address — changing it later invalidates the token.
- **Log shipping: stdout → Promtail → Loki.** Hatchet Go SDK does
  not have OTel auto-instrumentation in May 2026 (that's Python
  only). Container logs flow through the same Promtail config we
  already use for the app; the dashboard at `:8888` covers per-step
  observability that OTel traces would otherwise provide.
- Build *one* trivial 3-step workflow as a smoke test
  (`Clone → ListFiles → Done`). Confirm engine starts, worker
  registers, workflow durably resumes if killed mid-step. Bonus
  validation: 24-hour idle measurement to characterize the actual
  RAM/CPU footprint of `hatchet-lite` (no vendor-published figure
  exists [unverified]).
- Confirm Kamal deploy story works end-to-end (build, deploy,
  rollback, accessory boot).

### Phase 2 — Ship the first real AI workflow on Hatchet

Trigger: Phase 1 PoC is green.

- Build the actual feature (e.g. DeepWiki-style indexer) as a
  Hatchet workflow with proper step boundaries (Clone, Embed,
  RagIndex, GenerateDocs).
- Each step has its own retry policy + timeout. Set
  `Retries: 5` with exponential backoff on the embed step so paid
  embedding tokens are never re-spent after a successful
  per-file completion. Cheaper, idempotent steps (clone, index)
  use `Retries: 2-3` linear.
- Embedding step fans out per file. Watch for Postgres write
  amplification: ~5 transactions per task × N parallel embed
  steps. At our scale, fine; verify autovacuum on the `hatchet`
  database stays healthy.
- River keeps running for the existing email + export jobs —
  unchanged. **The two-system arrangement is the target state**,
  not a transition. See Phase 3 for the rationale.
- Production-bake the workflow for one cycle (1–2 weeks of real
  use) before moving on.

### Phase 3 — Confirm the two-system target state is right

Trigger: Phase 2 is stable in production.

The expected outcome here is **not consolidation onto Hatchet**.
Research and Hatchet's own founder on HN both confirm what we'd
suspect from looking at the code: a one-shot job is structurally a
River job, not a workflow. The 30-line River worker stays lighter
than the equivalent Hatchet code plus container + token + ports.
This phase exists to *verify* that hypothesis on our specific stack,
not to consolidate.

- Rewrite **one** River job — recommend `send_magic_link` because
  it is the smallest — as a 1-step Hatchet workflow. Keep the River
  version running on a feature flag; do not delete it.
- Compare side-by-side: lines of code, deploy story, monitoring
  surface, latency P99 (Hatchet's 5-Postgres-tx-per-task floor
  matters here). Be honest in the comparison; do not rationalize.
- **Expected outcome:** Hatchet feels heavier for the simple case.
  Confirmed empirically, keep River for one-shot stateless work,
  remove the experimental Hatchet send-magic-link, document the
  comparison in this file.
- **Only-if-surprised path:** If Hatchet 1-step DX is genuinely
  competitive with River (lines of code, deploy ergonomics, idle
  cost) — *and* the operational overhead of running both systems
  exceeds the savings of keeping River — reopen the consolidation
  question as Phase 4.

### Phase 4 — Consolidation (unlikely, conditional)

Trigger: Phase 3 surprises us by showing 1-step Hatchet is competitive
with River *and* two-system operational cost is real.

Default assumption: this phase never runs. Two queues, one Postgres,
separate concerns is the target state.

If it does run:

- Migrate in order of risk: `verification_email` → `2fa_otp` →
  `login_notification` → `data_export`.
- After each migration, verify SSE-progress integration for the
  export job still works (this is the most complex existing job).
- Once all River jobs are migrated and stable for one production
  cycle, remove River from `composition.Build`, `infra/kamal/`,
  and `.docs/`. Update this document to reflect consolidation
  complete.

---

## Kill criteria

Stop the migration and revisit if any of these fire during the
relevant phase:

- **Phase 1:** `hatchet-lite` cannot deploy cleanly on Kamal *or*
  the stdout-log shipping path to Loki via Promtail breaks *or* idle
  RAM/CPU footprint exceeds what fits on the current VPS SKU.
- **Phase 2:** The first real AI workflow has more orchestration
  pain than expected — e.g. fan-out cancellation does not propagate,
  durable replay loses inputs, the SDK forces structural changes we
  do not want.
- **Phase 3:** Already the expected outcome — Hatchet 1-step is
  heavier than River. No migration triggered. *Surprise* in this
  phase (Hatchet 1-step truly competitive) flips us toward Phase 4.
- **Any phase, Postgres health:**
  - Steady CPU on the cluster exceeds 30% during normal AI-workflow
    load. Hatchet's 5-transactions-per-task floor multiplied by
    fan-out can move this number fast.
  - Connection count exceeds budget. Hatchet defaults to
    `DATABASE_MAX_CONNS=50`; verify Postgres `max_connections`
    accommodates that plus the app pool plus headroom.
  - Autovacuum on the `hatchet` database falls behind. Tune per
    Hatchet's docs if we cross ~500 GB of run data (we won't, for
    a long time).

If the AI workflow drives sustained throughput past ~500 events/sec,
the whole question gets reopened — Hatchet's own benchmarks show DB
CPU climbing fast past 1000–2000 runs/sec, and at that point we'd
also be re-evaluating Temporal.

---

## What this document is not

- **Not a vendor comparison.** The detailed candidate writeups in
  issue #57 are correct and are not duplicated here.
- **Not a commitment to a specific timeline.** The trigger is the
  first concrete AI workflow proposal landing. That might be next
  month, or in six months. The plan is ready when it does.
- **Not a forever decision.** Watch Temporal and Restate; revisit
  the engine choice if Hatchet stalls or we outgrow it.

---

## Out of scope

- LLM-orchestration libraries (LangGraph, CrewAI). They run *inside*
  a workflow step. Pick those when we have something to orchestrate.
- No-code platforms (n8n, Make, Zapier). Wrong layer for us.
- Switching River for a different stateless job queue (Asynq,
  Faktory). Not motivated.

---

## Related

- Issue [#57](https://github.com/atilladeniz/next-go-pg/issues/57) —
  the spike issue this document closes
- [`.docs/river.md`](./river.md) — current job-queue stack
- [`.docs/background-jobs.md`](./background-jobs.md) — current job
  patterns

### Engine source / docs

- [Hatchet GitHub](https://github.com/hatchet-dev/hatchet) — engine
  + Go SDK source
- [`hatchet-lite` image docs](https://docs.hatchet.run/self-hosting/hatchet-lite) —
  single-container deployment, ports 7077 + 8888
- [Hatchet: Why Go for agents](https://hatchet.run/blog/go-agents) —
  Go-specific case for the engine
- [DBOS Transact Go SDK](https://github.com/dbos-inc/dbos-transact-golang) —
  the library-only alternative on the Nov 2026 revisit list

### Operational reality / production reports

- [Cynco — self-hosting Hatchet on AWS EC2](https://medium.com/@hazqeelafyq/self-host-hatchet-full-stack-app-on-aws) —
  one of the few public production reports for `hatchet-lite`,
  Postgres-only mode, sub-100 req/sec workload
- [HN: Absurd Workflows — Durable Execution with Just Postgres](https://news.ycombinator.com/item?id=45797228) —
  community discussion of DBOS Conductor paywall and Postgres
  write-amplification concerns
- [DBOS Conductor licensing](https://docs.dbos.dev/production/hosting-conductor) —
  the workflow UI paywall that rules out DBOS for OSS adoption today

### Use-case context

- [Cognition DeepWiki](https://docs.devin.ai/work-with-devin/deepwiki) —
  the kind of AI workflow this decision is sized for
- [deepwiki-open pipeline](https://deepwiki.com/AsyncFuncAI/deepwiki-open) —
  five-stage repository-processing pipeline as the canonical example
  of a multi-step LLM workflow with expensive intermediate state
