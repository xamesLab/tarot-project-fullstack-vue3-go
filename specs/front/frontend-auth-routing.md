# Frontend Auth + Routing (Progress)

Дата фиксации: 8 марта 2026.

## Что закрыто
- Auth flow:
  - страницы `login/register`,
  - login/register запросы в API,
  - хранение токена в `localStorage`,
  - `logout`,
  - hydration пользователя через `GET /api/me`.
- Route guards:
  - public-only (`/login`, `/register`),
  - private (`/dashboard`, `/settings`),
  - admin-only (`/admin`).

## Dev networking
- В `frontend/.env` задано `VITE_API_BASE_URL=/`.
- В `vite.config.ts` включен proxy:
  - `/api` -> `http://127.0.0.1:8080`.

Итог: в dev frontend отправляет запросы на свой origin (`/api/...`), а Vite проксирует их в backend.
