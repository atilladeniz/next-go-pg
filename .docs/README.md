# Technical Documentation (.docs)

LLM-friendly technical documentation for this project's tech stack.

**IMPORTANT:** Always check here FIRST before searching the internet!

## Structure

```
.docs/
├── README.md           # This file
├── nextjs.md           # Next.js 16 App Router
├── tanstack-query.md   # TanStack Query (React Query)
├── better-auth.md      # Better Auth
├── gorm.md             # GORM ORM
├── gorilla-mux.md      # Gorilla Mux Router
├── goca.md             # Goca CLI
├── orval.md            # Orval API Client Generator
├── shadcn.md           # shadcn/ui Components
├── tailwind.md         # Tailwind CSS 4
└── kamal-deploy.md     # Kamal Deployment (Docker)
```

## Usage

1. **Claude/LLM:** Read the relevant file before starting a task
2. **Developer:** Quick reference without internet

## Fetching Documentation

Use the `fetch-docs` command to download LLM-friendly documentation:

```bash
# Fetch docs from any URL
make fetch-docs url=https://tanstack.com/query/latest/docs

# With custom name
make fetch-docs url=https://nextjs.org/docs name=nextjs

# Examples
make fetch-docs url=https://orm.drizzle.team name=drizzle
make fetch-docs url=https://docs.stripe.com name=stripe
```

### How it works

The script tries these methods in order:

1. **llms.txt** - Checks if site has `/llms.txt` or `/llms-full.txt` (best quality)
2. **sitefetch** - Crawls entire documentation site (requires `bun i -g sitefetch`)
3. **Jina Reader** - Single page fallback (free, no auth needed)

### Requirements

```bash
# For full site crawling (recommended)
bun install -g sitefetch
```

### Options

```bash
# Only fetch single page (skip crawling)
./scripts/fetch-docs.sh https://example.com/page --single
```

## Tech Stack Overview

| Technology | Version | Docs |
|------------|---------|------|
| Next.js | 16 | `.docs/nextjs.md` |
| TanStack Query | 5 | `.docs/tanstack-query.md` |
| Better Auth | latest | `.docs/better-auth.md` |
| Go | 1.23 | `.docs/go.md` |
| GORM | 2 | `.docs/gorm.md` |
| Gorilla Mux | latest | `.docs/gorilla-mux.md` |
| Goca CLI | latest | `.docs/goca.md` |
| Orval | latest | `.docs/orval.md` |
| shadcn/ui | latest | `.docs/shadcn.md` |
| Tailwind CSS | 4 | `.docs/tailwind.md` |
| Kamal | 2.9.0 | `.docs/kamal-deploy.md` |

## Adding Docs Manually

If automatic fetching doesn't work:

```bash
# Check if site has llms.txt
curl https://example.com/llms.txt

# Or use Jina Reader directly
curl https://r.jina.ai/https://example.com/docs > .docs/example.md
```
