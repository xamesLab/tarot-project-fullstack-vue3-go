# NEBULA-TAROT frontend

Актуальный production SPA на React 18, TypeScript и Vite.

## Запуск

Требования: Node.js 20.19+ или 22.12+ и npm.

```bash
npm install
npm run dev
```

Routes:

- `/` — публичная карта дня;
- `/profile` — профиль, настройки и история авторизованного пользователя;
- `/spreads` — расклады авторизованного пользователя;
- `/admin` — keeper-only административный интерфейс;
- `/ui-kit` — development-only витрина компонентов.

## Команды

- `npm run dev` — Vite dev server;
- `npm run typecheck` — TypeScript check;
- `npm run lint` — ESLint без warnings;
- `npm test` — Vitest;
- `npm run test:deploy` — все тесты последовательно, один worker, лимиты теста и hook по 30 секунд, подробный вывод с длительностью и JS heap;
- `npm run build` — typecheck и production build;
- `npm run preview` — локальный preview production bundle;
- `npm run extract:references` — извлечение исходных prototype templates в ignored-файлы.

Команды Vite/Vitest явно используют `vite.config.ts`. Обычный `npm test` запускается максимум в двух worker-процессах с лимитом теста 10 секунд. Для VPS используется `npm run test:deploy`: один worker, без параллельного выполнения файлов, до 30 секунд на тест и hook. Ожидание отдельного асинхронного DOM-запроса в обоих режимах — до 5 секунд. Состав тестов, изоляция файлов и настоящие ленивые импорты сохраняются.

## Конфигурация API

```env
VITE_API_MODE=mock|http
VITE_AUTH_API_MODE=mock|http
VITE_API_URL=http://localhost:8888/api
```

- `dailyCard` всегда использует HTTP endpoints Go backend;
- `VITE_AUTH_API_MODE` выбирает transport для auth, profile, catalog, history, spreads и admin;
- `VITE_API_MODE` выбирает transport для generic readings и oracle facade;
- production использует относительный `VITE_API_URL=/api`, HTTP account transport и mock generic oracle/readings transport.

JWT хранится как `tarot_daily_token`; `/api/me` восстанавливает сессию после reload. Theme и language сохраняются как `nebula-tarot:theme` и `nebula-tarot:lang`.

## Архитектура

- `src/app` — router и providers;
- `src/api` — DTO, HTTP client, mock adapter и facade;
- `src/features` — пользовательские сценарии и feature state;
- `src/shared` — tokens, utilities и reusable UI;
- `src/widgets` — крупные композиционные блоки;
- `src/pages` — route-level страницы;
- `src/test` — unit, component и API-boundary tests.

Новый reusable UI сначала добавляется в `shared/ui` и демонстрируется на `/ui-kit`. `Arcana Home.html` и `Arcana UI Library.html` являются неизменяемыми визуальными references; их legacy-названия не являются брендом продукта.

Production frontend собирается на VPS скриптом [deploy.sh](../deploy.sh), который запускается от tarot-deploy. Тестовая конфигурация API находится в .env.test; production build использует .env.production. Подробности: [deployment guide](../specs/docs/DEPLOYMENT_GUIDE.md).
