## ADDED Requirements

### Requirement: Self-hosted Hatchet workflow engine runs in the dev stack

The system SHALL run a `hatchet-lite` workflow engine container as part of the development Docker Compose stack, providing durable workflow execution for AI agent features.

#### Scenario: Hatchet engine starts on `docker compose up`
- **WHEN** a developer runs `just dev` (which brings up the dev Docker Compose stack)
- **THEN** the `hatchet-lite` container starts within 60 s, with gRPC listening on internal port 7077 and the dashboard reachable at `http://localhost:8888`

#### Scenario: Hatchet uses a dedicated Postgres database
- **WHEN** the `hatchet-lite` container starts
- **THEN** it connects to a database named `hatchet` on the existing dev Postgres instance, runs its own migrations, and does NOT write to the application's `nextgopg` database

#### Scenario: Hatchet uses Postgres-only mode
- **WHEN** the `hatchet-lite` container starts
- **THEN** it operates without a RabbitMQ dependency; all queue state is in Postgres

### Requirement: Local Ollama LLM runtime is available to backend workers

The system SHALL run an Ollama container providing local LLM inference for AI workflow steps, eliminating cloud-LLM dependency during development.

#### Scenario: Ollama is reachable from backend workers
- **WHEN** the backend container makes an HTTP request to `http://ollama:11434/api/generate`
- **THEN** Ollama responds with model output using the configured model

#### Scenario: Default model is pulled on first start
- **WHEN** the Ollama container first starts and the configured `OLLAMA_MODEL` is not yet present in the persisted model volume
- **THEN** the model is pulled automatically before the service reports healthy, OR a clear log message instructs the operator to run `ollama pull <model>`

#### Scenario: Model selection is configurable
- **WHEN** a developer sets the `OLLAMA_MODEL` environment variable (e.g. `gemma4:e4b` for dev laptops, `gemma4:26b-a4b` for heavier workstations)
- **THEN** the Ollama service uses that model as the default for new requests

### Requirement: Workflow engine integration is encapsulated behind an application port

The backend SHALL define an application-layer port (`HatchetEnqueuer` or equivalent) that hides the workflow engine SDK from the application and HTTP layers.

#### Scenario: HTTP handler enqueues a workflow without importing the SDK
- **WHEN** the HTTP handler for an AI feature processes a request
- **THEN** it calls the application-layer port to enqueue a workflow run, and the HTTP handler package does NOT import `github.com/hatchet-dev/hatchet/sdks/go`

#### Scenario: Engine swap requires only infrastructure changes
- **WHEN** the workflow engine is ever swapped (e.g. to Temporal or DBOS)
- **THEN** only the infrastructure-layer adapter needs to change; the application use cases and HTTP handlers are unaffected

### Requirement: Workflow progress events flow through the existing platform SSE broker

The system SHALL publish workflow progress events to the platform SSE broker (`internal/platform/sse`) using a new event type so the frontend can observe runs without coupling to the workflow engine's own eventing.

#### Scenario: Progress event is published on step completion
- **WHEN** a workflow step completes (successfully or with failure)
- **THEN** an SSE event of type `ai-progress` is published to the platform broker with the run ID, step name, status, and step-specific data (file count, file index, etc.)

#### Scenario: Frontend subscribes to progress without engine API knowledge
- **WHEN** the frontend listens for `ai-progress` events via the existing SSE infrastructure
- **THEN** it receives all step-level updates for the current user's runs without making any direct calls to the Hatchet API

### Requirement: Workflow worker logs flow to the existing observability stack

The Hatchet engine and backend worker containers SHALL emit structured logs to stdout in a format consumable by Promtail, so they appear in the existing Loki/Grafana stack alongside application logs.

#### Scenario: Hatchet engine logs appear in Loki
- **WHEN** a workflow run completes and produces engine-level logs
- **THEN** those logs are picked up by Promtail from the container's stdout and queryable in Grafana with the `service` label set to `hatchet`

#### Scenario: Backend worker logs include the run ID
- **WHEN** the backend worker executes a workflow step
- **THEN** all log lines emitted from that step include the workflow run ID in the structured log context (`logger.WithContext(ctx).…`)

### Requirement: Workflow completion is tracked by a Prometheus counter

The system SHALL increment a Prometheus counter `ai_workflows_completed_total{status}` once per workflow run terminus, with `status` set to `success`, `failed`, or `cancelled`.

#### Scenario: Successful run increments success counter
- **WHEN** a workflow run reaches the final step and stores the result without error
- **THEN** `ai_workflows_completed_total{status="success"}` is incremented by 1

#### Scenario: Failed run increments failed counter
- **WHEN** a workflow run terminates with an error after exhausting retries on any step
- **THEN** `ai_workflows_completed_total{status="failed"}` is incremented by 1
