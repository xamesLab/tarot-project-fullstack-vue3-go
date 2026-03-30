# Architecture Summary

Last updated: March 30, 2026.

## System Style

Tarot Daily uses a straightforward monorepo architecture:

- frontend SPA
- backend REST API
- MongoDB database
- local development scripts and smoke tooling

The backend follows a layered structure:

- router
- handler
- service
- repository

Authentication is stateless and token-based.

## Frontend

### Stack

- Vue 3
- TypeScript
- Pinia
- Vue Router
- Vite

### Frontend Responsibilities

- authentication flow
- route guards
- user dashboard
- profile and settings pages
- admin shell and admin pages
- design-system reference via `/ui-kit`

### Current Frontend Runtime Notes

- `/ui-kit` is the only confirmed dev-only route in the current code
- `/dev` is planned, but not implemented yet
- route protection is handled in router guards
- visual changes are expected to align with the UI Kit
- background visuals are isolated in `AppBackground.vue`

## Backend

### Stack

- Go
- Gin
- MongoDB driver
- JWT

### Backend Responsibilities

- auth and token issuance
- user profile and settings handling
- daily reading generation and retrieval
- deck and theme management
- admin endpoints
- Gateway integration for async ask-flow
- callback verification and diagnostics

### Runtime Interfaces

- REST API under `/api`
- static deck uploads under `/uploads`
- static local card assets under `/assets/cards`

## Database

### Main Collections

- `users`
- `cards`
- `daily_readings`
- `deck_themes`
- `decks`
- `role_policies`
- `ask_requests`
- `admin_settings`

### Key Indexes

- `users.login` unique
- `daily_readings(userId, date)` unique
- `cards.code` unique
- `deck_themes.name` unique
- `role_policies.role` unique
- `ask_requests.requestId` unique
- `ask_requests(userId, createdAt)` non-unique lookup index

## Key Runtime Flows

### Auth Flow

Frontend sends login or register request -> backend validates and issues JWT -> frontend stores token and hydrates user state.

### Reading Flow

Frontend requests today’s reading -> backend resolves the user-local date from the saved timezone -> backend returns an existing reading or generates a new one -> frontend renders card and spread blocks.

### Admin Flow

Frontend enters admin routes -> backend checks JWT and admin role -> admin pages read or mutate protected resources.

### Async Ask Flow

Frontend calls `POST /api/ask` -> backend creates an `ask_request` and forwards the request to Gateway -> Gateway later calls backend callback -> frontend polls `GET /api/ask/:request_id`.

## Current Architectural Boundaries

- no finalized production deployment architecture
- no active external object storage architecture
- no active automated test architecture at the current stage because that work is intentionally paused
