# Nebula Tarot

`Nebula Tarot` — это MVP веб-платформы для ежедневных таро-раскладов с двуязычным интерфейсом, ролевой моделью доступа и административной панелью. Репозиторий находится в активной реализации: backend-ядро и dev-only `UI Kit` уже собраны, а спецификации в `specs/` остаются основным источником проектного контекста.

## Что входит в MVP

- ежедневная `Card of the Day` с детерминированной выдачей и сбросом в `03:00` по локальному часовому поясу пользователя
- случайный `3-card spread` с дневными лимитами по ролям
- регистрация, логин и профиль пользователя
- история раскладов с ограничением срока хранения по ролям
- выбор колоды и рубашки карт
- админ-панель для управления пользователями, колодами и ролевыми политиками
- поддержка двух языков: `English` и `Русский`
- светлая и тёмная темы
- задел под будущую AI-интеграцию через `Gateway stub`

## Роли пользователей

- `Seeker` — базовая роль для новых пользователей
- `Adept` — расширенные лимиты
- `Master` — ещё больше лимитов и истории
- `Keeper` — администраторский уровень с доступом в админ-панель

По умолчанию для MVP заданы такие ограничения для `3-card spread` и истории:

- `Seeker`: 3 расклада в день, история 1 день
- `Adept`: 5 раскладов в день, история 7 дней
- `Master`: 10 раскладов в день, история 15 дней
- `Keeper`: без практического лимита, история 30 дней

## Планируемый стек

### Frontend

- `Vue 3`
- `Vite`
- `Pinia`
- `Vue Router`
- `Axios`
- `Vue I18n`
- `SCSS`
- `VeeValidate`
- `date-fns`

### Backend

- `Go 1.22+`
- `Gin`
- `MongoDB 6.0+`
- `mongo-go-driver`
- `JWT`
- `bcrypt`
- `validator`
- `godotenv`

## Архитектура проекта

Целевая архитектура и текущая реализация организованы так:

- `frontend/` — SPA на `Vue 3`
- `backend/` — API на `Go + Gin`
- `MongoDB` — основное хранилище
- `backend/uploads/` — локальное файловое хранилище для колод и рубашек

Для локальной разработки в корне также есть:

- `start-dev.sh` — последовательный запуск backend, затем frontend
- `stop-dev.sh` — остановка локальных dev-процессов

## Что уже есть в репозитории

- [PRD](./specs/docs/PRD.md) — продуктовые требования, роли, пользовательские сценарии и scope MVP
- [TECHNICAL_SPEC](./specs/docs/TECHNICAL_SPEC.md) — целевая архитектура, API-контракты, схема MongoDB и требования к окружению
- [IMPLEMENTATION_ROADMAP](./specs/docs/IMPLEMENTATION_ROADMAP.md) — поэтапный план реализации на 6-8 недель
- [DESIGN_SYSTEM](./specs/docs/DESIGN_SYSTEM.md) — проектная дизайн-система и карта переиспользуемых frontend-компонентов
- [v0.1.0 plan](./specs/v0.1.0/README.md) — разбиение roadmap на исполнимые этапы
- `specs/temp/mock/tarot-cards.json` — исходные мок-данные карт для сидирования базы

## Текущий статус

Актуальный статус на `April 3, 2026`:

- продуктовые, технические и дизайн-системные документы уже зафиксированы как рабочий baseline
- `phase-0` и основная часть `phase-1` фактически реализованы
- backend уже содержит серверный каркас, auth, профиль пользователя и timezone/accounting-day логику
- frontend уже содержит базовый shell, auth flow, profile/settings forms, i18n/theme infrastructure и dev-only `UI Kit`
- `UI Kit v0.1.0` утверждён как source of truth для дальнейшей frontend-разработки
- `phase-2-core-frontend` завершён как стабильный frontend baseline
- светлая и тёмная темы уже переведены на общий design-system token layer и визуально согласованы между собой
- `phase-3-reading-features` реализован: backend reading API, dashboard, history и tarot-компоненты подключены end-to-end, включая persisted reveal-state карт в раскладах
- ручные тесты `phase-3` пройдены
- `phase-4-admin-panel` реализован на уровне расширенного backend/frontend baseline: keeper-only admin routes, user management, role policies, global settings, независимые каталоги `decks` и `card backs`, global defaults, архив для колод и рубашек и пользовательский override для `adept/master/keeper`
- для `phase-4` ещё не зафиксирована финальная ручная приёмка
- следующий основной шаг — manual acceptance для `phase-4`, затем `phase-5-polish-and-testing`

## Рекомендуемая последовательность старта

Если вы подключаетесь к проекту сейчас, удобнее идти так:

1. Прочитать [PRD](./specs/docs/PRD.md), [TECHNICAL_SPEC](./specs/docs/TECHNICAL_SPEC.md), [IMPLEMENTATION_ROADMAP](./specs/docs/IMPLEMENTATION_ROADMAP.md)
2. Для frontend отдельно изучить [DESIGN_SYSTEM](./specs/docs/DESIGN_SYSTEM.md) и dev-only [UI Kit](./frontend/src/pages/dev/UIKit.vue)
3. Поднять окружение: `Go 1.22+`, `Node.js 20 LTS`, `MongoDB 6.0+`
4. Использовать `./start-dev.sh` для локального запуска стека
5. Реализовывать задачи по фазам из [specs/v0.1.0](./specs/v0.1.0/README.md), начиная с ручной приёмки `phase-4-admin-panel`, затем переходить к `phase-5-polish-and-testing`

## Предполагаемое локальное окружение

В спецификации зафиксированы такие адреса для разработки:

- frontend: `http://localhost:5173`
- backend API: `http://localhost:8080/api`

Для локальной работы также предполагаются:

- self-hosted `MongoDB`
- SSH tunnel для локального доступа к базе
- отдельные `.env` файлы для фронтенда и бэкенда

## Важные особенности доменной логики

- `Card of the Day` зависит от `user_id + accounting_date + timezone`
- расчёт `accounting day` идёт не по полуночи, а с `03:00` до `02:59` следующего дня
- при смене часового пояса доступность дневной карты может измениться
- в новых раскладах карты изначально скрыты рубашкой и открываются по одной с сохранением состояния просмотра в самом `reading`
- история чтений должна удаляться через TTL-логику MongoDB на основе роли пользователя

## Документы, с которых лучше начинать

Если вы впервые заходите в проект, оптимальный порядок чтения такой:

1. [PRD](./specs/docs/PRD.md)
2. [TECHNICAL_SPEC](./specs/docs/TECHNICAL_SPEC.md)
3. [IMPLEMENTATION_ROADMAP](./specs/docs/IMPLEMENTATION_ROADMAP.md)
4. [DESIGN_SYSTEM](./specs/docs/DESIGN_SYSTEM.md)
5. [specs/v0.1.0](./specs/v0.1.0/README.md)

## Ближайшая цель репозитория

Закрыть ручную приёмку `phase-4-admin-panel`, затем перейти к `phase-5-polish-and-testing`, финальной приёмке и deployment.
