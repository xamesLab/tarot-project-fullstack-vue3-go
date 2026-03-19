# Архитектура этапа v0.1.2 (Frontend/Backend/Gateway)

Дата фиксации: 20 марта 2026.

## 1. Архитектурная цель
Добавить в backend асинхронный контур обработки пользовательского запроса через Gateway с сохранением обратной совместимости `v0.1.1`.

Контракт интеграции:
- `specs/version-0.1.2/INTEGRATION_INSTRUCTION.ru.md`

## 2. Изменения по слоям

## 2.1 Frontend
- Frontend работает по polling-модели:
  - `POST /api/ask` (создание задачи);
  - `GET /api/ask/:request_id` (опрос статуса до `success|error`).
- Для `v0.1.2` backend гарантирует стабильный контракт полей статуса и результата.

## 2.2 Backend
- Добавляются новые компоненты:
  - `ask` сервис оркестрации (create/poll/complete);
  - `gatewayclient` для `POST /request`;
  - callback verifier (allowlist/hmac);
  - repository для `ask_requests`.
- Добавляются новые endpoint’ы:
  - `POST /api/ask`;
  - `GET /api/ask/:request_id`;
  - `POST /api/gateway/callback`.
- `POST /api/ask` не блокируется ожиданием Agent и завершает запрос после постановки в Gateway.
- Callback выполняет идемпотентный upsert статуса/результата по `request_id`.

## 2.3 Deploy / Ops
- Добавить env-конфигурацию интеграции с Gateway:
  - `GATEWAY_BASE_URL` (для этапа `v0.1.2`: `https://api.xmlab.ru`);
  - `GATEWAY_REQUEST_PATH` (default `/request`);
  - `GATEWAY_TIMEOUT`;
  - `BACKEND_PUBLIC_BASE_URL` (для `callback_url`);
  - `GATEWAY_CALLBACK_HMAC_SECRET` и/или `GATEWAY_CALLBACK_ALLOWLIST`.
- Обновить `backend/.env.example` и `backend/README.md` для нового integration-flow.

## 3. Данные и миграции
- Добавляется коллекция `ask_requests` (или эквивалентная):
  - `requestId` (unique);
  - `userId`;
  - `text`;
  - `status`;
  - `resultText`;
  - `errorCode`, `errorMessage`;
  - `createdAt`, `updatedAt`, `completedAt`.
- Индексы:
  - unique `requestId`;
  - `userId + createdAt`.
- Миграции данных `v0.1.1` не требуются.

## 4. Риски и меры
- Риск: недоступность Gateway.
- Мера: явный timeout клиента + корректный маппинг ошибки в `POST /api/ask` + retry на стороне Gateway (вне backend).

- Риск: spoofed callback.
- Мера: верификация источника callback (allowlist/hmac), отклонение невалидных запросов.

- Риск: повторные callback по одному `request_id`.
- Мера: идемпотентное обновление, без дублирования и без деградации статуса.
