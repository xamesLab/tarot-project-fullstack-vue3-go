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
- Этап `v0.1.2` (backend Gateway integration): реализация завершена, валидация зафиксирована как `validation in progress`.
- Документы этапа:
  - `specs/version-0.1.2/requirements-v0.1.2.md`
  - `specs/version-0.1.2/architecture-v0.1.2.md`
  - `specs/version-0.1.2/steps.md`
  - `specs/version-0.1.2/testing-report-v0.1.2.md`
  - `specs/version-0.1.2/acceptance-criteria-check-v0.1.2.md`
  - `specs/version-0.1.2/final-status-v0.1.2.md`
- Этап `v0.1.3` (frontend UI Kit + Arcana background): выполнен.
- Документы этапа:
  - `specs/version-0.1.3/front/requirements.md`
  - `specs/version-0.1.3/front/step-by-step.md`
  - `specs/version-0.1.3/front/arcana-uikit-prompt.md`
  - `specs/version-0.1.3/front/arcana-background-prompt.md`
- Этап `v0.1.4`: запланирован как ближайший этап разработки.
- Документы этапа:
  - `specs/version-0.1.4/requirements-v0.1.4.md`
  - `specs/version-0.1.4/steps.md`

## Актуальный статус
- Базовый MVP: завершен и документально зафиксирован.
- `v0.1.2`: backend-реализация завершена; остается фактический интеграционный smoke ask-flow в окружении с доступным Gateway и callback security конфигурацией.
- `v0.1.3`: frontend-доработки по UI Kit, Arcana background и UX consistency завершены.
- `v0.1.4`: следующий рабочий этап, в который входят закрытие валидации `v0.1.2`, dev-tools `/dev` и cleanup smoke-данных.
- Все задачи, связанные с автотестами, отложены до отдельного явного запроса пользователя.
- Все задачи, связанные с файловым хранилищем / внешним storage для assets, отложены до отдельного явного запроса пользователя.
- Для UI-изменений действует правило `UI Kit First`, для фоновой системы действует `Background Freeze Rule`.

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
1. Выполнять чек-лист `specs/version-0.1.4/steps.md`.
2. Использовать `specs/version-0.1.4/requirements-v0.1.4.md` как актуальный scope ближайшего этапа.
3. Использовать `specs/main/next-development-stages.md` как верхнеуровневую сводку roadmap.
4. Не возвращаться к автотестам и файловому storage без отдельного явного запроса пользователя.

## Frontend UI правило
- Для компонентов, верстки и стилей действует правило `UI Kit First`: ориентироваться на `/ui-kit`.
- Документ правила: `frontend/UI_KIT_RULE.md`.

