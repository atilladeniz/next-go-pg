---
description: "Type: architecture, flow, data-model"
arguments:
  - name: topic
    description: What do you want to visualize? (e.g., "how users login", "database structure", "deployment setup")
    required: true
---

# Create System Design Concept

Erstelle ein neues Concept-Dokument mit Mermaid-Diagrammen in `.concepts/`.

## Workflow

1. **User beschreibt** was er visualisieren möchte: `$ARGUMENTS.topic`
2. **Du analysierst** und erkennst automatisch:
   - Welcher **Typ** passt (architecture, flow, data-model)
   - Welches **Mermaid-Diagramm** am besten geeignet ist
   - Einen passenden **Dateinamen** (kebab-case)
3. **Du erstellst** das Dokument mit Business Context + Mermaid

## Type Detection Guide

### → `architecture/`
**Wenn es um geht:**
- Wie Systeme/Services zusammenhängen
- Deployment, Server, Infrastructure
- Technische Komponenten und ihre Verbindungen
- "Big Picture" / Übersicht

**Mermaid:** C4Context, C4Container, flowchart mit Subgraphs

### → `flows/`
**Wenn es um geht:**
- Wie etwas abläuft (Prozess, Workflow)
- Wer mit wem kommuniziert (API Calls, Events)
- User Journey / UX Flow
- Zeitliche Abfolge von Aktionen

**Mermaid:** sequenceDiagram, flowchart, journey, stateDiagram-v2

### → `data-models/`
**Wenn es um geht:**
- Datenbank-Struktur, Tabellen, Relationen
- Entity/Domain Models
- Zustände eines Objekts (State Machine)
- Klassen und ihre Beziehungen

**Mermaid:** erDiagram, classDiagram, stateDiagram-v2

## Response Format

Nachdem du den Typ erkannt hast:

1. **Kurz bestätigen** was du verstanden hast
2. **Datei erstellen** in `.concepts/<type>/<name>.md`
3. **Inhalt:**
   - `## Business Context` - Was ist das Problem/Kontext
   - `## [Diagram]` - Mermaid Diagramm(e)
   - `## Details` - Tabelle mit Erklärungen (optional)

## Examples

| User sagt | Du erkennst | Output |
|-----------|-------------|--------|
| "wie funktioniert der login" | flow + sequenceDiagram | `.concepts/flows/login-flow.md` |
| "zeig mir die datenbank struktur" | data-model + erDiagram | `.concepts/data-models/database-schema.md` |
| "wie hängen frontend und backend zusammen" | architecture + C4Container | `.concepts/architecture/system-overview.md` |
| "was passiert wenn user bestellt" | flow + sequenceDiagram | `.concepts/flows/order-process.md` |
| "welche zustände hat eine bestellung" | data-model + stateDiagram | `.concepts/data-models/order-states.md` |
| "wie deployen wir" | architecture + flowchart | `.concepts/architecture/deployment.md` |

---

**User möchte visualisieren:** `$ARGUMENTS.topic`

Analysiere was der User braucht, wähle den passenden Typ und Diagramm-Art, und erstelle das Concept-Dokument.
