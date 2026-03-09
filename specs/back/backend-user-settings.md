# Backend User Settings (Step 11)

Дата фиксации: 8 марта 2026.

## 1. Цель шага
Реализовать настройки пользователя:
- `timezone`
- `deckThemeId`

## 2. Реализованные компоненты
- `backend/internal/auth/service.go`
  - добавлен `UpdateSettings`.
  - валидация таймзоны:
    - `UTC±HH:MM` (ограничения до `UTC+14:00`/`UTC-14:00`);
    - IANA (`Europe/Moscow`, `Asia/Tokyo` и т.д.);
    - пустое/`UTC` нормализуется в `UTC+00:00`.
  - проверка существования `deckThemeId` в `deck_themes`.
  - поддержка сброса темы через пустую строку.
- `backend/internal/repository/mongo/user_repository.go`
  - добавлен метод обновления настроек пользователя.
- `backend/internal/repository/mongo/deck_theme_repository.go`
  - доступ к `deck_themes` (поиск по `_id`, список).
- `backend/internal/transport/http/handlers/settings_handler.go`
  - handler для `PATCH /api/me/settings`.
- `backend/internal/transport/http/router.go`
  - подключен `PATCH /api/me/settings` в auth-зоне.

## 3. API
- `PATCH /api/me/settings` (Bearer JWT)

Пример payload:
```json
{
  "timezone": "Europe/Moscow",
  "deckThemeId": "67cc3751a6683d5a3a5e95b7"
}
```

Вариант сброса выбранной темы:
```json
{
  "deckThemeId": ""
}
```

Ответ:
- `200` с обновленным профилем пользователя.
- `400` при невалидных `timezone`/`deckThemeId`.
- `401` если пользователь из токена не найден.

## 4. Статус проверки
- Сборка backend проходит (`go test ./...`).
