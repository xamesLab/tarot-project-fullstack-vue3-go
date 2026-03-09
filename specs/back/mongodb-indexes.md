# MongoDB Index Plan (Step 6)

Дата фиксации: 8 марта 2026.

## 1. Цель
Зафиксировать обязательные индексы MVP и подготовить воспроизводимый скрипт их применения.

## 2. Обязательные индексы MVP

### `users`
- Index: `{ login: 1 }`
- Name: `uniq_users_login`
- Type: `unique`
- Назначение: исключить дублирование логинов.

### `daily_readings`
- Index: `{ userId: 1, date: 1 }`
- Name: `uniq_daily_readings_user_date`
- Type: `unique compound`
- Назначение: одна запись расклада на пользователя в пределах одной даты.

### `cards`
- Index: `{ code: 1 }`
- Name: `uniq_cards_code`
- Type: `unique`
- Назначение: стабильный уникальный код карты для ссылок в `assets` тем колод.

### `deck_themes`
- Index: `{ name: 1 }`
- Name: `uniq_deck_themes_name`
- Type: `unique`
- Назначение: исключить дубли названий тем колод в admin API.

### `role_policies`
- Index: `{ role: 1 }`
- Name: `uniq_role_policies_role`
- Type: `unique`
- Назначение: одна политика лимитов на роль.

## 3. Скрипты применения
- Mongo script: `deploy/mongo/create-indexes.js`
- PowerShell wrapper: `deploy/scripts/apply-mongo-indexes.ps1`

## 4. Команда запуска
```powershell
powershell -ExecutionPolicy Bypass -File deploy/scripts/apply-mongo-indexes.ps1
```

## 5. Фактическое применение
Индексы применены к живой БД `tarot_daily` 8 марта 2026 через Go-команду:
- `backend/cmd/indexes`

Контрольная проверка после применения:
- `users`: 1 doc
- `cards`: 78 docs
- `deck_themes`: 1 doc
- `daily_readings`: 0 docs
- `role_policies`: 0 docs (или больше, если роли уже добавлялись).
