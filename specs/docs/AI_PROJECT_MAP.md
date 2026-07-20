# Карта проекта Nebula Tarot для AI-агента

**Назначение:** навигация по актуальному коду, поиск источников истины и безопасное внесение изменений.

**Актуально на:** 21 июля 2026 года.

**Главное правило:** рабочий frontend находится в `frontend/` и написан на React + TypeScript. Каталог `frontend-legacy/` содержит завершённый Vue baseline `v0.1.x`; не читать и не изменять его, если задача явно не относится к legacy.

## 1. С чего начинать

Перед изменениями:

1. Прочитать корневой `README.md` и этот документ.
2. Проверить `git status --short`: в рабочем дереве могут быть чужие незавершённые изменения, их нельзя перезаписывать.
3. Найти ближайшие тесты изменяемого модуля.
4. Для UI свериться с `specs/docs/DESIGN_SYSTEM.md` и development route `/ui-kit`.
5. Для доменных и API-контрактов свериться с `specs/docs/TECHNICAL_SPEC.md`, но при расхождении считать текущий исполняемый код главным источником истины.

Приоритет источников истины:

1. актуальный код и тесты;
2. общие документы из `specs/docs/`;
3. `specs/v0.2.0/` как baseline перехода на React;
4. versioned-документы `v0.1.x` только как исторический контекст.

## 2. Карта верхнего уровня

| Путь | Что находится | Когда открывать |
|---|---|---|
| `frontend/` | Актуальный React SPA | Любые текущие UI, auth и Card of the Day задачи |
| `backend/` | Go/Gin API, MongoDB, OpenRouter | API, бизнес-правила, данные, AI, auth |
| `deploy.sh` | Локальная сборка и ручной production deploy | После commit и push в `main` |
| `specs/docs/` | Общая продуктовая и техническая документация | Перед архитектурными и UI-изменениями |
| `specs/v0.2.0/` | Scope и визуальные прототипы текущей итерации | Для проверки требований к React frontend |
| `scripts/`, `start-dev.sh`, `stop-dev.sh` | Локальный запуск и вспомогательные скрипты | Для full-stack запуска |
| `frontend-legacy/` | Старый Vue 3 клиент | Только по прямому указанию пользователя |
| `frontend/Arcana Home.html` | Неизменяемый визуальный референс Home | Для визуальной сверки |
| `frontend/Arcana UI Library.html` | Неизменяемый визуальный референс UI Kit | Для визуальной сверки компонентов |

Не использовать сгенерированные файлы как исходники: `frontend/dist/`, `frontend/vite.config.js`, `frontend/tailwind.config.js`, `frontend/*.d.ts`, backend binary `backend/tarot-app`.

## 3. Архитектура системы

```text
Browser
  -> React Router / AppProviders
  -> pages
  -> features + widgets
  -> typed API facade
       -> mock adapter (часть content API)
       -> HTTP client (/api)
            -> Gin middleware
            -> handlers
            -> services
            -> repositories
            -> MongoDB
            -> локальный SOCKS relay
            -> OpenRouter для AI
```

Frontend и backend имеют разные модели состояния:

- локальное клиентское состояние: Zustand;
- серверное/асинхронное состояние: TanStack Query;
- постоянная клиентская сессия и настройки: `localStorage`;
- постоянные доменные данные: MongoDB;
- изображения колод и рубашек: `backend/uploads/`, доступны через `/uploads`.

## 4. Актуальный frontend

### 4.1 Точки входа

| Путь | Ответственность |
|---|---|
| `frontend/src/main.tsx` | React root и подключение глобальных стилей |
| `frontend/src/app/App.tsx` | `AppProviders` + `RouterProvider` |
| `frontend/src/app/providers/AppProviders.tsx` | QueryClient, восстановление сессии, theme/language side effects, обработка unauthorized event |
| `frontend/src/app/router.tsx` | `/`, dev-only `/ui-kit`, fallback redirect |
| `frontend/src/pages/home/HomePage.tsx` | Композиция главной страницы без низкоуровневой бизнес-логики |

### 4.2 Слои frontend

