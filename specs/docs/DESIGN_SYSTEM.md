# Design System и компонентная архитектура для `Tarot Daily`

> Статус: принят как рабочий reference для `v0.1.0`  
> Source of truth: `frontend/src/pages/dev/UIKit.vue`  
> Контекст: `specs/v0.1.0/phase-2-core-frontend/pre-phase-2-core-frontend/README.md`

## Зачем нужен этот документ

Финальная версия `UI Kit v0.1.0` утверждается как визуальный и компонентный reference для дальнейшей frontend-разработки.  
Этот документ переводит текущую dev-only витрину в инженерную систему:

- определяет слои дизайн-системы
- разбивает `UI Kit` на переиспользуемые примитивы и составные компоненты
- задаёт целевую структуру `frontend/src`
- фиксирует правила переиспользования для `phase-2` и следующих этапов

## 1. Визуальная модель системы

`UI Kit v0.1.0` задаёт не просто набор красивых блоков, а конкретный визуальный язык:

- dark-first интерфейс с золотым акцентом
- высокий контраст и ритуальная, ornamental эстетика
- display-типографика для смысловых заголовков
- спокойная sans/body-типографика для формы и контента
- прямолинейная геометрия поверхностей без сложных вырезов
- визуальное разделение на:
  - фон
  - поверхность
  - акцент
  - контент
  - сакральные/орнаментальные маркеры

## 2. Слои дизайн-системы

Система должна строиться в 4 слоя.

### Layer A. Foundations

Это глобальные токены, которые не зависят от конкретного компонента:

- colors
- typography
- spacing
- sizing
- borders
- shadows
- motion
- z-index
- iconography / ornamental symbols

### Layer B. Primitives

Это самые маленькие визуальные блоки, не привязанные к бизнес-сущностям:

- `AppContainer`
- `Surface`
- `SectionHeader`
- `Divider`
- `Badge`
- `Tag`
- `IconSymbol`
- `Avatar`
- `Button`
- `Input`
- `Select`
- `Checkbox`
- `Radio`
- `Toggle`
- `SegmentedToggle`
- `Range`
- `Tabs`

### Layer C. Composites

Это составные UI-узлы, собираемые из примитивов:

- `AppHeader`
- `AppFooter`
- `UserChip`
- `ThemeLanguageControls`
- `TarotCardPreview`
- `ReadingCard`
- `FormField`
- `FormGroup`
- `SettingsControlGroup`
- `ModalDialog`
- `EmptyState`
- `LoadingPanel`

### Layer D. Screens / Feature Blocks

Это страницы и фичевые контейнеры:

- `LoginPage`
- `RegisterPage`
- `DashboardPage`
- `ProfilePage`
- `HistoryPage`
- `Admin*`

Правило: экран не должен реализовывать собственный визуальный язык. Он только композирует primitives и composites из design system.

## 3. Design tokens

На `April 1, 2026` token-layer уже вынесен из `UIKit.vue` в глобальные стили проекта: `frontend/src/styles/tokens.scss`, `frontend/src/styles/themes/_light.scss`, `frontend/src/styles/themes/_dark.scss` и `frontend/src/styles/components/*.scss`.

Текущий практический статус такой:

- app shell и базовые primitives уже используют общий design-system слой
- `UI Kit` остаётся visual regression reference и showcase для тех же primitives
- светлая тема больше не является побочным вариантом dark theme, а приведена к той же ornamental системе через parchment/brass palette

Дальнейшая задача этапов `phase-3+` не в повторной миграции токенов, а в расширении semantic roles без появления второго визуального языка.

### 3.1 Цветовые токены

Базовые semantic tokens:

