# Backend Updates Log v0.1.2

**Дата обновления:** 20 марта 2026  
**Версия:** v0.1.2  
**Статус:** ✅ Completed

---

## Изменения в рамках задачи "Табы для раскладов"

### Новые endpoint'ы

Добавлены 3 новых endpoint'а для генерации отдельных раскладов:

| Endpoint | Метод | Описание |
|----------|-------|----------|
| `/api/reading/generate/card-of-day` | POST | Генерация новой карты дня |
| `/api/reading/generate/spread-a` | POST | Генерация расклада A (3 карты) |
| `/api/reading/generate/spread-b` | POST | Генерация расклада B (5 карт) |

Все endpoint'ы требуют Bearer token и учитывают дневной лимит роли пользователя.

---

## Изменённые файлы

### 1. Сервисный слой (Service Layer)

**Файл:** `backend/internal/reading/service.go`

**Добавлено:**
```go
func (s *Service) GenerateCardOfDay(ctx context.Context, userID string) (*TodayView, error)
func (s *Service) GenerateSpreadA(ctx context.Context, userID string) (*TodayView, error)
func (s *Service) GenerateSpreadB(ctx context.Context, userID string) (*TodayView, error)
```

Каждый метод:
- Проверяет дневной лимит через `ensureDaily()`
- Вызывает соответствующий метод генерации
- Возвращает обновлённый `TodayView`

---

### 2. Бизнес-логика (Domain Logic)

**Файл:** `backend/internal/reading/daily.go`

**Добавлено:**
```go
func (s *Service) generateCardOfDayOnly(...) (*models.DailyReading, error)
func (s *Service) generateSpreadAOnly(...) (*models.DailyReading, error)
func (s *Service) generateSpreadBOnly(...) (*models.DailyReading, error)
```

Каждый метод:
- Генерирует один элемент (карту дня или расклад)
- Вызывает соответствующий метод repository для обновления

---

### 3. Репозиторий (Repository Layer)

**Файл:** `backend/internal/repository/mongo/daily_reading_repository.go`

**Добавлено:**
```go
func (r *DailyReadingRepository) UpdateCardOfDay(...) (*models.DailyReading, error)
func (r *DailyReadingRepository) UpdateSpreadA(...) (*models.DailyReading, error)
func (r *DailyReadingRepository) UpdateSpreadB(...) (*models.DailyReading, error)
```

Каждый метод:
- Обновляет только соответствующее поле в документе `daily_readings`
- Инкрементирует `spreadGenerationsCount`
- Обновляет `lastSpreadGeneratedAt` и `updatedAt`
- Поддерживает опциональное обновление `deckThemeId`

---

### 4. HTTP Handler

**Файл:** `backend/internal/transport/http/handlers/reading_handler.go`

**Добавлено:**
```go
func (h *ReadingHandler) GenerateCardOfDay(c *gin.Context)
func (h *ReadingHandler) GenerateSpreadA(c *gin.Context)
func (h *ReadingHandler) GenerateSpreadB(c *gin.Context)
```

Каждый handler:
- Вызывает соответствующий метод сервиса
- Обрабатывает ошибки через `mapError()`
- Возвращает JSON response с `TodayView`

---

### 5. Маршрутизация (Routes)

**Файл:** `backend/internal/transport/http/routes.go`

**Добавлено в `registerProtectedRoutes()`:**
```go
authProtected.POST("/reading/generate/card-of-day", deps.readingHandler.GenerateCardOfDay)
authProtected.POST("/reading/generate/spread-a", deps.readingHandler.GenerateSpreadA)
authProtected.POST("/reading/generate/spread-b", deps.readingHandler.GenerateSpreadB)
```

---

### 6. Документация

**Файл:** `backend/README.md`

**Добавлено:**
```markdown
- `POST /api/reading/generate/card-of-day` (Bearer token required, `v0.1.2`)
  - генерирует новую карту дня в рамках дневного лимита роли.
- `POST /api/reading/generate/spread-a` (Bearer token required, `v0.1.2`)
  - генерирует новый расклад A (3 карты) в рамках дневного лимита роли.
- `POST /api/reading/generate/spread-b` (Bearer token required, `v0.1.2`)
  - генерирует новый расклад B (5 карт) в рамках дневного лимита роли.
```

---

## Структура данных

### Коллекция `daily_readings`

**Обновляемые поля:**

| Поле | Тип | Описание |
|------|-----|----------|
| `cardOfDay` | `CardDraw` | Карта дня (обновляется при генерации карты дня) |
| `spreadA` | `[]CardDraw` | Расклад A, 3 карты (обновляется при генерации расклада A) |
| `spreadB` | `[]CardDraw` | Расклад B, 5 карт (обновляется при генерации расклада B) |
| `spreadGenerationsCount` | `int` | Счётчик генераций (инкрементируется при любой генерации) |
| `lastSpreadGeneratedAt` | `time.Time` | Время последней генерации |
| `updatedAt` | `time.Time` | Время обновления документа |
| `deckThemeId` | `ObjectID` | ID темы колоды (опционально) |

---

## Обработка ошибок

Все endpoint'ы возвращают стандартные ошибки:

| Код HTTP | Error Code | Описание |
|----------|------------|----------|
| 400 | `invalid_input` | Некорректные входные данные |
| 401 | `unauthorized` | Пользователь не найден |
| 403 | `forbidden` | Политика роли не найдена |
| 429 | `daily_spread_limit_reached` | Достигнут дневной лимит генераций |
| 503 | `cards_unavailable` | Каталог карт недоступен |

---

## Тестирование

### Ручное тестирование

```bash
# Получить карту дня
curl -X POST http://localhost:8080/api/reading/generate/card-of-day \
  -H "Authorization: Bearer <token>"

# Получить расклад A
curl -X POST http://localhost:8080/api/reading/generate/spread-a \
  -H "Authorization: Bearer <token>"

# Получить расклад B
curl -X POST http://localhost:8080/api/reading/generate/spread-b \
  -H "Authorization: Bearer <token>"
```

### Автоматические тесты

Unit-тесты для новых методов не добавлялись (используется существующая логика генерации).

---

## Зависимости

- MongoDB коллекция `daily_readings` (существующая)
- Сервис `roles.Service` для проверки дневного лимита
- Репозиторий `DailyReadingRepository` для обновления документов

---

## Обратная совместимость

- Существующий endpoint `POST /api/reading/generate` продолжает работать
- Все изменения обратно совместимы
- Frontend может использовать как новые endpoint'ы, так и старый

---

## Примечания

1. Каждая генерация (карты дня, расклада A, расклада B) инкрементирует общий счётчик `spreadGenerationsCount`
2. Дневной лимит применяется ко всем генерациям суммарно
3. При генерации обновляется только соответствующее поле, остальные данные сохраняются

---

## Связанные задачи

- Frontend: добавлены табы для переключения между раскладами
- Frontend: кнопка генерации на каждом табе (генерирует только активный расклад)
- Frontend: компонент `AppTabs.vue` для переключения вкладок

---

## Файлы для проверки

```
backend/
├── README.md                                          (обновлён)
├── internal/
│   ├── reading/
│   │   ├── service.go                                 (обновлён)
│   │   └── daily.go                                   (обновлён)
│   ├── repository/mongo/
│   │   └── daily_reading_repository.go                (обновлён)
│   └── transport/http/
│       ├── handlers/reading_handler.go                (обновлён)
│       └── routes.go                                  (обновлён)
```