```text
frontend/src/
  api/          DTO, HTTP client, mock adapter, единый API facade
  app/          bootstrap приложения, providers, router
  features/     пользовательские сценарии и состояние
  pages/        композиция маршрутов
  shared/       независимые стили, иконки, hooks/utilities и UI-компоненты
  widgets/      крупные визуальные блоки страниц
  test/         Vitest setup и интеграционные/smoke-тесты
```

Правила зависимостей:

- `shared/ui` не должен импортировать Zustand stores, API, router или feature-модули;
- `features` владеют сценарием и состоянием, но не всей разметкой страницы;
- `widgets` собирают визуальные блоки из `shared/ui` и feature-логики;
- `pages` соединяют готовые features/widgets;
- DTO и transport-контракты размещаются в `api`, а не внутри компонентов.

### 4.3 Где искать frontend-функциональность

| Задача | Основные файлы |
|---|---|
| API типы | `frontend/src/api/types.ts` |
| HTTP, JWT header, ошибки и URL изображений | `frontend/src/api/client.ts` |
| Выбор mock/HTTP transport | `frontend/src/api/index.ts` |
| Реальные HTTP endpoints | `frontend/src/api/http.ts` |
| Mock данные и ответы | `frontend/src/api/mock/` |
| Login/registration UI | `frontend/src/features/auth/AuthModal.tsx` |
| Auth session/state | `frontend/src/features/auth/store.ts` |
| Карта дня: guest cache/claim/resolve | `frontend/src/features/daily-card/bootstrap.ts` |
| React Query bootstrap карты | `frontend/src/features/daily-card/useDailyCardBootstrap.ts` |
| AI-толкование карты | `frontend/src/features/daily-card/useDailyCardInterpretation.ts` |
| Преобразование backend DTO для показа | `frontend/src/features/daily-card/presentation.ts` |
| Язык и тема | `frontend/src/features/preferences/store.ts` |
| i18n bootstrap и словари | `frontend/src/features/i18n/` |
| Oracle chat | `frontend/src/features/oracle-chat/OracleChat.tsx` |
| Карта и 3D flip | `frontend/src/shared/ui/TarotCard/` |
| Карта дня на сцене | `frontend/src/widgets/CardOfDayStage/` |
| Раскрываемые детали | `frontend/src/widgets/CardOfDayDetails/` |
| Фоновая анимация | `frontend/src/widgets/StarField/` |
| Header и auth entry points | `frontend/src/widgets/Header/` |
| Общие UI-компоненты | `frontend/src/shared/ui/` |
| SVG-иконки | `frontend/src/shared/icons/index.tsx` |

### 4.4 API facade и режимы транспорта

`frontend/src/api/index.ts` экспортирует один объект `api`:

- `api.readings` и `api.oracle` выбираются через `VITE_API_MODE`;
- `api.auth`, profile, catalog, history, spreads и admin выбираются через `VITE_AUTH_API_MODE`;
- `api.dailyCard` всегда использует текущий HTTP adapter;
- HTTP реализации `readings` и `oracle` пока намеренно недоступны и выбрасывают ошибку; их нельзя включать через `VITE_API_MODE=http` без реализации;
- базовый адрес берётся из `VITE_API_URL`, fallback — `/api`.

JWT хранится под ключом `tarot_daily_token`. `401`/`403` от authenticated-запроса вызывает событие `nebula:auth-unauthorized`; `AppProviders` выполняет logout.

### 4.5 Главный пользовательский поток: Card of the Day

```text
AppProviders.hydrateSession()
  -> useDailyCardBootstrap()
  -> bootstrapDailyCard(user)
       guest без cache -> POST /reading/guest/card-of-day
       guest с актуальным cache -> вернуть cache
       login с guest claim -> POST /reading/card-of-day/resolve с claim_token
       user без claim -> POST /reading/card-of-day/resolve
  -> CardOfDayStage
  -> пользователь переворачивает карту
  -> useDailyCardInterpretation()
       guest -> POST /reading/guest/card-of-day/interpretation
       user -> POST /reading/card-of-day/interpretation
  -> CardOfDayDetails
```

Guest cache: `nebula-tarot:guest-card-of-day:v3`. Смена accounting date инвалидирует cache. Старые ключи `v1` и `v2` удаляются при bootstrap.

Критические инварианты:

