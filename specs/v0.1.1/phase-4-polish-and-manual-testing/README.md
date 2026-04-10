# Этап 4. Полировка и финальная ручная приёмка

## Цель этапа

Довести `v0.1.1` до стабильного состояния без deployment-задач: проверить direct AI flow на основных сценариях, закрыть UX gaps и зафиксировать ручную приёмку версии.

## Что должно быть реализовано

### 1. UX polish

- проверить консистентность AI-блоков на dashboard и history
- поджать mobile layouts, если AI-слой раздувает интерфейс
- привести loading/error/success states к одному визуальному языку

### 2. Edge cases

- repeated asks
- refresh во время генерации
- частично открытый расклад
- выключенный AI feature flag
- достижение AI-лимита
- timeout провайдера
- невалидный payload от модели

### 3. Документация

- актуализировать root docs и roadmap при необходимости
- отметить manual acceptance для `v0.1.1`
- описать известные non-blocking ограничения версии

### 4. Финальная ручная приёмка

- user AI journey
- returning user AI journey
- admin AI controls journey
- RU/EN
- light/dark
- mobile/tablet/desktop

## Результат этапа

После этапа `v0.1.1` должна быть готова как завершённая non-deployment версия с прямым AI provider path и подтверждённой ручной приёмкой.

## Критерии готовности

- direct AI flow стабилен
- критические UX и backend bugs устранены
- ручные сценарии пройдены
- документация актуализирована
- deployment остаётся за пределами версии
