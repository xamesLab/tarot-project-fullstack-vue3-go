# Refactor Log v0.1.1

Дата фиксации: 9 марта 2026.

## 1. Цель

Зафиксировать выполненный не-функциональный рефакторинг проекта в версии `0.1.1`:
- улучшение читаемости кода;
- разбиение перегруженных файлов и пакетов;
- вынос повторяющейся логики в общие helper-ы;
- подготовка базы для безопасных следующих изменений.

Важно:
- API-контракты и пользовательский сценарий не менялись намеренно;
- функциональные багфиксы и визуальные изменения в этот журнал не включены, если они не были частью структурного рефакторинга.

## 2. Выполненный рефакторинг

### 2.1. HTTP transport layer

Сделано:
- добавлен общий helper-слой для хендлеров:
  - `internal/transport/http/handlers/helpers.go`
- в helper-слой вынесены:
  - `bindJSON(...)`
  - `currentUserID(...)`
  - `parseInt64Query(...)`
  - `respondByError(...)`
  - `respondInternalError(...)`

Эффект:
- убрано дублирование `ShouldBindJSON`, извлечения `auth_user_id` и типового error mapping;
- хендлеры стали короче и однообразнее по стилю ошибок.

Обновлены хендлеры:
- `auth_handler.go`
- `me_handler.go`
- `settings_handler.go`
- `reading_handler.go`
- `deck_themes_handler.go`
- `admin_settings_handler.go`
- `admin_roles_handler.go`
- `admin_users_handler.go`
- `decks_handler.go`
- `admin_decks_handler.go`

### 2.2. Router / bootstrap

Сделано:
- точка сборки HTTP-приложения разделена на отдельные файлы:
  - `internal/transport/http/router.go`
  - `internal/transport/http/bootstrap.go`
  - `internal/transport/http/routes.go`

Эффект:
- `router.go` больше не смешивает создание engine, DI и регистрацию всех маршрутов;
- структура запуска стала понятнее и проще для расширения.

### 2.3. Пакет `reading`

Сделано:
- перегруженный `internal/reading/service.go` разделён по обязанностям:
  - `service.go` — типы, конструктор, публичные entrypoints;
  - `daily.go` — сценарий дневного чтения и лимитов;
  - `draw.go` — логика генерации карт и раскладов;
  - `view.go` — сборка DTO и резолв ассетов.

Эффект:
- orchestration, random draw logic и view mapping больше не перемешаны в одном файле;
- пакет стало проще читать и тестировать по частям.

### 2.4. Пакет `decks`

Сделано в два этапа.

Этап 1:
- файловое сохранение и вычисление upload root вынесены из сервиса в:
  - `internal/decks/upload_storage.go`

Этап 2:
- пакет дополнительно разделён по обязанностям:
  - `service.go` — основной service API;
  - `upload.go` — upload back/card;
  - `view.go` — mapping deck -> view и резолв имён пользователей;
  - `assets.go` — helper-ы для card key, content type и статуса колоды.

Эффект:
- бизнес-логика колод отделена от работы с файловой системой и мелких asset helper-ов;
- код upload-сценариев стал изолированным и проще для проверки.

### 2.5. Пакет `auth`

Сделано:
- `internal/auth/service.go` разделён по обязанностям:
  - `service.go` — ошибки, DTO, struct `Service`, конструктор;
  - `account.go` — `Register`, `Login`, `Me`, `validateCredentials`, выдача токена;
  - `settings.go` — `UpdateSettings`;
  - `view.go` — `toUserView`, нормализация timezone/status.

Эффект:
- auth-модуль перестал быть монолитным;
- регистрация, логин, профиль и настройки теперь разделены логически.

### 2.6. Общие доменные helper-ы

Сделано:
- добавлен `internal/domain/models/user_helpers.go`

Вынесено в общее место:
- нормализация роли;
- валидация роли;
- fallback `displayName -> login`;
- нормализация статуса пользователя.

Эффект:
- убрано дублирование правил между `roles`, `adminusers`, `auth`, `decks`.

### 2.7. Общие input helper-ы

Сделано:
- добавлен `internal/platform/inpututil/input.go`

В helper вынесены:
- `RequireTrimmed(...)`
- `ParseObjectID(...)`
- `ParseOptionalObjectID(...)`

Эффект:
- убрана повторяющаяся логика `TrimSpace + ObjectIDFromHex`;
- сервисы используют единые правила валидации входа.

Пакеты, переведённые на helper:
- `auth`
- `decks`
- `deckthemes`
- `reading`

### 2.8. Общий env/config loader для CLI

Сделано:
- добавлен `internal/config/envfile.go`

Обновлены команды:
- `cmd/migrate011`
- `cmd/seed`
- `cmd/indexes`

Эффект:
- убран дублирующийся парсинг `.env` в утилитах;
- CLI-команды используют единый способ загрузки env-файла.

## 3. Добавленные тесты

Добавлены unit-тесты на вынесенную чистую логику:
- `internal/platform/timeutil/timezone_test.go` — уже существовал и сохранён как базовая опора;
- `internal/domain/models/user_helpers_test.go`
- `internal/platform/inpututil/input_test.go`
- `internal/reading/view_test.go`
- `internal/auth/account_test.go`
- `internal/auth/view_test.go`
- `internal/decks/assets_test.go`
- `internal/config/envfile_test.go`

Что покрыто:
- нормализация ролей и user display data;
- парсинг `ObjectID`;
- маппинг reading view;
- валидация auth credentials;
- нормализация user view;
- helper-ы для asset/status в `decks`;
- загрузка env-файла.

## 4. Проверка после рефакторинга

После каждого крупного этапа выполнялась проверка:

```bash
go test ./...
```

Результат:
- backend собирается;
- тесты проходят;
- целевые API-контракты не были сломаны рефакторингом.

## 5. Что ещё остаётся как следующий слой рефакторинга

Следующие кандидаты:
- разделить `internal/repository/mongo/user_repository.go` на CRUD и admin-list/filter helpers;
- добавить сервисные тесты для `auth.UpdateSettings`;
- добавить сервисные тесты для `decks.UploadBack` и `decks.UploadCard`;
- постепенно покрыть integration-тестами критические HTTP endpoints.