- accounting day начинается в `03:00` часового пояса пользователя;
- AI-толкование не запрашивается до раскрытия карты;
- guest claim должен переносить гостевую карту в authenticated reading;
- в query key должны оставаться user/guest identity, reading/card, claim token и language;
- при смене языка текст AI должен запрашиваться/кэшироваться отдельно.

### 4.6 UI и дизайн-система

Источники:

- tokens: `frontend/src/shared/styles/tokens.css`;
- fonts: `frontend/src/shared/styles/fonts.css`;
- animations: `frontend/src/shared/styles/animations.css`;
- глобальные и layout-стили: `frontend/src/shared/styles/index.css`;
- Home-specific стили: `frontend/src/pages/home/home.css`;
- component-specific стили карты: `frontend/src/shared/ui/TarotCard/tarot-card.css`;
- живой каталог компонентов: `/ui-kit` и `frontend/src/pages/ui-kit/`.

Правила безопасного изменения UI:

1. Не менять visual reference HTML.
2. Не вводить новый цвет, если существует CSS token.
3. Изменение общего компонента сначала проверить на Home и `/ui-kit`.
4. Сохранять существующие `className`, DOM hierarchy и accessibility attributes при чистом рефакторинге.
5. Новые пользовательские строки добавлять одновременно в `ru` и `en` locale JSON.
6. Интерактивные элементы должны работать с клавиатуры и иметь видимый `:focus-visible`.
7. Проверять dark/light, desktop/mobile и `prefers-reduced-motion`.

### 4.7 Клиентское хранение

| Ключ | Назначение |
|---|---|
| `tarot_daily_token` | JWT |
| `nebula-tarot:theme` | `dark` или `light` |
| `nebula-tarot:lang` | `ru` или `en` |
| `nebula-tarot:guest-card-of-day:v3` | Гостевая карта, timezone, accounting date и claim token |

Не переименовывать ключи без миграции: они являются частью совместимости между сессиями.

## 5. Backend

### 5.1 Точка входа и слои

`backend/main.go`:

1. загружает env и подключает MongoDB;
2. создаёт Gin router и middleware;
3. создаёт repositories, индексы и default settings;
4. собирает services и handlers;
5. регистрирует public/protected/admin routes;
6. запускает HTTP server и graceful shutdown.

```text
HTTP request
  -> middleware/     CORS, JWT auth, role check
  -> handlers/       bind/validate, HTTP status и JSON contract
  -> services/       бизнес-правила и orchestration
  -> repositories/   только MongoDB data access
  -> models/         persistence/API structures
```

Не переносить бизнес-правила в handlers или repositories. Handler отвечает за HTTP boundary; service — за решение; repository — за запрос к данным.

### 5.2 Где искать backend-функциональность

| Задача | Handler | Service | Repository/model |
|---|---|---|---|
| Auth | `handlers/auth.go` | `services/auth.go` | `repositories/user.go`, `models/auth.go`, `models/user.go` |
| Профиль/настройки | `handlers/user.go` | `services/user.go` | `repositories/user.go` |
| Карта дня и расклады | `handlers/reading.go` | `services/reading.go` | `repositories/reading.go`, `models/reading.go` |
| AI/ask | `handlers/ask.go` | `services/ask.go`, `services/ai_prompt.go` | `repositories/ask_request.go`, `repositories/ai_interpretation_cache.go` |
| Каталог колод/рубашек | `handlers/catalog.go` | `services/catalog.go` | `repositories/deck.go`, `repositories/card_back.go` |
| Admin | `handlers/admin.go` | `services/admin.go` | admin/policy/settings repositories и models |
| Ошибки HTTP boundary | `handlers/helpers.go` | `services/errors.go` | — |
| JWT/пароли/validation/time | — | — | `internal/utils/` |
| Upload и распаковка колоды | Admin handler/service | `services/admin.go` | `utils/deck_upload.go` |

Общие handler helpers находятся в `backend/internal/handlers/helpers.go`: JSON bind + validation, извлечение authenticated user id, единый error response. Не дублировать эту boundary-логику в новом handler.

### 5.3 Группы API

Все пути имеют prefix `/api`.

