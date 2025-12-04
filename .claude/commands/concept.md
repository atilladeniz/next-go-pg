---
description: Create a new system design concept with Mermaid diagrams
arguments:
  - name: name
    description: Name of the concept (e.g., "payment-flow", "notification-system")
    required: true
  - name: type
    description: "Type: architecture, flow, data-model"
    required: false
---

# Create System Design Concept

Erstelle ein neues Concept-Dokument mit Mermaid-Diagrammen in `.concepts/`.

## Instructions

1. **Analysiere** was der User visualisieren möchte
2. **Wähle** den passenden Diagram-Typ:
   - **architecture**: C4, Deployment, Infrastructure → `.concepts/architecture/`
   - **flow**: Sequence, Flowchart, State, User Journey → `.concepts/flows/`
   - **data-model**: ER, Class Diagram → `.concepts/data-models/`

3. **Erstelle** das Dokument mit:
   - Business Context (Textbeschreibung)
   - Mermaid Diagram(s)
   - Details/Tabellen wenn nötig

4. **Speichere** in `.concepts/<type>/<name>.md`

## Diagram Type Reference

| Use Case | Diagram Type | Mermaid |
|----------|-------------|---------|
| System Overview | C4 Context | `C4Context` |
| Service Architecture | C4 Container | `C4Container` |
| API Interaction | Sequence | `sequenceDiagram` |
| Process/Workflow | Flowchart | `flowchart TD` |
| Object States | State | `stateDiagram-v2` |
| Database Schema | ER Diagram | `erDiagram` |
| Class Structure | Class | `classDiagram` |
| User Experience | Journey | `journey` |
| Timeline | Gantt | `gantt` |
| Hierarchy | Mind Map | `mindmap` |

## Template

```markdown
# [Concept Name]

## Business Context

[Beschreibe die Business Logic und den Kontext]

## [Diagram Title]

\`\`\`mermaid
[Diagram Code]
\`\`\`

## Details

| Aspekt | Beschreibung |
|--------|--------------|
| ... | ... |

## Related Concepts

- [Link zu verwandten Concepts]
```

## Examples

### `/concept payment-flow flow`
→ Creates `.concepts/flows/payment-flow.md` with sequence/flowchart diagrams

### `/concept user-service architecture`
→ Creates `.concepts/architecture/user-service.md` with C4 diagrams

### `/concept order-model data-model`
→ Creates `.concepts/data-models/order-model.md` with ER diagram

## Argument: $ARGUMENTS.name

Name des Concepts: **{{ name | default: "[Frage nach Name]" }}**

## Argument: $ARGUMENTS.type

Typ: **{{ type | default: "auto-detect" }}**

---

Bitte erstelle jetzt das Concept-Dokument basierend auf dem Namen und Typ.
Frage nach Details falls nötig.
