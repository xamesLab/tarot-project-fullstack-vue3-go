# Backend Roles and Access (Step 10)

Дата фиксации: 8 марта 2026.

## 1. Цель шага
Реализовать backend-проверку ролей `user`/`admin` и ограничить admin-доступ.

## 2. Реализованные компоненты
- `backend/internal/transport/http/middleware/role.go`
  - middleware `RoleRequired(roles ...string)`
  - отказ `403 forbidden` при недостаточных правах.
- `backend/internal/transport/http/handlers/access_handler.go`
  - endpoint проверки админ-доступа.
- `backend/internal/transport/http/router.go`
  - `GET /api/admin/access-check`
  - защита цепочкой `AuthRequired` + `RoleRequired("admin")`.

## 3. API
- `GET /api/admin/access-check`
  - Требует Bearer JWT.
  - Доступен только роли `admin`.
  - Для роли `user` возвращает `403`.

## 4. Статус проверки
- Сборка backend проходит (`go test ./...`).