- `--ds-color-bg-page`
- `--ds-color-bg-page-alt`
- `--ds-color-bg-surface`
- `--ds-color-bg-surface-strong`
- `--ds-color-bg-elevated`
- `--ds-color-text-primary`
- `--ds-color-text-secondary`
- `--ds-color-text-muted`
- `--ds-color-accent-primary`
- `--ds-color-accent-light`
- `--ds-color-accent-dark`
- `--ds-color-border-subtle`
- `--ds-color-border-accent`
- `--ds-color-danger`
- `--ds-color-overlay`

Рекомендованное соответствие текущему `UI Kit`:

| Token | Value / source |
| --- | --- |
| `--ds-color-bg-page` | `#000000` |
| `--ds-color-bg-page-alt` | `#0D0D0D` |
| `--ds-color-bg-surface` | `#1A1A1A` |
| `--ds-color-bg-surface-strong` | `#111111` или ближайшая усиленная поверхность |
| `--ds-color-bg-elevated` | текущий elevated input/control background |
| `--ds-color-text-primary` | `#FFFFFF` |
| `--ds-color-text-secondary` | `#C4C4C4` |
| `--ds-color-text-muted` | `#6B6B6B` |
| `--ds-color-accent-primary` | `#D4AF37` |
| `--ds-color-accent-light` | `#F4E5B0` |
| `--ds-color-accent-dark` | `#9B7E1F` |
| `--ds-color-border-subtle` | текущий `var(--border-subtle)` |
| `--ds-color-border-accent` | текущий `var(--border-gold)` |
| `--ds-color-danger` | `#ff4444` |
| `--ds-color-overlay` | `rgba(0, 0, 0, 0.72)` |

### 3.2 Typography tokens

Нужны 3 семейства:

- `--ds-font-display`
- `--ds-font-body`
- `--ds-font-ui`

Типографические роли:

- `--ds-text-display-1`
- `--ds-text-display-2`
- `--ds-text-heading-1`
- `--ds-text-heading-2`
- `--ds-text-heading-3`
- `--ds-text-body-lg`
- `--ds-text-body-md`
- `--ds-text-body-sm`
- `--ds-text-label`
- `--ds-text-caption`

Правило использования:

- display/headings только для hero, section titles, tarot content headings
- body для интерпретаций, описаний, крупных текстов
- ui font для labels, controls, tabs, buttons, badges

### 3.3 Spacing tokens

Шкала должна быть простой и короткой:

- `--ds-space-2 = 8px`
- `--ds-space-3 = 12px`
- `--ds-space-4 = 16px`
- `--ds-space-5 = 20px`
- `--ds-space-6 = 24px`
- `--ds-space-8 = 32px`
- `--ds-space-10 = 40px`
- `--ds-space-12 = 48px`
- `--ds-space-15 = 60px`

### 3.4 Radius и stroke

Так как текущий UI кит уже ушёл от `clip-path`, форма должна держаться на границах и сетке:

- `--ds-radius-none = 0`
- `--ds-border-width = 1px`
- `--ds-border-width-strong = 2px`

Если позднее потребуется скругление, его нужно вводить отдельно и осознанно. Для `v0.1.0` базовая форма остаётся геометрически жёсткой.

### 3.5 Motion tokens

- `--ds-motion-fast = 160ms ease`
- `--ds-motion-base = 300ms ease`
- `--ds-motion-slow = 420ms ease`

Анимация допустима только для:

- hover/focus
- modal open/close
- tab content transition
- toggle thumb movement
- soft reveal контентных блоков

## 4. Базовые элементы из `UI Kit`

Ниже разбивка текущего `UI Kit` на инженерные примитивы.

### 4.1 Layout primitives

#### `AppContainer`

Назначение:

- ограничивает ширину контента
- задаёт горизонтальные отступы
- используется в header, page body, footer

Источник в `UI Kit`:

- `.container`
- `.header-container`

#### `Surface`

Назначение:

- базовая поверхность для карточек, контрол-групп, модалок, tabs-container

Варианты:

- `surface`
- `surface-elevated`
- `surface-ghost`
- `surface-bordered`

