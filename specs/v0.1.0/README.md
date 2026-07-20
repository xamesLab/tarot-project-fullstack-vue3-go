# План разработки v0.1.0

> **Архивный version baseline.** Файлы в этом каталоге фиксируют план и приёмку Vue-версии `v0.1.0`. Пути и deployment-команды отсюда не являются инструкцией для текущего React production. Актуальные операции: [`specs/docs/DEPLOYMENT_GUIDE.md`](../docs/DEPLOYMENT_GUIDE.md).

Этот каталог содержит поэтапный план разработки MVP `Nebula Tarot`, собранный на основе:

- `specs/docs/PRD.md`
- `specs/docs/TECHNICAL_SPEC.md`
- `specs/docs/IMPLEMENTATION_ROADMAP.md`
- `specs/docs/DESIGN_SYSTEM.md`

Цель структуры `v0.1.0` — превратить общую дорожную карту в удобные для исполнения пакеты работ. Каждый этап оформлен отдельной папкой и включает:

- `README.md` — описание цели, состава работ, ожидаемого результата и критериев готовности
- `MANUAL_TESTS.md` — ручные проверки, подтверждающие завершение этапа

## Этапы

1. `phase-0-project-setup` — подготовка окружения, каркаса проекта, UI Kit, i18n и seed-данных
2. `phase-1-core-backend` — аутентификация, профиль, базовые backend-слои и timezone-логика
3. `phase-2-core-frontend` — UI-основа, роутинг, auth-страницы, stores и layout
4. `phase-3-reading-features` — дневная карта, 3-карточный расклад, история и role-policy-логика
5. `phase-4-admin-panel` — административный контур, роли, политики, глобальные настройки, колоды и рубашки
6. `phase-5-polish-and-testing` — адаптивность, accessibility, оптимизация, ошибки и финальная ручная приёмка
7. `phase-6-deployment` — production deployment, smoke testing, мониторинг и запуск MVP

## Как использовать

1. Начинать с `phase-0-project-setup`
2. Для каждого этапа реализовывать задачи из локального `README.md`
3. После завершения реализации проходить `MANUAL_TESTS.md`
4. Переходить к следующему этапу только после выполнения критериев готовности

## Актуальный статус

На `April 10, 2026`:

- `phase-0-project-setup` — done
- `phase-1-core-backend` — done
- `phase-2-core-frontend` — done
- `phase-3-reading-features` — done
- `phase-4-admin-panel` — done
- `phase-5-polish-and-testing` — done
- версия `v0.1.0` готова как завершённый MVP baseline
- финальная ручная приёмка `v0.1.0` пройдена
- следующий отдельный этап — `phase-6-deployment`, если нужен production rollout

## Принципы исполнения

- не выходить за рамки MVP, зафиксированного в `PRD`
- держать frontend и backend синхронизированными по API-контрактам из `TECHNICAL_SPEC`
- для frontend-реализации опираться на `DESIGN_SYSTEM.md` и dev-only `UI Kit` как на source of truth
- использовать вертикальную разработку и раннюю ручную проверку
- учитывать двуязычность, темы и mobile-first поведение как обязательные условия готовности
