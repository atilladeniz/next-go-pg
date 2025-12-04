---
name: system-design
description: Create and manage system design diagrams with Mermaid. Use for architecture, flows, data models, and business logic visualization.
allowed-tools: Read, Write, Glob, Bash
---

# System Design Skill

Erstelle Mermaid-Diagramme für System Design und Architecture.

## Verzeichnis

```
.concepts/
├── architecture/     # C4, Deployment, Infrastructure
├── flows/            # Auth, Data, User Journey
└── data-models/      # ER, Class, State Diagrams
```

## Workflow

1. **Anfrage analysieren**: Was soll visualisiert werden?
2. **Diagram-Typ wählen**: Passend zum Use Case
3. **Business Logic dokumentieren**: Text vor Diagramm
4. **Mermaid erstellen**: Syntaktisch korrekt
5. **In .concepts/ speichern**: Richtige Kategorie

## Diagram Types

### Architecture
- **C4Context**: System-Übersicht mit Akteuren
- **C4Container**: Technische Container (Services)
- **Flowchart**: Deployment, Infrastructure

### Flows
- **Sequence**: API Calls, Service Interaktionen
- **Flowchart**: Prozesse, Entscheidungen
- **State**: Objekt-Zustände, State Machines
- **User Journey**: UX Flows

### Data Models
- **ER Diagram**: Datenbank-Relationen
- **Class Diagram**: Domain Models
- **Mind Map**: Hierarchien, Brainstorming

## Template

```markdown
# [Titel]

## Business Context

[Beschreibung der Business Logic]

## [Diagram Name]

\`\`\`mermaid
[Diagram Code]
\`\`\`

## Details

[Tabellen, Erklärungen, etc.]
```

## Mermaid Syntax Quick Reference

### Flowchart
```
flowchart TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Action]
    B -->|No| D[Other]
```

### Sequence
```
sequenceDiagram
    A->>B: Request
    B-->>A: Response
```

### ER Diagram
```
erDiagram
    ENTITY1 ||--o{ ENTITY2 : relationship
```

### State Diagram
```
stateDiagram-v2
    [*] --> State1
    State1 --> State2: event
```

### C4 Context
```
C4Context
    Person(user, "User")
    System(app, "App")
    Rel(user, app, "Uses")
```

## Datei-Naming

- `kebab-case.md`
- Beschreibender Name
- Kategorie-Ordner nutzen

Beispiele:
- `.concepts/architecture/deployment.md`
- `.concepts/flows/auth-flow.md`
- `.concepts/data-models/er-diagram.md`
