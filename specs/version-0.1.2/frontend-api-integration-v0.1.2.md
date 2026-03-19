# Frontend API Integration Guide v0.1.2

Дата: 20 марта 2026.
Цель: дать frontend-разработчику готовую инструкцию по подключению async flow `ask` в версии `v0.1.2`.

## 1. Что реализовано в API

В `v0.1.2` для фронтенда добавлен async-контур:
- `POST /api/ask` — создать задачу;
- `GET /api/ask/:request_id` — получить текущий статус/результат.
- `GET /api/gateway/health` — проверить доступность Gateway через backend.

Backend сам общается с Gateway (`https://api.xmlab.ru`) и принимает callback.
Frontend с Gateway напрямую не взаимодействует.

## 2. Базовые правила подключения

1. Все пользовательские запросы выполняются с `Authorization: Bearer <token>`.
2. Формат ошибок единый:
```json
{
  "error": "error_code",
  "message": "human readable text",
  "requestId": "optional-request-id"
}
```
3. Callback endpoint (`/api/gateway/callback`) — внутренний для Gateway, не используется фронтом.

## 3. Контракты endpoint-ов

## 3.1 POST /api/ask

Назначение:
- создать async-задачу;
- отправить ее в Gateway;
- вернуть `request_id` для дальнейшего polling.

Auth:
- требуется Bearer token.

Request:
```json
{
  "text": "Какая сегодня рекомендация?"
}
```

Response (`202 Accepted`):
```json
{
  "request_id": "req-67f7e1a9c0d2...",
  "status": "accepted|queued|processing|success|error",
  "result_text": "optional",
  "error_code": "optional",
  "error_message": "optional"
}
```

Примечание:
- обычно на старте приходит `accepted` или `queued`;
- в редких случаях может вернуться уже `error`, если dispatch в Gateway не удался.

Коды ошибок:
- `400 invalid_input` — пустой/некорректный `text`;
- `401 unauthorized` — нет/битый токен;
- `502 gateway_timeout|gateway_rejected|gateway_bad_response|gateway_unavailable`;
- `500 invalid_config` — backend не настроен для ask-интеграции.

## 3.2 GET /api/ask/:request_id

Назначение:
- получить текущее состояние задачи.

Auth:
- требуется Bearer token.

Path param:
- `request_id` из ответа `POST /api/ask`.

Response (`200 OK`):
```json
{
  "request_id": "req-67f7e1a9c0d2...",
  "status": "accepted|queued|processing|success|error",
  "result_text": "текст результата (для success)",
  "error_code": "код ошибки (для error)",
  "error_message": "сообщение ошибки (для error)"
}
```

Владение запросом:
- пользователь может читать только свои `request_id`;
- чужой или несуществующий `request_id` вернет `404 not_found`.

## 3.3 GET /api/gateway/health

Назначение:
- диагностика доступности Gateway из frontend без прямого вызова Gateway.

Auth:
- требуется Bearer token.

Response (`200 OK`, пример):
```json
{
  "status": "ok",
  "uptime": 81720.198966473
}
```

Что делает backend:
- вызывает `GET <GATEWAY_BASE_URL>/health` (для `v0.1.2`: `https://api.xmlab.ru/health`);
- проксирует валидный health response от Gateway.

Коды ошибок:
- `401 unauthorized` — нет/битый токен;
- `504 gateway_timeout` — timeout до Gateway;
- `502 gateway_bad_response|gateway_unavailable` — Gateway недоступен или вернул невалидный ответ.

## 4. Жизненный цикл статусов

Типичный переход:
1. `accepted`
2. `queued`
3. `processing` (опционально)
4. `success` или `error` (финальные статусы)

Финальные:
- `success`
- `error`

Только после финального статуса polling нужно останавливать.

## 5. Рекомендуемый polling для фронта

Рекомендуемый алгоритм:
1. Вызвать `POST /api/ask`.
2. Сохранить `request_id` в состоянии экрана.
3. Запустить polling `GET /api/ask/:request_id`:
   - интервал: 1.0–1.5 секунды;
   - max duration: 60–90 секунд.
4. Остановить polling при:
   - `status=success` (показать `result_text`);
   - `status=error` (показать `error_message`);
   - timeout на фронте (показать "Ответ задерживается, попробуйте позже").

