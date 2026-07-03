# Документация Nebula Tarot

Документы разделены на актуальные общие спецификации и versioned baselines. Для новой разработки frontend source of truth — раздел `v0.2.0` и актуальный `DESIGN_SYSTEM.md`.

## Актуальные документы

- [v0.2.0](./v0.2.0/README.md) — новый React frontend, нововведения, scope, ограничения и критерии приёмки
- [PRD](./docs/PRD.md) — продуктовый baseline и дополнение для `v0.2.0`
- [Technical Specification](./docs/TECHNICAL_SPEC.md) — backend-контракты и актуальная frontend-архитектура
- [Design System](./docs/DESIGN_SYSTEM.md) — компоненты, tokens и правила нового frontend
- [Implementation Roadmap](./docs/IMPLEMENTATION_ROADMAP.md) — исторический MVP roadmap и переход к `v0.2.0`
- [Telegram Mini App integration](./telegram-miniapp-integration.md) — отдельная интеграционная спецификация

## Версии

| Версия | Основной результат | Frontend |
|---|---|---|
| [v0.1.0](./v0.1.0/README.md) | Завершённый full-stack MVP | `frontend-legacy/`, Vue 3 |
| [v0.1.1](./v0.1.1/README.md) | Direct AI через OpenRouter | `frontend-legacy/`, Vue 3 |
| [v0.2.0](./v0.2.0/README.md) | Новый публичный experience и frontend foundation | `frontend/`, React + TypeScript |

## Правило интерпретации

Versioned-документы описывают состояние конкретной версии и не переписываются задним числом. Если старый документ ссылается на Vue, Pinia или `frontend/src/pages/dev/UIKit.vue`, эта ссылка относится к `v0.1.x`; после переноса соответствующий код находится в `frontend-legacy/`. Для `v0.2.0` используются React, Zustand, React Query и route `/ui-kit` из `frontend/src/pages/ui-kit/`.
