# Design System Nebula Tarot

**Актуальная версия:** `v0.2.0`

**Source of truth:** `frontend/src/shared/` и development route `/ui-kit`

**Проверено:** 21 июля 2026 года

Этот документ описывает дизайн-систему нового React frontend. Дизайн-система Vue-клиента `v0.1.x` сохранена вместе с кодом в `frontend-legacy/` и больше не является source of truth для `frontend/`.

## 1. Принципы

- Visual source — неизменяемые прототипы `frontend/Arcana Home.html` и `frontend/Arcana UI Library.html`.
- Product brand во всех пользовательских текстах — `NEBULA-TAROT`; имя Arcana сохраняется только у исходных файлов.
- CSS custom properties — единственный источник theme colors.
- Один и тот же компонент используется на Home и в `/ui-kit`; showcase не дублирует production markup.
- `shared/ui` принимает props и не зависит от Zustand, i18next, router или API.
- Feature/widget может владеть состоянием, переводами и domain composition.
- Сложная графика и 3D behavior остаются в component CSS; Tailwind используется там, где он не ухудшает fidelity.
- Любой интерактивный элемент должен иметь keyboard path и видимый `:focus-visible`.

## 2. Visual foundations

### 2.1 Цветовые tokens

Source: `frontend/src/shared/styles/tokens.css`.

| Token | Dark | Light | Назначение |
|---|---|---|---|
| `--gold` | `#c9a96e` | `#8a6a2c` | основной accent |
| `--gold-light` | `#e8d5a3` | `#a07f3a` | светлый accent и звёзды |
| `--gold-dim` | `#8a6f47` | `#6b5526` | вторичный accent |
| `--ink` | `#0a0a0f` | `#f0e8d9` | фон страницы |
| `--ink-mid` | `#13131c` | `#fcf8f1` | surfaces |
| `--ink-soft` | `#1c1c2a` | `#e9decb` | elevated/control surfaces |
| `--smoke` | `#2a2a3d` | `#dfd1bb` | дополнительный surface |
| `--mist` | `#e8e2d9` | `#221c2d` | основной текст |
| `--mist-dim` | `#9b9189` | `#564e40` | muted text |
| `--cream` | `#f5f0e8` | `#171120` | контрастный текст |
| `--line` | gold 12% | gold 28% | границы |
| `--line-soft` | gold 10% | gold 18% | мягкие разделители |

Dark theme задаётся в `:root`. Light theme переопределяет тот же контракт в `:root.light`. Компоненты не должны определять отдельную палитру hex-значений, кроме неизменяемой иллюстрации карты и намеренно контрастных цветов внутри gradients.

### 2.2 Typography

Шрифты self-hosted через `@fontsource`:

- `Cinzel` — wordmark, headings, roman numerals и section labels
- `Cormorant Garamond` — body copy, interpretations, form controls и card titles
- `Raleway` — navigation, controls, metadata и buttons

Подключаются только используемые weights из `frontend/src/shared/styles/fonts.css`. Для кириллицы должен сохраняться корректный fallback, поскольку Cinzel не содержит полного Cyrillic набора.

### 2.3 Motion

Source: `frontend/src/shared/styles/animations.css`.

- `twinkle` — star field
- `float` — ambient orbs
- `spin` — zodiac wheel, полный оборот за 120 секунд
- card flip — `.85s` с `transform-style: preserve-3d`
- card tilt — pointer-only hook `useTilt`
- modal — opacity + translate/scale transition

`prefers-reduced-motion: reduce` отключает ambient animations и tilt, сокращает переходы. Motion не должен быть единственным способом передать состояние.

## 3. Слои

```text
tokens/fonts/animations
        ↓
shared/icons + shared/lib
        ↓
shared/ui
        ↓
features + widgets
        ↓
pages
```

### Foundations

- `shared/styles/tokens.css` — theme contract
- `shared/styles/fonts.css` — font loading
- `shared/styles/animations.css` — shared motion
- `shared/styles/index.css` — reset, Tailwind layers и глобальные component classes

### Shared UI

| Группа | Компоненты | Основные состояния |
|---|---|---|
| Actions | `Button` | primary, outline, ghost; md/sm; disabled; icon-only |
| Labels | `Chip`, `Badge` | outline/solid; gold/soft/dim; optional dot |
| Forms | `Field`, `TextareaField` | default, focus, error, disabled browser state |
| Controls | `Toggle`, `Checkbox`, `RadioCard`, `SegmentedControl` | on/off, checked, selected |
| Identity | `Avatar`, `UserPill`, `Logo` | sm/md/lg, ring, optional tagline |
| Navigation | `NavBar` | default, hover, active |
| Structure | `Ornament`, `Divider`, `SectionLabel` | star/cross, solid/dashed |
| Tarot | `TarotCard` | face/back, flippable, tilt |
| Content | `AspectBlock`, `CornerPanel`, `Tooltip` | sphere type, decorated panel, top/below tooltip |
| Overlay | `Modal` | open, Escape/overlay close, focus trap |

