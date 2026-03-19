# Tarot Daily (MVP)

Монорепозиторий проекта Tarot Daily.

## Репозиторий
- `specs/` - требования, архитектура, план работ.
- `frontend/` - клиентское приложение (Vue 3 + Pinia + Router).
- `backend/` - серверное API (Go + Gin + MongoDB).
- `deploy/` - локальная инфраструктура и dev-скрипты запуска.

## Статус
- Шаг 1 (Scope Freeze): выполнен.
- Шаг 2 (структура репозитория и базовая документация): выполнен.
- Шаг 3 (локальное окружение разработки): выполнен.
- Шаг 4 (конфигурация через `.env` и `.env.example`): выполнен.
- Шаг 5 (модель данных MongoDB): выполнен.
- Шаг 6 (индексы MongoDB): выполнен.
- Шаг 7 (seed-данные): выполнен.
- Шаг 8 (базовый каркас backend): выполнен.
- Шаг 9 (auth-модуль): выполнен.
- Шаг 10 (роли и проверка прав доступа): выполнен.
- Backend role policies (назначение ролей и лимиты генерации): выполнен.
- Шаг 11 (настройки пользователя: timezone + deck theme): выполнен.
- Шаг 12 (endpoint списка доступных тем колод): выполнен.
- Шаг 13 (ядро генерации ежедневного расклада): выполнен.
- Шаг 14 (`GET /api/reading/today`): выполнен.
- Шаг 15 (backend admin API тем колод): выполнен.
- Шаг frontend bootstrap (выполнен досрочно): Vue 3 + Pinia + Router + базовый layout.
- Шаг frontend Dashboard (интеграция reading API): выполнен.
- Шаг frontend Settings (таймзона + выбор темы): выполнен.
- Шаг frontend Admin (темы + роли + лимиты): выполнен.
- Шаг интеграции frontend/backend (контракты + smoke): выполнен.
- Шаг UI states (errors/loading/empty): выполнен.
- Шаг unit + smoke тестирования: выполнен.
- Шаг стабилизации релиз-кандидата: выполнен.
- Шаг финальной документации запуска/использования: выполнен.
- Шаг ручной mobile-first проверки: выполнен.
- Шаг проверки acceptance criteria: выполнен.
- Финальный статус MVP: зафиксирован.

## Post-MVP статус
- Этап `v0.1.1`: выполнен.
- Итоги:
  - `specs/version-0.1.1/final-status-v0.1.1.md`
  - `specs/version-0.1.1/testing-report-v0.1.1.md`
  - `specs/version-0.1.1/acceptance-criteria-check-v0.1.1.md`
- Этап `v0.1.2`: запущен (`in progress`).
- Документы этапа:
  - `specs/version-0.1.2/requirements-v0.1.2.md`
  - `specs/version-0.1.2/architecture-v0.1.2.md`
  - `specs/version-0.1.2/steps.md`
  - `specs/version-0.1.2/testing-report-v0.1.2.md`
  - `specs/version-0.1.2/acceptance-criteria-check-v0.1.2.md`
  - `specs/version-0.1.2/final-status-v0.1.2.md`

## Текущее ограничение MVP
Границы MVP и критерии приемки зафиксированы в:
- `specs/main/mvp-scope-freeze.md`
- `specs/back/mongodb-data-model.md`
- `specs/back/mongodb-indexes.md`
- `specs/back/seed-strategy.md`
- `specs/back/backend-bootstrap.md`
- `specs/back/auth-module.md`
- `specs/back/backend-roles-access.md`
- `specs/back/backend-role-policies.md`
- `specs/back/backend-user-settings.md`
- `specs/back/backend-deck-themes-list.md`
- `specs/back/backend-reading-engine.md`
- `specs/back/backend-admin-deck-themes.md`
- `specs/front/frontend-bootstrap.md`
- `specs/front/frontend-auth-routing.md`
- `specs/front/frontend-dashboard-reading.md`
- `specs/front/frontend-settings.md`
- `specs/front/frontend-admin.md`
- `specs/front/frontend-backend-integration-smoke.md`
- `specs/front/frontend-ui-states.md`
- `specs/version-0.0.1/testing-report.md`
- `specs/main/final-runbook.md`
- `specs/main/acceptance-criteria-check.md`
- `specs/main/mobile-first-manual-check.md`
- `specs/main/mvp-final-status.md`

## Ближайшие шаги
1. Выполнять чек-лист `specs/version-0.1.2/steps.md`.

