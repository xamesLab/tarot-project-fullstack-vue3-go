# Frontend Admin Integration

Дата фиксации: 9 марта 2026.

## 1. Что реализовано
- Страница `/admin` расширена и теперь покрывает:
  - управление темами колод;
  - управление role policies (лимиты генераций по ролям);
  - назначение роли пользователю по логину.

## 2. API-интеграция
- Темы:
  - `GET /api/admin/deck-themes`
  - `POST /api/admin/deck-themes`
- Роли и лимиты:
  - `GET /api/admin/roles`
  - `POST /api/admin/roles`
  - `PATCH /api/admin/users/role`

## 3. Обновленные frontend-файлы
- `frontend/src/services/admin.api.ts`
- `frontend/src/pages/AdminPage.vue`

## 4. UI поведение
- Добавлен единый refresh.
- Добавлены отдельные формы:
  - создание темы;
  - сохранение политики роли (`dailySpreadLimit`, `-1` для unlimited);
  - назначение роли пользователю.
- Добавлены статусные сообщения и обработка ошибок API.

## 5. Проверка
- `npm run build` проходит успешно (`vue-tsc --noEmit` + `vite build`).
