# AI Agent Development Plan
## Phase 2: Core Frontend

## Mandatory Documentation Sync

Before implementing the frontend shell, the agent must review:

1. `specs/docs/PRD.md`
2. `specs/docs/TECHNICAL_SPEC.md`
3. `specs/docs/IMPLEMENTATION_ROADMAP.md`
4. `specs/v0.1.0/phase-2-core-frontend/README.md`
5. `specs/v0.1.0/phase-2-core-frontend/DESIGN_SYSTEM.md`
6. `specs/v0.1.0/phase-2-core-frontend/MANUAL_TESTS.md`

The agent must verify route structure, component inventory, auth flows, and i18n expectations before editing the frontend.

If there is any overlap between central project docs and the local phase design-system document, implementation decisions for this phase must follow `specs/v0.1.0/phase-2-core-frontend/DESIGN_SYSTEM.md`.

## Mission

Build the application shell that users interact with before tarot features are added: routing, auth pages, shared components, app stores, and dashboard placeholders.

## Step-by-Step Execution

### Step 1. Inspect the frontend scaffold

- verify Vite entry files
- verify router bootstrap
- verify theme and i18n setup from Phase 0
- identify what is already reusable versus what still needs to be created

### Step 2. Confirm backend contracts

- inspect the actual auth and profile endpoints from Phase 1
- align frontend services with real payloads and responses
- do not guess API shapes if the backend differs from the docs

### Step 3. Build the API client

- create the Axios instance
- add base URL configuration
- add request interceptor for JWT
- add response handling for `401` and `403`
- centralize API error normalization where possible

### Step 4. Build service modules

- `services/auth.js`
- `services/user.js`
- placeholder `services/reading.js` if needed by dashboard scaffolding

Keep services thin and route all stateful logic into stores.

### Step 5. Build global stores

- implement `auth` store
- implement `app` store
- keep `readings` store as a Phase 3-ready placeholder
- persist only the values that are actually intended to survive reloads

### Step 6. Build routing and guards

- add public routes
- add protected routes
- add admin route group
- add dev-only route
- enforce auth guard behavior
- enforce keeper-only behavior for admin routes

### Step 7. Build shared UI primitives

- `Button`
- `Input`
- `Select`
- `Modal`

These primitives should:

- respect the phase design-system tokens and component map
- be reusable by auth pages and later tarot/admin pages
- be visible in the UI Kit

### Step 8. Build layout shell

- `Header`
- `Footer`
- optional reusable app container structure
- language switcher
- theme toggle
- logout action
- conditional admin navigation entry

### Step 9. Build composables

- `useTheme`
- `useAuth`

Keep composables small and focused on view-facing convenience logic.

### Step 10. Implement auth pages

- `Login.vue`
- `Register.vue`
- form validation
- loading states
- API error states
- success redirects
- localized labels and messages

### Step 11. Implement dashboard skeleton

- authenticated landing page
- greeting section
- placeholder areas for card-of-day and spread features
- visible empty/loading states
- avoid prematurely implementing reading logic here

### Step 12. Verify key flows

- unauthenticated user gets redirected to `/login`
- registered user is redirected to `/`
- persisted session behavior works as intended
- theme and language changes are visible in shell UI
- `/ui-kit` shows the created base components

### Step 13. Validate against manual acceptance

- run through the phase manual checklist
- verify both EN and RU
- verify that routing, guards, and auth flow behave predictably

## AI Agent Guardrails

- do not implement real tarot reading generation in this phase
- do not add admin business pages yet
- do not create one-off UI patterns that bypass the shared component layer
- do not introduce frontend state that contradicts backend truth
- do not resolve design ambiguity by inventing a third UI direction; prefer `specs/v0.1.0/phase-2-core-frontend/DESIGN_SYSTEM.md`

## Completion Standard

This phase is complete only when a user can register, log in, land on a protected dashboard shell, and navigate a stable base application layout.