Источник:

- `.card`
- `.control-group`
- `.tabs-container`
- `.user-profile`
- `.modal`

#### `SectionHeader`

Назначение:

- единый заголовок секции

Состав:

- `eyebrow` по необходимости
- `title`
- optional `subtitle` / `actions`

Источник:

- `.section-header`
- `.section-title`
- hero header

#### `Divider`

Варианты:

- `simple`
- `ornamental`
- `footer`

Источник:

- `.ornamental-divider`
- `.footer-divider`

#### `IconSymbol`

Назначение:

- декоративный символический элемент

Допустимые значения в `v0.1.0`:

- `star`
- `moon`
- `diamond`
- `circle`

Важно: это не general icon library. Это ограниченный ornamental set.

### 4.2 Typography primitives

#### `DisplayText`

Для hero и больших смысловых заголовков.

#### `Heading`

Варианты:

- `h1`
- `h2`
- `h3`
- `h4`

#### `BodyText`

Варианты:

- `lg`
- `md`
- `sm`

#### `LabelText`

Для:

- field labels
- tabs
- badge text
- button text
- metadata

### 4.3 Interactive primitives

#### `Button`

Это один компонент с пропсами:

- `variant: primary | secondary | ghost | outlined | text`
- `size: sm | md | lg`
- `fullWidth`
- `iconOnly`
- `leadingIcon`
- `trailingIcon`
- `starSides`
- `loading`
- `disabled`

В `UI Kit` это уже видно в:

- `.btn-primary`
- `.btn-secondary`
- `.btn-ghost`
- `.btn-outlined`
- `.btn-text`
- `.btn-small`
- `.btn-large`
- `.btn-full`
- `.btn-icon`
- `.btn-icon-text`
- `.btn-star-sides`

Вывод: не нужно плодить отдельные компоненты `PrimaryButton`, `GhostButton`, `TextButton`. Нужен один `BaseButton`.

#### `Input`

Пропсы:

- `modelValue`
- `label`
- `placeholder`
- `type`
- `error`
- `hint`
- `disabled`

Варианты:

- `text`
- `email`
- `password`

#### `Select`

Пропсы:

- `modelValue`
- `label`
- `options`
- `placeholder`
- `error`
- `disabled`

#### `Checkbox`

Использовать как controlled component с label slot.

#### `RadioGroup`

Нужны:

- `RadioGroup`
- `RadioOption`

Не отдельные разрозненные радио-кнопки по месту.

#### `Toggle`

В проекте нужны 2 типа toggle.

##### `SwitchToggle`

Обычный on/off:

- dark mode in settings
- sound
- show hints
- notifications

Источник:

- `.toggle-switch`

##### `SegmentedToggle`

Бинарный segmented control:

- `EN/RU`
- `Light/Dark`
- в будущем возможны `Grid/List`, `Yes/No`

Источник:

- `.header-binary-toggle`
- блок `Final Toggles`

Это важный компонент. Его нужно вынести отдельно, не оставлять как private-стиль шапки.

#### `RangeSlider`

Для настройки интенсивности/размеров/прозрачности.

Пропсы:

- `modelValue`
- `min`
- `max`
- `step`
- `label`
- `showValue`

#### `Tabs`

Нужен один базовый компонент `Tabs` с variant API:

- `variant: underline | pills | boxed`

Состав:

- `Tabs`
- `TabsList`
- `TabsTrigger`
- `TabsPanel`

Источник:

- `.tabs-header`
- `.tab-button`
- `.tabs-pills`
- `.tabs-boxed`

### 4.4 Informational primitives

#### `Badge`

Для небольшого статуса:

- `Today`
- `3 of 5`
- `Major`
- role/status labels

#### `Tag`

Для множественных тематических меток:

- `Freedom`
- `Adventure`
- `Intuition`

#### `Avatar`

Минимальный reusable primitive для пользователя.

#### `UserChip`

