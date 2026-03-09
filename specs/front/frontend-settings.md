# Frontend Settings Integration

Дата фиксации: 9 марта 2026.

## 1. Что реализовано
- Страница `Settings` теперь:
  - инициализируется из текущего профиля пользователя (`authStore.user`);
  - загружает доступные темы колоды;
  - сохраняет `timezone` и `deckThemeId` через `PATCH /api/me/settings`.
- После сохранения обновляется профиль в `authStore`, чтобы новые настройки сразу использовались в приложении.

## 2. Важная правка по сбросу темы
- При выборе `Default` в UI теперь отправляется `deckThemeId: ""`.
- Это корректно сбрасывает `deckThemeId` на backend.
- Ранее отправлялось `undefined`, из-за чего backend не менял поле.

## 3. Измененные файлы
- `frontend/src/services/settings.api.ts`
- `frontend/src/stores/settings.store.ts`
- `frontend/src/stores/auth.store.ts`
- `frontend/src/pages/SettingsPage.vue`

## 4. Проверка
- `npm run build` проходит успешно (`vue-tsc --noEmit` + `vite build`).
