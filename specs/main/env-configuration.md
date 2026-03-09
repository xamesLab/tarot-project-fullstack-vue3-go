# Конфигурация окружения (Step 4)

Дата фиксации: 8 марта 2026.

## Цель шага
Перевести конфигурацию проекта на `.env`-подход и зафиксировать шаблоны переменных окружения для frontend и backend.

## Реализовано
- Добавлен шаблон backend: `backend/.env.example`
- Добавлен локальный backend env: `backend/.env`
- Добавлен шаблон frontend: `frontend/.env.example`
- Добавлен локальный frontend env: `frontend/.env`
- Добавлены правила игнорирования env-файлов: `.gitignore`

## Переменные backend
- `APP_PORT`
- `APP_ENV`
- `JWT_SECRET`
- `MONGO_HOST`
- `MONGO_PORT`
- `MONGO_URI`
- `MONGO_DB_NAME`
- `MONGO_APP_USER`
- `MONGO_AUTH_SOURCE`

## Переменные frontend
- `VITE_API_BASE_URL`
- `VITE_APP_NAME`

## Принятые значения для локальной разработки
- Backend API: `http://localhost:8080`
- Mongo host: `127.0.0.1`
- Mongo port: `27017`
- Mongo DB name: `tarot_daily`
- Mongo app user: `tarot_daily_app`
- Mongo authSource: `tarot_daily`
- Mongo URI template: `mongodb://tarot_daily_app:<set-password>@127.0.0.1:27017/tarot_daily?authSource=tarot_daily`

## Примечание по безопасности
Секреты хранятся только в `.env`. Файлы `.env.example` используются как шаблоны без чувствительных значений.
