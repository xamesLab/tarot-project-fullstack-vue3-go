# Auth Module (Step 9)

Дата фиксации: 8 марта 2026.

## 1. Цель шага
Реализовать MVP-аутентификацию:
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/me`

С требованиями:
- JWT-сессия,
- хэширование паролей (bcrypt),
- роль пользователя по умолчанию `user`.

## 2. Реализованные компоненты

### Domain
- `internal/domain/models/user.go`

### Repository
- `internal/repository/mongo/user_repository.go`
  - create user
  - find by login
  - find by id
  - duplicate-key detector

### JWT
- `internal/platform/jwt/token.go`
  - генерация токена (HS256)
  - парсинг и валидация токена
  - claims: `uid`, `login`, `role`, `exp`

### Auth service
- `internal/auth/service.go`
  - register (валидация, bcrypt hash, create, issue token)
  - login (find, bcrypt compare, issue token)
  - me (по `uid` из токена)

### HTTP handlers + middleware
- `internal/transport/http/handlers/auth_handler.go`
- `internal/transport/http/handlers/me_handler.go`
- `internal/transport/http/middleware/auth.go`

### Router wiring
- `internal/transport/http/router.go`
  - `POST /api/auth/register`
  - `POST /api/auth/login`
  - `GET /api/me` (Bearer required)

## 3. Конфиг
Добавлен параметр:
- `JWT_TTL` (default `24h`)

Обновлены:
- `backend/.env`
- `backend/.env.example`

## 4. Интеграционная проверка
Проверка выполнена на реальной MongoDB через запущенный API:
- `health/live` -> `ok`
- `register` -> создан новый пользователь
- `login` -> возвращен JWT
- `me` -> возвращены профиль и роль `user`

## 5. Ограничения текущего шага
- Logout endpoint и refresh token не реализованы (не требуются MVP шагом 9).
