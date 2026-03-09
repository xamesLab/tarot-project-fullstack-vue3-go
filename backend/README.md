# Backend

Папка для серверного API MVP и этапа `v0.1.1`.

Планируемый стек:
- Go
- Gin
- MongoDB

План работ:
- `specs/back/implementation-plan-backend.md`

## Environment
- Шаблон: `backend/.env.example`
- Локальный файл: `backend/.env`

Переменные:
- `APP_PORT` - порт backend API.
- `APP_ENV` - режим приложения (`dev` для MVP).
- `JWT_SECRET` - секрет подписи JWT.
- `JWT_TTL` - время жизни JWT (например, `24h`).
- `MONGO_HOST` - хост MongoDB для dev-туннеля.
- `MONGO_PORT` - порт MongoDB для dev-туннеля.
- `MONGO_DB_NAME` - имя базы данных.
- `MONGO_APP_USER` - пользователь приложения MongoDB.
- `MONGO_AUTH_SOURCE` - authSource для пользователя приложения.
- `MONGO_URI` - строка подключения к MongoDB (пароль указывать только локально в `.env`).
- `HTTP_READ_HEADER_TIMEOUT`, `HTTP_READ_TIMEOUT`, `HTTP_WRITE_TIMEOUT`, `HTTP_IDLE_TIMEOUT` - HTTP таймауты.
- `SHUTDOWN_TIMEOUT` - таймаут graceful shutdown.
- `MONGO_CONNECT_TIMEOUT`, `MONGO_PING_TIMEOUT` - таймауты подключения и ping MongoDB.

## Run API (Backend Skeleton)
Запуск:
- `$env:GO111MODULE='on'; go run ./cmd/api`

Быстрый запуск одной командой:
- PowerShell (из папки `backend`): `./dev.ps1`
- Bash (из папки `backend`): `./dev.sh`

Запуск bash:
`export GO111MODULE=on`
`export MONGO_URI='mongodb://tarot_daily_app:331707tD@127.0.0.1:27017/tarot_daily?authSource=tarot_daily'`
`"/c/Program Files/Go/bin/go.exe" run ./cmd/api`

Health endpoints:
- `GET /api/health/live`
- `GET /api/health/ready`

Static assets:
- `GET /uploads/*` - загруженные в админке изображения колод.
- `GET /assets/cards/*` - локальный каталог изображений `specs/temp/card-image`.

