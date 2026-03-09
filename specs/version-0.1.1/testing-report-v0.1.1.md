# Testing Report v0.1.1

Дата проверки: 9 марта 2026.

## 1. Backend checks
- Команда: `go test ./...` (из `backend/`)
- Результат: `PASS`

## 2. Frontend checks
- Команда: `npm run build` (из `frontend/`)
- Результат: `PASS`

## 3. API smoke v0.1.1
- Скрипт: `deploy/scripts/smoke-v0.1.1-api.ps1`
- Запуск:
  - `powershell -ExecutionPolicy Bypass -File deploy/scripts/smoke-v0.1.1-api.ps1 -AdminLogin admin -AdminPassword admin12345`
- Покрытие:
  - user register/login;
  - `GET/PATCH /api/me/profile`;
  - стабильность `GET /api/reading/today`;
  - проверка `imageUrl` и `backImageUrl` в reading payload;
  - admin users exact search и фильтры;
  - admin decks: create/list/get;
  - upload `back` и upload 78 карт;
  - автоматический переход deck `process -> done`;
  - `GET /api/decks/available` возвращает только `done`;
  - запрет назначения `preferredDeckId` для `process` deck;
  - разрешение назначения `preferredDeckId` для `done` deck.
- Результат: `PASS` (`Smoke v0.1.1 passed`).

## 4. Frontend scenario checks (locale/theme + card flip)
- Скрипт: `deploy/scripts/check-frontend-v0.1.1.ps1`
- Запуск:
  - `powershell -ExecutionPolicy Bypass -File deploy/scripts/check-frontend-v0.1.1.ps1`
- Покрытие:
  - сохранение `locale/theme` в localStorage;
  - наличие `ru/en` и `light/dark` переключателей;
  - наличие profile API integration (`/api/me/profile`, `/api/decks/available`);
  - flip-механика карточки (`closed -> click -> reveal front`).
- Результат: `PASS`.

## 5. Regression note
- Дополнительно выполнена коррекция БД для `v0.1.1`:
  - команда: `go run ./cmd/migrate011 --env-file .env --apply`
  - итог: каталог `cards` приведен к 78 картам из `specs/temp/mock/tarot-cards.json`, legacy cards удалены, `daily_readings` очищены.

## 6. Post-refactor verification
- После структурного рефакторинга backend повторно выполнялась базовая проверка сборки и тестов:
  - команда: `go test ./...` (из `backend/`)
  - результат: `PASS`
- Добавлены unit-тесты на вынесенные helper-ы и view-logic:
  - `internal/domain/models/user_helpers_test.go`
  - `internal/platform/inpututil/input_test.go`
  - `internal/config/envfile_test.go`
  - `internal/reading/view_test.go`
  - `internal/auth/account_test.go`
  - `internal/auth/view_test.go`
  - `internal/decks/assets_test.go`
- Детальный журнал структурных изменений: `specs/version-0.1.1/refactor-log.md`
