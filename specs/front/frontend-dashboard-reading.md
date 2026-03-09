# Frontend Dashboard Reading Integration

Дата фиксации: 9 марта 2026.

## 1. Что реализовано
- Dashboard интегрирован с backend reading API:
  - `GET /api/reading/today`
  - `POST /api/reading/generate`
- На экране отображаются:
  - `cardOfDay`,
  - `spreadA`,
  - `spreadB`,
  - счетчики лимитов генерации (`used/limit/remaining`),
  - флаг возможности генерации (`canGenerateMore`).
- Добавлена кнопка генерации нового комплекта раскладов.
- Добавлена обработка лимита генерации (`429` -> понятное сообщение).

## 2. Измененные frontend-файлы
- `frontend/src/types/reading.ts`
  - расширен контракт ответа reading API.
- `frontend/src/services/reading.api.ts`
  - добавлен `generateSpreads`.
- `frontend/src/stores/reading.store.ts`
  - добавлено состояние `generating`;
  - добавлен action `generate`.
- `frontend/src/pages/DashboardPage.vue`
  - обновлен UI Dashboard:
    - кнопки reload/generate,
    - блок информации о лимитах,
    - сообщения состояния.

## 3. Проверка
- `npm run build` проходит успешно (`vue-tsc --noEmit` + `vite build`).
