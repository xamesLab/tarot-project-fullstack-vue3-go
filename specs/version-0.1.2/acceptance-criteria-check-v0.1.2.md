# Acceptance Criteria Check v0.1.2

Дата проверки: 20 марта 2026.
Источник критериев: `specs/version-0.1.2/requirements-v0.1.2.md` (раздел 4).

## Статус по критериям
1. `POST /api/ask` принимает валидный `text` и возвращает `request_id` + `accepted`: `PASS`
2. Backend отправляет задачу в Gateway в контрактном формате (`text`, `callback_url`, `request_id`, `metadata`): `PASS`
3. `POST /api/gateway/callback` принимает `success|error` и идемпотентно обновляет статус по `request_id`: `PASS`
4. `GET /api/ask/:request_id` отдает корректный статус/результат/ошибку пользователю: `PASS`
5. Callback защищен (allowlist и/или HMAC), невалидные callback отклоняются: `PASS`
6. Интеграционный smoke ask-flow проходит, `go test ./...` проходит: `PARTIAL`

## Доказательная база
- Реализация endpoint и интеграции:
  - `backend/internal/transport/http/routes.go`
  - `backend/internal/transport/http/bootstrap.go`
  - `backend/internal/transport/http/handlers/ask_handler.go`
  - `backend/internal/ask/service.go`
  - `backend/internal/gatewayclient/client.go`
  - `backend/internal/transport/http/middleware/gateway_callback_verify.go`
- Unit-тесты:
  - `backend/internal/ask/service_test.go`
  - `backend/internal/gatewayclient/client_test.go`
  - `backend/internal/transport/http/middleware/gateway_callback_verify_test.go`
- Команда:
  - `go test ./...` (из `backend/`) — `PASS`
- Smoke-скрипт:
  - `deploy/scripts/smoke-v0.1.2-ask.ps1` — добавлен
  - фактический запуск в текущей сессии: `NOT RUN`

## Вывод
- Этап `v0.1.2`: `IN PROGRESS`
