---
name: system-design
description: Create and manage system design diagrams with Mermaid. Use for architecture, flows, data models, and business logic visualization.
allowed-tools: Read, Write, Glob, Bash
---

# System Design Skill

Create Mermaid diagrams for System Design and Architecture.

## Directory

```
.concepts/
├── architecture/     # C4, Deployment, Infrastructure
├── flows/            # Auth, Data, User Journey
└── data-models/      # ER, Class, State Diagrams
```

## Workflow

1. **Analyze request**: What needs to be visualized?
2. **Choose diagram type**: Appropriate for the use case
3. **Document business logic**: Text before diagram
4. **Create Mermaid**: Syntactically correct
5. **Save in .concepts/**: Correct category

## Diagram Types

### Architecture
- **C4Context**: System overview with actors
- **C4Container**: Technical containers (services)
- **Flowchart**: Deployment, Infrastructure

### Flows
- **Sequence**: API calls, service interactions
- **Flowchart**: Processes, decisions
- **State**: Object states, state machines
- **User Journey**: UX flows

### Data Models
- **ER Diagram**: Database relations
- **Class Diagram**: Domain models
- **Mind Map**: Hierarchies, brainstorming

## Template

```markdown
# [Title]

## Business Context

[Description of business logic]

## [Diagram Name]

\`\`\`mermaid
[Diagram Code]
\`\`\`

## Details

[Tables, explanations, etc.]
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

## File Naming

- `kebab-case.md`
- Descriptive name
- Use category folders

Examples:
- `.concepts/architecture/deployment.md`
- `.concepts/flows/auth-flow.md`
- `.concepts/data-models/er-diagram.md`
