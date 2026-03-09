# MongoDB Data Model (Step 5)

Дата фиксации: 8 марта 2026.

## 1. Цель
Зафиксировать модель данных MVP для коллекций:
- `users`
- `cards`
- `daily_readings`
- `deck_themes`

## 2. Общие правила
- Формат даты в `daily_readings.date`: `YYYY-MM-DD`.
- Все ссылки между коллекциями храним как `ObjectId`.
- Временные поля (`createdAt`, `updatedAt`, `generatedAt`) храним в UTC (`ISODate`).
- Для MVP допускается хранение таймзоны в виде фиксированного offset (`UTC+03:00`) или IANA (`Europe/Moscow`).

## 3. Коллекция `users`

Поля:
- `_id`: `ObjectId`, required.
- `login`: `string`, required, unique, длина 3..64.
- `passwordHash`: `string`, required.
- `role`: `string`, required, default `user` (MVP роли: `user`, `admin`; поддерживаются кастомные роли через `role_policies`).
- `timezone`: `string`, required, default `UTC+00:00`.
- `deckThemeId`: `ObjectId`, optional.
- `createdAt`: `date`, required.
- `updatedAt`: `date`, required.

Пример:
```json
{
  "_id": "ObjectId('...')",
  "login": "alice",
  "passwordHash": "$2a$10$...",
  "role": "user",
  "timezone": "UTC+03:00",
  "deckThemeId": "ObjectId('...')",
  "createdAt": "2026-03-08T10:00:00Z",
  "updatedAt": "2026-03-08T10:00:00Z"
}
```

## 4. Коллекция `cards`

Поля:
- `_id`: `ObjectId`, required.
- `code`: `string`, required, unique (например, `major_fool`, `minor_cups_ace`).
- `name`: `string`, required.
- `arcanaType`: `string`, required, enum: `major | minor`.
- `description`: `string`, required.
- `interpretationOpen`: `string`, required.
- `interpretationClosed`: `string`, required.

Пример:
```json
{
  "_id": "ObjectId('...')",
  "code": "major_fool",
  "name": "The Fool",
  "arcanaType": "major",
  "description": "Start of a new journey.",
  "interpretationOpen": "Openness to new possibilities.",
  "interpretationClosed": "Fear of uncertainty."
}
```

## 5. Коллекция `deck_themes`

Поля:
- `_id`: `ObjectId`, required.
- `name`: `string`, required, unique.
- `isPlaceholder`: `bool`, required, default `true`.
- `assets`: `object`, required, формат: `cardCode -> imageUrl`.
- `createdBy`: `ObjectId`, required (admin user id).
- `createdAt`: `date`, required.

Пример:
```json
{
  "_id": "ObjectId('...')",
  "name": "Default Classic",
  "isPlaceholder": true,
  "assets": {
    "major_fool": "/themes/default/major_fool.jpg"
  },
  "createdBy": "ObjectId('...')",
  "createdAt": "2026-03-08T10:00:00Z"
}
```

## 6. Коллекция `daily_readings`

Поля:
- `_id`: `ObjectId`, required.
- `userId`: `ObjectId`, required.
- `date`: `string`, required, формат `YYYY-MM-DD` (по таймзоне пользователя).
- `timezoneUsed`: `string`, required.
- `deckThemeId`: `ObjectId`, optional.
- `cardOfDay`: `object`, required:
  - `cardId`: `ObjectId`, required.
  - `state`: `string`, required, enum: `open | closed`.
- `spreadA`: `array<object>`, required, ровно 3 элемента:
  - `cardId`: `ObjectId`, required.
  - `state`: `string`, required, enum: `open | closed`.
- `spreadB`: `array<object>`, required, ровно 5 элементов:
  - `cardId`: `ObjectId`, required.
  - `state`: `string`, required, enum: `open | closed`.
- `spreadGenerationsCount`: `int`, required, default `0`.
- `generatedAt`: `date`, required.
- `lastSpreadGeneratedAt`: `date`, optional.
- `updatedAt`: `date`, required.

Пример:
```json
{
  "_id": "ObjectId('...')",
  "userId": "ObjectId('...')",
  "date": "2026-03-08",
  "timezoneUsed": "UTC+03:00",
  "deckThemeId": "ObjectId('...')",
  "cardOfDay": {
    "cardId": "ObjectId('...')",
    "state": "open"
  },
  "spreadA": [
    { "cardId": "ObjectId('...')", "state": "closed" },
    { "cardId": "ObjectId('...')", "state": "open" },
    { "cardId": "ObjectId('...')", "state": "open" }
  ],
  "spreadB": [
    { "cardId": "ObjectId('...')", "state": "closed" },
    { "cardId": "ObjectId('...')", "state": "open" },
    { "cardId": "ObjectId('...')", "state": "open" },
    { "cardId": "ObjectId('...')", "state": "closed" },
    { "cardId": "ObjectId('...')", "state": "open" }
  ],
  "spreadGenerationsCount": 1,
  "generatedAt": "2026-03-08T10:00:00Z"
}
```

## 7. Коллекция `role_policies`

Поля:
- `_id`: `ObjectId`, required.
- `role`: `string`, required, unique.
- `dailySpreadLimit`: `int`, required (`-1` означает безлимит).
- `createdAt`: `date`, required.
- `updatedAt`: `date`, required.

Пример:
```json
{
  "_id": "ObjectId('...')",
  "role": "vip",
  "dailySpreadLimit": 25,
  "createdAt": "2026-03-09T08:00:00Z",
  "updatedAt": "2026-03-09T08:00:00Z"
}
```

## 8. Контрактные ограничения (уровень сервиса)
- Для `spreadA` не допускаются повторяющиеся `cardId` внутри массива.
- Для `spreadB` не допускаются повторяющиеся `cardId` внутри массива.
- Повторы между `cardOfDay`, `spreadA`, `spreadB` в один день разрешены (MVP-упрощение).
- Для одного пользователя на одну дату должна существовать ровно одна запись `daily_readings`.
- `cardOfDay` фиксируется на бизнес-сутки пользователя (сброс в `03:00` по таймзоне пользователя).
- Лимит `spread`-генераций задается по роли пользователя.

## 9. Решения, переданные в следующий шаг (индексы)
Индексы будут физически созданы на шаге 6, но здесь фиксируются требования:
- `users.login` unique.
- `daily_readings(userId, date)` unique compound.
- `cards.code` unique.
- `role_policies.role` unique.