Составной элемент из:

- `Avatar`
- `user name`
- `secondary meta`

Источник:

- `.user-profile`

## 5. Составные компоненты проекта

### 5.1 Layout-компоненты

#### `AppHeader`

Состав:

- `BrandMark`
- `ThemeLanguageControls`
- `UserChip` или auth actions

Режимы:

- `public`
- `authenticated`
- `admin`

#### `AppFooter`

Состав:

- ornamental divider
- supporting copy
- legal/meta text

#### `PageSection`

Нужен reusable section wrapper с:

- title
- optional subtitle
- optional actions
- content slot

### 5.2 Auth и form-компоненты

#### `FormField`

Объединяет:

- `label`
- control
- hint
- error

Это лучше, чем повторять label/error руками в каждой форме.

#### `AuthCard`

Использовать на `Login` и `Register`.

Состав:

- `Surface`
- `SectionHeader`
- form content
- actions/footer links

### 5.3 Контентные tarot-компоненты

#### `TarotCardPreview`

Состав:

- surface frame
- back ornament
- symbol
- optional face image later

Это отдельный reusable visual primitive, потому что он встретится:

- в card of the day
- в spread preview
- в history
- в admin deck management

#### `ReadingCard`

Состав:

- `Surface`
- `Badge`
- `TarotCardPreview`
- `Heading`
- `BodyText`
- `Tag` list
- action buttons

#### `SpreadSummaryCard`

Для dashboard/history. По сути расширение `ReadingCard`, а не новый визуальный стиль.

### 5.4 Feedback-компоненты

#### `ModalDialog`

Состав:

- overlay
- surface
- title
- body
- actions

Варианты:

- `default`
- `confirm`
- `danger`

#### `EmptyState`

Для первого визита, пустой истории, отсутствия данных.

Состав:

- ornamental icon
- title
- description
- primary action
- optional secondary action

#### `LoadingPanel`

Не generic spinner-only блок. В этом проекте лучше использовать surface-based skeleton/loading card.

## 6. Карта компонентов для `phase-2`

Минимальный обязательный набор компонентов для реализации `phase-2-core-frontend`:

### Base

- `BaseButton.vue`
- `BaseInput.vue`
- `BaseSelect.vue`
- `BaseCheckbox.vue`
- `BaseRadioGroup.vue`
- `BaseSwitchToggle.vue`
- `BaseSegmentedToggle.vue`
- `BaseTabs.vue`
- `BaseBadge.vue`
- `BaseTag.vue`
- `BaseAvatar.vue`
- `BaseSurface.vue`
- `BaseDivider.vue`
- `BaseModal.vue`

### Form

- `FormField.vue`
- `AuthCard.vue`

### Layout

- `AppHeader.vue`
- `AppFooter.vue`
- `PageSection.vue`
- `PageContainer.vue`

### Domain

- `UserChip.vue`
- `TarotCardPreview.vue`
- `ReadingCard.vue`
- `EmptyStateCard.vue`
- `LoadingCard.vue`

## 7. Рекомендуемая структура `frontend/src`

```text
frontend/src/
  components/
    base/
      BaseAvatar.vue
      BaseBadge.vue
      BaseButton.vue
      BaseCheckbox.vue
      BaseDivider.vue
      BaseInput.vue
      BaseModal.vue
      BaseRadioGroup.vue
      BaseSegmentedToggle.vue
      BaseSelect.vue
      BaseSurface.vue
      BaseSwitchToggle.vue
      BaseTabs.vue
      BaseTag.vue
    form/
      AuthCard.vue
      FormField.vue
    layout/
      AppFooter.vue
      AppHeader.vue
      PageContainer.vue
      PageSection.vue
    domain/
      ReadingCard.vue
      TarotCardPreview.vue
      UserChip.vue
    feedback/
      EmptyStateCard.vue
      LoadingCard.vue
  composables/
    useTheme.js
    useLocale.js
  styles/
    tokens.scss
    themes/
      _dark.scss
      _light.scss
    foundations/
      _reset.scss
      _typography.scss
      _motion.scss
    components/
      _button.scss
      _field.scss
      _surface.scss
      _tabs.scss
```

