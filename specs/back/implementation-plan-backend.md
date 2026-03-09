# Implementation Plan: Backend (Go + Gin + MongoDB)

## 1. Objective
Build an MVP backend for Tarot Daily with:
- auth and role-based access
- daily reading generation
- timezone-aware day boundaries
- deck theme management
- admin-only endpoints

## 2. Scope
- In scope: REST API, MongoDB schema/repositories, auth, reading generation logic, admin access control.
- Out of scope: production-grade deployment, autoscaling, advanced observability.

## 3. Infrastructure Note (Important)
MongoDB is already installed on a remote VPS.  
Development connection must use an SSH tunnel.

## 4. Environment and Connectivity Plan

### 4.1 Environment Variables
- `APP_PORT`
- `JWT_SECRET`
- `MONGO_URI` (local tunnel endpoint, example: `mongodb://127.0.0.1:27017/tarot_daily`)
- `MONGO_DB_NAME`
- `APP_ENV` (`dev`)

### 4.2 SSH Tunnel Workflow (Dev)
1. Open tunnel from local machine to VPS:
   - `ssh -L 27017:127.0.0.1:27017 <user>@<vps-host>`
2. Keep tunnel session active while backend runs.
3. Backend connects via local forwarded port using `MONGO_URI`.
4. Add a short troubleshooting section in README:
   - connection refused
   - authentication errors
   - port conflicts

## 5. Backend Architecture
- Router layer (Gin routes + middleware)
- Handler layer (HTTP)
- Service layer (business logic)
- Repository layer (MongoDB operations)
- Domain models / DTOs

Suggested modules:
- `auth`
- `users`
- `readings`
- `cards`
- `deckthemes`
- `admin`

## 6. Data and Index Plan

Collections:
- `users`
- `cards`
- `daily_readings`
- `deck_themes`

Indexes:
- `users.login` unique
- `daily_readings` unique compound: (`userId`, `date`)
- Optional: `deck_themes.name` unique

Seed data:
- Seed 78 cards (minimum fields from PRD).
- Seed default placeholder deck theme.
- Seed one admin account for local testing.

## 7. Delivery Phases

### Phase 1: Project Bootstrap
- Initialize Go module and Gin app.
- Add config loader and env validation.
- Add Mongo connection manager with startup health check.
- Add graceful shutdown.

### Phase 2: Auth and User Profile
- Implement:
  - `POST /api/auth/register`
  - `POST /api/auth/login`
  - `GET /api/me`
- Password hashing via bcrypt/argon2.
- JWT issuing and verification middleware.
- Add role field (`user` default, `admin` manual assignment).

### Phase 3: Settings and Theme Selection
- Implement:
  - `PATCH /api/me/settings` (timezone, deckTheme)
  - `GET /api/deck-themes`
- Validate timezone format and fallback to `UTC+00:00`.
- Ensure selected deck theme exists.

### Phase 4: Daily Reading Engine
- Implement:
  - `GET /api/reading/today`
- Logic:
  - compute user-local date from stored timezone
  - return existing reading for date if present
  - otherwise generate and save
- Generation rules:
  - card of day: 1 card
  - spread A: 3 cards (no duplicates inside spread)
  - spread B: 5 cards (no duplicates inside spread)
  - allow duplicates across different blocks
  - assign random state `open` or `closed` per card
  - pick interpretation by state
- Resolve image URL from selected `deckTheme`.

### Phase 5: Admin API
- Implement admin middleware (`role == admin`).
- Implement:
  - `GET /api/admin/deck-themes`
  - `POST /api/admin/deck-themes`
- Validate payload for deck theme creation (MVP minimal schema).
- Return 403 for non-admin users.

### Phase 6: Validation and Error Handling
- Standard error response format.
- Request validation for all input payloads.
- Clear HTTP status policy:
  - 400 validation errors
  - 401 unauthorized
  - 403 forbidden
  - 404 not found
  - 500 internal errors

### Phase 7: Testing
- Unit tests:
  - auth service
  - reading generation service
  - timezone date resolution
- Handler tests for key endpoints.
- Integration smoke tests against tunneled MongoDB.

### Phase 8: Dev Documentation
- Add README sections:
  - local run steps
  - env vars
  - SSH tunnel command and verification
  - seed command
  - API endpoint list

## 8. Security Baseline
- Hash passwords only, never store plaintext.
- Protect private endpoints with JWT middleware.
- Protect admin endpoints with role-check middleware.
- Sanitize and validate all inputs.
- Do not log secrets or raw tokens.

## 9. Definition of Done (Backend)
- All MVP endpoints implemented and tested at basic level.
- Daily reading logic matches PRD rules, including `open/closed` states.
- Timezone-aware day calculation works with default fallback `UTC+00:00`.
- Admin endpoints are accessible only to `admin`.
- Backend reliably works with MongoDB over SSH tunnel in local dev.
