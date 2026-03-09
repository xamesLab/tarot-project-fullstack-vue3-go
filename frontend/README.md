# Frontend

Папка для клиентского приложения MVP.

Планируемый стек:
- Vue 3
- Pinia
- Vue Router

План работ:
- `specs/front/implementation-plan-frontend.md`

## Environment
- Шаблон: `frontend/.env.example`
- Локальный файл: `frontend/.env`

Переменные:
- `VITE_API_BASE_URL` - базовый URL backend API (`/` в dev через Vite proxy).
- `VITE_APP_NAME` - отображаемое имя приложения.

## Frontend Bootstrap
Стек:
- Vue 3
- Pinia
- Vue Router
- Vite + TypeScript

Команды:
- `npm install`
- `npm run dev`
- `npm run build`

Dev proxy:
- `vite.config.ts` проксирует `/api` -> `http://127.0.0.1:8080`.
- Это позволяет фронтенду работать без CORS-настроек backend в dev.

Реализованный каркас:
- страницы: `login`, `register`, `dashboard`, `profile`, `settings`, `admin/*`, `not-found`
- роутер и guards (public/private/admin)
- базовые stores: `auth`, `reading`, `settings`, `ui-preferences`
- API-слой: `auth.api`, `reading.api`, `settings.api`, `profile.api`, `admin.api`

Дополнительно для `v0.1.1`:
- страница `profile` (`/profile`) с `displayName/timezone/preferredDeckId`;
- API-слой `profile.api` (`/api/me/profile`, `/api/decks/available`);
- card flip UI: карта по умолчанию закрыта и раскрывается по клику.