## 8. Правила переиспользования

### 8.1 Что разрешено

- собирать страницы только из базовых и составных компонентов
- добавлять новые variants в существующие компоненты
- расширять tokens и semantic aliases без ломки существующих ролей

### 8.2 Что запрещено

- рисовать новый локальный стиль кнопок внутри страницы
- добавлять page-specific инпуты вне `BaseInput`/`FormField`
- дублировать tab/toggle/button разметку вручную
- создавать второй набор theme variables вне design-system слоя
- оставлять ornamental стили только внутри `UIKit.vue`, если они уже стали проектным стандартом

## 9. Матрица состояний

Для каждого интерактивного компонента обязателен одинаковый state-contract:

- `default`
- `hover`
- `focus-visible`
- `active`
- `disabled`
- `error`, где применимо
- `loading`, где применимо

Первой волной это обязательно для:

- `BaseButton`
- `BaseInput`
- `BaseSelect`
- `BaseSwitchToggle`
- `BaseSegmentedToggle`
- `BaseTabs`
- `BaseModal`

## 10. Миграция от текущего состояния

Основной migration risk, существовавший в начале `phase-2`, в базовом объёме уже закрыт:

- tokens и theme variables вынесены в общий слой
- app shell переведён на reusable primitives
- `UI Kit` показывает phase-2 base components как часть общей системы

Оставшийся риск для следующих этапов уже другой: не откатиться обратно к page-specific стилям по мере добавления tarot features и admin screens.

### Практическое продолжение после миграции

1. Добавлять tarot/admin UI только поверх существующих primitives и composites.
2. Расширять semantic tokens, если появятся новые состояния или сущности.
3. Держать `UI Kit` и реальные страницы синхронизированными.
4. Не возвращаться к локальным page-level кнопкам, полям, табам и toggle-паттернам.

Ключевой принцип: `UI Kit` должен стать витриной системы, а не её отдельной реализацией.

## 11. Практическое решение для ближайших этапов

Для `phase-2-core-frontend` предлагается такой порядок реализации:

1. `tokens + themes`
2. `BaseSurface`, `BaseButton`, `BaseInput`, `BaseSelect`
3. `BaseSwitchToggle`, `BaseSegmentedToggle`, `BaseTabs`
4. `FormField`, `AuthCard`
5. `AppHeader`, `AppFooter`, `PageContainer`, `PageSection`
6. `UserChip`, `TarotCardPreview`, `ReadingCard`
7. затем только страницы `Login`, `Register`, `Dashboard`

Такой порядок уменьшает дублирование и даёт быстрый выход на рабочий `phase-2`.

## 12. Definition of done для дизайн-системы `v0.1.0`

Система считается реально внедрённой, когда:

- `UI Kit` больше не хранит уникальные production-стили только внутри себя
- базовые токены вынесены в глобальные стили
- страницы auth/dashboard используют единые компоненты
- theme/language switchers используют `BaseSegmentedToggle`
- dashboard cards используют `BaseSurface` + `ReadingCard`
- новый UI в проекте нельзя сделать, не опираясь на существующие primitives

На `April 1, 2026` этот baseline для `phase-2` считается достигнутым. Следующие этапы должны развивать систему, а не переизобретать её.

## 13. Главный вывод

`UI Kit v0.1.0` уже достаточно зрелый, чтобы быть не просто dev-страницей, а архитектурной основой frontend.

Инженерное решение для проекта:

- один набор tokens
- один набор base-компонентов
- один набор layout/composite-компонентов
- `UI Kit` как showcase и regression reference
- все следующие страницы строятся только поверх этой системы
