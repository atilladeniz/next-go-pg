<!-- BEGIN:nextjs-agent-rules -->

# Next.js: ALWAYS read docs before coding

Before any Next.js work, find and read the relevant doc in `node_modules/next/dist/docs/`. Your training data is outdated — the docs are the source of truth.

<!-- END:nextjs-agent-rules -->

<!-- BEGIN:openspec-rules -->

# OpenSpec: spec-driven workflow

This repo uses [OpenSpec](https://github.com/Fission-AI/OpenSpec) for spec-driven feature work. Before implementing a non-trivial change, write a proposal under `openspec/changes/<name>/` and follow the artifact order: `proposal.md` → `design.md` → `tasks.md` → implementation.

Slash commands (Claude Code):

| Command | Purpose |
|---------|---------|
| `/opsx:propose <description>` | Create a change and generate all artifacts in one step |
| `/opsx:new <name>` | Create an empty change directory |
| `/opsx:continue` | Generate the next missing artifact for the current change |
| `/opsx:apply` | Implement the tasks listed in `tasks.md` |
| `/opsx:verify` | Check that implementation matches artifacts before archiving |
| `/opsx:archive` | Move a completed change into `openspec/changes/archive/` |
| `/opsx:explore` | Think-mode: investigate ideas before committing to a change |

**Before implementing any feature**: read `openspec/changes/<feature>/proposal.md` (and `design.md` if present). The spec is the contract — execute against it, not assumptions.

Prerequisite (one-time per machine):

```bash
npm install -g @fission-ai/openspec@latest
```

See `.docs/openspec.md` for the verb cheatsheet.

<!-- END:openspec-rules -->
