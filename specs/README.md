# Документация Nebula Tarot

Документы разделены на действующие общие/операционные документы и versioned baselines. Для текущего состояния source of truth — код и тесты, затем документы из `specs/docs/`.

## Актуальные документы

- [Карта проекта для AI-агента](./docs/AI_PROJECT_MAP.md) — где искать код, как проходят основные потоки и как безопасно работать с проектом
- [Production deployment](./docs/DEPLOYMENT_GUIDE.md) — действующий путь от `git push` до обновления production
- [OpenRouter relay-прокси](./docs/OPENROUTER_PROXY.md) — схема, эксплуатация, проверки и откат без чувствительных данных
- [v0.2.0](./v0.2.0/README.md) — новый React frontend, нововведения, scope, ограничения и критерии приёмки
- [PRD](./docs/PRD.md) — продуктовый baseline и дополнение для `v0.2.0`
- [Technical Specification](./docs/TECHNICAL_SPEC.md) — backend-контракты и актуальная frontend-архитектура
- [Design System](./docs/DESIGN_SYSTEM.md) — компоненты, tokens и правила нового frontend
- [Implementation Roadmap](./docs/IMPLEMENTATION_ROADMAP.md) — исторический MVP roadmap и переход к `v0.2.0`
- [Telegram Mini App integration](./telegram-miniapp-integration.md) — отдельная интеграционная спецификация

Telegram-документ был составлен для Vue baseline и требует React re-baseline перед реализацией; это явно отмечено в самом файле.

## Версии

| Версия | Основной результат | Frontend |
|---|---|---|
| [v0.1.0](./v0.1.0/README.md) | Завершённый full-stack MVP | `frontend-legacy/`, Vue 3 |
| [v0.1.1](./v0.1.1/README.md) | Direct AI через OpenRouter | `frontend-legacy/`, Vue 3 |
| [v0.2.0](./v0.2.0/README.md) | Новый публичный experience и React foundation | `frontend/`, React + TypeScript |

## Правило интерпретации

Versioned-документы описывают состояние конкретной версии и не переписываются задним числом. Если старый документ ссылается на Vue, Pinia или `frontend/src/pages/dev/UIKit.vue`, эта ссылка относится к `v0.1.x`; после переноса соответствующий код находится в `frontend-legacy/`. Для текущего frontend используются React, Zustand, React Query и route `/ui-kit` из `frontend/src/pages/ui-kit/`.

Фазовые `AI_AGENT_PLAN.md`, `MANUAL_TESTS.md` и README под `specs/v0.x/` проверены как исторические артефакты. Они не являются инструкцией по production deployment или текущей серверной конфигурации.
