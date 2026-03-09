# Backend Deck Themes List (Step 12)

Дата фиксации: 8 марта 2026.

## 1. Цель шага
Реализовать endpoint списка доступных тем колод для пользовательских настроек.

## 2. Реализованные компоненты
- `backend/internal/deckthemes/service.go`
  - сервис получения списка тем из MongoDB.
  - формирует API-вид с полями `id`, `name`.
- `backend/internal/transport/http/handlers/deck_themes_handler.go`
  - HTTP handler `GET /api/deck-themes`.
- `backend/internal/transport/http/router.go`
  - подключен маршрут `GET /api/deck-themes` в auth-зоне (`AuthRequired`).

## 3. API
- `GET /api/deck-themes`
  - Требует Bearer JWT.
  - Ответ: массив тем для выбора в UI.

Пример ответа:
```json
[
  { "id": "67cc3751a6683d5a3a5e95b7", "name": "Default Classic" }
]
```

## 4. Статус проверки
- Сборка backend проходит (`go test ./...`).
