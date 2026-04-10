# AI Agent Development Plan
## Phase 4: Polish and Manual Testing

## Mandatory Documentation Sync

Before implementation, review:

1. `specs/v0.1.1/README.md`
2. `specs/v0.1.1/phase-4-polish-and-manual-testing/README.md`
3. `specs/v0.1.1/phase-4-polish-and-manual-testing/MANUAL_TESTS.md`
4. all previous `v0.1.1` phase README files

## Mission

Stabilize the direct AI release candidate, close major UX gaps, and complete the final manual acceptance for `v0.1.1`.

## Step-by-Step Execution

### Step 1. Audit the complete AI journey

- new user flow
- returning user flow
- admin governance flow
- mobile and localization flow

### Step 2. Fix the highest-risk issues

- broken or confusing loading states
- repeated request bugs
- unstable restore behavior after reload
- mobile layout regressions

### Step 3. Run the manual checklist end-to-end

- execute the documented scenarios
- record blocking vs non-blocking issues
- retest after fixes

### Step 4. Update docs

- capture release readiness
- note any accepted limitations
- keep deployment explicitly out of scope

## Guardrails

- do not treat deployment as part of this phase
- do not mark the version done until the manual flows pass
- do not hide unresolved blockers behind documentation wording

## Completion Standard

This phase is complete only when `v0.1.1` feels coherent, stable, and manually accepted as a non-deployment version with direct AI integration.
