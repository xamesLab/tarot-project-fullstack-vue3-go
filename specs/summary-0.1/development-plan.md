# Development Plan by Major Stages

Last updated: March 30, 2026.

## Stage 0: MVP Foundation

Delivered:

- monorepo structure
- backend and frontend bootstrap
- auth
- daily reading flow
- settings
- admin area baseline
- smoke scripts and local runbook

Outcome:

- locally reproducible product baseline

## Stage `v0.1.1`

Delivered:

- profile and deck improvements
- richer admin flows
- related documentation and validation artifacts

Outcome:

- stronger product completeness after MVP

## Stage `v0.1.2`

Delivered:

- backend Gateway integration
- async ask flow
- callback processing
- callback security
- polling model for frontend

Outcome:

- backend contract for Gateway-driven async requests exists

Still open:

- execute a real integration smoke run and update `v0.1.2` status documents with the actual result

## Stage `v0.1.3`

Delivered:

- `/ui-kit`
- dev-only registration for the UI Kit route
- Arcana background system
- UI Kit-first governance
- mobile design harmonization for the reference page

Outcome:

- frontend now has a visual reference workflow and stronger UI guardrails

## Planned Stage `v0.1.4`

Documented scope:

- complete real validation of `v0.1.2`
- add `/dev` as a dev-tools index
- improve smoke-data cleanup and reproducibility

Important reality check:

- this stage is planned in documentation
- its major items are not fully implemented in the current repository yet

## Recommended Execution Order

1. Prepare a runtime environment where backend, MongoDB, Gateway access, and callback security are configured correctly.
2. Run the real `v0.1.2` ask-flow smoke and capture the result.
3. Update `v0.1.2` final validation documents.
4. Add `/dev` only after validation work is closed or at least unblocked.
5. Add smoke cleanup tooling for repeatable local and smoke cycles.

## Explicitly Deferred Work

- automated tests
- visual regression automation
- external file storage migration

These should stay out of scope until the user explicitly reopens them.
