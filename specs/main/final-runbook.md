# Final Runbook (Local MVP)

Дата фиксации: 9 марта 2026.

## 1. Предусловия
- Установлены `Go`, `Node.js`, `npm`, `ssh`.
- Заполнены `backend/.env` и `frontend/.env` (если используется).
- Поднят SSH-туннель к MongoDB:
  - `ssh -L 27017:d51673dc67a4.vps.myjino.ru:27017 -p 49401 root@d51673dc67a4.vps.myjino.ru`

## 2. Подготовка БД (backend)
Из `backend/`:
```powershell
$env:GO111MODULE='on'
go run ./cmd/indexes --env-file .env --stats
go run ./cmd/seed --env-file .env --apply --admin-login admin --admin-password <admin-password>
```

## 3. Запуск backend API
Из `backend/`:
```powershell
$env:GO111MODULE='on'
go run ./cmd/api
```
API по умолчанию: `http://127.0.0.1:8080`.

## 4. Запуск frontend
Из `frontend/`:
```powershell
npm install
npm run dev
```
UI по умолчанию: `http://127.0.0.1:5173`.

## 5. Быстрые проверки
- Backend unit/build:
```powershell
cd backend
$env:GO111MODULE='on'
go test ./...
```
- Frontend build/typecheck:
```powershell
cd frontend
npm run build
```
- Backend smoke:
```powershell
powershell -ExecutionPolicy Bypass -File deploy/scripts/smoke-backend-api.ps1 -ApiBaseUrl http://127.0.0.1:8080 -AdminLogin admin -AdminPassword <admin-password>
```

## 6. Базовый user/admin flow
1. Зарегистрировать обычного пользователя.
2. Проверить `Dashboard` (`cardOfDay`, `spreadA`, `spreadB`, лимиты).
3. Проверить `Settings` (timezone + theme).
4. Войти под `admin`.
5. Проверить `Admin`:
   - создать тему,
   - создать/обновить policy роли,
   - назначить роль пользователю.
