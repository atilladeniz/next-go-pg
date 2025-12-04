# Technical Documentation (.docs)

LLM-friendly technical documentation for this project's tech stack.

**WICHTIG:** Immer ZUERST hier nachschauen bevor im Internet recherchiert wird!

## Struktur

```
.docs/
├── README.md           # Diese Datei
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

## Verwendung

1. **Claude/LLM:** Lese die relevante Datei bevor du eine Aufgabe startest
2. **Entwickler:** Schnelle Referenz ohne Internet

## Docs hinzufuegen

Wenn du eine neue Technologie verwendest oder Docs aktualisieren willst:

```bash
# Docs von llms.txt holen (falls verfuegbar)
curl https://example.com/llms.txt > .docs/example.md

# Oder manuell erstellen mit den wichtigsten Patterns
```

## Tech Stack Uebersicht

| Technologie | Version | Docs |
|-------------|---------|------|
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
