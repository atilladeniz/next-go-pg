# OpenSpec Cheatsheet

[OpenSpec](https://github.com/Fission-AI/OpenSpec) is the spec-driven workflow for this repo. Use it before any non-trivial feature: a written proposal beats an ad-hoc implementation.

## Prerequisite

One-time install per machine:

```bash
npm install -g @fission-ai/openspec@latest
```

Confirm: `openspec --version` (≥ 1.3.x).

## Layout

```
openspec/
├── changes/
│   ├── <change-name>/       # Active change in progress
│   │   ├── README.md        # Short title + description
│   │   ├── proposal.md      # What & why
│   │   ├── design.md        # How (architecture, trade-offs)
│   │   ├── tasks.md         # Implementation checklist
│   │   └── specs/           # Delta specs (added/modified capabilities)
│   └── archive/             # Completed changes
└── specs/                   # Current canonical specs (the source of truth)
```

## Slash Commands (Claude Code)

| Command | When to use |
|---------|-------------|
| `/opsx:propose "<idea>"` | Quick start — creates change and all artifacts in one step |
| `/opsx:new <kebab-name>` | Just the directory; fill artifacts manually |
| `/opsx:explore` | Think-mode before committing to a change |
| `/opsx:continue` | Generate the next missing artifact |
| `/opsx:apply` | Execute the tasks in `tasks.md` |
| `/opsx:verify` | Sanity check before archiving |
| `/opsx:archive` | Move completed change into `archive/` and sync `specs/` |
| `/opsx:sync` | Apply delta specs to main specs without archiving |

## Manual CLI (when slash commands aren't available)

```bash
openspec new change <name>            # Create change directory
openspec list                          # List active changes (--specs for specs)
openspec list --specs                  # List canonical specs
openspec status --change <name>        # Show artifact completion state
openspec validate                      # Validate all changes/specs
openspec validate <change>             # Validate a specific change
openspec show <change-or-spec>         # Print a change/spec
openspec archive <change>              # Archive a completed change
openspec view                          # Interactive dashboard (TUI)
```

## Workflow

1. **Propose** — `/opsx:propose "add billing portal"`. Generates `proposal.md`, `design.md`, `tasks.md`.
2. **Review** — open `openspec/changes/<name>/` and edit artifacts. The spec is the contract.
3. **Apply** — `/opsx:apply` walks through `tasks.md` and implements each step.
4. **Verify** — `/opsx:verify` reads the artifacts and checks the implementation matches.
5. **Archive** — `/opsx:archive` moves the change into `archive/` and updates `openspec/specs/`.

## Rules for this repo

- **Before touching code** for a feature listed in `openspec/changes/`, read its `proposal.md` (and `design.md` if present). Execute against the spec, not assumptions.
- Use OpenSpec for: new features (3+ steps), architectural decisions, anything spanning multiple files.
- Skip OpenSpec for: small fixes, refactors, dependency bumps. Those go straight into a PR.
- Change names use kebab-case: `add-billing-portal`, `streamline-cmd-binaries`.

## Related

- `AGENTS.md` — repo-root agent rules (loaded by Claude Code automatically)
- `.claude/CLAUDE.md` — project context (auto-imports AGENTS.md)
- [OpenSpec GitHub](https://github.com/Fission-AI/OpenSpec) — upstream docs
