# Testing Report v0.1.2

Дата проверки: 20 марта 2026.
Статус: `IN PROGRESS`.

## 1. Backend checks
- Команда: `go test ./...` (из `backend/`)
- Результат: `PASS`
- Примечание:
  - добавлены и проходят unit-тесты:
    - `backend/internal/ask/service_test.go`
    - `backend/internal/gatewayclient/client_test.go`
    - `backend/internal/transport/http/middleware/gateway_callback_verify_test.go`

## 2. Gateway client checks
- Что проверяем:
  - корректная сборка payload для `POST https://api.xmlab.ru/request`;
  - timeout/retry поведение на уровне клиента backend;
  - маппинг ошибок Gateway в backend-ошибки.
- Результат: `PASS` (unit)

## 3. API contract checks (Frontend <-> Backend)
- Endpoint: `POST /api/ask`
  - валидация `text`;
  - ответ `request_id + accepted`.
- Endpoint: `GET /api/ask/:request_id`
  - статусы `accepted|queued|processing|success|error`;
  - корректная выдача результата/ошибки.
- Результат: `PASS` (code + unit on ask service)

## 4. Callback checks (Gateway -> Backend)
- Endpoint: `POST /api/gateway/callback`
  - `success` callback;
  - `error` callback;
  - idempotency при повторном callback;
  - проверка подписи/allowlist.
- Результат: `PASS` (unit)
- Примечание:
  - callback verification вынесен в middleware `GatewayCallbackVerify`;
  - idempotency проверяется тестами `CompleteSuccess/CompleteError`.

## 5. Smoke ask-flow v0.1.2
- Сценарий:
  - создать запрос через `POST /api/ask`;
  - выполнить callback;
  - получить финальный статус через `GET /api/ask/:request_id`.
- Скрипт:
  - `deploy/scripts/smoke-v0.1.2-ask.ps1`
- Результат запуска в текущей сессии: `NOT RUN`
- Причина:
  - для фактического прогона требуется поднятый backend + доступный Gateway + настроенная callback security конфигурация.
