# Frontend v0.1.3: UI Kit, Background, UX Consistency

**Версия:** v0.1.3  
**Дата:** 21 марта 2026  
**Статус:** Completed

---

## 1. Контекст

В рамках `v0.1.3` на фронтенде выполнен комплекс доработок:
1. Dev-only страница `/ui-kit` как эталон визуальной системы.
2. Фоновая система Arcana через отдельный компонент `AppBackground.vue`.
3. Приведение `mode-switcher` к стилю UI Kit (section 6).
4. Улучшение мобильной гармонии UI Kit (типографика, ритм, карточки, таблицы).
5. Закрепление правил для агентов: `UI Kit First` и `Background Freeze Rule`.

---

## 2. Цели

1. Поддерживать единый визуальный стандарт через `/ui-kit`.
2. Сохранить безопасную изоляцию внутренних страниц (dev-only маршруты).
3. Сделать фон более атмосферным, но не мешающим читаемости и взаимодействию.
4. Улучшить восприятие UI Kit на мобильных экранах.
5. Зафиксировать проектные правила для дальнейшей разработки.

---

## 3. Scope

В рамках задачи:
- роутинг dev-only (`/ui-kit`);
- standalone страница `frontend/src/pages/dev/UIKit.vue`;
- фоновой слой `frontend/src/components/AppBackground.vue`;
- адаптация `ModeSwitcher`/`ModeToggle` под UI Kit;
- mobile-доработки UI Kit;
- документация и агентские правила.

Вне рамок:
- полная миграция всего продуктового интерфейса на UI Kit;
- e2e-набор по визуальным регрессиям;
- переработка backend API.

---

## 4. Реализованные требования

### 4.1 Dev-only UI Kit

1. Маршрут `/ui-kit` доступен только в dev.
2. В production маршрут не регистрируется.
3. Guard дополнительно блокирует `meta.devOnly` вне dev.

### 4.2 Arcana Background

1. Добавлен фиксированный декоративный фоновой слой `AppBackground.vue`.
2. Фон не перехватывает взаимодействие (`pointer-events: none`).
3. Поддерживаются light/dark режимы через `isLight` prop из `App.vue`.
4. Декоративные элементы: blooms, mandala, constellations, glyphs, frame, vignette.

### 4.3 Mode Switcher Consistency

1. Языковой переключатель — segmented control в стиле UI Kit section 6.
2. Theme toggle — круглая кнопка 34x34 в стиле UI Kit section 6.
3. Символы взяты из UI Kit; `Icon`-компонент сохранен как fallback для SVG.

### 4.4 Mobile Harmonization (UI Kit)

1. Добавлены дополнительные адаптивные корректировки (`<=760px`, `<=480px`).
2. Ключевые размеры переведены на `clamp(...)`.
3. Улучшены ритм отступов, масштаб типографики, размеры карточек и таблиц на мобильных.

### 4.5 Agent Governance

1. Добавлены правила:
   - `frontend/UI_KIT_RULE.md`
   - `AGENTS.md`
   - `frontend/AGENTS.md`
2. Закреплено обязательное правило `UI Kit First`.
3. Закреплено правило `Background Freeze Rule` (фон не трогать без явного запроса пользователя).

---

## 5. Техническая реализация

Роутинг и guards:
- `frontend/src/router/dev-tools.routes.ts`
- `frontend/src/router/index.ts`
- `frontend/src/router/guards.ts`

UI Kit и адаптивность:
- `frontend/src/pages/dev/UIKit.vue`

Фоновая система:
- `frontend/src/components/AppBackground.vue`
- `frontend/src/App.vue`
- `frontend/src/assets/styles/base.css` (очистка конфликтующих старых фоновых слоев)

Mode switcher:
- `frontend/src/components/ModeSwitcher.vue`
- `frontend/src/components/ModeToggle.vue`

Документация/правила:
- `frontend/README.md`
- `frontend/UI_KIT_RULE.md`
- `AGENTS.md`
- `frontend/AGENTS.md`
- `specs/version-0.1.3/front/requirements.md`
- `specs/version-0.1.3/front/step-by-step.md`

---

## 6. Критерии приемки

- [x] `/ui-kit` доступен в dev и закрыт в production.
- [x] UI Kit реализован как эталонная страница с 11 секциями.
- [x] Фон подключен отдельным компонентом и не мешает UX.
- [x] `mode-switcher` приведен к визуальному языку UI Kit.
- [x] Мобильная версия UI Kit выглядит гармонично на узких экранах.
- [x] Добавлены и задокументированы агентские правила UI/фон.
- [x] TypeScript и production build проходят.

---

## 7. Проверки

Выполнено:
1. `npm run build` — успешно.
2. Проверка production bundle — dev-route `/ui-kit` не попадает в prod-регистрацию.

---

## 8. Будущее развитие

1. Добавить dev-index (`/dev`) как каталог внутренних страниц.
2. Добавить визуальные snapshot/e2e тесты для UI Kit.
3. Постепенно выравнивать продуктовые страницы по токенам и паттернам UI Kit.
