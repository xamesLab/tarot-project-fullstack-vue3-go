# Архитектура этапа v0.1.1

## 1. Архитектурная цель
Реализовать новые функции без пересборки ядра MVP:
- переиспользовать текущий стек `Vue 3 + Pinia + Vue Router` и `Go + Gin + MongoDB`;
- расширить domain-модель и API;
- выделить изолированный admin-контур на уровне маршрутов, layout и прав доступа.

## 2. Изменения по слоям

## 2.1 Frontend
- Добавляется global-store режимов UI:
- `locale: ru|en`;
- `theme: light|dark`.
- Добавляется i18n-слой для UI-строк.
- Добавляется отдельный layout для админки (`AdminShell`) с отдельной навигацией.
- Пользовательский layout и admin layout не смешивают меню.
- Добавляются новые страницы:
- `ProfilePage` (основное приложение);
- `AdminUsersPage`;
- `AdminDecksPage`;
- `AdminSettingsPage`.
- Компоненты карточек переводятся на модель `faceDown -> revealOnClick`.

## 2.2 Backend
- Расширяется профиль пользователя:
- `displayName`;
- `timezone`;
- `status: active|banned|freeze`;
- `preferredDeckId`.
- Добавляется/расширяется домен колод:
- сущность `deck` с владельцем, статусом и набором изображений.
- Добавляются admin API:
- пользователи: list/filter/sort/search/pagination;
- колоды: create/list/get/update/upload assets;
- общие настройки (каркас + базовые параметры).
- Сохраняется role-based доступ (`admin` для админских API).

## 2.3 Хранилище (MongoDB + assets)
- MongoDB хранит метаданные колод и пользователей.
- Изображения колод для этапа `0.1.1` хранятся локально в файловой системе backend.
- В БД сохраняются ссылки/пути к файлам (`backImagePath`, `cards[id].imagePath`).
- Источник начальных тестовых артефактов:
- `specs/temp/card-image/*.jpg`;
- `specs/temp/mock/tarot-cards.json`.

## 3. Целевая модель данных

### 3.1 `users`
- `_id`
- `email/login`
- `passwordHash`
- `role`
- `status` (`active|banned|freeze`)
- `displayName`
- `timezone`
- `preferredDeckId` (nullable)
- `createdAt`
- `updatedAt`

### 3.2 `decks`
- `_id`
- `name`
- `status` (`process|done`)
- `createdByUserId`
- `backImagePath` (nullable)
- `cards` (map: `cardId -> imagePath`)
- `cardsUploadedCount`
- `requiredCardsCount` (`78`)
- `createdAt`
- `updatedAt`

### 3.3 `card_catalog` (опционально, если нужно вынести из seed)
- `id`
- `key`
- `name`
- `arcana`
- `type`
- `interpretation`

## 4. Индексы
- `users`: `role`, `status`, `createdAt`, `displayName`.
- `decks`: `status`, `createdByUserId`, `createdAt`.
- Для точного поиска по имени:
- индекс по `displayName` (с согласованной стратегией регистра).

## 5. API-контракты (MVP+)

### 5.1 Public/User
- `GET /api/me/profile`
- `PATCH /api/me/profile` (`displayName`, `timezone`, `preferredDeckId`)
- `GET /api/decks/available` (только `done`)
- `GET /api/reading/today` (возвращает карты выбранной колоды)

### 5.2 Admin: users
- `GET /api/admin/users?page=1&pageSize=20&role=&status=&sort=registeredAt:desc&name=`
- Поддержка точного поиска по `name`.

### 5.3 Admin: decks
- `POST /api/admin/decks` (создание колоды)
- `GET /api/admin/decks` (список колод)
- `GET /api/admin/decks/{deckId}` (детали + карта загруженных изображений)
- `POST /api/admin/decks/{deckId}/assets/back` (загрузка рубашки)
- `POST /api/admin/decks/{deckId}/assets/cards/{cardKey}` (загрузка карты, `cardKey=22..77|100..121`)
- Статус пересчитывается автоматически после каждой загрузки; отдельный `finalize` endpoint не требуется.

### 5.4 Admin: settings
- `GET /api/admin/settings`
- `PATCH /api/admin/settings`

## 6. Бизнес-правила
- Админка доступна только по прямому URL `/admin`, без ссылок из main-навигации.
- Статус deck:
- `process`, если нет `back` или загружено < 78 карт;
- `done`, если есть `back` и ровно 78 карт.
- В профиле пользователь может выбрать только deck со статусом `done`.
- В user-flow карта сначала закрыта (рубашка), контент скрыт.
- По клику карта раскрывается и показывает фронт + описание.
- Поиск пользователей по имени в админке выполняется как точное регистронезависимое совпадение.
- Локализация `en` применяется только к UI-слою; данные карт не переводятся на этапе `0.1.1`.

## 7. Потоки

### 7.1 Admin decks flow
1. Админ создает deck (имя).
2. Загружает `back` и карты частями.
3. Backend после каждого сохранения пересчитывает статус.
4. При полной загрузке статус автоматически `done`.

### 7.2 User profile + reading flow
1. Пользователь открывает профиль и выбирает `done` deck.
2. Выбор сохраняется в `preferredDeckId`.
3. При получении расклада backend возвращает ссылки карт выбранной колоды + метаданные карты.
4. Frontend рисует рубашку до клика, после клика раскрывает карту.

## 8. Риски и меры
- Риск: несоответствие `id` карты и файла изображения.
- Мера: строгая валидация `cardKey` (`22..77`, `100..121`) и health-check deck completeness.
- Риск: рост объема локального storage.
- Мера: ограничение размера файла + cleanup-процедуры + подготовка к выносу в object storage.
- Риск: различия локали контента карт и UI.
- Мера: на этапе `0.1.1` локализуем UI, контент карт оставляем как в каталоге.
