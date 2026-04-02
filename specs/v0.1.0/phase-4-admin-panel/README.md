# Этап 4. Административная панель

## Цель этапа

Реализовать административный контур для роли `keeper`: управление пользователями, ролевыми политиками, глобальными настройками и колодами.

## Текущее состояние

На `April 2, 2026` этап реализован как рабочий backend/frontend baseline:

- keeper-only admin routes подключены
- admin dashboard, users, role policies и settings работают через реальные API
- отдельные каталоги `decks` и `card backs` подключены к backend
- deck creation через ZIP upload колоды, валидацию структуры и извлечение в `backend/uploads/decks/{deck_id}/` реализован
- card back creation через отдельную загрузку изображения и сохранение в `backend/uploads/card-backs/{card_back_id}/` реализован
- global `default deck` и `default back` настраиваются в admin settings
- пользовательские overrides в профиле для `selected_deck_id` и `selected_back_id` реализованы для ролей `adept`, `master` и `keeper`

Финальная ручная приёмка по `MANUAL_TESTS.md` ещё не зафиксирована.

## Что должно быть реализовано

### 1. Backend admin endpoints

- dashboard metrics
- список пользователей с поиском и фильтрами
- изменение роли пользователя
- отключение аккаунта
- просмотр и изменение `role_policies`
- просмотр и изменение `admin_settings`
- проверка admin access

### 2. Управление колодами

- создание deck metadata
- загрузка ZIP-архива колоды
- создание card back metadata
- отдельная загрузка изображения рубашки
- валидация структуры
- проверка состава `78` карт
- проверка форматов, размера и aspect ratio
- извлечение в `uploads/decks/{deck_id}/`
- сохранение рубашки в `uploads/card-backs/{card_back_id}/`
- смена статуса `processing` -> `done`
- отдельные списки доступных колод и рубашек
- установка default deck и default back для всех пользователей, если у пользователя нет собственного выбора
- пользовательский выбор deck/back в профиле только для `adept`, `master` и `keeper`

### 3. Frontend admin shell

- `AdminLayout`
- sidebar navigation
- breadcrumbs
- admin guards
- показ ссылки на admin-раздел только для `keeper`

### 4. Admin pages

- `AdminDashboard`
- `Users`
- `Decks`
- `RolePolicies`
- `Settings`
- сервисный слой `admin.js`
- store `admin.js`

## Результат этапа

После этапа проект получает рабочий административный интерфейс, достаточный для управления пользователями, лимитами, настройками и контентом колод.

Отдельно фиксируем важное правило доменной логики:

- в админке `keeper` управляет колодами и рубашками независимо
- default deck и default back настраиваются глобально
- эти значения используются всеми пользователями по умолчанию
- пользовательские override-выборы допускаются только для ролей `adept`, `master` и `keeper`

## Основные риски этапа

- ошибки в проверке роли `keeper`
- небезопасная обработка ZIP-файлов
- невалидные колоды и частично загруженные файлы
- массовое влияние изменений role policy на реальных пользователей

## Критерии готовности

- admin endpoints доступны только `keeper`
- список пользователей и изменение ролей работают
- role policies редактируются и сохраняются
- global settings редактируются и сохраняются
- deck upload проходит валидацию и создаёт рабочую колоду
- card back upload проходит валидацию и создаёт рабочую рубашку
- admin умеет раздельно управлять deck/back defaults
- профиль позволяет `adept`, `master` и `keeper` выбрать deck/back overrides независимо
- admin pages доступны на EN и RU
