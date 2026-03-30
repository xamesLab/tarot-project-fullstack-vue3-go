# Environment and Local Runbook

Last updated: March 30, 2026.

## Required Tools

- Go
- Node.js
- npm
- ssh
- PowerShell

## Environment Files

### Backend

- template: `backend/.env.example`
- local file: `backend/.env`

Important variables:

- `APP_PORT`
- `APP_ENV`
- `JWT_SECRET`
- `JWT_TTL`
- `MONGO_URI`
- `MONGO_DB_NAME`
- `GATEWAY_BASE_URL`
- `GATEWAY_REQUEST_PATH`
- `GATEWAY_TIMEOUT`
- `BACKEND_PUBLIC_BASE_URL`
- `GATEWAY_CALLBACK_HMAC_SECRET`
- `GATEWAY_CALLBACK_ALLOWLIST`

### Frontend

- template: `frontend/.env.example`
- local file: `frontend/.env`

Important variables:

- `VITE_API_BASE_URL`
- `VITE_APP_NAME`

## Local Run Sequence

### 1. Prepare SSH tunnel to MongoDB

```powershell
ssh -L 27017:d51673dc67a4.vps.myjino.ru:27017 -p 49401 root@d51673dc67a4.vps.myjino.ru
```

### 2. Prepare backend database and seed data

From `backend/`:

```powershell
$env:GO111MODULE='on'
go run ./cmd/indexes --mongo-uri "mongodb://tarot_daily_app:<password>@127.0.0.1:27017/tarot_daily?authSource=tarot_daily" --mongo-db tarot_daily --stats
go run ./cmd/seed --env-file .env --apply --admin-login admin --admin-password <admin-password>
```

### 3. Start backend

From `backend/`:

```powershell
$env:GO111MODULE='on'
go run ./cmd/api
```

Alternative shortcuts:

- `backend/dev.ps1`
- `backend/dev.sh`

Default backend URL:

- `http://127.0.0.1:8080`

### 4. Start frontend

From `frontend/`:

```powershell
npm install
npm run dev
```

Default frontend URL:

- `http://127.0.0.1:5173`

## Basic Validation Commands

### Backend tests

```powershell
cd backend
$env:GO111MODULE='on'
go test ./...
```

### Frontend build

```powershell
cd frontend
npm run build
```

### Backend smoke

```powershell
powershell -ExecutionPolicy Bypass -File deploy/scripts/smoke-backend-api.ps1 -ApiBaseUrl http://127.0.0.1:8080 -AdminLogin admin -AdminPassword <admin-password>
```

### `v0.1.1` frontend and backend smoke

```powershell
powershell -ExecutionPolicy Bypass -File deploy/scripts/smoke-v0.1.1-api.ps1 -AdminLogin admin -AdminPassword <admin-password>
```

### `v0.1.2` ask-flow smoke

```powershell
powershell -ExecutionPolicy Bypass -File deploy/scripts/smoke-v0.1.2-ask.ps1 -ApiBaseUrl http://127.0.0.1:8080 -CallbackHmacSecret <secret>
```

Notes:

- backend must already be running
- Gateway must be reachable
- callback verification must be configured

## Important Current Limitation

The planned smoke cleanup tooling is not present in the repository yet. If smoke data must be removed today, it requires manual cleanup or new implementation work.