| Доступ | Маршруты |
|---|---|
| Public health | `GET /health`, `/ai/health`, `/gateway/health` |
| Public auth | `POST /auth/register`, `/auth/login` |
| Public guest reading | `POST /reading/guest/card-of-day` |
| Public guest AI | `POST /reading/guest/card-of-day/interpretation` |
| Authenticated user | `GET /me`, `PATCH /me/profile`, `PATCH /me/settings` |
| Authenticated readings | `/reading/card-of-day/resolve`, `/reading/generate/*`, reveal, history |
| Authenticated AI | `POST /ask`, `GET /ask/:requestId`, user card interpretation |
| Authenticated catalog | `GET /decks/available`, `/card-backs/available` |
| Keeper only | `/admin/*` |

Полный и точный список путей искать в методах `RegisterRoutes` и `RegisterPublicRoutes` внутри `backend/internal/handlers/`. Регистрация групп и middleware находится в `backend/main.go`.

### 5.4 MongoDB

| Collection | Repository | Назначение |
|---|---|---|
| `users` | `user.go` | пользователи, роли, настройки, login metadata |
| `cards` | `card.go` | каталог карт |
| `daily_readings` | `reading.go` | карта дня и расклады |
| `decks` | `deck.go` | колоды |
| `card_backs` | `card_back.go` | рубашки |
| `role_policies` | `policy.go` | лимиты и retention по ролям |
| `admin_settings` | `admin_settings.go` | глобальные настройки |
| `ask_requests` | `ask_request.go` | состояние AI-запросов |
| `ai_interpretation_cache` | `ai_interpretation_cache.go` | кэш AI-толкований |

Индексы создаются при старте в `main.go` и/или seed utility. При изменении uniqueness, TTL или query pattern обновлять repository indexes и тесты одновременно.

### 5.5 Ключевые доменные правила

- Card of the Day определяется пользователем/гостем, timezone и accounting date.
- Accounting day: `03:00`–`02:59`, реализация времени — `backend/internal/utils/time.go` и frontend counterpart в daily-card bootstrap.
- Роль определяет дневные лимиты, retention и доступные настройки.
- Публичная регистрация контролируется `admin_settings.registration_enabled`.
- Роль по умолчанию из admin settings не может автоматически выдать `keeper`.
- Reveal state хранится в reading и влияет на доступность AI-интерпретации.
- Guest claim подписан JWT secret и должен проверяться до привязки карты к пользователю.
- Ошибки сервисов должны оставаться машинно-читаемыми (`validation_failed`, `invalid_credentials`, `invalid_guest_claim` и т. п.).
- `services/ask.go` большой по объёму; перед изменением определить конкретный поток: guest interpretation, user interpretation, request reuse, provider call, cache или limit.

### 5.6 Конфигурация backend

Шаблон: `backend/.env.example`. Обязательная переменная для запуска — `MONGO_URI`.

Основные группы:

- app/server: `APP_ENV`, `APP_PORT`, HTTP и shutdown timeouts;
- auth: `JWT_SECRET`, `JWT_TTL`;
- MongoDB: `MONGO_URI`, `MONGO_DB_NAME`, connect/ping timeouts;
- browser access/assets: `FRONTEND_URL`, `BACKEND_PUBLIC_BASE_URL`;
- direct AI: `AI_DIRECT_ENABLED`, `OPENROUTER_*`;
- production egress: `HTTPS_PROXY`, `NO_PROXY`;
- legacy gateway stub: `GATEWAY_*`.

Секреты не добавлять в Git, тестовые ответы и логи.

### 5.7 Production runtime

Production использует immutable releases:

```text
/opt/tarot-app/releases/<commit-sha>
/opt/tarot-app/current -> releases/<commit-sha>
```

Ключевые внешние runtime-ресурсы:

- `/etc/tarot-app/app.env` — env и секреты, не входит в Git/release;
- `/var/lib/tarot-app/uploads` — persistent uploads;
- `tarot-app.service` — backend под `tarot-user`;
- `tarot-openrouter-relay.service` — отдельный SSH/SOCKS relay для OpenRouter;
- MongoDB — отдельный сервис, обычный deploy его не изменяет;
- `gtwapi` — вне deploy flow и не перезапускается.

После push production обновляется вручную из корня репозитория:

```bash
bash deploy.sh
```

Действующие runbooks:

- `specs/docs/DEPLOYMENT_GUIDE.md`;
- `specs/docs/OPENROUTER_PROXY.md`.

