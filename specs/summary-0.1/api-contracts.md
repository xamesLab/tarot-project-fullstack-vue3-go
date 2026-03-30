# API Contracts Summary

Last updated: March 30, 2026.

This file is a practical API overview based on the current backend routes and related docs. It is not a Swagger replacement.

## Auth

### `POST /api/auth/register`

Creates a user and returns a token plus user payload.

Typical payload:

```json
{
  "login": "reader",
  "password": "strong-password"
}
```

### `POST /api/auth/login`

Logs a user in and returns a token plus user payload.

### `GET /api/me`

Bearer token required.

Returns the current user profile.

## User Profile and Settings

### `GET /api/me/profile`

Bearer token required.

Alias-style profile endpoint used by the frontend.

### `PATCH /api/me/profile`

Bearer token required.

### `PATCH /api/me/settings`

Bearer token required.

Typical writable fields:

- `displayName`
- `timezone`
- `preferredDeckId`
- `deckThemeId`

### `GET /api/deck-themes`

Bearer token required.

Returns legacy deck theme options.

### `GET /api/decks/available`

Bearer token required.

Returns user-selectable decks in `done` status.

## Reading

### `GET /api/reading/today`

Bearer token required.

Returns the reading for the current user-local day. If no reading exists yet, backend generates and stores one.

Typical response shape:

```json
{
  "date": "2026-03-29",
  "timezoneUsed": "Europe/Moscow",
  "cardOfDay": {},
  "spreadA": [],
  "spreadB": [],
  "spreadGenerationsUsed": 2,
  "spreadGenerationsLimit": 10,
  "spreadGenerationsRemaining": 8,
  "canGenerateMore": true
}
```

### `POST /api/reading/generate`

Bearer token required.

Generates the combined spread package subject to role-based limits.

### `POST /api/reading/generate/card-of-day`

Bearer token required.

Generates only the card-of-day portion.

### `POST /api/reading/generate/spread-a`

Bearer token required.

Generates only spread A.

### `POST /api/reading/generate/spread-b`

Bearer token required.

Generates only spread B.

## Gateway Integration (`v0.1.2`)

### `POST /api/ask`

Bearer token required.

Creates an async request routed to Gateway.

Payload:

```json
{
  "text": "What should I focus on today?"
}
```

Typical response:

```json
{
  "request_id": "req-123",
  "status": "accepted"
}
```

### `GET /api/ask/:request_id`

Bearer token required.

Returns the current async request state.

Expected statuses:

- `accepted`
- `queued`
- `processing`
- `success`
- `error`

Possible result fields:

- `result_text`
- `error_code`
- `error_message`

### `POST /api/gateway/callback`

No Bearer token.

Protected by callback verification rules:

- allowlist and/or
- HMAC signature

Supported callback statuses:

- `success`
- `error`

Success payload shape:

```json
{
  "request_id": "req-123",
  "status": "success",
  "text": "Your answer"
}
```

Error payload shape:

```json
{
  "request_id": "req-123",
  "status": "error",
  "error": "agent_timeout",
  "message": "Agent timed out"
}
```

### `GET /api/gateway/health`

Bearer token required.

Backend performs a health check against the configured Gateway.

## Admin

All admin endpoints require:

- Bearer token
- user role `admin`

### Access and roles

- `GET /api/admin/access-check`
- `GET /api/admin/roles`
- `POST /api/admin/roles`
- `PATCH /api/admin/users/role`
- `GET /api/admin/users`

### Deck themes

- `GET /api/admin/deck-themes`
- `POST /api/admin/deck-themes`

### Decks and uploads

- `GET /api/admin/decks`
- `POST /api/admin/decks`
- `GET /api/admin/decks/:deckId`
- `POST /api/admin/decks/:deckId/assets/back`
- `POST /api/admin/decks/:deckId/assets/cards/:cardKey`

Upload notes:

- multipart form-data
- formats: `jpeg`, `png`, `webp`
- file size limit: `8MB`
- valid card keys: `22..77`, `100..121`

### General admin settings

- `GET /api/admin/settings`
- `PATCH /api/admin/settings`

Typical settings fields:

- `projectTitle`
- `supportEmail`
- `maintenanceMode`
- `dayResetHour`
- `maxUploadMb`
