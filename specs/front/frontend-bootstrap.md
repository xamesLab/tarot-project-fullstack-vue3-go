# Frontend Bootstrap (Skipped Step: Frontend Skeleton)

Дата фиксации: 8 марта 2026.

## 1. Цель
Подготовить рабочий каркас frontend раньше запланированной очереди шагов:
- инициализация проекта Vue 3,
- подключение Pinia и Vue Router,
- базовая структура страниц/компонентов/stores/services.

## 2. Что реализовано

### Infra
- `frontend/package.json`
- `frontend/vite.config.ts`
- `frontend/tsconfig.json`
- `frontend/index.html`
- `frontend/src/main.ts`
- `frontend/src/App.vue`

### Styles
- `src/assets/styles/variables.css`
- `src/assets/styles/base.css`

### Router
- `src/router/index.ts`
- `src/router/guards.ts`

Маршруты:
- `/login`
- `/register`
- `/dashboard`
- `/settings`
- `/admin`
- `/:pathMatch(.*)*`

### Pages
- `src/pages/LoginPage.vue`
- `src/pages/RegisterPage.vue`
- `src/pages/DashboardPage.vue`
- `src/pages/SettingsPage.vue`
- `src/pages/AdminPage.vue`
- `src/pages/NotFoundPage.vue`

### Components
- `src/components/AppHeader.vue`
- `src/components/CardView.vue`
- `src/components/SpreadView.vue`
- `src/components/ThemePicker.vue`
- `src/components/BaseLoader.vue`

### State
- `src/stores/auth.store.ts`
- `src/stores/reading.store.ts`
- `src/stores/settings.store.ts`

### API Services
- `src/services/apiClient.ts`
- `src/services/auth.api.ts`
- `src/services/reading.api.ts`
- `src/services/settings.api.ts`
- `src/services/admin.api.ts`

### Types
- `src/types/auth.ts`
- `src/types/reading.ts`
- `src/types/deck-theme.ts`

## 3. Проверка
- `npm install` выполнен успешно.
- `npm run build` выполнен успешно.
- Сборка `vite` + `vue-tsc` проходит без ошибок.
