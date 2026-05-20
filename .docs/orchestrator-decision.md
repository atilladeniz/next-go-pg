# Workflow Orchestrator Decision (Spike for #57)

**Date:** 2026-05-20
**Status:** Decision recorded — implementation deferred.
**Decision:** Stay on River for everything we run today. Adopt Hatchet
**only when** we have a concrete AI workflow that needs durable
multi-step execution. Do **not** preemptively swap River out.

This document is the deliverable of issue #57 — a decision-support
artifact, not a spec. It will move into git history once #57 is
closed; revisit if we genuinely decide to build AI agent workflows.

---

## Why this question exists

We use [River](./river.md) for background jobs today:

```
backend/internal/notifications/infrastructure/jobs/   (4 email kinds)
backend/internal/exports/infrastructure/jobs/          (data export with SSE progress)
```

River is a Postgres-native job queue. It is great at **enqueue → work →
done**. It is not a workflow orchestrator: multi-step chains, durable
resume after crash mid-step, per-step retry policies, human-in-the-loop
pauses, fan-out/fan-in over LLM calls — all hand-rolled if you try them
on River.

If we add AI agent features (multi-step LLM chains, tool-use loops,
human-in-the-loop, durable resume), we will eventually want a real
workflow engine. The question is: do we **swap now** (replace River
with Hatchet preemptively) or **defer** (keep River; pick the engine
when the first real AI workflow lands)?

---

## Candidate landscape (snapshot, May 2026)

| Tool | Language | Persistence | Self-host story | AI workflows | Footprint | Verdict |
|------|----------|-------------|-----------------|--------------|-----------|---------|
| **[Hatchet](https://hatchet.run)** | Go-native, Go SDK | Postgres | Single binary; runs on Kamal/Docker; ~7.2k stars; active | First-class | One extra service (engine + worker) | **Primary candidate when we need it** |
| **[Temporal](https://temporal.io)** | Polyglot (Go SDK is first-class) | Cassandra or Postgres | Heavy cluster (history/matching/frontend) or Temporal Cloud ($$); battle-tested | First-class | Significant infra | Strong if we outgrow Hatchet; today over-spec |
| **[Inngest](https://inngest.com)** | TS-first, Go SDK exists | Hosted (self-host possible) | SaaS happy path | First-class (AgentKit) | Low if SaaS, higher if self-host | Wrong language profile — our backend is Go |
| **[Restate](https://restate.dev)** | Rust core | Bundled storage | Single binary | First-class | Light | Too young; small ecosystem |

The detailed candidate writeups in the original #57 issue are accurate
and not repeated here — focus is on what to do, not what exists.

---

## What changed our recommendation

The issue text leans toward **consolidate onto Hatchet, retire River**.
After looking at the actual state of the repo, we are recommending the
**opposite** for now. Three reasons:

### 1. We have no AI workflows yet — only stateless jobs

Every job we run today is a one-shot enqueue → process → done. None of
the Hatchet-specific primitives (durable resume mid-step, fan-out/fan-in,
human-in-the-loop, per-step retry) buy us anything for the current
workload. Adopting Hatchet today is paying a price for capability we
won't use.

### 2. River is working

Four job kinds in two bounded contexts, SSE-integrated export progress,
synchronous email fallback on River unavailability, golden tests, full
DDD integration (#58 refactor). It boots, runs, and handles failures
correctly. Replacing a working system is only worth the cost when the
replacement has a concrete win.

### 3. "Either River alone or Hatchet alone" is the right framing — but
the cut-over moment is **when AI workflows arrive**, not now

The issue correctly identifies that running River + Hatchet in parallel
is "the worst of both worlds". The framing is right; the timing is
wrong. We are not yet at the fork — we are still at the part of the
roadmap where River is sufficient.

The disciplined move:

- **Today:** keep River. Do nothing.
- **When the first real AI workflow lands** (multi-step, durable
  resume, per-step retry policies actually needed): build *that one
  feature* on Hatchet, alongside River, as an *explicit, time-boxed
  transition state*.
- **Then:** validate Hatchet's ergonomics for a trivial 1-step job
  (the critical question — rewrite `send_magic_link` as a 1-step
  Hatchet workflow, confirm DX is not painful).
  - **If ergonomics are fine:** migrate the rest of River's jobs to
    Hatchet, remove River.
  - **If 1-step Hatchet feels heavyweight:** accept the two-system
    debt knowingly, scope Hatchet to AI workflows only, keep River
    for stateless jobs.

This is the same migration path the issue describes, just gated on a
concrete AI use-case landing instead of preemptive.

---

## What we are committing to

**Now:**

- No change to `backend/internal/notifications/infrastructure/jobs/`
  or `backend/internal/exports/infrastructure/jobs/`.
- No new infrastructure on the Kamal box.
- Issue #57 closed with this document as the deliverable.

**When the first AI workflow proposal lands:**

- Open a new spike issue: *"PoC: First AI workflow on Hatchet,
  alongside River"*.
- Deploy Hatchet engine + one worker on the existing Kamal host.
- Build that one workflow end-to-end. Confirm:
  - Hatchet engine self-hosts cleanly on Kamal (1 accessory in
    `infra/kamal/deploy.yml`).
  - Postgres load is acceptable on the shared DB, or Hatchet gets
    its own DB.
  - OTel collector output plugs into our existing Loki/Grafana stack
    (`infra/loki/`, `infra/grafana/`).
  - Go SDK ergonomics for a *trivial* 1-step job are not painful.
- After that workflow ships and is stable for at least one production
  cycle, decide on full River-to-Hatchet migration.

**Kill criteria for the eventual migration:**

- Hatchet 1-step workflow ergonomics are visibly worse than River's
  current ~30-line worker pattern → keep River, scope Hatchet to
  multi-step only.
- Postgres load doubles unacceptably when Hatchet shares the DB →
  either give Hatchet its own DB, or revisit.
- Hatchet's bundled OTel collector cannot be redirected to our Loki
  stack → keep River for production observability continuity.

---

## What this document is not

- **Not a vendor comparison.** The original issue covers that
  adequately; rehashing those tables in code review adds nothing.
- **Not a rejection of Hatchet.** The recommendation is "adopt at the
  right moment", not "don't adopt".
- **Not a forever decision.** The decision is "the cut-over trigger is
  a concrete AI workflow"; the trigger may fire next month or never.
  Revisit when context changes.

---

## Out of scope

- LLM-orchestration libraries (LangGraph, CrewAI). They run *inside* a
  workflow step. We pick those when we have something to orchestrate.
- No-code platforms (n8n, Make, Zapier). Wrong layer for us.
- Switching River for a different stateless job queue (Asynq, Faktory,
  etc.). Not motivated — River works.

---

## Related

- Issue [#57](https://github.com/atilladeniz/next-go-pg/issues/57) —
  the spike issue this document closes
- [`.docs/river.md`](./river.md) — current job-queue stack
- [`.docs/background-jobs.md`](./background-jobs.md) — current job
  patterns
- [`.docs/auth-architecture.md`](./auth-architecture.md) — sibling
  spike on auth architecture (issue #47), not yet written
