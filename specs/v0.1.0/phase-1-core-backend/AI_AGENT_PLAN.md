# AI Agent Development Plan
## Phase 1: Core Backend

## Mandatory Documentation Sync

Before writing backend code, the agent must review:

1. `specs/docs/PRD.md`
2. `specs/docs/TECHNICAL_SPEC.md`
3. `specs/docs/IMPLEMENTATION_ROADMAP.md`
4. `specs/v0.1.0/phase-1-core-backend/README.md`
5. `specs/v0.1.0/phase-1-core-backend/MANUAL_TESTS.md`

The agent must confirm that endpoint names, payloads, validation rules, and user fields match the documentation exactly before implementation starts.

## Mission

Implement the backend foundation required for authentication, user profile management, JWT-based access control, and timezone-aware accounting-day calculations.

## Step-by-Step Execution

### Step 1. Reconcile API and data model contracts

- extract the required request and response shapes for auth and profile endpoints
- confirm user schema fields from the technical specification
- confirm password, username, and timezone validation rules
- confirm which fields are editable in MVP and which are not

### Step 2. Inspect the Phase 0 backend scaffold

- verify server boot flow
- verify config loading pattern
- verify MongoDB connection strategy
- identify missing packages and extension points

### Step 3. Implement core models

- create user-related models
- create request DTOs for register, login, profile updates, and settings updates
- create JWT claims model
- keep field names and validation tags consistent with the docs

### Step 4. Implement utility layer

- password hashing and verification
- JWT creation and validation
- timezone validation helper
- accounting-date helper
- next-reset-time helper
- any shared validation helpers needed by handlers and services

### Step 5. Implement repositories

- build the user repository
- add the card repository only to the extent Phase 1 requires it
- ensure uniqueness checks for email and username are supported
- keep repository methods narrow and aligned with the roadmap

### Step 6. Implement services

- `Register`
- `Login`
- `LoadUser`
- `GetUserProfile`
- `UpdateProfile`
- `UpdateSettings`

The agent should keep business rules in services, not in handlers.

### Step 7. Implement middleware

- auth middleware for JWT validation
- role middleware for protected role-based access
- CORS middleware
- standard unauthorized and forbidden response behavior

### Step 8. Implement handlers and routes

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/me`
- `PATCH /api/me/profile`
- `PATCH /api/me/settings`

Handlers should:

- validate input
- call the correct service
- return documented response shapes
- avoid embedding business logic

### Step 9. Wire route protection

- protect `/api/me` endpoints with auth middleware
- confirm role middleware is reusable for later phases
- avoid exposing admin-only concerns here

### Step 10. Verify edge cases

- duplicate email
- duplicate username
- weak password
- invalid timezone
- malformed token
- missing token
- boundary behavior around `02:59` and `03:00`

### Step 11. Run manual validation

- test register and login end-to-end
- test profile retrieval and updates
- test settings updates
- test unauthorized access rejection
- verify token claims include `user_id`, `role`, and `timezone`

### Step 12. Cross-check against acceptance

- walk the phase manual tests one by one
- compare real behavior to docs
- record any doc mismatch instead of silently diverging

## AI Agent Guardrails

- do not implement reading generation yet
- do not add admin endpoints in this phase
- do not invent fields outside the documented user model
- do not silently alter API names or payload shapes

## Completion Standard

This phase is complete only when the frontend can reliably authenticate, load the current user, update profile settings, and trust the backend timezone logic.

