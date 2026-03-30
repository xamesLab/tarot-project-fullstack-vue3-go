# Current Project Status

Last updated: March 30, 2026.

## Executive Summary

The project is beyond MVP. The repository contains a working backend API, a working Vue frontend, local development scripts, and several post-MVP improvements. The codebase is usable for manual development and validation, but some roadmap items described in planning documents are still not implemented.

## Confirmed Delivered Milestones

### MVP

Delivered and documented:

- backend skeleton
- auth and role checks
- user settings
- daily reading generation
- admin endpoints
- frontend bootstrap and core user flows
- local runbook and acceptance checks

### `v0.1.1`

Delivered and documented as a completed stage.

### `v0.1.2`

Implemented:

- async backend-to-Gateway ask flow
- Gateway callback handling
- callback verification via allowlist and/or HMAC
- Gateway health endpoint
- ask-request persistence and status polling

Validation status:

- implementation is documented as completed
- backend tests exist in the repository
- the real ask-flow smoke run is still pending in a properly configured runtime environment

### `v0.1.3`

Implemented in code:

- dev-only `/ui-kit` page
- `AppBackground.vue`
- UI Kit-based visual governance
- dev-only route registration pattern
- mode switcher alignment and mobile refinements for the design reference page

## What Is Present in Code Today

- frontend dev-only route: `/ui-kit`
- admin pages: users, decks, settings
- backend static asset serving from `/uploads` and `/assets/cards`
- Gateway ask-flow endpoints and callback endpoint
- smoke scripts for backend, `v0.1.1`, and `v0.1.2`

## What Is Planned but Not Yet Implemented

These items appear in `v0.1.4` planning documents, but they are not present in the current repository state:

- `/dev` index page for internal developer tools
- smoke cleanup command or script
- updated `v0.1.2` final validation artifacts based on a successful real smoke run

## Current Project Constraints

- background system must not be changed unless explicitly requested
- UI work must follow `UI Kit First`
- automated tests are postponed until the user explicitly asks for them
- external file storage migration is postponed until the user explicitly asks for it

## Practical Interpretation

The safest way to treat the project today is:

- a working local MVP-plus product
- a documented codebase with a clear visual governance model
- a project whose nearest real work is integration validation and dev hygiene
- not a project where automation or storage migration should be advanced right now
