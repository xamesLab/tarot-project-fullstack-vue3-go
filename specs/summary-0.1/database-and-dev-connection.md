# Database and Development Connection Guide

Last updated: March 30, 2026.

## Database Model

The project uses MongoDB.

Important collections referenced in code and docs:

- `users`
- `cards`
- `daily_readings`
- `deck_themes`
- `decks`
- `role_policies`
- `ask_requests`
- `admin_settings`

## Local Development Connection Model

The repository documentation assumes MongoDB is reachable through an SSH tunnel from the developer machine.

## Standard SSH Tunnel Example

```powershell
ssh -L 27017:d51673dc67a4.vps.myjino.ru:27017 -p 49401 root@d51673dc67a4.vps.myjino.ru
```

This exposes MongoDB locally on:

- host: `127.0.0.1`
- port: `27017`

## Typical Application Mongo URI

```text
mongodb://tarot_daily_app:<set-password>@127.0.0.1:27017/tarot_daily?authSource=tarot_daily
```

## Backend Variables Used for MongoDB

- `MONGO_URI`
- `MONGO_DB_NAME`
- `MONGO_HOST`
- `MONGO_PORT`
- `MONGO_APP_USER`
- `MONGO_AUTH_SOURCE`
- `MONGO_CONNECT_TIMEOUT`
- `MONGO_PING_TIMEOUT`

## Database Bootstrap and Validation Commands

Run from `backend/`.

### 1. Optional connectivity check

```powershell
$env:GO111MODULE='on'
go run ./cmd/mongocheck -uris "mongodb://tarot_daily_app:<password>@127.0.0.1:27017/tarot_daily?authSource=tarot_daily"
```

### 2. Create or rotate the MongoDB app user

```powershell
$env:GO111MODULE='on'
go run ./cmd/dbinit -app-pass <password>
```

### 3. Ensure indexes exist

```powershell
$env:GO111MODULE='on'
go run ./cmd/indexes --mongo-uri "mongodb://tarot_daily_app:<password>@127.0.0.1:27017/tarot_daily?authSource=tarot_daily" --mongo-db tarot_daily --stats
```

### 4. Seed application data

```powershell
$env:GO111MODULE='on'
go run ./cmd/seed --env-file .env --apply --admin-login admin --admin-password <admin-password>
```

## What These Commands Do

### `cmd/mongocheck`

Tests connectivity to one or more Mongo URIs.

### `cmd/dbinit`

Ensures the application database exists and creates or updates the MongoDB app user.

### `cmd/indexes`

Ensures major indexes exist, including:

- unique user login
- unique daily reading per user/date
- ask-request indexes
- deck and role-policy indexes

### `cmd/seed`

Seeds:

- 78 tarot cards
- default deck theme data
- admin account

## Storage Notes

- uploaded deck assets are currently stored on local backend filesystem storage and served from `/uploads`
- static card assets are available from `/assets/cards`
- external asset storage migration is not active right now and is explicitly deferred