Auth endpoints:
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/me` (Bearer token required)
- `GET /api/me/profile` (Bearer token required)
  - alias для `/api/me` в контексте `v0.1.1`.
- `PATCH /api/me/settings` (Bearer token required, MVP-совместимость)
  - payload: `displayName` (опционально), `timezone` (IANA или `UTC±HH:MM`), `preferredDeckId` (только `done` deck), `deckThemeId` (legacy).
- `PATCH /api/me/profile` (Bearer token required)
  - alias для `/api/me/settings` в контексте `v0.1.1`.
- `GET /api/decks/available` (Bearer token required)
  - список колод со статусом `done` в формате `{ "id": "...", "name": "..." }`.
- `GET /api/deck-themes` (Bearer token required)
  - legacy endpoint тем колод MVP (`deck_themes`).
- `GET /api/reading/today` (Bearer token required)
  - возвращает фиксированную карту дня (одна на пользователя в сутки).
  - суточная граница фиксирована на `03:00` по таймзоне пользователя.
  - при первом открытии за сутки также формируется первый комплект раскладов.
  - в карточках возвращаются `imageUrl` (лицевая сторона) и `backImageUrl` (рубашка).
- `POST /api/reading/generate` (Bearer token required)
  - генерирует новый комплект `spreadA` + `spreadB` в рамках дневного лимита роли.
  - `user` по умолчанию: `10` генераций/сутки.
  - `admin`: без лимита.

Role-based check:
- `GET /api/admin/access-check` (Bearer + role `admin`)
- `GET /api/admin/deck-themes` (Bearer + role `admin`)
  - список тем колод.
- `POST /api/admin/deck-themes` (Bearer + role `admin`)
  - создание темы колоды.
  - payload: `{ "name": "Aurora Theme" }`
- `GET /api/admin/decks` (Bearer + role `admin`)
  - список колод (`process|done`) с информацией о создателе и количестве загруженных карт.
- `POST /api/admin/decks` (Bearer + role `admin`)
  - создание колоды.
  - payload: `{ "name": "Classic Tarot" }`
- `GET /api/admin/decks/{deckId}` (Bearer + role `admin`)
  - детали колоды, список загруженных карт и `backImageUrl`.
- `POST /api/admin/decks/{deckId}/assets/back` (Bearer + role `admin`)
  - загрузка рубашки.
  - multipart form-data: `file=<image>`
- `POST /api/admin/decks/{deckId}/assets/cards/{cardKey}` (Bearer + role `admin`)
  - загрузка изображения карты по ключу.
  - multipart form-data: `file=<image>`
  - статус колоды пересчитывается автоматически после каждой загрузки.
  - поддерживаемые форматы: `jpeg`, `png`, `webp`.
  - лимит размера файла: `8MB`.
  - `cardKey` валидируется по фиксированному каталогу 78 карт (`22..77`, `100..121`).
  - повторная загрузка уже существующего `cardKey` запрещена (уникальность в рамках колоды).
  - загрузка сверх 78 карт блокируется.
- `GET /api/admin/roles` (Bearer + role `admin`)
  - список ролей и лимитов генерации.
- `POST /api/admin/roles` (Bearer + role `admin`)
  - upsert лимита по роли.
  - payload: `{ "role": "vip", "dailySpreadLimit": 25 }` (`-1` = без лимита).
- `PATCH /api/admin/users/role` (Bearer + role `admin`)
  - назначение роли пользователю.
  - payload: `{ "login": "alice", "role": "vip" }`
- `GET /api/admin/users` (Bearer + role `admin`)
  - список пользователей для админки.
  - query:
    - `page` (по умолчанию `1`)
    - `pageSize` (по умолчанию `20`)
    - `role` (опционально, фильтр)
    - `status` (`active|banned|freeze`)
    - `sort` (`registeredAt:desc` или `registeredAt:asc`)
    - `name` (точное регистронезависимое совпадение по `displayName`, fallback на `login`)

## Seed Command (Go)
Команда: `backend/cmd/seed`

Безопасный режим (по умолчанию, без записи в БД):
- `go run ./cmd/seed --env-file .env`

Dry-run + экспорт payload:
- `go run ./cmd/seed --env-file .env --export-json ../deploy/mongo/seed-payload.json`

Применение в MongoDB:
- `go run ./cmd/seed --env-file .env --apply --admin-login admin --admin-password <strong-password>`

Подключение ssh тунеля для dev:
- `ssh -L 27017:d51673dc67a4.vps.myjino.ru:27017 -p 49401 root@d51673dc67a4.vps.myjino.ru`

Если в системе включен GOPATH-режим:
- `$env:GO111MODULE='on'; go run ./cmd/seed --env-file .env`

Что делает seed:
- готовит 78 карт (`cards`);
- upsert default темы колоды (`deck_themes`);
- upsert admin пользователя (`users`) с bcrypt-хэшем пароля.

## Migration v0.1.1 (Go)
Команда: `backend/cmd/migrate011`

Назначение:
- синхронизирует каталог `cards` с `specs/temp/mock/tarot-cards.json`;
- удаляет legacy-карты вне диапазона `22..77`, `100..121`;
- очищает `daily_readings` (чтобы убрать ссылки на старые `cardId`);
- обновляет assets default deck theme на `/assets/cards/{id}.jpg`.

Запуск:
- dry-run:
  - `go run ./cmd/migrate011 --env-file .env`
- apply:
  - `go run ./cmd/migrate011 --env-file .env --apply`

## Indexes Command (Go)
Команда: `backend/cmd/indexes`

Применение индексов:
- `go run ./cmd/indexes --mongo-uri "mongodb://<user>:<password>@127.0.0.1:27017/tarot_daily?authSource=tarot_daily" --mongo-db tarot_daily --stats`

Что делает:
- создает/подтверждает индексы:
  - `users.login` unique,
  - `daily_readings(userId,date)` unique,
  - `cards.code` unique,
  - `deck_themes.name` unique,
  - `role_policies.role` unique;
- выводит статистику документов по ключевым коллекциям.
