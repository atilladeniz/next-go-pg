## ADDED Requirements

### Requirement: Authenticated user can request a repository summary

The system SHALL allow an authenticated user to submit a public Git repository URL and receive an asynchronous summarization run that produces per-file and repo-level summaries.

#### Scenario: User submits a valid repo URL
- **WHEN** an authenticated user POSTs `{"repoURL": "https://github.com/owner/repo"}` to `/api/v1/ai/summarize-repo`
- **THEN** the system responds with HTTP 202, returning a JSON body containing the run ID and an initial status `pending`

#### Scenario: User submits an invalid repo URL
- **WHEN** an authenticated user POSTs `{"repoURL": "not-a-url"}` to `/api/v1/ai/summarize-repo`
- **THEN** the system responds with HTTP 400 and a descriptive error message; no workflow run is enqueued

#### Scenario: Unauthenticated request is rejected
- **WHEN** an unauthenticated request hits `/api/v1/ai/summarize-repo`
- **THEN** the system responds with HTTP 401; no workflow run is enqueued

### Requirement: Summarization workflow has five discrete durable steps

The summarization workflow SHALL consist of five steps, each independently checkpointed in the workflow engine, with per-step retry policies tuned to the cost and failure mode of each step.

#### Scenario: Clone step retries 3 times on transient failure
- **WHEN** the clone step fails with a transient network error
- **THEN** the workflow engine retries up to 3 times with linear backoff before marking the run as failed

#### Scenario: SummarizeFile retries 5 times with exponential backoff
- **WHEN** the SummarizeFile step fails (e.g. Ollama is temporarily unavailable)
- **THEN** the workflow engine retries up to 5 times with exponential backoff

#### Scenario: SummarizeFile fans out per-file
- **WHEN** the Traverse step yields N relevant files
- **THEN** the workflow engine launches N SummarizeFile sub-tasks in parallel, with each sub-task independently checkpointed

#### Scenario: Mid-run crash resumes from the last completed step
- **WHEN** the backend worker container is killed during execution of the Aggregate step
- **THEN** after the worker restarts, the workflow resumes from Aggregate without re-running Clone, Traverse, or any already-completed SummarizeFile sub-tasks

### Requirement: Workflow progress is streamed to the frontend via SSE

The system SHALL publish progress events for each workflow step transition over the existing SSE channel, so the frontend can update its UI in real time without polling.

#### Scenario: User sees per-file progress
- **WHEN** SummarizeFile sub-tasks complete one by one
- **THEN** the frontend receives `ai-progress` SSE events containing `{step: "summarize_file", fileIndex: N, fileCount: M, filename: "..."}` for each completion and updates a progress bar

#### Scenario: User sees final completion event
- **WHEN** the Store step completes successfully
- **THEN** the frontend receives a final `ai-progress` event with `{step: "store", status: "completed"}` and re-fetches the summary

### Requirement: Summary results are persisted in a new aggregate

The system SHALL persist the workflow output as a `RepoSummary` aggregate in the application's primary database, with one row per run, containing the repo URL, status, run timestamps, per-file summaries, and the final repo-level summary.

#### Scenario: User retrieves their own summary
- **WHEN** an authenticated user GETs `/api/v1/ai/summaries/{id}` for a run they own
- **THEN** the system responds with HTTP 200 and the full `RepoSummary` payload including all per-file summaries and the final repo summary

#### Scenario: User cannot access another user's summary
- **WHEN** an authenticated user GETs `/api/v1/ai/summaries/{id}` for a run owned by another user
- **THEN** the system responds with HTTP 404 (not 403) to avoid leaking existence of other users' runs

#### Scenario: Per-file summary value objects enforce non-empty filename
- **WHEN** a `FileSummary` is constructed via the domain factory
- **THEN** an empty filename returns a domain error and the value object is not created

### Requirement: Frontend page renders the summarization flow

The frontend SHALL provide a protected page at `/ai/summarize` where authenticated users can submit a repo URL and observe live progress.

#### Scenario: Anonymous visitor is redirected to login
- **WHEN** an unauthenticated user navigates to `/ai/summarize`
- **THEN** the frontend redirects to `/login` (consistent with other protected pages)

#### Scenario: Authenticated user can submit and see progress
- **WHEN** an authenticated user enters a repo URL and clicks "Zusammenfassen"
- **THEN** the frontend submits the request, shows a progress UI that updates from SSE events, and renders the final summary when the run completes

#### Scenario: Progress UI shows step name in German
- **WHEN** a workflow step transitions
- **THEN** the German-language step label (e.g. "Repository klonen…", "Dateien analysieren…", "Zusammenfassung erstellen…") is displayed in the UI

### Requirement: Workflow auto-cancels on user disconnect

The system SHALL propagate request cancellation from the HTTP context through the workflow context so that an in-progress run is cancellable when the user closes the page or revokes the session.

#### Scenario: Frontend closes the SSE connection
- **WHEN** the frontend SSE connection for a run closes (user closes the page)
- **THEN** if no other client is subscribed to the run's progress, the workflow's context is cancelled and any in-progress step exits via `ctx.Err() == context.Canceled`

#### Scenario: Cancellation increments the cancelled counter
- **WHEN** a run is cancelled mid-execution
- **THEN** `ai_workflows_completed_total{status="cancelled"}` is incremented by 1 and the `RepoSummary.Status` value is updated to `cancelled`
