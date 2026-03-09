# Backend Role Policies and Assignment

Дата фиксации: 9 марта 2026.

## 1. Цель
Добавить управляемые роли и дневные лимиты генерации раскладов:
- создание/обновление политики роли;
- назначение роли пользователю;
- применение лимита в `POST /api/reading/generate`.

## 2. Новая коллекция
`role_policies`:
- `_id`
- `role` (unique)
- `dailySpreadLimit` (`-1` означает без лимита)
- `createdAt`
- `updatedAt`

## 3. API (admin only)
- `GET /api/admin/roles`
  - список ролей и лимитов.
- `POST /api/admin/roles`
  - upsert политики роли.
  - пример: `{ "role": "vip", "dailySpreadLimit": 25 }`.
- `PATCH /api/admin/users/role`
  - назначение роли пользователю по логину.
  - пример: `{ "login": "alice", "role": "vip" }`.

## 4. Реализованные компоненты
- `backend/internal/domain/models/role_policy.go`
- `backend/internal/repository/mongo/role_policy_repository.go`
- `backend/internal/roles/service.go`
- `backend/internal/transport/http/handlers/admin_roles_handler.go`
- `backend/internal/repository/mongo/user_repository.go` (обновление роли пользователя)
- `backend/internal/transport/http/router.go` (роуты admin API)

## 5. Дефолтные политики (fallback)
- `user` -> `10`
- `admin` -> `-1` (без лимита)

Если пользователь получает кастомную роль, но политика роли не задана, генерация раскладов блокируется до создания политики.
