# Final Status v0.1.1

Дата фиксации: 9 марта 2026.

## Итог
- Статус этапа: `v0.1.1 завершен`.
- Шаги чек-листа: `21/21` (`specs/version-0.1.1/steps.md`).

## Что подтверждено
- Functional scope `v0.1.1` реализован:
  - locale/theme режимы;
  - изолированная админка (`users/decks/settings`);
  - профиль пользователя (`name/timezone/preferredDeckId`);
  - каталог карт и маппинг изображений;
  - card flip UI.
- API smoke `v0.1.1`: `PASS`.
- Frontend checks `v0.1.1`: `PASS`.
- Backend/Frontend build checks: `PASS`.
- После функционального завершения `v0.1.1` выполнен структурный рефакторинг backend без изменения целевого scope:
  - вынесены общие HTTP helper-ы;
  - разделены router/bootstrap/routes;
  - декомпозированы пакеты `reading`, `decks`, `auth`;
  - вынесены shared helper-ы `user_helpers`, `inpututil`, `envfile`;
  - добавлены unit-тесты на вынесенную чистую логику.

## Сопутствующие документы
- Требования: `specs/version-0.1.1/requirements-v0.1.1.md`
- Архитектура: `specs/version-0.1.1/architecture-v0.1.1.md`
- Отчет проверок: `specs/version-0.1.1/testing-report-v0.1.1.md`
- Критерии приемки: `specs/version-0.1.1/acceptance-criteria-check-v0.1.1.md`
- Журнал рефакторинга: `specs/version-0.1.1/refactor-log.md`
