# AI Agent Development Plan
## Phase 0: AI Foundation

## Mandatory Documentation Sync

Before implementation, review:

1. `specs/docs/PRD.md`
2. `specs/docs/TECHNICAL_SPEC.md`
3. `specs/docs/IMPLEMENTATION_ROADMAP.md`
4. `specs/v0.1.1/README.md`
5. `specs/v0.1.1/phase-0-ai-foundation/README.md`
6. `specs/v0.1.1/phase-0-ai-foundation/MANUAL_TESTS.md`

## Mission

Prepare the backend and data model for direct AI integration through `OpenRouter`, while preserving a future fallback-friendly provider architecture.

## Step-by-Step Execution

### Step 1. Inspect current gateway and ask-request stubs

- confirm what already exists in handlers, models, and docs
- identify which parts are placeholders only
- preserve route continuity where possible

### Step 2. Design provider abstraction

- define the minimal interface for AI generation
- keep it compatible with both direct provider and future gateway fallback
- avoid leaking transport concerns into handlers

### Step 3. Extend config layer

- add direct AI config fields to backend config
- load secret and model name from env
- keep secrets out of logs and API responses

### Step 4. Update ask-request model

- add provider and model metadata
- add request lifecycle timestamps and status
- add usage/error fields needed for observability

### Step 5. Define endpoint behavior

- settle final request/response contract for `/api/ask`
- decide reuse and idempotency rules
- document authorization checks

### Step 6. Validate against manual checklist

- confirm config behavior with missing secrets
- confirm no secret values are exposed
- confirm the architecture still allows future gateway fallback

## Guardrails

- do not hardcode secret values
- do not expose AI keys to frontend
- do not remove gateway code paths outright
- do not introduce deployment concerns in this phase

## Completion Standard

This phase is complete only when the codebase is structurally ready for direct AI calls and the future gateway fallback remains possible.
