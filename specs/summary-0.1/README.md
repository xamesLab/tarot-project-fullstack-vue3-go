# Tarot Daily Summary 0.1

Last updated: March 30, 2026.

This folder is an English handoff pack for the current repository state of Tarot Daily. It is meant to help a new contributor understand what the product is, what is already implemented, what is only planned, how the codebase is organized, and which constraints currently shape development.

## Recommended Reading Order

1. `agent-project-brief.md`
2. `PRD.md`
3. `current-status.md`
4. `development-plan.md`
5. `architecture.md`
6. `domain-model.md`
7. `api-contracts.md`
8. `database-and-dev-connection.md`
9. `environment-and-runbook.md`
10. `project-structure.md`
11. `ui-and-development-rules.md`
12. `open-items-and-deferred-work.md`

## UI Deep Dive

For a transfer-focused UI Kit handoff, use:

- `UI/README.md`
- `UI/ui-kit-source-of-truth.md`
- `UI/visual-tokens-and-typography.md`
- `UI/component-specification.md`
- `UI/layout-responsive-and-interaction.md`
- `UI/agent-transfer-playbook.md`

## What This Summary Covers

- product scope and intended user value
- actual implementation status across MVP and versioned stages
- backend, frontend, database, and integration architecture
- high-level API contracts
- MongoDB and local development connection model
- local setup and validation commands
- repository structure
- UI and workflow constraints
- open items, planned work, and explicitly deferred work

## Important Notes About Accuracy

This summary was aligned to the current repository contents, not only to roadmap documents. In particular:

- `/ui-kit` exists today and is wired as a dev-only route.
- `/dev` is described in planning documents, but it is not implemented in the current codebase yet.
- smoke cleanup tooling is planned in `v0.1.4`, but no cleanup command or script exists in the repository today.
- automated test work is intentionally paused until the user explicitly asks for it.
- external file storage migration is also intentionally paused until the user explicitly asks for it.
