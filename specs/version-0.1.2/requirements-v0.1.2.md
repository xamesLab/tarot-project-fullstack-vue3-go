# Требования этапа v0.1.2 (Gateway Integration)

Дата фиксации: 20 марта 2026.

## 1. Цель этапа
Доработать backend `Tarot Daily` для асинхронной интеграции `Frontend -> Backend -> Gateway -> Agent` и формализовать контракт взаимодействия frontend с новым backend flow.

Источник интеграционного контракта:
- `specs/version-0.1.2/INTEGRATION_INSTRUCTION.ru.md`

Базовый URL Gateway для этапа `v0.1.2`:
- `https://api.xmlab.ru`

## 2. Функциональный scope

### 2.1 Frontend -> Backend API
- Добавить endpoint `POST /api/ask`:
  - вход: `text`;
  - валидация входа;
  - генерация/нормализация `request_id`;
  - отправка задания в Gateway;
  - ответ frontend в формате `accepted`.
- Добавить endpoint `GET /api/ask/:request_id`:
  - возврат текущего статуса задачи (`accepted|queued|processing|success|error`);
  - возврат результата/ошибки после callback.

### 2.2 Backend -> Gateway integration
- Реализовать HTTP-клиент вызова `POST /request` в Gateway.
- Для `v0.1.2` использовать `GATEWAY_BASE_URL=https://api.xmlab.ru`.
- Передавать обязательные поля: `text`, `callback_url`, `request_id`, `metadata.user_id`.
- Обрабатывать ответы/ошибки Gateway с маппингом на backend-ошибки API.

### 2.3 Gateway -> Backend callback
- Добавить endpoint `POST /api/gateway/callback`.
- Поддержать callback статусы `success` и `error`.
- Реализовать идемпотентное обновление результата по `request_id`.
- Ввести проверку подлинности callback (allowlist и/или HMAC подпись).

### 2.4 Хранение статусов и результатов запросов
- Добавить persistence слой для асинхронных запросов (коллекция/хранилище `ask_requests`).
- Хранить минимум:
  - `request_id`;
  - `user_id`;
  - `text`;
  - `status`;
  - `result_text`;
  - `error_code`, `error_message`;
  - `created_at`, `updated_at`, `completed_at`.
- Добавить необходимые индексы:
  - unique по `request_id`;
  - индекс для пользовательского доступа (`user_id`, `created_at`).

### 2.5 Контракт для frontend
- Зафиксировать polling flow для frontend:
  - `POST /api/ask` -> получение `request_id`;
  - периодический `GET /api/ask/:request_id` до финального статуса.
- Обновить backend-документацию endpoint’ов и примеров payload/response.

## 3. Нефункциональные требования
- Не ломать существующие API-контракты `v0.1.1` (регрессии недопустимы).
- `POST /api/ask` должен быть быстрым и не ждать ответа Agent.
- Callback-обработчик должен быть идемпотентным.
- Ошибки Gateway/Agent не должны приводить к падению backend процесса.
- Новые конфигурации и команды должны быть документированы в `backend/README.md`.

## 4. Критерии приемки
- `POST /api/ask` возвращает `request_id` и статус `accepted` для валидного запроса.
- Backend отправляет задачу в Gateway в контрактном формате.
- По каждому `request_id` backend принимает callback и сохраняет финальный статус (`success|error`).
- `GET /api/ask/:request_id` отдает корректный статус и данные результата/ошибки.
- Callback endpoint защищен (allowlist и/или HMAC) и обрабатывает повторные callback без дублирования данных.
- Базовые backend-проверки проходят (`go test ./...`), интеграционный smoke по ask-flow проходит.

## 5. Ограничения этапа
- На этапе `v0.1.2` не внедряется WebSocket/SSE push (только polling для frontend).
- Этап не включает разработку Gateway и Agent, только интеграцию backend с их контрактами.
