# Чек-лист реализации v0.1.2 (без привязки ко времени)

- [x] Шаг 1. Зафиксировать контракт интеграции `Frontend -> Backend -> Gateway` в документах `v0.1.2` (source: `INTEGRATION_INSTRUCTION.ru.md`).

- [x] Шаг 2. Спроектировать модель `ask_requests` и индексы (`requestId` unique, `userId+createdAt`).
- [x] Шаг 3. Реализовать repository-слой для `ask_requests` (create/get/update status/result/error).

- [x] Шаг 4. Добавить конфигурацию интеграции в backend (`GATEWAY_BASE_URL=https://api.xmlab.ru`, `GATEWAY_REQUEST_PATH`, `GATEWAY_TIMEOUT`, `BACKEND_PUBLIC_BASE_URL`, callback security vars).
- [x] Шаг 5. Обновить `backend/.env.example` под новые переменные.

- [x] Шаг 6. Реализовать `gatewayclient` для `POST /request` с таймаутом и обработкой ошибок.
- [x] Шаг 7. Реализовать сервис оркестрации `ask` (создание запроса, вызов Gateway, чтение статуса).

- [x] Шаг 8. Добавить endpoint `POST /api/ask`:
  - auth-required;
  - валидация payload;
  - генерация `request_id`;
  - dispatch в Gateway;
  - ответ `accepted`.

- [x] Шаг 9. Добавить endpoint `GET /api/ask/:request_id`:
  - проверка доступа пользователя к своему запросу;
  - возврат `status/result/error`.

- [x] Шаг 10. Добавить endpoint `POST /api/gateway/callback`:
  - проверка подлинности callback (allowlist/hmac);
  - idempotent update запроса;
  - поддержка `success|error`.

- [x] Шаг 11. Добавить middleware/helper для callback verification и единый маппинг ошибок интеграции.

- [x] Шаг 12. Обновить роутинг, bootstrap и DI-контейнер backend под новые сервисы/handler’ы.

- [x] Шаг 13. Написать unit-тесты:
  - сервис `ask`;
  - callback idempotency;
  - callback verification;
  - gateway client error mapping.

- [x] Шаг 14. Добавить интеграционный smoke-скрипт `ask-flow`:
  - `POST /api/ask`;
  - имитация callback `success`;
  - проверка `GET /api/ask/:request_id`.

- [x] Шаг 15. Подготовить сценарий error-flow:
  - callback `status=error`;
  - корректная выдача ошибки frontend через polling endpoint.

- [x] Шаг 16. Обновить backend документацию endpoint’ов и env:
  - `backend/README.md`;
  - примеры запросов/ответов для frontend.

- [x] Шаг 17. Обновить `testing-report-v0.1.2.md` по фактическим проверкам.
- [x] Шаг 18. Обновить `acceptance-criteria-check-v0.1.2.md` и зафиксировать `final-status-v0.1.2.md`.
