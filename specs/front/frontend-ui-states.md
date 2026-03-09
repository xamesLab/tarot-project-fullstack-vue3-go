# Frontend UI States (errors/loading/empty)

Дата фиксации: 9 марта 2026.

## 1. Что закрыто
- Добавлены и унифицированы состояния загрузки, ошибок и пустых данных на ключевых экранах:
  - `Dashboard`
  - `Settings`
  - `Admin`
  - `Login`
  - `Register`

## 2. Ключевые изменения
- `Dashboard`:
  - loading при загрузке `today`;
  - ошибка API;
  - empty-state при отсутствии данных;
  - отдельное состояние генерации раскладов.
- `Settings`:
  - loading списка тем;
  - отображение ошибок сохранения;
  - success-сообщение после сохранения.
- `Admin`:
  - loading по разделам (themes/roles);
  - success/error сообщения действий;
  - блокировка кнопок на время фоновых операций.
- `Login`/`Register`:
  - loading-label на кнопках submit.

## 3. Проверка
- `npm run build` выполняется успешно.
