# Agent Project Brief

Last updated: March 30, 2026.

This is the fastest reliable entry point for a coding agent joining the Tarot Daily repository. It summarizes the project without the dedicated UI transfer pack.

## Project Essence

Tarot Daily is a tarot-themed entertainment web app. A logged-in user gets:

- one card of the day
- spread A with 3 cards
- spread B with 5 cards
- deterministic daily results based on the user-local date and timezone

The product also includes:

- authentication
- user profile and settings
- admin area
- deck and deck-theme management
- async Gateway-backed ask flow added in `v0.1.2`

## Current Reality

The project is beyond MVP and is usable locally.

Implemented and present in code:

- backend API in Go + Gin
- frontend SPA in Vue 3 + Pinia + Router
- MongoDB persistence
- auth and admin flows
- reading generation flow
- deck uploads served from backend local storage
- static fallback card assets
- dev-only `/ui-kit`
- Gateway ask flow with callback verification

Planned in docs but not yet implemented in code:

- `/dev` index page
- smoke cleanup tooling
- final real-world validation closure for `v0.1.2`

## Active Product Status

Stages:

- MVP: delivered
- `v0.1.1`: delivered
- `v0.1.2`: implemented, but still documented as `validation in progress`
- `v0.1.3`: delivered
- `v0.1.4`: planned as the current roadmap stage

The most important unfinished real work is:

- run the real `v0.1.2` ask-flow smoke in a properly configured environment
- update `v0.1.2` final status documents based on the actual result

## Architecture in One View

Repository shape:

- `backend/` for API and domain logic
- `frontend/` for the SPA
- `deploy/` for smoke scripts and dev helpers
- `specs/` for requirements, status, and roadmap docs

Backend architecture:

- router
- handler
- service
- repository

Frontend architecture:

- pages
- components
- stores
- services
- router guards

Runtime model:

- frontend talks to `/api`
- backend serves uploaded assets from `/uploads`
- backend serves static local card assets from `/assets/cards`

## Core Domain Model

Main entities:

- `User`: auth identity, role, status, timezone, preferred deck
- `Card`: tarot catalog card with open and closed interpretations
- `DailyReading`: one user-day reading record with card of day, spread A, spread B
- `DeckTheme`: legacy theme-based asset mapping
- `Deck`: newer uploadable/admin-managed deck model with `process` and `done` states
- `RolePolicy`: daily spread limit per role
- `AskRequest`: async Gateway-backed request and status record
- `AdminSettings`: global operational settings

Important modeling nuance:

- the project currently mixes legacy `deck_theme` compatibility with newer `deck` flows

## Main User and Admin Flows

User flows:

- register and login
- load current user
- open dashboard and receive today's reading
- generate spreads within role limits
- edit profile and settings
- choose deck or legacy theme

Admin flows:

- access protected admin routes
- inspect users
- manage role policies and user roles
- create decks
- upload deck back images and card assets
- edit global admin settings

Gateway flow:

- frontend calls `POST /api/ask`
- backend creates `ask_request` and sends request to Gateway
- Gateway calls `POST /api/gateway/callback`
- frontend polls `GET /api/ask/:request_id`

## API Surface You Should Assume Exists

Auth and profile:

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/me`
- `GET /api/me/profile`
- `PATCH /api/me/profile`
- `PATCH /api/me/settings`

User content:

- `GET /api/deck-themes`
- `GET /api/decks/available`
- `GET /api/reading/today`
- `POST /api/reading/generate`
- `POST /api/reading/generate/card-of-day`
- `POST /api/reading/generate/spread-a`
- `POST /api/reading/generate/spread-b`

Gateway:

- `POST /api/ask`
- `GET /api/ask/:request_id`
- `POST /api/gateway/callback`
- `GET /api/gateway/health`

Admin:

- `GET /api/admin/access-check`
- `GET/POST /api/admin/roles`
- `PATCH /api/admin/users/role`
- `GET /api/admin/users`
- `GET/POST /api/admin/deck-themes`
- `GET/POST /api/admin/decks`
- `GET /api/admin/decks/:deckId`
- `POST /api/admin/decks/:deckId/assets/back`
- `POST /api/admin/decks/:deckId/assets/cards/:cardKey`
- `GET/PATCH /api/admin/settings`

## Local Environment and Run Model

Tools:

- Go
- Node.js
- npm
- ssh
- PowerShell

Typical local flow:

1. Open SSH tunnel to MongoDB.
2. Prepare Mongo indexes and seed data.
3. Start backend from `backend/`.
4. Start frontend from `frontend/`.

Useful commands:

- backend tests: `go test ./...`
- frontend validation: `npm run build`
- backend smoke: `deploy/scripts/smoke-backend-api.ps1`
- `v0.1.1` smoke: `deploy/scripts/smoke-v0.1.1-api.ps1`
- `v0.1.2` smoke: `deploy/scripts/smoke-v0.1.2-ask.ps1`

Important limitation:

- smoke cleanup tooling is not implemented yet

## Development Constraints

Treat these as project rules, not suggestions:

- UI work must follow `UI Kit First`
- background system must not be changed unless explicitly requested
- `/ui-kit` exists and is the current visual source of truth
- `/dev` should not be assumed to exist
- automated tests are deferred until explicit user request
- external file storage migration is deferred until explicit user request

## Practical Guidance for a New Agent

Assume:

- the repository is operational
- documentation matters
- manual validation is currently preferred over expanding automation
- planned work must be distinguished from already-delivered work

Do not assume:

- `v0.1.2` validation is fully closed
- `/dev` exists
- cleanup utilities exist
- storage migration is active
- Playwright or other automation is part of the current active scope

## Where to Read Next

For full context after this brief:

- `current-status.md`
- `architecture.md`
- `api-contracts.md`
- `environment-and-runbook.md`
- `open-items-and-deferred-work.md`

For UI-specific transfer work:

- use the separate `UI/` folder
