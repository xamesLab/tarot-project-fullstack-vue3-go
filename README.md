# Nebula Tarot

`Nebula Tarot` — full-stack tarot-платформа с Go API, MongoDB и двумя поколениями frontend. Актуальная frontend-итерация — `v0.2.0`: новый публичный Card of the Day experience на React и TypeScript. Завершённый Vue-клиент `v0.1.x` сохранён в `frontend-legacy/` как функциональный full-stack baseline.

## Структура репозитория

| Путь | Назначение | Статус |
|---|---|---|
| `frontend/` | Новый SPA на React + TypeScript | Актуальная разработка `v0.2.0` |
| `frontend-legacy/` | Предыдущий SPA на Vue 3 | Legacy baseline `v0.1.x` |
| `backend/` | REST API на Go + Gin, MongoDB, JWT, OpenRouter | Действующий backend |
| `specs/` | Общие спецификации и документы по версиям | [Индекс документации](./specs/README.md) |

## Что нового в v0.2.0

- новый визуальный язык `NEBULA-TAROT`, перенесённый из `Arcana Home.html` и `Arcana UI Library.html`
- публичная страница `/` с картой дня, 3D-переворотом карты, толкованием по четырём сферам и посланием дня
- mock oracle chat с асинхронными ответами
- регистрационный modal для закрытых разделов с локальной валидацией и mock submit
- RU/EN и dark/light theme с сохранением выбора в `localStorage`
- dev-only `/ui-kit`, собранный из тех же shared-компонентов, что и Home
- feature-oriented архитектура: `app`, `api`, `features`, `shared`, `widgets`, `pages`
- типизированная граница API с переключателем `VITE_API_MODE=mock|http`

Полное описание изменений и границ версии: [specs/v0.2.0](./specs/v0.2.0/README.md).

## Стек

### Актуальный frontend

- React 18, TypeScript 5, Vite 5
- React Router, Zustand, TanStack Query
- i18next + react-i18next
- react-hook-form + zod
- Tailwind CSS, CSS custom properties и компонентные CSS-стили
- Vitest + Testing Library

### Backend

- Go 1.23, Gin
- MongoDB и официальный Go driver
- JWT, bcrypt, validator
- direct AI integration через OpenRouter

### Legacy frontend

- Vue 3, Pinia, Vue Router, Vue I18n
- Axios, VeeValidate, SCSS

## Быстрый старт нового frontend

Требования: Node.js 20 LTS и npm.

```bash
cd frontend
npm install
npm run dev
```

Откройте `http://localhost:5173/`. В development также доступен `http://localhost:5173/ui-kit`; в production route не регистрируется. Файлы `.env.development` и `.env.production` уже выбирают mock API.

Проверки:

```bash
npm run typecheck
npm run lint
npm test
npm run build
```

## Запуск backend и legacy-клиента

Этот режим нужен для существующих auth, profile, readings, history и admin flow, пока HTTP-адаптер нового frontend не реализован.

```bash
cp backend/.env.example backend/.env
cd backend
go mod download
go run .
```

В другом терминале:

```bash
cp frontend-legacy/.env.example frontend-legacy/.env
cd frontend-legacy
npm install
npm run dev
```

Backend по умолчанию доступен на `http://localhost:8888/api`, оба frontend используют порт `5173`, поэтому одновременно их следует запускать на разных портах.

Корневой `start-dev.sh` по-прежнему запускает `backend/` вместе с актуальным `frontend/`, но новый клиент останется в mock mode и не будет обращаться к поднятому API. Для проверки прежнего end-to-end flow используйте `frontend-legacy/`.

## API-режим нового frontend

`frontend/src/api/index.ts` выбирает реализацию по `VITE_API_MODE`:

- `mock` — рабочий режим `v0.2.0`; данные карты дня, oracle и регистрация обслуживаются локальными typed fixtures
- `http` — зарезервирован, но пока намеренно возвращает ошибку `HTTP API is not implemented`

Новый клиент поэтому можно разрабатывать автономно. Интеграция с Go API — следующий отдельный этап; наличие backend в репозитории не означает, что `frontend/` уже использует его.

## Текущий статус версий

- `v0.1.0` — завершённый MVP baseline: auth, профиль, карта дня, 3-card spread, история, каталоги и admin panel
- `v0.1.1` — direct AI flow через OpenRouter и административные настройки AI
- `v0.2.0` — новая frontend-архитектура и публичный Card of the Day experience; работает через mock API

## Документация

Начинать рекомендуется в таком порядке:

1. [v0.2.0 — изменения, scope и ограничения](./specs/v0.2.0/README.md)
2. [Technical Specification](./specs/docs/TECHNICAL_SPEC.md)
3. [Design System](./specs/docs/DESIGN_SYSTEM.md)
4. [Product Requirements](./specs/docs/PRD.md)
5. [Implementation Roadmap](./specs/docs/IMPLEMENTATION_ROADMAP.md)
6. [Frontend implementation plan](./frontend/IMPLEMENTATION_PLAN.md)

## Важные доменные правила legacy/full-stack baseline

- Card of the Day зависит от `user_id + accounting_date + timezone`
- accounting day длится с `03:00` до `02:59` следующего дня
- 3-card spread ограничивается ролью пользователя
- состояние раскрытия карт сохраняется в reading
- срок хранения истории зависит от роли и реализуется через TTL-логику MongoDB

Эти правила остаются частью продукта, но в mock-oriented Home версии `v0.2.0` ещё не подключены к реальному пользовательскому состоянию.
