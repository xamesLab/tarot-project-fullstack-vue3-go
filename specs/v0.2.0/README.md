# Nebula Tarot v0.2.0

`v0.2.0` — новая frontend-итерация проекта. Она заменяет визуальную и архитектурную основу клиента, но пока не заменяет backend-connected возможности `v0.1.x`. Старый Vue frontend перенесён в `frontend-legacy/`, новый React frontend расположен в `frontend/`.

**Статус:** implemented frontend baseline, mock API mode

**Дата фиксации документа:** 3 июля 2026

**Исходный план:** [`frontend/IMPLEMENTATION_PLAN.md`](../../frontend/IMPLEMENTATION_PLAN.md)

## Цели версии

- воспроизвести Home и UI Library из исходных Arcana-прототипов в бренде `NEBULA-TAROT`
- открыть Card of the Day без обязательной авторизации
- отделить UI от транспорта данных через стабильный typed API contract
- заложить масштабируемую feature-oriented структуру
- сделать shared UI Kit исполняемым контрактом дизайн-системы
- сохранить RU-first локализацию, английский язык и две темы

## Нововведения

### Новый Home

Route `/` стал публичной ритуальной страницей карты дня. В неё входят:

- фон из 70 генерируемых звёзд, орбов и созвездий
- sticky header с навигацией, user pill, RU/EN switch и theme toggle
- moon/date strip
- карта Верховной Жрицы с 3D flip по click, Enter или Space
- ключевые слова и краткое толкование
- раскрываемые сферы Love, Career, Health и Spirit
- Message of the Day в corner-ornament panel
- oracle chat с локализованным приветствием и mock-ответом
- адаптивное поведение и `prefers-reduced-motion`

### Регистрационный flow

Все четыре пункта верхней навигации пока считаются gated actions и открывают единый registration modal. Форма содержит name, email, password и необязательную birth date, валидируется через `react-hook-form + zod`, отправляется в mock API и показывает success state.

Это UX-заглушка, а не реальная регистрация: JWT, backend session и восстановление пользователя после reload в `v0.2.0` не подключены.

### UI Kit как runtime-контракт

Route `/ui-kit` демонстрирует реальные компоненты `shared/ui`, используемые Home. Route регистрируется лениво только при `import.meta.env.DEV`; production catch-all перенаправляет на `/`.

Библиотека включает:

- Button, Chip, Badge, Field и TextareaField
- Toggle, Checkbox, RadioCard и SegmentedControl
- Avatar, UserPill, NavBar и Logo
- Ornament, Divider и SectionLabel
- TarotCard с flip/tilt behavior
- AspectBlock, CornerPanel и Tooltip
- доступный Modal с focus trap, Escape/overlay close, scroll lock и restore focus

### Темы и локализация

- default language — `ru`, fallback — `ru`; поддерживаются `ru` и `en`
- default theme — `dark`; поддерживаются `dark` и `light`
- значения хранятся в `nebula-tarot:lang` и `nebula-tarot:theme`
- `<html lang>`, `data-theme` и класс `.light` синхронизируются в `AppProviders`
- UI-строки лежат в namespace `common` и `home`; локализованный контент reading поступает из API fixtures
- смена языка обновляет seeded greeting oracle chat

### Новая архитектура

```text
frontend/src/
├── api/           # DTO, public API facade, mock и будущий HTTP adapter
├── app/           # App, router и providers
├── features/      # auth, preferences, i18n, oracle chat
├── pages/         # Home и dev-only UI Kit
├── shared/        # styles, icons, hooks и controlled UI
├── widgets/       # page-level композиция
└── test/          # smoke setup
```

Правила зависимостей:

- `shared/ui` не импортирует store, router или page-level данные
- `features` владеют пользовательским поведением и локальным состоянием
- `widgets` собирают domain/page blocks
- server state хранится в TanStack Query, preferences — в Zustand
- данные UI запрашиваются только через `api`, а не напрямую из fixtures

### Typed API boundary

Публичный интерфейс `NebulaTarotApi` содержит:

```ts
api.readings.getDaily(language)
api.oracle.ask(question, language)
api.auth.register(payload)
```

`VITE_API_MODE=mock|http` выбирает реализацию в одном месте. В текущей версии development и production используют `mock`; HTTP adapter оставлен явной заглушкой и не должен считаться интегрированным.

### Обновлённый toolchain

| Область | v0.1.x | v0.2.0 |
|---|---|---|
| UI | Vue 3 | React 18 + TypeScript |
| State | Pinia | Zustand + TanStack Query |
| Routing | Vue Router | React Router data router |
| i18n | Vue I18n | i18next + react-i18next |
| Forms | VeeValidate | react-hook-form + zod |
| Styles | SCSS | Tailwind + CSS variables + CSS |
| Tests | ручной baseline | Vitest + Testing Library smoke tests |

## Scope версии

### Реализовано

- Home и dev-only UI Kit
- shared component library и design tokens
- dark/light, RU/EN и persistence настроек
- typed mock readings, oracle и registration APIs
- loading/error behavior для daily reading
- keyboard-accessible card flip и modal focus management
- typecheck, lint, unit smoke test и production build scripts

### Не входит в v0.2.0

- реальная связь нового frontend с Go API
- JWT auth, login, profile и session bootstrap
- настоящие Spreads, Horoscope и Learn routes
- history, deck/card-back selection и admin panel
- реальный LLM oracle streaming
- SSR/SSG, payments, PWA и production deployment

Эти возможности доступны в разной степени в `frontend-legacy/` или остаются следующими этапами миграции.

## Запуск и проверка

```bash
cd frontend
npm install
npm run dev
```

Проверка качества:

```bash
npm run typecheck
npm run lint
npm test
npm run build
```

## Критерии приёмки

- `/` открывается без backend и авторизации
- Home корректно работает на русском и английском, в dark и light themes
- выбор языка и темы сохраняется после reload
- карта переворачивается мышью и клавиатурой
- interpretation раскрывается и сворачивается
- oracle chat добавляет вопрос и локализованный mock response
- gated nav открывает форму регистрации; validation и success state работают
- `/ui-kit` доступен в development и отсутствует как route в production
- typecheck, lint, tests и build завершаются успешно

## Следующий этап

Главный технический шаг после `v0.2.0` — реализовать `frontend/src/api/http.ts`, согласовать DTO с действующими Go endpoints и последовательно перенести auth/session, reading history, spreads, catalog/profile и admin flows из `frontend-legacy/`. До этого новый и legacy frontend следует воспринимать как два разных runtime-клиента, а не как взаимозаменяемые сборки.
