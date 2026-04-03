# AI Agent Development Plan
## Phase 0: Project Setup

## Mandatory Documentation Sync

Before changing any code or creating any project structure, the agent must review and align with these source documents:

1. `specs/docs/PRD.md`
2. `specs/docs/TECHNICAL_SPEC.md`
3. `specs/docs/IMPLEMENTATION_ROADMAP.md`
4. `specs/v0.1.0/phase-0-project-setup/README.md`
5. `specs/v0.1.0/phase-0-project-setup/MANUAL_TESTS.md`

The agent should treat the docs above as the source of truth. If there is a conflict, the agent must pause, document the conflict, and avoid inventing undocumented behavior.

## Mission

Create the minimal but solid project foundation required for all later phases: backend scaffold, frontend scaffold, local configuration, UI Kit entry point, i18n bootstrap, theming bootstrap, and database seed path.

## Delivery Strategy

Work in thin vertical slices. After each slice, verify that the repo is still internally consistent and that the next slice has the prerequisites it needs.

## Step-by-Step Execution

### Step 1. Inspect the current repository state

- verify whether `backend/` already exists
- verify whether `frontend/` already exists
- inspect root-level docs and any existing setup artifacts
- check for an existing `.gitignore`
- check whether any prior bootstrap work already happened

### Step 2. Establish the target directory structure

- create the `backend/` root structure described in `TECHNICAL_SPEC.md`
- create the `frontend/` root structure described in `TECHNICAL_SPEC.md`
- avoid adding feature-specific code that belongs to later phases
- keep placeholders minimal but meaningful

### Step 3. Bootstrap backend runtime

- initialize the Go module
- add the minimal dependency set required by the roadmap
- create the server entry point
- add a minimal config loader pattern
- wire a basic Gin server that can start locally
- keep the implementation intentionally simple and phase-appropriate

### Step 4. Bootstrap frontend runtime

- initialize the Vite + Vue 3 application
- add the required phase-0 dependencies
- create the entry point, router root, app root, and empty page/component structure
- ensure the app renders successfully before any UI depth is added

### Step 5. Prepare local configuration

- create `.env.example` files or equivalent configuration guidance
- define backend environment variables expected by the technical spec
- define `VITE_API_URL` for the frontend
- make sure runtime-only files are ignored by Git

### Step 6. Create the UI Kit shell

- add the dev-only `UIKit.vue` page
- add the `/ui-kit` route guarded to dev usage
- create visible placeholder sections for typography, colors, buttons, inputs, cards, and modals
- do not overbuild components that belong to Phase 2

### Step 7. Add theme infrastructure

- create SCSS variables and theme files
- define the core CSS custom properties for light and dark themes
- implement a minimal theme toggle mechanism
- apply theme classes or data attributes in a way that later phases can extend cleanly

### Step 8. Add i18n infrastructure

- install and wire Vue I18n
- create empty or minimal `en.json` and `ru.json`
- add the language switcher scaffold
- ensure the app can render in both languages without structural issues

### Step 9. Prepare database seed path

- inspect `specs/temp/mock/tarot-cards.json`
- create the seed script location and basic loading flow
- map the source card JSON into the target MongoDB card schema
- create seed logic for role policies, admin settings, and the default deck
- include required indexes from the technical specification

### Step 10. Add basic documentation where needed

- document setup assumptions that are necessary for future agents
- keep documentation focused on actual repository behavior
- avoid documenting features that do not yet exist

### Step 11. Run local verification

- verify backend starts
- verify frontend starts
- verify `/ui-kit` is reachable in dev mode
- verify theme switching works at a basic level
- verify language switching works at a basic level
- verify the seed path is ready and, if possible, runnable

### Step 12. Validate against the manual acceptance checklist

- walk through `MANUAL_TESTS.md`
- note anything intentionally deferred
- do not claim completion if MongoDB seed flow or bootstrap runtime is still broken

## AI Agent Guardrails

- do not begin Phase 1 auth work in this phase
- do not add production-grade abstractions prematurely
- do not create fake endpoints that imply Phase 1 is done
- do not implement tarot business logic yet

## Completion Standard

This phase is complete only when the repository is ready for feature implementation, not merely when folders exist.

