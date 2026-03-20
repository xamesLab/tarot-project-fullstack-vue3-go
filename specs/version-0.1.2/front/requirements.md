# Frontend Task: Gateway Health Check в карточке «Карта дня»

**Версия:** v0.1.2  
**Дата:** 20 марта 2026  
**Статус:** ✅ Completed (с модальным окном)

---

## 1. Контекст задачи

В рамках версии `v0.1.2` backend предоставляет endpoint для проверки доступности Gateway:

- `GET /api/gateway/health` — проксирует запрос к `https://api.xmlab.ru/health`

Этот endpoint используется для диагностики: frontend не должен обращаться к Gateway напрямую.

**Текущее состояние:**
- На DashboardPage отображается карточка «Карта дня» через компонент `SpreadView`
- Кнопка проверки здоровья Gateway отсутствует

**Целевое состояние:**
- В карточке «Карта дня» добавлена кнопка «Проверить Gateway»
- По клику выполняется запрос к `/api/gateway/health`
- Результат выводится в консоль браузера (для отладки)

---

## 2. Требования

### 2.1 Функциональные требования

1. **Кнопка проверки здоровья**
   - Расположение: в карточке «Карта дня» (рядом с заголовком или в actions-row)
   - Текст кнопки: «Проверить Gateway» (с локализацией)
   - Состояние disabled во время запроса

2. **API-запрос**
   - Метод: `GET /api/gateway/health`
   - Auth: требуется Bearer token
   - Обработка ошибок: корректный маппинг статусов (401, 502, 504)

3. **Модальное окно**
   - Отображение результата в модальном окне (не в консоли)
   - Компонент переиспользуемый (можно использовать на других страницах)
   - Требования к дизайну:
     - Ширина и высота: не более 70% экрана
     - Фон: затемнённый overlay
     - Кнопка закрытия: крестик в верхнем правом углу
     - Заголовок: по центру в верхней части
     - Поддержка слотов: header, default, footer (опционально)
     - Закрытие по клику на overlay и по ESC
   - Отображение статуса:
     - Успех: зелёная иконка ✓, сообщение о времени работы
     - Ошибка: красная иконка !, сообщение об ошибке

4. **Вывод результата**
   - Успех (`200 OK`): модальное окно с заголовком «Gateway доступен», uptime в формате «Xч Yм»
   - Ошибка: модальное окно с заголовком «Ошибка проверки Gateway», текст ошибки

### 2.2 Технические требования

1. **Архитектура**
   -遵循现有项目约定 (follow existing conventions)
   - API-клиент: создать/расширить `services/gateway.api.ts`
   - Store: расширить `stores/reading.store.ts` или создать отдельный `stores/gateway.store.ts`
   - Компонент: добавить кнопку в `SpreadView.vue` или `DashboardPage.vue`

2. **Типизация**
   - Использовать TypeScript
   - Добавить типы в `types/reading.ts` или `types/gateway.ts`

3. **Локализация**
   - Добавить ключи в i18n словари (ru/en)

4. **Обработка ошибок**
   - Использовать существующий `ApiError` класс
   - Не показывать пользователю alert/toast (только консоль для v0.1.2)

---

## 3. Контракты API

### GET /api/gateway/health

**Request:**
```http
GET /api/gateway/health
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "status": "ok",
  "uptime": 81720.198966473
}
```

**Error responses:**
```json
// 401 Unauthorized
{
  "error": "unauthorized",
  "message": "Invalid or missing token"
}

// 504 Gateway Timeout
{
  "error": "gateway_timeout",
  "message": "Backend timed out while connecting to Gateway"
}

// 502 Bad Gateway
{
  "error": "gateway_bad_response",
  "message": "Gateway returned invalid response"
}
```

---

## 4. Типы данных (TypeScript)

```ts
export interface GatewayHealthResponse {
  status: "ok";
  uptime: number;
}

export interface GatewayHealthError {
  error: string;
  message: string;
  requestId?: string;
}
```

---

## 5. Критерии приемки (Acceptance Criteria)

- [ ] Кнопка «Проверить Gateway» отображается в карточке «Карта дня»
- [ ] По клику выполняется `GET /api/gateway/health` с Authorization header
- [ ] Результат отображается в модальном окне (не в консоли)
- [ ] Модальное окно соответствует требованиям:
  - [ ] Ширина/высота не более 70%
  - [ ] Фон затемнён
  - [ ] Кнопка-крестик в верхнем правом углу
  - [ ] Заголовок по центру
  - [ ] Закрытие по клику на overlay
  - [ ] Закрытие по ESC
- [ ] Успешный ответ: зелёная иконка, заголовок «Gateway доступен», uptime
- [ ] Ошибка: красная иконка, заголовок «Ошибка проверки Gateway», текст ошибки
- [ ] Кнопка блокируется на время запроса (disabled state)
- [ ] Код соответствует существующим конвенциям проекта
- [ ] TypeScript компилируется без ошибок
- [ ] Локализация добавлена для ru/en
- [ ] Компонент AppModal переиспользуемый

---

## 6. Риски и заметки

1. **Auth token:** Убедиться, что токен доступен в момент вызова (пользователь авторизован)
2. **CORS:** Backend должен корректно проксировать запросы к Gateway
3. **Будущее расширение:** В следующих версиях можно добавить UI-отображение статуса (toast/indicator)

---

## 7. Ссылки

- `specs/version-0.1.2/frontend-api-integration-v0.1.2.md` — раздел 3.3 про `/api/gateway/health`
- `specs/version-0.1.2/INTEGRATION_INSTRUCTION.ru.md` — общая архитектура
- `frontend/src/services/reading.api.ts` — пример API-клиента
- `frontend/src/stores/reading.store.ts` — пример store
