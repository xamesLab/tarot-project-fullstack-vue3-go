# AI Agent Development Plan
## Phase 3: Reading Features

## Mandatory Documentation Sync

Before implementing tarot business logic, the agent must review:

1. `specs/docs/PRD.md`
2. `specs/docs/TECHNICAL_SPEC.md`
3. `specs/docs/IMPLEMENTATION_ROADMAP.md`
4. `specs/v0.1.0/phase-3-reading-features/README.md`
5. `specs/v0.1.0/phase-3-reading-features/MANUAL_TESTS.md`

The agent must specifically verify the documented rules for:

- deterministic card-of-the-day generation
- accounting day reset at `03:00` local time
- role-based spread limits
- history retention behavior

## Mission

Implement the core tarot experience of the MVP across backend and frontend.

## Step-by-Step Execution

### Step 1. Reconfirm domain rules before coding

- review card-of-day determinism requirements
- review three-card spread constraints
- review retention policy behavior
- confirm route names and expected responses

### Step 2. Inspect available card and policy data

- verify seeded cards exist and match the expected schema
- verify seeded `role_policies` exist
- verify default deck assumptions
- identify whether repository access patterns already support the new logic

### Step 3. Implement reading models and repository

- create reading entities
- create card-in-reading structures
- implement persistence for generated readings
- implement history retrieval with sorting and filtering support
- keep TTL and retention concerns aligned with the data model

### Step 4. Implement card-of-day service logic

- compute accounting date from timezone
- check whether a card already exists for the same user, date, and reading type
- compute deterministic seed from `user_id + accounting_date + timezone`
- select exactly one card
- save and return the reading
- expose next reset timing when needed by UI

### Step 5. Implement three-card spread service logic

- compute accounting date
- fetch user role policy
- count already generated spreads for the accounting day
- reject generation when the limit is reached
- select three unique cards
- persist the reading
- return remaining count information

### Step 6. Implement history service logic

- fetch user history
- support filter by reading type
- support pagination
- respect retention expectations from the user role

### Step 7. Implement backend endpoints

- `POST /api/reading/generate/card-of-day`
- `POST /api/reading/generate/spread-a` or the final route name used by the implementation
- `GET /api/reading/history`

Do not rename routes casually. If the implementation uses a different final route than the roadmap draft, document the reason.

### Step 8. Implement tarot UI components

- `TarotCard`
- `CardOfDay`
- `ThreeCardSpread`
- `ReadingHistory`

These components should support:

- localized content
- loading and error states
- responsive layout
- future deck-image extensibility

### Step 9. Extend reading services and stores

- add frontend API calls for generation and history
- extend the `readings` store with fetch and generate actions
- expose computed state for availability and limits
- keep client-side state consistent with backend responses

### Step 10. Implement dashboard behavior

- load existing card-of-day state
- generate card-of-day on demand
- generate spreads on demand
- show remaining spread count
- handle first-time and already-generated cases

### Step 11. Implement history page

- render reading list
- support filter controls
- support pagination
- display retention-related messaging

### Step 12. Validate edge cases

- repeated card-of-day requests on the same accounting day
- exact spread-limit exhaustion
- timezone change effects
- `02:59` versus `03:00`
- localization of tarot data and errors

### Step 13. Run acceptance walkthrough

- execute the manual phase checklist
- verify behavior for at least the `seeker` role
- verify history and limits reflect role policies correctly

## AI Agent Guardrails

- do not introduce payment or monetization logic
- do not implement AI interpretation features
- do not add reversed cards or advanced spread semantics
- do not bypass documented role-policy checks in the UI

## Completion Standard

This phase is complete only when the end user can complete the core product journey: get the daily card, generate allowed spreads, and review history in a localized UI.