Рекомендованные UI-состояния:
- `idle` — до отправки;
- `submitting` — пока `POST /api/ask`;
- `waiting` — polling активен;
- `success` — вывод результата;
- `failed` — вывод ошибки;
- `timeout` — локальный таймаут фронта.

## 6. Готовые TypeScript типы

```ts
export type AskStatus = "accepted" | "queued" | "processing" | "success" | "error";

export interface AskResponse {
  request_id: string;
  status: AskStatus;
  result_text?: string;
  error_code?: string;
  error_message?: string;
}

export interface ApiErrorBody {
  error: string;
  message: string;
  requestId?: string;
}
```

## 7. Пример API-клиента (fetch)

```ts
const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "http://127.0.0.1:8080";

function authHeaders(token: string): HeadersInit {
  return {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };
}

async function parseOrThrow(res: Response) {
  if (!res.ok) {
    const err = (await res.json()) as { error: string; message: string; requestId?: string };
    throw new Error(`${err.error}: ${err.message}`);
  }
  return res.json();
}

export async function createAsk(token: string, text: string) {
  const res = await fetch(`${API_BASE}/api/ask`, {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify({ text }),
  });
  return (await parseOrThrow(res)) as {
    request_id: string;
    status: "accepted" | "queued" | "processing" | "success" | "error";
    result_text?: string;
    error_code?: string;
    error_message?: string;
  };
}

export async function getAsk(token: string, requestId: string) {
  const res = await fetch(`${API_BASE}/api/ask/${encodeURIComponent(requestId)}`, {
    method: "GET",
    headers: authHeaders(token),
  });
  return (await parseOrThrow(res)) as {
    request_id: string;
    status: "accepted" | "queued" | "processing" | "success" | "error";
    result_text?: string;
    error_code?: string;
    error_message?: string;
  };
}
```

## 8. Пример polling-утилиты

```ts
export async function waitAskFinal(
  token: string,
  requestId: string,
  opts?: { intervalMs?: number; timeoutMs?: number }
) {
  const intervalMs = opts?.intervalMs ?? 1200;
  const timeoutMs = opts?.timeoutMs ?? 90000;
  const startedAt = Date.now();

  while (true) {
    const state = await getAsk(token, requestId);
    if (state.status === "success" || state.status === "error") {
      return state;
    }
    if (Date.now() - startedAt > timeoutMs) {
      throw new Error("ask_poll_timeout");
    }
    await new Promise((r) => setTimeout(r, intervalMs));
  }
}
```

## 9. Минимальный UX-скрипт на экране чата/запроса

1. Пользователь вводит текст.
2. Кнопка "Отправить":
   - дизейбл при пустом тексте;
   - вызывает `createAsk`.
3. Если `createAsk` вернул финальный статус:
   - сразу отрисовать финал.
4. Иначе:
   - запустить `waitAskFinal`.
5. На `success`:
   - отрисовать `result_text`.
6. На `error`:
   - отрисовать `error_message` + кнопку "Повторить".

## 10. Частые причины ошибок и что делать

1. `401 unauthorized`
- токен отсутствует/просрочен;
- проверить refresh/login flow на фронте.

2. `400 invalid_input`
- пустой или некорректный `text`;
- добавить фронтовую валидацию trim + min length.

3. `502 gateway_*`
- backend не смог корректно отправить/получить ответ от Gateway;
- показывать пользователю retry action.

4. `500 invalid_config`
- backend не настроен (`BACKEND_PUBLIC_BASE_URL`, callback security);
- эскалировать на backend/devops.

## 11. Чек-лист интеграции для frontend

- Реализован API client для `POST /api/ask` и `GET /api/ask/:request_id`.
- Реализован polling с остановкой по `success|error`.
- Покрыты UI-состояния: `submitting`, `waiting`, `success`, `failed`, `timeout`.
- Единообразно отображаются backend error codes/messages.
- На странице есть action "повторить запрос" после `error/timeout`.

## 12. Важно про callback

`POST /api/gateway/callback` не вызывается из frontend и не должен попадать в frontend SDK.
Это внутренний интеграционный endpoint backend <-> gateway.
