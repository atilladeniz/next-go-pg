---
description: "Type: architecture, flow, data-model"
arguments:
  - name: topic
    description: What do you want to visualize? (e.g., "how users login", "database structure", "deployment setup")
    required: true
---

# Create System Design Concept

Create a new concept document with Mermaid diagrams in `.concepts/`.

## Workflow

1. **User describes** what they want to visualize: `$ARGUMENTS.topic`
2. **You analyze** and automatically detect:
   - Which **type** fits (architecture, flow, data-model)
   - Which **Mermaid diagram** is best suited
   - An appropriate **filename** (kebab-case)
3. **You create** the document with Business Context + Mermaid

## Type Detection Guide

### → `architecture/`
**When it's about:**
- How systems/services connect
- Deployment, server, infrastructure
- Technical components and their connections
- "Big picture" / overview

**Mermaid:** C4Context, C4Container, flowchart with subgraphs

### → `flows/`
**When it's about:**
- How something works (process, workflow)
- Who communicates with whom (API calls, events)
- User journey / UX flow
- Chronological sequence of actions

**Mermaid:** sequenceDiagram, flowchart, journey, stateDiagram-v2

### → `data-models/`
**When it's about:**
- Database structure, tables, relations
- Entity/domain models
- States of an object (state machine)
- Classes and their relationships

**Mermaid:** erDiagram, classDiagram, stateDiagram-v2

## Response Format

After detecting the type:

1. **Briefly confirm** what you understood
2. **Create file** in `.concepts/<type>/<name>.md`
3. **Content:**
   - `## Business Context` - What is the problem/context
   - `## [Diagram]` - Mermaid diagram(s)
   - `## Details` - Table with explanations (optional)

## Examples

| User says | You detect | Output |
|-----------|------------|--------|
| "how does the login work" | flow + sequenceDiagram | `.concepts/flows/login-flow.md` |
| "show me the database structure" | data-model + erDiagram | `.concepts/data-models/database-schema.md` |
| "how do frontend and backend connect" | architecture + C4Container | `.concepts/architecture/system-overview.md` |
| "what happens when user orders" | flow + sequenceDiagram | `.concepts/flows/order-process.md` |
| "what states does an order have" | data-model + stateDiagram | `.concepts/data-models/order-states.md` |
| "how do we deploy" | architecture + flowchart | `.concepts/architecture/deployment.md` |

---

**User wants to visualize:** `$ARGUMENTS.topic`

Analyze what the user needs, choose the appropriate type and diagram style, and create the concept document.
