# Project Structure Summary

Last updated: March 30, 2026.

## Repository Layout

```text
tarot_project/
  backend/
  frontend/
  deploy/
  specs/
  README.md
  AGENTS.md
```

## `backend/`

Main areas:

- `cmd/`
  - runnable commands such as `api`, `dbinit`, `indexes`, `migrate011`, `mongocheck`, `seed`
- `internal/`
  - application code split by domain modules and platform helpers
- `storage/uploads/`
  - current local storage for uploaded deck assets
- `.env.example`
- `README.md`

Important backend subareas:

- `internal/auth`
- `internal/reading`
- `internal/decks`
- `internal/ask`
- `internal/adminusers`
- `internal/adminsettings`
- `internal/transport/http`
- `internal/repository/mongo`
- `internal/platform/*`

## `frontend/`

Main areas:

- `src/pages`
- `src/components`
- `src/stores`
- `src/services`
- `src/router`
- `src/assets/styles`
- `.env.example`
- `README.md`

Important frontend pages currently present:

- `LoginPage.vue`
- `RegisterPage.vue`
- `DashboardPage.vue`
- `ProfilePage.vue`
- `SettingsPage.vue`
- `admin/AdminUsersPage.vue`
- `admin/AdminDecksPage.vue`
- `admin/AdminSettingsPage.vue`
- `dev/UIKit.vue`

Important clarification:

- `dev/DevIndexPage.vue` is not present in the current repository state

## `deploy/`

Contains operational scripts and local development helpers, including:

- backend smoke scripts
- `v0.1.1` and `v0.1.2` smoke scripts
- environment preflight check
- MongoDB index helper wrapper
- seed payload example

## `specs/`

Contains:

- product-level documents
- backend and frontend implementation specs
- versioned stage documents
- roadmap and backlog
- this summary pack

## Practical Recommendation

When onboarding into the repository:

1. read `specs/summary-0.1/README.md`
2. read the summary PRD and current status
3. inspect `README.md`, `backend/README.md`, and `frontend/README.md`
4. move into version-specific documents only when working on a specific stage
