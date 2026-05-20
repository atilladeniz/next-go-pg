# Workflow Orchestrator Decision (Spike for #57)

**Date:** 2026-05-20
**Status:** Decision recorded — implementation gated on the first
concrete AI workflow proposal.
**Decision:** **Adopt Hatchet** as the workflow engine for AI agent
features. Keep River for the existing stateless email + export jobs
*during the transition*; migrate them after the first AI workflow
proves the stack on Kamal. **Do not start two systems in parallel as
a permanent state.**

This document is the deliverable of issue #57. It is a *decision-support
artifact* — the implementation will land in follow-up issues, gated on
the first concrete AI workflow proposal (e.g. a DeepWiki-style code
indexer, an agentic researcher, a multi-step LLM pipeline).

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
- **Self-host story is single-binary + Postgres.** Drops cleanly into
  `infra/kamal/deploy.yml` as an accessory, like Loki and Grafana
  already are.
- **Bundled observability** (log ingestor + OTel collector) can pipe
  into our existing Loki / Grafana stack (`infra/loki/`,
  `infra/grafana/`).
- **No per-workflow spawn limit** (Temporal caps at 51,200 spawned
  tasks per workflow — fine for now, but a known ceiling).

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

---

## Migration plan

The migration runs in four phases. Each phase is independently
shippable; we can pause between phases if a phase reveals something
unexpected.

### Phase 1 — PoC alongside River (week of the first AI workflow proposal)

Trigger: an AI feature lands in the roadmap (e.g. issue
*"Add code-indexer workflow"*).

- Deploy Hatchet engine + one worker as a Kamal accessory in
  `infra/kamal/deploy.yml`. Use a dedicated Postgres database to
  avoid sharing the app DB on a brand-new system. Revisit DB-sharing
  in Phase 4.
- Wire Hatchet's OTel collector into the existing Loki / Grafana
  stack. If that fails, the migration pauses for an observability
  decision before we ship the first workflow.
- Build *one* trivial 3-step workflow as a smoke test
  (`Clone → ListFiles → Done`). Confirm engine starts, worker
  registers, workflow durably resumes if killed mid-step.
- Confirm Kamal deploy story works end-to-end (build, deploy,
  rollback, accessory boot).

### Phase 2 — Ship the first real AI workflow on Hatchet

Trigger: Phase 1 PoC is green.

- Build the actual feature (e.g. DeepWiki-style indexer) as a
  Hatchet workflow with proper step boundaries (Clone, Embed,
  RagIndex, GenerateDocs).
- Each step has its own retry policy + timeout. Embedding step
  fans out per file.
- River keeps running for the existing email + export jobs —
  unchanged. **No two-system-as-target state**, this is an
  *explicit transition state*.
- Production-bake the workflow for one cycle (1–2 weeks of real
  use) before moving on.

### Phase 3 — Validate Hatchet ergonomics for trivial 1-step jobs (kill criterion)

Trigger: Phase 2 is stable in production.

- Rewrite **one** River job — recommend `send_magic_link` because
  it is the smallest — as a 1-step Hatchet workflow.
- Compare side-by-side: lines of code, deploy story, monitoring
  surface, latency P99. Be honest in the comparison; do not
  rationalize.
- **Decision gate:**
  - If 1-step DX is fine → proceed to Phase 4.
  - If 1-step Hatchet feels visibly heavier than the current
    ~30-line River worker → **stop**. Scope Hatchet to AI
    workflows only, keep River for stateless jobs, accept
    the two-system debt with eyes open.

### Phase 4 — Migrate remaining River jobs (only if Phase 3 passes)

Trigger: Phase 3 confirms Hatchet ergonomics are acceptable for the
simple case.

- Migrate in order of risk: `verification_email` → `2fa_otp` →
  `login_notification` → `data_export`.
- After each migration, verify SSE-progress integration for the
  export job still works (this is the most complex existing job).
- Once all River jobs are migrated and stable for one production
  cycle, remove River from `composition.Build`, `infra/kamal/`,
  and `.docs/`. Update `.docs/orchestrator-decision.md` to reflect
  consolidation complete.

---

## Kill criteria

Stop the migration and revisit if any of these fire during the
relevant phase:

- **Phase 1:** Hatchet engine cannot deploy cleanly on Kamal *or*
  its OTel collector cannot be wired into our Loki stack.
- **Phase 2:** The first real AI workflow has more orchestration
  pain than expected — e.g. fan-out cancellation does not propagate,
  durable replay loses inputs, the SDK forces structural changes we
  do not want.
- **Phase 3:** 1-step workflow ergonomics are visibly worse than
  River — see kill clause in the phase description.
- **Any phase:** Postgres load on the shared cluster (if we share) or
  on the dedicated Hatchet DB exceeds what our single-replica Kamal
  can handle.

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
- [Hatchet GitHub](https://github.com/hatchet-dev/hatchet) — engine
  + Go SDK source
- [Hatchet: Why Go for agents](https://hatchet.run/blog/go-agents) —
  Go-specific case for the engine
- [Cognition DeepWiki](https://docs.devin.ai/work-with-devin/deepwiki) —
  the kind of AI workflow this decision is sized for
