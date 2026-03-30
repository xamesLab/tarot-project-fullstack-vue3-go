# Open Items and Deferred Work

Last updated: March 30, 2026.

This file separates real open work from intentionally postponed work.

## Active Open Items

### 1. Close `v0.1.2` real validation

What is missing:

- run the actual `v0.1.2` ask-flow smoke in a properly configured environment
- confirm Gateway reachability
- confirm callback verification configuration
- update final validation documents with the real outcome

Why it matters:

- the implementation exists, but the stage is still documented as `validation in progress`

### 2. Implement planned `v0.1.4` dev hygiene work

Planned but not yet delivered:

- `/dev` developer index page
- smoke cleanup tooling

Why it matters:

- roadmap and summary docs should continue to distinguish between planned work and implemented work

## Deferred by Explicit Product Decision

### Automated tests

Deferred until the user explicitly requests them:

- E2E automation
- visual regression automation
- broader automated UI test coverage

Practical meaning:

- do not treat test automation as the active stage
- do not expand or repair an automation layer unless requested

### External file storage

Deferred until the user explicitly requests it:

- object storage migration
- S3-compatible upload architecture
- asset URL migration strategy

Practical meaning:

- current backend-served local storage is the active model

## Medium-Priority Future Work

Documented future candidates:

- MongoDB backup and restore runbook
- changelog and release-notes discipline

These are not blocked forever, but they are not the current main track either.
