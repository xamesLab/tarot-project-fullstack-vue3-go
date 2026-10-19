# Nebula Tarot

`Nebula Tarot` — production full-stack tarot-платформа с Go API, MongoDB и React frontend. Предыдущий Vue 3 клиент сохранён в `frontend-legacy/` как исторический baseline `v0.1.x`.

**Текущее состояние на 21 июля 2026 года:** React frontend развёрнут на `https://nebula-tarot.ru`; публичная карта дня и AI-интерпретация работают через Go backend и OpenRouter.

## Структура репозитория

| Путь | Назначение | Статус |
|---|---|---|
| `frontend/` | React 18 + TypeScript SPA | Актуальный production frontend |
| `frontend-legacy/` | Предыдущий Vue 3 SPA | Архивный full-stack baseline `v0.1.x` |
| `backend/` | REST API на Go + Gin, MongoDB, JWT, OpenRouter | Действующий backend |
| `deploy.sh` | Сборка и production deployment на VPS | Действующий deploy entry point |
| `specs/` | Общая, versioned и operational документация | [Индекс](./specs/README.md) |

## Актуальный frontend

- публичный Card of the Day с guest claim и переносом карты в аккаунт;
- AI-толкование по четырём сферам;
- login/registration и восстановление JWT-сессии;
- protected routes `/profile` и `/spreads`;
- keeper-only `/admin`;
- RU/EN, dark/light theme и адаптивный интерфейс;
- development-only `/ui-kit`;
- feature-oriented слои `app`, `api`, `features`, `shared`, `widgets`, `pages`.

Транспорт разделён по возможностям:

- `dailyCard` всегда использует реальный HTTP API;
- auth, profile, catalogs, history, spreads и admin используют transport из `VITE_AUTH_API_MODE`;
- generic readings и oracle chat используют transport из `VITE_API_MODE`;
- production сейчас использует HTTP для account-функций и mock для generic oracle/readings facade.

## Стек

### Frontend

- React 18, TypeScript 5, Vite 7;
- React Router, Zustand, TanStack Query;
- i18next + react-i18next;
- react-hook-form + zod;
- CSS custom properties, Tailwind CSS и компонентные CSS-стили;
- Vitest + Testing Library.

### Backend

- Go 1.23, Gin;
- MongoDB и официальный Go driver;
- JWT, bcrypt, validator;
- direct OpenRouter integration;
- стандартный Go HTTP transport через ограниченный SSH/SOCKS relay в production.

## Локальный запуск

Требования: Go `1.23.x`, Node.js 20.19+ или 22.12+, npm и доступная MongoDB.

```bash
cp backend/.env.example backend/.env
./start-dev.sh
```

`start-dev.sh` defaults to `DEV_DATA_MODE=auto`: it tries to start the Go API and uses the frontend mock API when MongoDB is not configured or cannot be reached. To choose explicitly:

```bash
DEV_DATA_MODE=mock ./start-dev.sh # Start the UI with mock data
DEV_DATA_MODE=mongo ./start-dev.sh # Require the Go API and MongoDB
```

Frontend mock adapters and fixture data live in [`mocks/`](./mocks/README.md). Changes made in mock mode are not written to MongoDB.

По умолчанию:

- frontend: `http://localhost:5173`;
- backend: `http://localhost:8888`;
- API: `http://localhost:8888/api`;
- UI Kit: `http://localhost:5173/ui-kit` только в development.

Раздельный запуск:

```bash
cd backend
go mod download
go run .
```

```bash
cd frontend
npm install
npm run dev
```

## Проверки

```bash
cd frontend
npm run typecheck
npm run lint
npm test
npm run build
```

```bash
cd backend
go test ./...
go vet ./...
go build ./...
```

## Production deployment

Команда git push публикует изменения в GitHub, но сама по себе не обновляет production. Деплой запускается на VPS: подключитесь к серверу, переключитесь на пользователя tarot-deploy и запустите скрипт из отдельного клона репозитория:

    sudo -iu tarot-deploy
    cd /home/tarot-deploy/nebula-tarot
    bash deploy.sh

Скрипт получает точный commit из origin/main, проверяет и собирает backend и frontend на VPS, затем вызывает существующий installer. Installer проверяет checksum, активирует релиз, запускает smoke tests и при сбое выполняет rollback. Локальная машина не собирает и не загружает файлы релиза. Production env, MongoDB, uploads, TLS, gtwapi, VPN и OpenRouter relay при обычном деплое не заменяются.

Для первоначальной настройки нужен read-only GitHub deploy key пользователя tarot-deploy и клон приватного репозитория. Команды настройки, деплоя, проверки и отката описаны в [Production deployment](./specs/docs/DEPLOYMENT_GUIDE.md).

## Документация

1. [Индекс документации](./specs/README.md)
2. [Карта проекта](./specs/docs/AI_PROJECT_MAP.md)
3. [Production deployment](./specs/docs/DEPLOYMENT_GUIDE.md)
4. [OpenRouter relay-прокси](./specs/docs/OPENROUTER_PROXY.md)
5. [Technical Specification](./specs/docs/TECHNICAL_SPEC.md)
6. [Design System](./specs/docs/DESIGN_SYSTEM.md)
7. [Product Requirements](./specs/docs/PRD.md)

## Исторические версии

- `v0.1.0` — Vue full-stack MVP baseline;
- `v0.1.1` — direct OpenRouter flow поверх Vue baseline;
- `v0.2.0` — переход на React и новый визуальный foundation;
- текущее состояние включает post-`v0.2.0` интеграцию React frontend с production backend.

Versioned-документы фиксируют контекст своей версии и не должны использоваться как operational runbook. Для текущего кода приоритет имеют код и тесты, затем `specs/docs/`.
