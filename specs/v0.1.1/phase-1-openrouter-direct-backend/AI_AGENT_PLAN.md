# AI Agent Development Plan
## Phase 1: OpenRouter Direct Backend

## Mandatory Documentation Sync

Before implementation, review:

1. `specs/docs/TECHNICAL_SPEC.md`
2. `specs/v0.1.1/README.md`
3. `specs/v0.1.1/phase-0-ai-foundation/README.md`
4. `specs/v0.1.1/phase-1-openrouter-direct-backend/README.md`
5. `specs/v0.1.1/phase-1-openrouter-direct-backend/MANUAL_TESTS.md`

## Mission

Turn the AI stub into a real backend integration using direct `OpenRouter` calls.

## Step-by-Step Execution

### Step 1. Inspect current handler/service layout

- locate current `/api/ask` stub behavior
- determine whether health or callback routes need coexistence
- preserve compatible route names where reasonable

### Step 2. Implement provider client

- add HTTP client with timeout and header injection
- support bearer auth and content-type
- optionally inject attribution headers from config

### Step 3. Build prompt and response handling

- convert reading data into provider messages
- include revealed cards and reading type in the prompt
- include spread positions for three-card readings
- explicitly request a compact interpretation only
- constrain the intended response length to `100` words maximum
- capture result text, usage, model, and status
- validate and normalize the model response

### Step 4. Persist request lifecycle

- create/update `ask_requests`
- mark `pending -> processing -> completed/failed`
- store safe diagnostics for failed calls

### Step 5. Expose endpoints

- wire `POST /api/ask`
- wire `GET /api/ask/:request_id`
- enforce ownership and auth checks

### Step 6. Add AI health diagnostics

- verify config presence
- verify provider reachability at a minimal level
- keep responses secret-safe

### Step 7. Run manual acceptance

- execute happy path
- validate auth failures
- validate provider failure handling
- validate repeated request behavior

## Guardrails

- do not move secret handling to frontend
- do not silently swallow provider failures
- do not make gateway removal a prerequisite for direct AI
- do not depend on deployment setup in this phase

## Completion Standard

This phase is complete only when the backend can request, persist, and return AI interpretations through `OpenRouter` in a production-like local flow.
