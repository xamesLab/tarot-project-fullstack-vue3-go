# Архитектура проекта Tarot Daily (MVP)

## 1. Цель архитектуры
Обеспечить простой, расширяемый и поддерживаемый MVP для ежедневных таро-раскладов с:
- авторизацией и ролями (`user`/`admin`);
- генерацией и хранением ежедневного результата по таймзоне пользователя;
- выбором темы колоды пользователем;
- админ-управлением темами колод.

## 2. Архитектурный стиль
- Frontend + Backend + Database.
- Клиент-серверное взаимодействие по REST API.
- Разделение backend по слоям: `router -> handler -> service -> repository`.
- Stateless API (аутентификация через JWT).

## 3. Технологический стек
- Frontend: Vue 3, Pinia, Vue Router.
- Backend: Go, Gin.
- Database: MongoDB.
- Аутентификация: JWT.

## 4. Логическая схема компонентов
1. Frontend (SPA):
- страницы: login/register/dashboard/settings/admin;
- хранение состояния (auth, reading, settings);
- route guards для приватных и admin-маршрутов.

2. Backend API:
- публичные endpoint'ы auth;
- приватные endpoint'ы профиля/раскладов;
- admin endpoint'ы для управления темами;
- middleware: auth, role-check, validation, error handling.

3. MongoDB:
- `users`, `cards`, `daily_readings`, `deck_themes`;
- уникальные индексы: `users.login`, `daily_readings(userId, date)`.

## 5. Ключевые доменные модули backend
- `auth`: регистрация, логин, выдача JWT.
- `users`: профиль и настройки (таймзона, deckTheme).
- `readings`: генерация/выдача ежедневного расклада.
- `cards`: справочник карт (78+).
- `deckthemes`: пользовательский список тем колод.
- `admin`: операции управления темами колод.

## 6. Потоки данных
1. Регистрация/логин:
- frontend -> `/api/auth/register|login` -> backend -> MongoDB;
- backend возвращает JWT, frontend сохраняет токен.

2. Получение расклада на сегодня:
- frontend -> `GET /api/reading/today` (JWT);
- backend определяет дату по `timezone` пользователя;
- если запись в `daily_readings` за дату есть, возвращает ее;
- иначе генерирует, сохраняет, возвращает результат.

3. Настройки пользователя:
- frontend -> `PATCH /api/me/settings` -> backend валидирует и сохраняет.

4. Админ-функции:
- frontend(admin) -> `/api/admin/deck-themes`;
- backend проверяет роль `admin`, затем выполняет CRUD MVP-уровня.

## 7. Правила генерации ежедневного контента
- Card of Day: 1 карта.
- Spread A: 3 карты, без повторов внутри расклада.
- Spread B: 5 карт, без повторов внутри расклада.
- Повторы между блоками (Card of Day / A / B) допускаются в MVP.
- Для каждой карты случайно выбирается `open`/`closed`.
- Интерпретация выбирается по состоянию карты.

## 8. Безопасность (baseline)
- Хэширование паролей (bcrypt/argon2).
- JWT middleware для приватных endpoint'ов.
- Role-based middleware для admin endpoint'ов.
- Валидация входных данных и единый формат ошибок.

## 9. Нефункциональные требования
- Mobile-first UI.
- Время загрузки ключевого экрана: 2-3 секунды в dev при типовом объеме данных.
- Простота кодовой базы для небольшой команды.

## 10. Расширяемость после MVP
- Подключение AI-интерпретаций через feature flag.
- Развитие админки до полноценного CMS-модуля.
- Отображение истории раскладов в UI.
