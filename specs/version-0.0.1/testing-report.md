# Testing Report (Unit + Smoke)

Дата фиксации: 9 марта 2026.

## 1. Unit
- Добавлены unit-тесты timezone-логики:
  - `backend/internal/platform/timeutil/timezone_test.go`
- Покрыто:
  - расчет бизнес-даты с reset `03:00`,
  - нормализация IANA/UTC offset,
  - обработка невалидной таймзоны.

## 2. Build checks
- Backend:
  - `go test ./...` проходит.
- Frontend:
  - `npm run build` проходит (`vue-tsc --noEmit` + `vite build`).

## 3. Smoke
- Добавлен повторяемый smoke-скрипт:
  - `deploy/scripts/smoke-backend-api.ps1`
- Сценарий проверяет:
  - register/login/me,
  - settings update,
  - stable card-of-day,
  - spread generation,
  - admin access + roles API.
- Фактический запуск smoke:
  - результат: `Smoke passed`.
