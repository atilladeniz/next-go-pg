---
name: tech-docs
description: Access local technical documentation before searching the internet. Use FIRST when researching any tech stack question.
allowed-tools: Read, Glob
---

# Technical Documentation Skill

**IMMER ZUERST `.docs/` durchsuchen** bevor im Internet recherchiert wird!

## Dokumentations-Verzeichnis

```
.docs/
├── README.md           # Uebersicht
├── nextjs.md           # Next.js 16 App Router
├── tanstack-query.md   # TanStack Query / React Query
├── better-auth.md      # Better Auth
├── gorm.md             # GORM ORM
├── gorilla-mux.md      # Gorilla Mux Router
├── goca.md             # Goca CLI
├── orval.md            # Orval API Client Generator
├── shadcn.md           # shadcn/ui Components
└── tailwind.md         # Tailwind CSS 4
```

## Workflow

1. **Frage zu einer Technologie?**
   → Zuerst `.docs/<tech>.md` lesen

2. **Docs nicht vorhanden?**
   → Dann erst im Internet recherchieren
   → Ergebnisse in `.docs/` speichern fuer spaeter

3. **Docs veraltet?**
   → Aktualisieren mit neuesten Infos

## Beispiel

```
User: "Wie funktioniert HydrationBoundary in TanStack Query?"

1. Read .docs/tanstack-query.md
2. Falls nicht vorhanden/unvollstaendig:
   - WebFetch von TanStack Query Docs
   - Relevante Info in .docs/tanstack-query.md speichern
3. Frage beantworten
```

## Docs hinzufuegen

Viele Libraries bieten `llms.txt` an:

```bash
# Beispiel: TanStack Query
curl https://tanstack.com/query/latest/llms.txt > .docs/tanstack-query.md

# Oder manuell die wichtigsten Patterns dokumentieren
```

## Prioritaet

1. `.docs/` (lokal, schnell, projekt-spezifisch)
2. `.claude/skills/` (projekt-spezifische Patterns)
3. Internet (nur wenn lokal nicht vorhanden)
