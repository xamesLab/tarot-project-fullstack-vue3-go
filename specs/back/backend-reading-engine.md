# Backend Reading Engine (Step 13-14)

Дата фиксации: 9 марта 2026.

## 1. Что реализовано
- Генерация и хранение `карты дня` c правилом:
  - одна карта на пользователя в рамках бизнес-суток;
  - граница бизнес-суток: `03:00` по таймзоне пользователя.
- Генерация раскладов:
  - `spreadA` (3 карты, без повторов внутри расклада);
  - `spreadB` (5 карт, без повторов внутри расклада);
  - состояния карт: `open` / `closed`.
- `GET /api/reading/today`:
  - возвращает фиксированную карту дня;
  - на первом открытии суток автогенерирует первый комплект раскладов.
- `POST /api/reading/generate`:
  - генерирует новый комплект раскладов;
  - учитывает дневной лимит генераций по роли.

## 2. Лимиты по ролям в генерации раскладов
- `user` по умолчанию: `10` генераций в сутки.
- `admin` по умолчанию: без лимита (`-1`).
- Лимит берется из `role_policies`, если политика роли задана.

## 3. Новые/обновленные backend-компоненты
- Domain:
  - `backend/internal/domain/models/card.go`
  - `backend/internal/domain/models/daily_reading.go`
- Timezone utility:
  - `backend/internal/platform/timeutil/timezone.go`
- Repositories:
  - `backend/internal/repository/mongo/card_repository.go`
  - `backend/internal/repository/mongo/daily_reading_repository.go`
- Service:
  - `backend/internal/reading/service.go`
- HTTP:
  - `backend/internal/transport/http/handlers/reading_handler.go`
  - `backend/internal/transport/http/router.go`

## 4. Формат ответа `GET /api/reading/today`
Ответ включает:
- `date`, `timezoneUsed`;
- `cardOfDay`;
- `spreadA`, `spreadB`;
- счетчики лимитов:
  - `spreadGenerationsUsed`,
  - `spreadGenerationsLimit` (если лимит ограничен),
  - `spreadGenerationsRemaining` (если лимит ограничен),
  - `canGenerateMore`.

## 5. Проверка
Smoke-проверка на реальной MongoDB через SSH-туннель подтвердила:
- повторный `GET /api/reading/today` возвращает ту же `cardOfDay`;
- для роли `user` лимит `10` блокируется на 11-й попытке (`429`);
- после смены роли и лимита генерация пересчитывается.
