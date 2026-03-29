# Back-log (current post-MVP)

Дата актуализации: 29 марта 2026.

## Цель
Список следующих задач после закрытия MVP, backend-этапа `v0.1.2` и frontend-этапа `v0.1.3`.

## Задачи

1. Этап `v0.1.4`
- Актуальный scope этапа:
  - `specs/version-0.1.4/requirements-v0.1.4.md`
  - `specs/version-0.1.4/steps.md`
- Статус: `planned`
- Приоритет: `high`

2. Закрытие валидации `v0.1.2`
- Выполнить фактический smoke ask-flow:
  - `deploy/scripts/smoke-v0.1.2-ask.ps1`;
  - проверить callback security конфигурацию;
  - обновить финальные статусные документы `v0.1.2`.
- Статус: `todo`
- Приоритет: `high`

3. Dev tools index и развитие UI Kit workflow
- Добавить `/dev` как каталог внутренних страниц.
- Подготовить площадку для дальнейших dev-only инструментов.
- Начать выравнивание продуктовых экранов по токенам и паттернам `/ui-kit`.
- Статус: `todo`
- Приоритет: `high`

4. E2E UI-тесты ключевых сценариев
- Статус: `deferred until explicit user request`
- Приоритет: `paused`
- Добавить Playwright/Cypress сценарии:
  - card flip (закрытая карта -> клик -> открытая карта);
  - профиль (`/profile`) и выбор `done` колоды;
  - админка: users/decks/settings.

5. Cleanup скрипт тестовых данных smoke
- Добавить утилиту очистки тестовых сущностей:
  - пользователи `v011*`, `smoke*`;
  - тестовые колоды, созданные smoke-скриптами.
- Статус: `todo`
- Приоритет: `high`

6. Визуальные snapshot-проверки для UI Kit
- Статус: `deferred until explicit user request`
- Приоритет: `paused`
- Добавить визуальные snapshot/e2e проверки для `/ui-kit`.
- Зафиксировать базовую защиту от регрессий после введения Arcana background и UI Kit.

7. Внешнее хранилище изображений
- Статус: `deferred until explicit user request`
- Приоритет: `paused`
- Подготовить вынос assets из локального FS в object storage (S3-совместимое).
- Добавить стратегию миграции и обратной совместимости URL.

8. Changelog и release notes
- Добавить релизный changelog для `v0.1.1`.
- Зафиксировать версию backend/frontend в релизной документации.
- Статус: `todo`
- Приоритет: `medium`

9. Backup/restore runbook MongoDB
- Добавить инструкции резервного копирования и восстановления.
- Включить проверку восстановления на тестовом стенде.
- Статус: `todo`
- Приоритет: `medium`
