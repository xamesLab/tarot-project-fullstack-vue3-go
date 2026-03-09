# Backend Admin Deck Themes API

Дата фиксации: 9 марта 2026.

## 1. Цель
Реализовать admin API для просмотра и добавления тем колод.

## 2. Endpoint'ы
- `GET /api/admin/deck-themes`
  - доступ только `admin`;
  - возвращает массив тем в формате `{ id, name }`.
- `POST /api/admin/deck-themes`
  - доступ только `admin`;
  - payload: `{ "name": "Aurora Theme" }`;
  - создает новую тему с `isPlaceholder=true`, пустым `assets`, `createdBy` из admin JWT.

## 3. Валидация и ошибки
- длина имени темы: `2..128`.
- дубли имени блокируются (`409 theme_exists`).

## 4. Реализованные компоненты
- `backend/internal/deckthemes/service.go`
  - добавлен `Create`.
- `backend/internal/repository/mongo/deck_theme_repository.go`
  - добавлен `Create`.
- `backend/internal/transport/http/handlers/deck_themes_handler.go`
  - добавлен `Create`.
- `backend/internal/transport/http/router.go`
  - добавлены admin routes `/api/admin/deck-themes`.
- `backend/cmd/indexes/main.go`
  - добавлен unique index `deck_themes.name`.

## 5. Проверка
Smoke-проверка подтвердила:
- `admin` может создать тему и видеть ее в списке;
- повторное создание темы с тем же именем дает `409`;
- `user` получает `403` на `GET /api/admin/deck-themes`.
