# Nebula Tarot v0.2.0

`v0.2.0` — зафиксированный baseline перехода проекта с Vue на React. Старый Vue frontend перенесён в `frontend-legacy/`, новый React frontend расположен в `frontend/`.

**Статус baseline:** HTTP auth + mock content

**Текущее дополнение:** post-baseline HTTP integration и production deployment выполнены

**Дата фиксации документа:** 3 июля 2026

**Проверено:** 21 июля 2026

**Исходный план:** [`frontend/IMPLEMENTATION_PLAN.md`](../../frontend/IMPLEMENTATION_PLAN.md)

После фиксации baseline React-клиент был расширен: Card of the Day и AI-интерпретация всегда используют Go HTTP API; добавлены profile, spreads и admin routes; account-related adapters подключены к backend; клиент развёрнут в production. Ниже сохранён исходный scope `v0.2.0`, поэтому утверждения о mock-only Home следует читать как состояние на дату фиксации, а не как текущую operational конфигурацию.

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

Card of the Day остаётся публичной. Guest pill открывает login, а Spreads, Horoscope и Learn для гостя открывают registration. Единый modal переключается между режимами без закрытия.

Регистрация использует `email`, необязательный `username`, `password` и обязательный IANA timezone. Login/registration подключены к Go API, JWT сохраняется в `localStorage`, а профиль восстанавливается через `/api/me`. В шапке отображаются Guest или реальные username/role и доступен logout.

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
├── api/           # DTO, hybrid facade, mock content и HTTP auth
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
api.auth.login(payload)
api.auth.getCurrentUser()
```

Transport выбирается по доменам: `VITE_API_MODE=mock` оставляет readings/oracle на fixtures, а `VITE_AUTH_API_MODE=http` подключает auth к backend. Это позволяет мигрировать API вертикальными slices, не ломая Home.

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

## Исходный scope версии

### Реализовано

- Home и dev-only UI Kit
- shared component library и design tokens
- dark/light, RU/EN и persistence настроек
- typed mock readings/oracle и HTTP login/registration/session APIs
- loading/error behavior для daily reading
- keyboard-accessible card flip и modal focus management
- typecheck, lint, unit smoke test и production build scripts

### Не входило в исходный v0.2.0 baseline

- реальная связь нового frontend с Go API
- profile/settings screen и остальные protected routes
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

## Исходные критерии приёмки

- `/` открывается без backend и авторизации
- Home корректно работает на русском и английском, в dark и light themes
- выбор языка и темы сохраняется после reload
- карта переворачивается мышью и клавиатурой
- interpretation раскрывается и сворачивается
- oracle chat добавляет вопрос и локализованный mock response
- guest pill открывает login, gated nav — registration; validation, session restore и logout работают
- `/ui-kit` доступен в development и отсутствует как route в production
- typecheck, lint, tests и build завершаются успешно

## Следующий этап на дату baseline

Исходным следующим шагом была реализация `frontend/src/api/http.ts` и перенос auth/session, reading history, spreads, catalog/profile и admin flows. Большая часть этого переноса выполнена. Текущие незакрытые направления — HTTP implementation generic readings/oracle facade, дальнейшая продуктовая разработка и CI/CD automation. Для текущего состояния использовать корневой README, `specs/docs/AI_PROJECT_MAP.md` и operational runbooks.