## 6. Тесты и команды

### Frontend

```bash
cd frontend
npm install
npm run typecheck
npm run lint
npm test
npm run build
```

Тесты находятся в `frontend/src/test/`. Для компонентов использовать Testing Library и проверять поведение/доступность, а не внутреннюю реализацию. API boundary уже покрывается mocked `fetch` и typed fixtures.

### Backend

```bash
cd backend
go test ./...
go vet ./...
go build ./...
```

Если локальное окружение принудительно отключает Go modules, запускать с `GO111MODULE=on`. Unit-тесты расположены рядом с пакетами в `*_test.go`; большинство service-тестов используют fake repositories и не требуют MongoDB.

### Full-stack development

```bash
cp backend/.env.example backend/.env
./start-dev.sh
```

По умолчанию frontend — `http://localhost:5173`, backend — `http://localhost:8888`, API — `http://localhost:8888/api`.

## 7. Алгоритм работы AI-агента

### Для исправления бага

1. Воспроизвести баг существующим тестом или минимальной read-only проверкой.
2. Проследить поток от page/widget до feature/API, затем handler/service/repository.
3. Исправлять на слое, который владеет нарушенным правилом.
4. Добавить regression test рядом с владельцем правила.
5. Запустить targeted tests, затем полный набор соответствующей части проекта.

### Для новой frontend-функции

1. Добавить/уточнить DTO в `api/types.ts`.
2. Реализовать transport в `api/http.ts` или mock adapter.
3. Асинхронный сценарий разместить в feature hook/store.
4. Переиспользуемый внешний вид разместить в `shared/ui` и показать в `/ui-kit`.
5. Собрать страницу через widget/page без дублирования.
6. Добавить ru/en строки и тесты.

### Для нового backend endpoint

1. Модель запроса/ответа — `models`.
2. Data access — repository interface + implementation.
3. Бизнес-правило — service с тестом на fake repository.
4. HTTP bind/validate/status — handler через общие helpers.
5. Route и подходящий middleware — `RegisterRoutes` + `main.go`.
6. Если меняется API contract, синхронно обновить frontend DTO/adapter и документацию.

### Для рефакторинга

1. Зафиксировать зелёный baseline.
2. Не смешивать рефакторинг с изменением поведения без явной причины.
3. Сохранять public API, JSON fields, storage keys, query keys, DOM/CSS contracts.
4. Делать небольшие логические изменения и проверять после каждого блока.
5. В конце просмотреть `git diff --check`, `git diff --stat` и убедиться, что legacy/сгенерированные файлы не затронуты.

## 8. Типичные ошибки

- Искать актуальный UI во Vue-каталоге вместо `frontend/`.
- Считать `VITE_API_MODE=http` полностью реализованным: HTTP readings/oracle пока отсутствуют.
- Считать `git push` автоматическим деплоем: GitHub Actions пока не настроен, нужен `bash deploy.sh`.
- Копировать production env, адрес relay-сервера или ключевой материал в Git/документацию.
- Вызывать AI до reveal или терять language в query/cache key.
- Считать календарную полночь границей карты дня вместо `03:00`.
- Добавлять fetch напрямую в React component вместо `api` + feature hook.
- Хранить server state в Zustand вместо TanStack Query.
- Дублировать handler error/bind/auth logic вместо `handlers/helpers.go`.
- Возвращать внутреннюю ошибку MongoDB/OpenRouter пользователю.
- Менять hex-цвета отдельно от theme tokens.
- Обновлять только один locale.
- Перезаписывать существующие незакоммиченные изменения.
- Запускать форматирование или массовую замену на всём репозитории без проверки diff.

## 9. Чек-лист завершения задачи

- Изменён правильный слой и нет ненужного дублирования.
- `frontend-legacy/` не затронут без явного требования.
- API, storage, i18n и UI contracts сохранены либо осознанно обновлены вместе.
- Добавлены или обновлены тесты для нового поведения.
- Frontend: typecheck, lint, tests, build.
- Backend: tests, vet, build.
- Для UI проверены Home и `/ui-kit`, dark/light и responsive состояния.
- `git diff --check` чистый; в diff нет секретов, build artifacts и посторонних файлов.
