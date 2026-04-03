# Этап 0. Подготовка проекта

> **Статус:** выполнено  
> **Кратко:** backend и frontend scaffold созданы, MongoDB seed выполнен, локальный backend проверен, frontend зависимости установлены, `build` прошёл, `/ui-kit` доступен в dev-режиме.

## Цель этапа

Подготовить рабочую основу проекта: окружение разработки, стартовую структуру репозитория, базовый backend и frontend каркас, UI Kit, i18n и начальное наполнение базы данными карт и ролевых политик.

## Что должно быть реализовано

### 1. Окружение разработки

- установить `Go 1.22+`
- установить `Node.js 20 LTS`
- подготовить подключение к `MongoDB 6.0+`
- создать локальную или удалённую базу `tarot_daily`
- зафиксировать переменные окружения для frontend и backend

### 2. Инициализация backend

- создать папку `backend/`
- инициализировать Go-модуль
- создать базовую слоистую структуру:
  - `config/`
  - `internal/middleware/`
  - `internal/handlers/`
  - `internal/services/`
  - `internal/repositories/`
  - `internal/models/`
  - `internal/utils/`
  - `uploads/`
- поднять базовый сервер `Gin` на `localhost:8080`
- настроить `.gitignore` для `uploads/` и `.env`

### 3. Инициализация frontend

- создать папку `frontend/`
- инициализировать проект на `Vite + Vue 3`
- подготовить структуру:
  - `src/router/`
  - `src/stores/`
  - `src/services/`
  - `src/pages/`
  - `src/components/`
  - `src/composables/`
  - `src/locales/`
  - `src/styles/`
  - `src/utils/`
- настроить dev-сервер на `localhost:5173`
- подключить `.env` с `VITE_API_URL`

### 4. UI Kit и темы

- создать dev-страницу `frontend/src/pages/dev/UIKit.vue`
- добавить dev-route `/ui-kit`
- подготовить секции UI Kit:
  - Typography
  - Colors
  - Buttons
  - Inputs
  - Cards
  - Modals
- внедрить CSS variables и SCSS theme files
- реализовать переключение светлой и тёмной темы

### 5. i18n

- подключить `Vue I18n`
- создать `frontend/src/locales/en.json`
- создать `frontend/src/locales/ru.json`
- реализовать переключатель языка
- убедиться, что архитектура готова к поэтапному наполнению переводов

### 6. Seed-данные и схема БД

- создать seed-скрипт
- использовать `specs/temp/mock/tarot-cards.json` как источник данных
- загрузить в БД `78` карт
- создать стартовые документы:
  - `role_policies`
  - `admin_settings`
  - `decks` для дефолтной колоды
- создать обязательные индексы

## Результат этапа

После завершения этапа у команды должен быть репозиторий, в котором:

- backend стартует локально
- frontend стартует локально
- MongoDB подключена и заполнена базовыми данными
- доступны переключение темы и языка
- есть стартовая страница UI Kit

## Артефакты этапа

- каркас `backend/`
- каркас `frontend/`
- стартовые `.env.example` или эквивалентная документация по конфигурации
- seed-скрипт
- загруженные коллекции `cards`, `role_policies`, `admin_settings`, `decks`

## Зависимости

- утверждённые документы из `specs/docs`
- доступ к MongoDB
- мок-данные карт

## Критерии готовности

- backend работает на `localhost:8080`
- frontend работает на `localhost:5173`
- база доступна и содержит сид-данные
- `/ui-kit` доступен в dev-режиме
- темы переключаются
- языки переключаются
