# AI Agent Development Plan
## Phase 5: Polish and Testing

## Mandatory Documentation Sync

Before polishing or fixing behavior, the agent must review:

1. `specs/docs/PRD.md`
2. `specs/docs/TECHNICAL_SPEC.md`
3. `specs/docs/IMPLEMENTATION_ROADMAP.md`
4. `specs/v0.1.0/phase-5-polish-and-testing/README.md`
5. `specs/v0.1.0/phase-5-polish-and-testing/MANUAL_TESTS.md`

The agent must compare the current implementation against documented MVP expectations before calling anything "done".

## Mission

Bring the product from feature-complete to release-ready by tightening UX quality, resilience, responsiveness, accessibility, and operational confidence.

## Step-by-Step Execution

### Step 1. Audit the current product state

- inspect all implemented screens
- compare implemented behavior to documented behavior
- list visible inconsistencies and functional gaps
- separate critical defects from polish issues

### Step 2. Audit UI consistency

- compare spacing, typography, color usage, radius, and shadows across pages
- verify UI Kit reflects the real component library
- normalize duplicate patterns

### Step 3. Audit responsive behavior

- test dashboard on mobile
- test history on mobile
- test admin pages on mobile and tablet
- test hidden/revealed card states on dashboard across breakpoints
- test admin archive views on mobile and tablet
- fix layout breakpoints and overflow issues

### Step 4. Audit accessibility

- keyboard navigation
- modal close behavior
- focus visibility
- input labeling
- color contrast
- readable validation and error states
- keyboard reveal for hidden cards

### Step 5. Audit performance

- inspect large image usage
- confirm lazy loading where appropriate
- confirm code splitting strategy
- measure important page and API response times
- fix obvious bottlenecks first

### Step 6. Audit error handling

- network disconnects
- backend `500` responses
- unauthorized and forbidden flows
- malformed requests
- empty states
- stale or deleted resource handling
- archived resource fallback behavior

### Step 7. Validate domain edge cases

- timezone changes with an active daily card
- exact spread limit boundary
- invalid deck uploads
- `0_meta.json` in uploaded deck archives
- repeated reveal of an already opened card
- access to archived or deleted decks/card backs
- accounting-day edge behavior

### Step 8. Fix defects in priority order

- critical
- high
- medium
- low

The agent should keep a defect list and close issues in a traceable way.

### Step 9. Re-run full manual journeys

- new user journey
- returning user journey
- admin journey
- theme journey
- i18n journey
- browser sanity checks

### Step 10. Document residual risk

- record what was verified
- record what could not be verified
- record any non-blocking issues accepted for MVP

## AI Agent Guardrails

- do not use polish work to introduce new scope
- do not hide broken behavior behind cosmetic changes
- do not claim accessibility compliance without actual verification
- do not treat unresolved critical bugs as acceptable for release
- do not break persisted reveal-state while polishing dashboard UI

## Completion Standard

This phase is complete only when the implemented product feels coherent, stable, and trustworthy across the documented MVP journeys.
