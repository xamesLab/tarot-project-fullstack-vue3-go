# Frontend/Backend Integration Smoke

Дата фиксации: 9 марта 2026.

## 1. Цель
Подтвердить рабочую интеграцию frontend-контрактов с backend API по критичному user/admin flow.

## 2. Проверенные сценарии
- Регистрация и `GET /api/me`.
- `GET /api/deck-themes`.
- `PATCH /api/me/settings` (таймзона + выбор темы).
- `GET /api/reading/today` (стабильная `cardOfDay` в рамках суток).
- `POST /api/reading/generate`.
- Admin:
  - `POST /api/admin/roles`,
  - `PATCH /api/admin/users/role`,
  - `POST /api/admin/deck-themes`.

## 3. Результаты smoke
- `USER_REGISTER_OK=it1773007060123`
- `SETTINGS_SAVED_TZ=UTC+03:00`
- `CARD_STABLE=True`
- `LIMIT_USER=10`
- `ASSIGNED_ROLE=vip`
- `LIMIT_AFTER_ASSIGN=25`
- `ADMIN_THEME_CREATE=Smoke 1773007064318`

## 4. Дополнительно
- Frontend type-check + build успешно:
  - `npm run build` (`vue-tsc --noEmit`, `vite build`).
