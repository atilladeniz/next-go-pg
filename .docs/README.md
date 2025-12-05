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

## Adding Docs

When using a new technology or updating docs:

```bash
# Fetch docs from llms.txt (if available)
curl https://example.com/llms.txt > .docs/example.md

# Or create manually with the most important patterns
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
| **Kamal** | 2.9.0 | `.docs/kamal-deploy.md` |
