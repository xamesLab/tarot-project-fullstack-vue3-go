# Структура папок и файлов проекта (MVP)

Ниже предложена целевая структура монорепозитория для Tarot Daily на базе подготовленных спецификаций.

## 1. Корень проекта
```text
pro/
  specs/
  frontend/
  backend/
  deploy/
  .env.example
  README.md
```

## 2. Папка `specs/`
```text
specs/
  main/
    PRD.md
    PRD-ru.md
    architecture.md
    project-structure.md
    mvp-scope-freeze.md
    steps.md
  front/
    frontend-bootstrap.md
    frontend-auth-routing.md
    frontend-dashboard-reading.md
    frontend-settings.md
    frontend-admin.md
    implementation-plan-frontend.md
  back/
    backend-bootstrap.md
    auth-module.md
    backend-reading-engine.md
    backend-role-policies.md
    implementation-plan-backend.md
    plan.txt
  temp/
    *.jpg
```

## 3. Frontend (`frontend/`)
```text
frontend/
  src/
    main.ts
    App.vue
    pages/
      LoginPage.vue
      RegisterPage.vue
      DashboardPage.vue
      SettingsPage.vue
      AdminPage.vue
    components/
      AppHeader.vue
      CardView.vue
      SpreadView.vue
      ThemePicker.vue
      BaseLoader.vue
    stores/
      auth.store.ts
      reading.store.ts
      settings.store.ts
    services/
      apiClient.ts
      auth.api.ts
      reading.api.ts
      settings.api.ts
      admin.api.ts
    router/
      index.ts
      guards.ts
    types/
      auth.ts
      reading.ts
      deck-theme.ts
    assets/
      styles/
        base.css
        variables.css
  public/
  package.json
  vite.config.ts
  tsconfig.json
  .env.example
```

## 4. Backend (`backend/`)
```text
backend/
  cmd/
    api/
      main.go
  internal/
    config/
      config.go
    transport/
      http/
        router.go
        middleware/
          auth.go
          role.go
          error_handler.go
        handlers/
          auth_handler.go
          user_handler.go
          reading_handler.go
          deck_theme_handler.go
          admin_handler.go
    domain/
      models/
        user.go
        card.go
        daily_reading.go
        deck_theme.go
      dto/
        auth_dto.go
        reading_dto.go
        settings_dto.go
    service/
      auth_service.go
      user_service.go
      reading_service.go
      deck_theme_service.go
      admin_service.go
    repository/
      mongo/
        user_repository.go
        card_repository.go
        daily_reading_repository.go
        deck_theme_repository.go
    platform/
      mongo/
        client.go
        indexes.go
      jwt/
        token.go
      password/
        hasher.go
    seed/
      seed_cards.go
      seed_deck_themes.go
      seed_admin.go
  migrations/
  go.mod
  go.sum
  .env.example
```

## 5. Инфраструктура и запуск (`deploy/`)
```text
deploy/
  docker-compose.yml
  mongo/
    init.js
```

## 6. Файлы окружения
`backend/.env.example`
- `APP_PORT`
- `JWT_SECRET`
- `MONGO_URI`
- `MONGO_DB_NAME`
- `APP_ENV`

`frontend/.env.example`
- `VITE_API_BASE_URL`

## 7. Примечания по структуре
- `internal/` в backend ограничивает использование внутренних пакетов только сервисом.
- Все HTTP-контракты проходят через `handlers` и `dto`.
- Бизнес-логика остается в `service`, доступ к данным только через `repository`.
- Структура фронтенда разделена по ответственности: `pages`, `components`, `stores`, `services`.
