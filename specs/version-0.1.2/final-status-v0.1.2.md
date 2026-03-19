# Final Status v0.1.2

Дата фиксации старта: 20 марта 2026.

## Итог
- Статус этапа: `v0.1.2 implementation completed, validation in progress`.
- Шаги чек-листа: `18/18` (`specs/version-0.1.2/steps.md`).

## Что реализовано
- Backend интеграция с Gateway:
  - `gatewayclient` (`POST /request`, timeout/error handling);
  - сервис оркестрации `ask` (create/get/complete success/complete error).
- API контур для frontend:
  - `POST /api/ask`;
  - `GET /api/ask/:request_id`.
- Callback контур:
  - `POST /api/gateway/callback`;
  - idempotency на уровне `ask` сервиса;
  - callback verification middleware (allowlist/hmac).
- Тестовое покрытие `v0.1.2`:
  - unit-тесты `ask`, `gatewayclient`, callback verification middleware.
- Smoke артефакты:
  - добавлен скрипт `deploy/scripts/smoke-v0.1.2-ask.ps1` (success + error flow).

## Текущий статус валидации
- `go test ./...` (backend): `PASS`.
- Интеграционный smoke `v0.1.2` в текущей сессии: `NOT RUN` (требуется поднятый backend + доступный Gateway + callback security конфигурация).

## Целевой scope этапа
- Интеграция backend с Gateway по async-контракту;
- API для frontend: `POST /api/ask`, `GET /api/ask/:request_id`;
- callback endpoint: `POST /api/gateway/callback` с idempotency;
- безопасность callback (allowlist/hmac);
- интеграционный smoke ask-flow и обновленная backend-документация.

## Сопутствующие документы
- Требования: `specs/version-0.1.2/requirements-v0.1.2.md`
- Архитектура: `specs/version-0.1.2/architecture-v0.1.2.md`
- Отчет проверок: `specs/version-0.1.2/testing-report-v0.1.2.md`
- Критерии приемки: `specs/version-0.1.2/acceptance-criteria-check-v0.1.2.md`
