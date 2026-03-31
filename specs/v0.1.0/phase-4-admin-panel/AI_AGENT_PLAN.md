# AI Agent Development Plan
## Phase 4: Admin Panel

## Mandatory Documentation Sync

Before implementing admin functionality, the agent must review:

1. `specs/docs/PRD.md`
2. `specs/docs/TECHNICAL_SPEC.md`
3. `specs/docs/IMPLEMENTATION_ROADMAP.md`
4. `specs/v0.1.0/phase-4-admin-panel/README.md`
5. `specs/v0.1.0/phase-4-admin-panel/MANUAL_TESTS.md`

The agent must re-check:

- keeper-only access rules
- editable role policy fields
- editable global settings fields
- expected deck upload validation rules

## Mission

Build the administrative tools required to operate the MVP safely and manage its users, policies, settings, and decks.

## Step-by-Step Execution

### Step 1. Confirm admin scope

- list the exact admin endpoints required
- confirm that only `keeper` has admin access
- confirm the allowed operations on users, roles, settings, and decks
- verify that no out-of-scope admin features are added

### Step 2. Inspect the current backend and frontend shells

- verify route middleware is ready for keeper-only protection
- verify existing stores and router patterns can support admin sections
- identify reusable table, modal, and form primitives

### Step 3. Implement backend admin services

- dashboard metrics aggregation
- user listing with filters and pagination
- role change logic
- account disable logic
- role policy retrieval and update logic
- global settings retrieval and update logic

### Step 4. Implement backend admin endpoints

- access-check endpoint
- dashboard endpoint
- users endpoint
- user role update endpoint
- roles retrieval and update endpoints
- settings retrieval and update endpoints

Every admin route must be protected with keeper-only middleware.

### Step 5. Implement deck management backend

- create deck repository and service logic
- support deck metadata creation
- support ZIP upload handling
- validate ZIP structure and image constraints
- extract valid archives into the deck storage area
- move deck status through the documented lifecycle
- clean up temporary artifacts after processing

### Step 6. Build admin frontend shell

- add `AdminLayout`
- add sidebar navigation
- add breadcrumbs or a comparable orientation pattern
- expose admin navigation only to keeper users

### Step 7. Build admin pages

- `AdminDashboard`
- `Users`
- `RolePolicies`
- `Settings`
- `Decks`

Pages should use shared UI primitives instead of ad hoc markup.

### Step 8. Build admin data layer

- `services/admin.js`
- `stores/admin.js`
- page-level loading and error management
- optimistic updates only when safe and easy to reason about

### Step 9. Validate deck upload behavior

- valid archive path
- invalid archive path
- missing cards
- wrong file formats
- wrong image dimensions or size
- cleanup of rejected uploads

### Step 10. Validate role-sensitive behavior

- non-keeper users cannot access admin routes
- role changes affect downstream permissions as expected
- policy changes affect later reading behavior

### Step 11. Run manual acceptance

- execute all admin manual test flows
- test EN and RU
- test both successful and rejected operations

## AI Agent Guardrails

- do not make admin endpoints accessible to non-keeper roles
- do not skip file validation for deck uploads
- do not expose destructive actions without explicit confirmation in the UI
- do not silently mutate user-facing business rules outside documented policy fields

## Completion Standard

This phase is complete only when a keeper user can safely operate the system through the admin UI and all admin actions are correctly enforced by the backend.

