# AI Agent Development Plan
## Phase 2: AI Reading Experience

## Mandatory Documentation Sync

Before implementation, review:

1. `specs/docs/DESIGN_SYSTEM.md`
2. `specs/v0.1.1/README.md`
3. `specs/v0.1.1/phase-2-ai-reading-experience/README.md`
4. `specs/v0.1.1/phase-2-ai-reading-experience/MANUAL_TESTS.md`

## Mission

Integrate AI interpretation into the existing reading UX without breaking the current dashboard/history experience or UI-system alignment.

## Step-by-Step Execution

### Step 1. Inspect current reading screens

- identify where fully revealed reading state is already available
- identify how the existing `Meaning` / `Значение расклада` block is currently populated
- prefer reuse of existing surface/header components

### Step 2. Add frontend service/store support

- implement `/api/ask` client calls
- keep per-reading request state
- preserve completed results across reloads when possible

### Step 3. Implement UI states

- idle state with CTA
- loading state
- completed state with rendered AI text inside the existing `Meaning` block
- error state with retry

### Step 4. Keep the UX restrained

- use the existing `Meaning` block instead of inventing a second interpretation container
- do not overload mobile with oversized AI panels
- keep all UI aligned with the design system

### Step 5. Validate manual flows

- card of day flow
- three-card spread flow
- reload persistence
- RU/EN and mobile/desktop sanity

## Guardrails

- do not expose provider internals directly in UI copy
- do not make AI available before the product rule allows it
- do not add a separate visual language outside the current UI kit

## Completion Standard

This phase is complete only when AI interpretation feels native to the product and remains stable across reading types and breakpoints.
