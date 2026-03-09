# Deploy (Dev)

Папка для локальной инфраструктуры разработки.

Содержит:
- `docker-compose.yml` для локального запуска сервисов;
- вспомогательные файлы инициализации MongoDB;
- инструкции по dev-запуску.

## Проверка окружения
- Скрипт preflight: `deploy/scripts/check-dev-env.ps1`
- Запуск:
  - `powershell -ExecutionPolicy Bypass -File deploy/scripts/check-dev-env.ps1`

## Индексы MongoDB
- Скрипт индексов: `deploy/mongo/create-indexes.js`
- PowerShell-обертка применения: `deploy/scripts/apply-mongo-indexes.ps1`
- Запуск:
  - `powershell -ExecutionPolicy Bypass -File deploy/scripts/apply-mongo-indexes.ps1`

## Seed Payload
- Генерируется Go-командой `backend/cmd/seed`.
- Пример файла: `deploy/mongo/seed-payload.json`.

## Smoke API (Backend)
- Скрипт smoke: `deploy/scripts/smoke-backend-api.ps1`
- Назначение:
  - проверка критичного user/admin flow через HTTP API.
- Перед запуском:
  - backend API должен быть запущен (например, `http://127.0.0.1:8080`);
  - admin пользователь должен существовать.
- Запуск:
  - `powershell -ExecutionPolicy Bypass -File deploy/scripts/smoke-backend-api.ps1`
  - с параметрами:
    - `-ApiBaseUrl http://127.0.0.1:8080`
    - `-AdminLogin admin`
    - `-AdminPassword <password>`

## Smoke API (v0.1.1)
- Скрипт smoke: `deploy/scripts/smoke-v0.1.1-api.ps1`
- Назначение:
  - проверка profile/admin-decks/admin-users flow этапа `v0.1.1`;
  - проверка ограничений `preferredDeckId` (`process` запрещен, `done` разрешен);
  - проверка `imageUrl/backImageUrl` в reading payload и deck `process->done`.
- Запуск:
  - `powershell -ExecutionPolicy Bypass -File deploy/scripts/smoke-v0.1.1-api.ps1 -AdminLogin admin -AdminPassword <password>`

## Frontend checks (v0.1.1)
- Скрипт: `deploy/scripts/check-frontend-v0.1.1.ps1`
- Назначение:
  - проверка `locale/theme` persistence;
  - проверка profile integration (`/api/me/profile`, `/api/decks/available`);
  - проверка card flip mechanics по коду фронтенда.
- Запуск:
  - `powershell -ExecutionPolicy Bypass -File deploy/scripts/check-frontend-v0.1.1.ps1`
