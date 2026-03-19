# Формализованная инструкция: цепочка Frontend -> Backend -> Gateway -> Agent

## 1. Цель

Документ задает единый контракт взаимодействия между компонентами:

- Frontend (клиент)
- Backend (бизнес-API)
- Gateway (асинхронный маршрутизатор)
- Agent (AI-обработчик)

Базовый URL Gateway для этапа `v0.1.2`:
- `https://api.xmlab.ru`

## 2. Роли компонентов

1. Frontend:
- отправляет пользовательский запрос в Backend;
- получает статус/результат из Backend.

2. Backend:
- принимает запрос от Frontend;
- формирует и отправляет задачу в Gateway;
- принимает callback от Gateway;
- сохраняет результат и отдает его Frontend.

3. Gateway:
- валидирует входящий запрос;
- асинхронно обрабатывает через очередь и воркеры;
- вызывает Agent;
- доставляет результат в `callback_url`.

4. Agent:
- синхронно обрабатывает текст;
- возвращает результат или ошибку.

## 3. Нормативные требования (MUST/SHOULD)

1. Backend MUST передавать в Gateway валидный `callback_url`, доступный из сети Gateway.
2. Gateway MUST отвечать на `POST /request` немедленно (`202 Accepted`), без ожидания Agent.
3. Gateway MUST вызывать Agent с таймаутом и retry согласно конфигурации.
4. Gateway MUST отправлять итог в callback даже при ошибке Agent (`status: "error"`).
5. Backend MUST проверять подлинность callback-запросов (минимум: allowlist источника или подпись/HMAC).
6. Frontend SHOULD получать результат через polling/WebSocket/SSE, а не ждать долгий синхронный ответ.

## 4. Контракты эндпоинтов

### 4.1 Frontend -> Backend

Endpoint (пример):
- `POST /api/ask`

Request body:
```json
{
  "text": "Какая погода в Москве?"
}
```

Response body (пример):
```json
{
  "request_id": "req-123",
  "status": "accepted"
}
```

### 4.2 Backend -> Gateway

Endpoint:
- `POST https://api.xmlab.ru/request`

Request body:
```json
{
  "text": "Какая погода в Москве?",
  "callback_url": "https://backend.example.com/api/gateway/callback",
  "request_id": "req-123",
  "metadata": {
    "user_id": "u-42"
  }
}
```

Response body:
```json
{
  "request_id": "req-123",
  "status": "queued"
}
```

### 4.3 Gateway -> Agent

Endpoint (из `AGENT_URL`):
- `POST /process`

Request body:
```json
{
  "request_id": "req-123",
  "text": "Какая погода в Москве?",
  "metadata": {
    "user_id": "u-42"
  }
}
```

Success response (пример):
```json
{
  "text": "Сейчас +12, облачно."
}
```

### 4.4 Gateway -> Backend (callback)

Endpoint (из `callback_url`):
- `POST /api/gateway/callback`

Success callback:
```json
{
  "request_id": "req-123",
  "status": "success",
  "text": "Сейчас +12, облачно.",
  "agent_latency_ms": 120,
  "delivered_at": "2026-03-18T19:10:00.000Z"
}
```

Error callback:
```json
{
  "request_id": "req-123",
  "status": "error",
  "error": "agent_timeout",
  "message": "Agent did not respond after 3 attempt(s): timeout of 500ms exceeded",
  "agent_latency_ms": 0,
  "delivered_at": "2026-03-18T19:10:00.000Z"
}
```

## 5. Порядок выполнения (sequence)

1. Frontend -> Backend: `POST /api/ask`.
2. Backend генерирует/принимает `request_id`.
3. Backend -> Gateway: `POST /request`.
4. Gateway -> Backend: `202 {request_id,status:"queued"}`.
5. Gateway (worker) -> Agent: `POST AGENT_URL`.
6. Agent -> Gateway: success/error.
7. Gateway -> Backend: callback с финальным статусом.
8. Backend сохраняет результат и предоставляет его Frontend.

## 6. Обработка ошибок

1. Ошибка валидации во входе Gateway:
- HTTP `400`
- `error: "validation_error"`.

2. Ошибка middleware pipeline:
- HTTP `400`
- `error: "pipeline_error"`.

3. Timeout/ошибка Agent:
- Gateway логирует retry попытки;
- Backend получает callback с `status: "error"`.

4. Недоступный `callback_url`:
- Gateway делает retry (`CALLBACK_RETRY_COUNT`);
- после исчерпания логирует `Callback delivery exhausted`;
- процесс не падает.

## 7. Безопасность

1. Admin API защищается заголовком:
- `X-Admin-Token: <ADMIN_TOKEN>`.

2. Рекомендовано для callback:
- подпись `X-Signature` (HMAC SHA-256) и проверка на Backend;
- idempotency по `request_id` (повторный callback не должен дублировать запись).

3. Продакшн:
- только HTTPS между внешними компонентами;
- секреты только в `.env`/secret manager.

## 8. Наблюдаемость и мониторинг

1. Логировать как минимум:
- `request_id`, `status`, `attempt`, `error`, `agent_latency_ms`.

2. Контрольные метрики:
- `processedTotal`, `errorTotal`, `queueDepth`, uptime.

3. Health-check:
- `GET /health` должен стабильно отвечать `200`.

## 9. Критерии приемки интеграции

1. Backend получает `202` от Gateway для валидных запросов.
2. Для каждого `request_id` приходит callback (success/error).
3. При отказах Agent сервис не падает, retry работает.
4. При dead callback воркеры продолжают обработку других задач.
5. Admin API недоступен без токена и доступен с валидным токеном.

## 10. Рекомендуемая минимальная реализация на стороне Backend

1. Таблица/хранилище запросов: `request_id`, `status`, `result`, `error`, timestamps.
2. Endpoint `POST /api/ask`:
- валидация `text`;
- отправка в Gateway;
- возврат `request_id` клиенту.

3. Endpoint `POST /api/gateway/callback`:
- валидация подписи/источника;
- idempotent update по `request_id`;
- сохранение финального статуса.

4. Endpoint чтения статуса для Frontend:
- `GET /api/ask/:request_id` (или push через WebSocket/SSE).
