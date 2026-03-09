# Seed Strategy (Step 7, Safe Mode)

Дата фиксации: 8 марта 2026.

## 1. Цель шага
Подготовить воспроизводимый seed-процесс через Go-команду для:
- `cards` (78 карт),
- `deck_themes` (default placeholder theme),
- `users` (test admin).

## 2. Реализация
- Команда: `backend/cmd/seed`
- Основной файл: `backend/cmd/seed/main.go`
- Зависимости:
  - `go.mongodb.org/mongo-driver`
  - `golang.org/x/crypto/bcrypt`

## 3. Режимы запуска

### Safe Mode (default)
Без записи в БД, только генерация и проверка состава seed.
```powershell
go run ./cmd/seed --env-file .env
```

### Safe Mode + payload export
```powershell
go run ./cmd/seed --env-file .env --export-json ../deploy/mongo/seed-payload.json
```

### Apply Mode (запись в MongoDB)
```powershell
go run ./cmd/seed --env-file .env --apply --admin-login admin --admin-password <strong-password>
```

## 4. Idempotency
Применение организовано через `upsert`:
- `cards` по `code`,
- `deck_themes` по `name`,
- `users` (admin) по `login`.

Повторный запуск не должен плодить дубликаты.

## 5. Фактический результат в этой сессии
- Seed-команда успешно выполнена в safe mode.
- Сформирован payload: `deploy/mongo/seed-payload.json`.
- Дополнительно seed применен к живой БД `tarot_daily` (`--apply`):
  - `cards`: 78 вставок,
  - `deck_themes`: default theme создана,
  - `users`: test admin создан.
