# AI Agent Development Plan
## Pre-Phase 2: UI Kit Foundation

## Mandatory Documentation Sync

Before changing the UI Kit or any frontend design-system code, the agent must review:

1. `specs/docs/PRD.md`
2. `specs/docs/TECHNICAL_SPEC.md`
3. `specs/docs/IMPLEMENTATION_ROADMAP.md`
4. `specs/v0.1.0/phase-2-core-frontend/pre-phase-2-core-frontend/README.md`
5. `specs/v0.1.0/phase-2-core-frontend/pre-phase-2-core-frontend/MANUAL_TESTS.md`
6. `specs/UI/UIKit.vue`

The agent must treat `specs/UI/UIKit.vue` as the source template for this task. The goal is not to invent a new UI Kit from scratch, but to adapt and implement it inside the actual frontend codebase.

## Mission

Turn the current UI Kit route into a real design-system and component-library reference page that future frontend work must follow by default.

## Step-by-Step Execution

### Step 1. Reconcile the design source

- inspect `specs/UI/UIKit.vue` in full
- identify its sections, component patterns, and visual language
- map its structure onto the current frontend architecture
- list what can be reused and what must be adapted

### Step 2. Inspect the current UI Kit implementation

- review the current dev route and current `UIKit.vue`
- compare the existing implementation with the source template
- identify missing sections, missing states, and visual mismatches

### Step 3. Define the UI Kit responsibilities

The implemented page must act as:

- a design-system homepage
- a component showcase
- a visual token reference
- a source of truth for future frontend work

The agent must preserve this role while keeping the page autonomous.

### Step 4. Implement the UI Kit page

- adapt the source template from `specs/UI/UIKit.vue`
- preserve its overall design intent
- structure the page into clear sections
- ensure the page is maintainable inside the actual frontend app

### Step 5. Make the page autonomous

- do not depend on auth state
- do not depend on API responses
- do not require real user data
- use static or demo data where necessary

### Step 6. Enforce dev-only behavior

- expose the route only in development mode
- prevent access in production mode
- ensure accidental production exposure does not happen

### Step 7. Ensure future component usefulness

- showcase tokens and primitives in a way that later pages can reuse
- include meaningful component states
- avoid turning the page into a decorative mockup only

### Step 8. Verify practical coverage

The page should visibly cover:

- typography
- colors
- buttons
- form controls
- cards and surfaces
- modal pattern
- layout/system patterns
- key component states

### Step 9. Update project agent guidance

- add a project-level agent note that frontend development must follow the UI Kit by default
- allow deviations only when the user explicitly asks for a different direction

### Step 10. Validate manually

- run the manual checklist for this pre-phase
- verify dev access, production restriction, and autonomy
- confirm the page is genuinely useful as a frontend reference

## AI Agent Guardrails

- do not make the UI Kit auth-protected
- do not make the UI Kit depend on backend calls
- do not replace the provided source template with an unrelated design
- do not build future pages in styles that ignore the UI Kit unless the user explicitly asks for it

## Completion Standard

This pre-phase is complete only when the UI Kit is a real, autonomous, dev-only design-system reference and the project agent guidance explicitly tells future frontend work to follow it.

