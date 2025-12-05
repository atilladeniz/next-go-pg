# FSD Migration Plan - Next.js Frontend

## Ziel

Migration der aktuellen Frontend-Struktur zu Feature-Sliced Design (FSD) mit Steiger-Linting.

## Entscheidungen

- **Entities**: Ja, mit User Entity (zukunftssicher)
- **Aliases**: Layer-Aliases (`@shared/`, `@features/`, `@widgets/`, `@entities/`)
- **Linting**: Recommended (alle Steiger-Regeln)

---

## Aktuelle Struktur (IST)

```text
frontend/src/
├── api/                    # Orval-generiert
│   ├── endpoints/          # React Query Hooks
│   ├── models/             # TypeScript Types
│   └── custom-fetch.ts     # Fetch Wrapper
├── app/                    # Next.js App Router
│   ├── (auth)/             # Login, Register
│   ├── (protected)/        # Dashboard, API-Test
│   └── api/auth/           # Better Auth Handler
├── components/             # ALLES gemischt!
│   ├── ui/                 # shadcn/ui (60+ Komponenten)
│   ├── header.tsx          # App-spezifisch
│   ├── mode-toggle.tsx     # Feature
│   ├── providers.tsx       # App-Config
│   └── theme-provider.tsx  # App-Config
├── hooks/                  # Global hooks
│   ├── use-sse.ts          # SSE für Stats
│   ├── use-auth-sync.ts    # Cross-Tab Auth
│   └── use-mobile.ts       # Responsive
└── lib/                    # Utilities
    ├── auth.ts             # Better Auth Server
    ├── auth-client.ts      # Better Auth Client
    ├── auth-server.ts      # Session Helper
    ├── get-query-client.ts # TanStack Query
    └── utils.ts            # cn() helper
```

---

## Ziel-Struktur (SOLL) - FSD

```text
frontend/src/
├── app/                        # Next.js App Router (BLEIBT)
│   ├── (auth)/
│   ├── (protected)/
│   ├── api/auth/
│   ├── layout.tsx
│   └── globals.css
│
├── widgets/                    # Composite UI Blocks
│   └── header/
│       ├── ui/
│       │   └── header.tsx
│       └── index.ts
│
├── features/                   # User Interactions
│   ├── auth/
│   │   ├── ui/
│   │   │   ├── login-form.tsx
│   │   │   └── register-form.tsx
│   │   ├── model/
│   │   │   └── use-auth-sync.ts
│   │   └── index.ts
│   ├── theme-toggle/
│   │   ├── ui/
│   │   │   └── mode-toggle.tsx
│   │   └── index.ts
│   └── stats/
│       ├── ui/
│       │   └── stats-grid.tsx
│       ├── model/
│       │   └── use-sse.ts
│       └── index.ts
│
├── entities/                   # Business Objects
│   └── user/
│       ├── ui/
│       │   └── user-info.tsx   # User display component
│       ├── model/
│       │   └── types.ts        # User type definition
│       └── index.ts
│
└── shared/                     # Wiederverwendbar
    ├── ui/                     # shadcn/ui Komponenten
    │   ├── button.tsx
    │   ├── card.tsx
    │   └── index.ts            # Public API
    ├── api/                    # Orval-generiert
    │   ├── endpoints/
    │   ├── models/
    │   ├── custom-fetch.ts
    │   └── index.ts
    ├── lib/                    # Utilities
    │   ├── auth/
    │   │   ├── auth.ts
    │   │   ├── auth-client.ts
    │   │   ├── auth-server.ts
    │   │   └── index.ts
    │   ├── query-client.ts
    │   ├── utils.ts
    │   └── index.ts
    └── config/                 # App-weite Config
        ├── providers.tsx
        └── index.ts
```

---

## Migrations-Schritte

### Phase 1: shared/ Layer aufbauen (Basis)

1. `shared/ui/` - shadcn/ui Komponenten verschieben
   - `components/ui/*` → `shared/ui/*`
   - `shared/ui/index.ts` erstellen (Public API)

2. `shared/lib/` - Utilities gruppieren
   - `lib/utils.ts` → `shared/lib/utils.ts`
   - `lib/get-query-client.ts` → `shared/lib/query-client.ts`
   - Auth-Files gruppieren in `shared/lib/auth/`

3. `shared/api/` - Orval-Output verschieben
   - `api/*` → `shared/api/*`
   - Orval Config anpassen

4. `shared/config/` - App-Konfiguration
   - `components/providers.tsx` → `shared/config/providers.tsx`
   - `components/theme-provider.tsx` → `shared/config/theme-provider.tsx`

### Phase 2: entities/ Layer aufbauen

5. `entities/user/` - User Business Object
   - User Type Definition erstellen
   - User-Info Component für Display

### Phase 3: features/ Layer aufbauen

6. `features/auth/` - Login/Register Logik
   - Login/Register Forms aus Pages extrahieren
   - `use-auth-sync.ts` → `features/auth/model/`

7. `features/theme-toggle/` - Dark Mode
   - `mode-toggle.tsx` → `features/theme-toggle/ui/`

8. `features/stats/` - Dashboard Stats
   - `stats-grid.tsx` → `features/stats/ui/`
   - `use-sse.ts` → `features/stats/model/`

### Phase 4: widgets/ Layer aufbauen

9. `widgets/header/` - App Header
    - `header.tsx` → `widgets/header/ui/`

### Phase 5: TypeScript Paths & Linting

10. `tsconfig.json` - Path Aliases erweitern

```json
{
  "paths": {
    "@/*": ["./src/*"],
    "@shared/*": ["./src/shared/*"],
    "@entities/*": ["./src/entities/*"],
    "@features/*": ["./src/features/*"],
    "@widgets/*": ["./src/widgets/*"]
  }
}
```

11. Steiger installieren - FSD Linting

```bash
bun add -D steiger @feature-sliced/steiger-plugin
```

12. `steiger.config.js` erstellen (recommended config)

13. Biome Config anpassen - Neue Pfade zu ignores

### Phase 6: App-Layer Imports aktualisieren

14. Alle Imports in `app/` aktualisieren
    - Von `@/components/ui/*` → `@shared/ui`
    - Von `@/lib/*` → `@shared/lib`
    - Von `@/hooks/*` → `@features/*/model`

---

## Import-Regeln (Steiger)

```text
app/        → widgets, features, entities, shared
widgets/    → features, entities, shared
features/   → entities, shared
entities/   → shared
shared/     → (nur externe libs)
```

Verboten:

- features → features (Cross-Imports)
- shared → features (Upward-Imports)
- entities → features (Upward-Imports)

---

## Orval Config Anpassung

```typescript
// orval.config.ts
export default defineConfig({
  api: {
    output: {
      target: "./src/shared/api/endpoints",
      schemas: "./src/shared/api/models",
      override: {
        mutator: {
          path: "./src/shared/api/custom-fetch.ts",
          name: "customFetch",
        },
      },
    },
  },
})
```

---

## Risiken & Mitigationen

| Risiko | Mitigation |
|--------|------------|
| Viele Import-Änderungen | TypeScript zeigt alle Fehler sofort |
| Orval-Pfade ändern sich | Einmal `make api` nach Migration |
| Breaking Changes | Feature-Branch, inkrementelle Migration |

---

## Nächste Schritte nach Approval

1. Phase 1-6 sequentiell durchführen
2. Nach jeder Phase: `bun run typecheck` zum Verifizieren
3. Am Ende: `make api` für Orval-Regeneration
4. Steiger-Lint ausführen und Fehler fixen
