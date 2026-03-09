# Backend Bootstrap (Step 8)

Дата фиксации: 8 марта 2026.

## 1. Цель шага
Реализовать базовый каркас backend:
- config loader,
- router,
- middleware,
- error handling,
- healthcheck,
- graceful shutdown.

## 2. Реализованные компоненты

### Entry point
- `backend/cmd/api/main.go`
- Логика:
  - загрузка конфигурации из `.env`,
  - подключение к MongoDB + startup ping,
  - запуск HTTP-сервера,
  - корректное завершение по `SIGINT/SIGTERM`,
  - graceful shutdown HTTP + Mongo disconnect.

### Config
- `backend/internal/config/config.go`
- Что покрыто:
  - чтение `.env`,
  - fallback по значениям,
  - валидация обязательных переменных,
  - duration-поля для таймаутов.

### Mongo platform
- `backend/internal/platform/mongo/client.go`
- Что покрыто:
  - connect,
  - startup ping,
  - database handle,
  - disconnect.

### HTTP transport
- `backend/internal/transport/http/router.go`
- Middleware:
  - request id (`X-Request-ID`),
  - panic recovery в JSON-формате,
  - единый error handler.
- Handlers:
  - `GET /` базовый статус сервиса,
  - `GET /api/health/live`,
  - `GET /api/health/ready`.

### Response
- `backend/internal/transport/http/response/json.go`
- Единый JSON-формат ошибки:
  - `error`,
  - `message`,
  - `requestId` (если доступен).

## 3. Проверка шага
- `go mod tidy` выполнен.
- `go test ./...` выполнен успешно (сборка всех пакетов прошла).

## 4. Ограничения на момент шага
- Реализован только инфраструктурный каркас, без бизнес-endpoint’ов auth/readings/admin.
