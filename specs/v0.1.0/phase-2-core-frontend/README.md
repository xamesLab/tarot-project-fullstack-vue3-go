# Этап 2. Базовый frontend

## Цель этапа

Собрать пользовательский каркас приложения: сервисный слой запросов, глобальные stores, роутинг, базовые UI-компоненты, auth-страницы и skeleton dashboard.

## Что должно быть реализовано

### Артефакты дизайна

- разработка должна в первую очередь опираться на [specs/docs/DESIGN_SYSTEM.md](../../docs/DESIGN_SYSTEM.md)
- локальный [DESIGN_SYSTEM.md](./DESIGN_SYSTEM.md) используется как фазовый артефакт и рабочая копия для `phase-2`
- если между центральным docs и локальным phase-2 design-system есть пересечения или расхождения, реализация этого этапа должна следовать локальному [DESIGN_SYSTEM.md](./DESIGN_SYSTEM.md)
- `frontend/src/pages/dev/UIKit.vue` считается visual reference для `v0.1.0`
- новые страницы и компоненты не должны вводить второй параллельный UI-язык

### 1. API слой

- общий `Axios` client
- автоматическая подстановка JWT
- обработка `401` и `403`
- базовые сервисы:
  - `auth`
  - `user`
  - placeholder для `reading`

### 2. Глобальное состояние

- `auth` store
- `app` store
- `readings` store в виде заготовки
- persistence для темы, языка и токена, если это утверждено реализацией

### 3. Роутинг

- публичные маршруты: `/login`, `/register`
- защищённые маршруты: `/`, `/profile`, `/history`
- admin route group: `/admin/*`
- dev-only route: `/ui-kit`
- guard на авторизацию
- guard на роль `keeper`

### 4. Базовые UI-компоненты

- `Button`
- `Input`
- `Select`
- `Modal`
- `Header`
- `Footer`
- composables:
  - `useTheme`
  - `useAuth`

### 5. Страницы аутентификации

- `Login.vue`
- `Register.vue`
- валидация форм
- отображение ошибок API и ошибок формы
- редирект после успеха
- поддержка EN и RU

### 6. Каркас dashboard

- приветствие пользователя
- placeholders для `Card of the Day`
- placeholders для `3-card spread`
- загрузочное состояние
- пустое состояние для первого визита

## Результат этапа

После этапа пользователь уже может открыть приложение, зарегистрироваться, войти, попасть на защищённый dashboard и пользоваться базовой оболочкой интерфейса.

## Текущее состояние реализации

На `April 1, 2026` этап собран как рабочий baseline для перехода к `phase-3`.

- реализованы `Axios` client, `auth/user/reading` services и глобальные stores `auth/app/readings`
- добавлены router guards для публичных, приватных, admin и dev-only маршрутов
- собраны базовые primitives и layout shell по локальному `DESIGN_SYSTEM.md`
- реализованы `Login`, `Register`, `Dashboard`, `Profile`, `History` и admin placeholders
- `Profile` подключён к `PATCH /api/me/profile`, настройки темы и языка к `PATCH /api/me/settings`
- `UI Kit` показывает вынесенные phase-2 base components как reusable showcase
- production build frontend проходит и подтверждает целостность shell-слоя

## Зависимости

- завершённый backend этапа 1
- рабочие auth/profile endpoints

## Критерии готовности

- auth-страницы работают end-to-end
- router guards защищают приватные маршруты
- тема и язык переключаются
- UI Kit содержит базовые компоненты
- skeleton dashboard отображается после логина
