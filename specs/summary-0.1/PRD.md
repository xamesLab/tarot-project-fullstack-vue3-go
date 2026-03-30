# Product Requirements Document

Last updated: March 30, 2026.

## Product Name

Tarot Daily

## Product Goal

Tarot Daily is an entertainment web application in the tarot domain. The product gives a logged-in user a compact mystical experience centered around a daily tarot session, personal preferences, and a lightweight account area.

## Core Value Proposition

- one card of the day
- two spread formats
- deterministic daily result per user and per timezone-aware date
- personal account and preference management
- admin area for content and operational management

## Primary Audience

- women and girls interested in mysticism and tarot
- short repeat visits, usually once or twice per day
- entertainment-first usage rather than professional divination tooling

## Product Scope

### User-facing scope

- register and log in with login and password
- keep a token-based authenticated session
- open a dashboard with:
  - card of the day
  - spread A
  - spread B
- reveal cards in the UI
- manage timezone and preferred deck or theme settings
- get the same daily result within the same user-local day

### Admin-facing scope

- access a protected admin area
- manage users and user roles
- manage decks and upload deck assets
- manage global admin settings

## Daily Content Rules

- card of the day: 1 card
- spread A: 3 cards
- spread B: 5 cards
- no duplicates inside a single spread
- duplicates between different blocks are allowed
- each card has `open` or `closed` state
- interpretation depends on state
- results are fixed for the user-local day

## Non-Goals at the Current Stage

- production infrastructure design
- subscriptions or payments
- full editorial CMS
- advanced observability platform work
- external object storage rollout
- automated test expansion unless the user explicitly requests it

## Delivered Product Evolution

### MVP

Delivered:

- backend and frontend bootstrap
- auth and role checks
- settings and daily reading flow
- admin area baseline
- local runbook and smoke scripts
- mobile-first baseline

### Post-MVP `v0.1.1`

Delivered:

- richer profile and deck flows
- admin user, deck, and settings screens
- related validation and smoke artifacts

### Post-MVP `v0.1.2`

Implemented:

- async Gateway integration
- `POST /api/ask`
- `GET /api/ask/:request_id`
- `POST /api/gateway/callback`
- callback verification via allowlist and/or HMAC

Still pending:

- real integration smoke validation in a fully configured environment

### Post-MVP `v0.1.3`

Delivered:

- dev-only `/ui-kit`
- Arcana background system
- UI Kit-first governance
- mode switcher alignment
- mobile UI harmonization for the design reference page

## Current Product Phase

### Planned `v0.1.4`

Approved scope in documentation:

- close real validation of `v0.1.2`
- add a dev index page at `/dev`
- improve smoke-data cleanup and reproducibility

Important note:

- these items are planned, but they are not fully implemented in the current codebase yet

### Explicitly Deferred Work

- all automated test work
- visual regression automation
- external file storage migration