### Features

- `preferences` — theme/language store and persistence
- `i18n` — i18next resources and initialization
- `auth` — registration modal and mock user state
- `oracle-chat` — message state and API mutation

### Widgets

- `Header`
- `DateStrip`
- `CardOfDayStage`
- `Interpretation`
- `BackgroundDecor`

Widgets собирают shared primitives в крупные блоки Home и могут читать translations/store или принимать API DTO.

## 4. Component contracts

### Button

Использовать `variant="primary"` для главного действия, `outline` для равнозначного secondary action, `ghost` для отмены или малозаметного действия. Label должен быть глаголом; icon-only button обязан иметь `aria-label`.

### Controls

`Toggle` использует `role="switch"`, `Checkbox` — `role="checkbox"`, `RadioCard` — `role="radio"`, `SegmentedControl` — group с `aria-pressed`. Эти компоненты controlled: состояние всегда передаёт владелец feature/widget.

### TarotCard

- `flippable` включает click и keyboard activation через Enter/Space
- `aria-pressed` отражает flipped state
- `tilt` применяется только для showcase/hover-capable устройств
- face и back имеют одинаковую геометрию и скрывают обратную сторону через `backface-visibility`

### Modal

Обязательные свойства: `open`, `onClose`, `title`. Реализация через portal должна:

- выставлять `role="dialog"` и `aria-modal="true"`
- блокировать body scroll
- переводить focus на первый control или `initialFocusRef`
- удерживать Tab внутри panel
- закрываться по Escape и click на overlay
- возвращать focus к исходному trigger

## 5. Page composition

### Home `/`

```text
BackgroundDecor
└── Header
    └── Page
        ├── DateStrip
        ├── SectionLabel
        └── Card stage
            ├── CardOfDayStage
            ├── Interpretation
            └── OracleChat
└── Footer
└── RegistrationModal
```

### UI Kit `/ui-kit`

UI Kit показывает foundations, controls, tarot components, decorated content и modal behavior. Он существует только в development и загружается через lazy import. Новая visual primitive считается частью системы только после появления здесь и повторного использования в реальном screen/widget.

## 6. Responsive rules

- минимальная поддерживаемая ширина — `320px`
- основной content container — до `1180px`
- ниже `900px` скрываются desktop navigation, constellation decor и user pill
- horizontal page padding уменьшается с `40px` до `20px`
- ниже `480px` скрывается logo tagline и уплотняются language switch/modal paddings

Новые компоненты должны быть usable на touch и не зависеть от hover. Hover enhancement допустим только как необязательный дополнительный эффект.

## 7. Accessibility

- visible gold focus ring для button, link, input, textarea и custom tabindex controls
- decorative SVG получают `aria-hidden`
- динамический chat log использует `aria-live="polite"`
- ошибки Field связываются с input через `aria-describedby`
- theme, language и icon actions имеют локализованные accessible names
- contrast проверяется отдельно для обеих тем

## 8. Добавление нового элемента

1. Проверить, можно ли собрать задачу из существующих tokens/components.
2. Добавить недостающие tokens сразу в dark и light blocks.
3. Создать controlled component в `shared/ui/<Component>/` и barrel export.
4. Добавить keyboard и accessible semantics.
5. Продемонстрировать состояния на `/ui-kit`.
6. Только после этого использовать компонент в feature/widget/page.
7. Запустить `npm run typecheck`, `npm run lint`, `npm test` и `npm run build`.

## 9. Запрещённые отклонения

- импорт store, router, API или i18next из `shared/ui`
- page-specific копии shared markup
- новая палитра вне token contract
- дублирование inline SVG между страницами
- animation без reduced-motion behavior
- кликабельный `div` без keyboard semantics
- production-ссылка на `/ui-kit`

## 10. Definition of done

Компонент готов, когда:

- использует tokens и обе темы
- имеет все необходимые visual/interactive states
- работает мышью, touch и клавиатурой
- показан в UI Kit, если является reusable primitive
- не содержит domain/store/router зависимости в shared layer
- не добавляет непереводимые пользовательские строки
- проходит typecheck, lint, tests и production build
