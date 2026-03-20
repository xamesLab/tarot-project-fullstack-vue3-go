# Пошаговый план реализации: Gateway Health Check + Табы

**Версия:** v0.1.2
**Оценка времени:** 4-5 часов
**Статус:** ✅ Выполнено

---

## Часть 1: Gateway Health Check (✅ Выполнено)

Шаги 1-14 выполнены. Модальное окно работает.

---

## Часть 2: Табы для раскладов (✅ Выполнено)

### Реализовано

1. **Переиспользуемый компонент табов** — `AppTabs.vue` ✅
2. **3 вкладки на Dashboard:** ✅
   - Карта дня (первая, активна по умолчанию)
   - Расклад A (3 карты)
   - Расклад B (5 карт)
3. **Кнопка генерации** — отдельная на каждом табе ✅
4. **Кнопка «Проверить Gateway»** — остаётся в actions-row ✅
5. **Дизайн:** горизонтальные табы сверху ✅
6. **Backend endpoint'ы:** ✅
   - `POST /api/reading/generate/card-of-day`
   - `POST /api/reading/generate/spread-a`
   - `POST /api/reading/generate/spread-b`

### Требования

1. **Переиспользуемый компонент табов** — `AppTabs.vue`
2. **3 вкладки на Dashboard:**
   - Карта дня (первая, активна по умолчанию)
   - Расклад A (3 карты)
   - Расклад B (5 карт)
3. **Кнопка генерации** — отдельная на каждом табе (генерирует только активный расклад)
4. **Кнопка «Проверить Gateway»** — остаётся в actions-row (не меняется)
5. **Дизайн:** горизонтальные табы сверху
6. **Сохранение состояния:** не требуется

### План реализации

#### Шаг 15: Создать переиспользуемый компонент AppTabs.vue

**Файл:** `frontend/src/components/AppTabs.vue`

**Props:**
- `tabs`: массив объектов `{ id: string, label: string }`
- `modelValue`: id активного таба (v-model)

**Emits:**
- `update:modelValue`: при смене таба

**Функционал:**
- Горизонтальные табы сверху
- Подсветка активного таба
- Переключение по клику
- Поддержка клавиатуры (Left/Right arrows)

#### Шаг 16: Обновить reading.store.ts

**Изменения:**
- Добавить метод `generateCardOfDay(token)` — генерация только карты дня
- Добавить метод `generateSpreadA(token)` — генерация только расклада A
- Добавить метод `generateSpreadB(token)` — генерация только расклада B
- Метод `generate()` оставить для обратной совместимости (генерирует все)

#### Шаг 17: Обновить DashboardPage.vue

**Изменения:**
- Обернуть 3 блока в компонент `AppTabs`
- Добавить 3 таба: `cardOfDay`, `spreadA`, `spreadB`
- Показывать только контент активного таба
- Добавить 3 кнопки «Сгенерировать» (по одной на таб)
- Кнопка «Проверить Gateway» остаётся в actions-row

#### Шаг 18: Добавить ключи локализации

**Файл:** `frontend/src/i18n/messages.ts`

```ts
// ru
tabCardOfDay: "Карта дня",
tabSpreadA: "Расклад A",
tabSpreadB: "Расклад B",
generateCardOfDay: "Сгенерировать карту дня",
generateSpreadA: "Сгенерировать расклад A",
generateSpreadB: "Сгенерировать расклад B",

// en
tabCardOfDay: "Card of the Day",
tabSpreadA: "Spread A",
tabSpreadB: "Spread B",
generateCardOfDay: "Generate Card of the Day",
generateSpreadA: "Generate Spread A",
generateSpreadB: "Generate Spread B",
```

#### Шаг 19: Финальная проверка

- [ ] Компонент AppTabs переиспользуемый
- [ ] 3 вкладки отображаются корректно
- [ ] Карта дня открыта по умолчанию
- [ ] Кнопка генерации на каждом табе своя
- [ ] Генерируется только активный расклад
- [ ] Кнопка «Проверить Gateway» на месте
- [ ] TypeScript компилируется
- [ ] Сборка проходит

---

## Шаг 1: Добавить типы для Gateway Health

**Файл:** `frontend/src/types/reading.ts` или `frontend/src/types/gateway.ts`

```ts
// Добавить в конец файла
export interface GatewayHealthResponse {
  status: "ok";
  uptime: number;
}
```

---

## Шаг 2: Создать API-клиент для Gateway

**Файл:** `frontend/src/services/gateway.api.ts`

```ts
import type { GatewayHealthResponse } from "../types/gateway";
import { request } from "./apiClient";

export function getGatewayHealth(token: string): Promise<GatewayHealthResponse> {
  return request<GatewayHealthResponse>("/api/gateway/health", { token });
}
```

---

## Шаг 3: Добавить ключи локализации

**Файлы:**
- `frontend/src/i18n/locales/ru.json`
- `frontend/src/i18n/locales/en.json`

**ru.json:**
```json
{
  "gatewayHealthCheck": "Проверить Gateway",
  "gatewayHealthChecking": "Проверка...",
  "gatewayHealthSuccess": "Gateway доступен",
  "gatewayHealthError": "Ошибка проверки Gateway"
}
```

**en.json:**
```json
{
  "gatewayHealthCheck": "Check Gateway",
  "gatewayHealthChecking": "Checking...",
  "gatewayHealthSuccess": "Gateway is available",
  "gatewayHealthError": "Gateway check failed"
}
```

---

## Шаг 4: Расширить reading store (опционально)

**Файл:** `frontend/src/stores/reading.store.ts`

Добавить метод для проверки здоровья:

```ts
import { getGatewayHealth } from "../services/gateway.api";
import type { GatewayHealthResponse } from "../types/gateway";

// В теле defineStore:
const gatewayHealthLoading = ref(false);

async function checkGatewayHealth(token: string): Promise<GatewayHealthResponse> {
  gatewayHealthLoading.value = true;
  try {
    const result = await getGatewayHealth(token);
    console.log("Gateway health:", result);
    return result;
  } catch (err) {
    console.error("Gateway health error:", err);
    throw err;
  } finally {
    gatewayHealthLoading.value = false;
  }
}

// Вернуть из store:
return {
  // ...existing
  gatewayHealthLoading,
  checkGatewayHealth
};
```

---

## Шаг 5: Добавить кнопку в DashboardPage

**Файл:** `frontend/src/pages/DashboardPage.vue`

**Изменения в template:**
```vue
<div class="actions-row">
  <!-- Новая кнопка -->
  <button
    class="btn btn-secondary"
    type="button"
    @click="checkGateway"
    :disabled="readingStore.gatewayHealthLoading"
  >
    {{ readingStore.gatewayHealthLoading ? t("gatewayHealthChecking") : t("gatewayHealthCheck") }}
  </button>
  
  <button class="btn btn-secondary" type="button" @click="reload" :disabled="readingStore.loading">
    {{ readingStore.loading ? t("dashboardReloading") : t("dashboardReload") }}
  </button>
  <button class="btn" type="button" @click="generate" :disabled="generateDisabled">
    {{ readingStore.generating ? t("dashboardGenerating") : t("dashboardGenerate") }}
  </button>
</div>
```

**Изменения в script:**
```ts
async function checkGateway() {
  if (!authStore.token) return;
  try {
    await readingStore.checkGatewayHealth(authStore.token);
  } catch (err) {
    // Ошибка уже залогирована в store
  }
}
```

---

## Шаг 6: Альтернативный вариант (без изменения store)

Если не хотим менять store, можно вызвать API напрямую в компоненте:

**Файл:** `frontend/src/pages/DashboardPage.vue`

```ts
import { getGatewayHealth } from "../services/gateway.api";
import type { GatewayHealthResponse } from "../types/gateway";

const gatewayHealthLoading = ref(false);

async function checkGateway() {
  if (!authStore.token) return;
  gatewayHealthLoading.value = true;
  try {
    const result = await getGatewayHealth(authStore.token);
    console.log("Gateway health:", result);
  } catch (err) {
    console.error("Gateway health error:", err);
  } finally {
    gatewayHealthLoading.value = false;
  }
}
```

---

## Шаг 7: Проверка работы

1. Запустить frontend:
   ```bash
   cd frontend
   npm run dev
   ```

2. Авторизоваться в приложении

3. Открыть Dashboard

4. Нажать кнопку «Проверить Gateway»

5. Открыть DevTools Console (F12)

6. Убедиться, что в консоли появляется:
   - Успех: `Gateway health: { status: "ok", uptime: 12345.67 }`
   - Или ошибка: `Gateway health error: ...`

---

## Шаг 8: Финальная проверка

- [ ] TypeScript компилируется: `npm run type-check`
- [ ] ESLint не ругается: `npm run lint`
- [ ] Кнопка отображается и кликабельна
- [ ] Запрос уходит с правильным Authorization header
- [ ] Ответ выводится в консоль
- [ ] Кнопка блокируется во время запроса

---

## Структура файлов после изменений

```
frontend/src/
├── types/
│   ├── reading.ts          (или gateway.ts)
│   └── ...
├── services/
│   ├── gateway.api.ts      (новый)
│   ├── reading.api.ts
│   └── ...
├── stores/
│   ├── reading.store.ts    (расширен)
│   └── ...
├── pages/
│   └── DashboardPage.vue   (изменён)
└── i18n/locales/
    ├── ru.json             (дополнен)
    └── en.json             (дополнен)
```

---

## Рекомендации

1. **Предпочтительный вариант:** Шаг 5 (с изменением store) — лучше для переиспользования
2. **Альтернатива:** Шаг 6 (без store) — быстрее, но менее масштабируемо
3. **Будущее:** Вынести статус Gateway в отдельный store для глобального состояния

---

## Шаг 9: Создать переиспользуемый компонент модального окна

**Файл:** `frontend/src/components/AppModal.vue`

**Требования:**
- Ширина и высота: не более 70% экрана (`max-width: 70vw`, `max-height: 70vh`)
- Фон: затемнённый overlay (`rgba(0, 0, 0, 0.5)`)
- Кнопка закрытия: крестик в верхнем правом углу
- Заголовок: по центру в верхней части
- Поддержка слотов: `header` (заголовок), `default` (контент), `footer` (опционально)
- Закрытие по клику на overlay и по ESC

```vue
<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="modal-overlay" @click="onOverlayClick">
        <div class="modal-container" role="dialog" aria-modal="true" @click.stop>
          <button class="modal-close" type="button" @click="close" aria-label="Close">
            <span class="modal-close-icon">&times;</span>
          </button>
          <header v-if="$slots.header" class="modal-header">
            <slot name="header" />
          </header>
          <div class="modal-content">
            <slot />
          </div>
          <footer v-if="$slots.footer" class="modal-footer">
            <slot name="footer" />
          </footer>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from "vue";

const props = defineProps<{
  modelValue: boolean;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
}>();

function close() {
  emit("update:modelValue", false);
}

function onOverlayClick() {
  close();
}

function handleEscapeKey(event: KeyboardEvent) {
  if (event.key === "Escape" && props.modelValue) {
    close();
  }
}

onMounted(() => {
  document.addEventListener("keydown", handleEscapeKey);
});

onUnmounted(() => {
  document.removeEventListener("keydown", handleEscapeKey);
});
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.modal-container {
  position: relative;
  background: var(--surface-elevated, #fff);
  border-radius: var(--radius, 12px);
  max-width: 70vw;
  max-height: 70vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.modal-close {
  position: absolute;
  top: 12px;
  right: 12px;
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.modal-close:hover {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 50%;
}

.modal-close-icon {
  font-size: 1.5rem;
  line-height: 1;
  color: var(--text-muted, #666);
}

.modal-header {
  text-align: center;
  padding: 20px 24px 12px;
  border-bottom: 1px solid var(--border, rgba(0, 0, 0, 0.1));
}

.modal-header h2 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.modal-content {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--border, rgba(0, 0, 0, 0.1));
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

/* Transition */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.95);
}

@media (max-width: 768px) {
  .modal-container {
    max-width: 90vw;
    max-height: 80vh;
  }
}
</style>
```

---

## Шаг 10: Добавить типы для статуса Gateway в модальное окно

**Файл:** `frontend/src/types/reading.ts` (дополнение)

```ts
export interface GatewayHealthStatus {
  status: "success" | "error";
  title: string;
  message: string;
  uptime?: number;
}
```

---

## Шаг 11: Расширить reading.store.ts для хранения статуса

**Файл:** `frontend/src/stores/reading.store.ts`

Добавить состояние для модального окна и результата:

```ts
import type { GatewayHealthResponse, GatewayHealthStatus } from "../types/reading";

// В теле defineStore:
const isHealthModalOpen = ref(false);
const healthStatus = ref<GatewayHealthStatus | null>(null);

async function checkGatewayHealth(token: string): Promise<GatewayHealthResponse> {
  gatewayHealthLoading.value = true;
  try {
    const result = await getGatewayHealth(token);
    healthStatus.value = {
      status: "success",
      title: "Gateway доступен",
      message: `Время работы: ${formatUptime(result.uptime)}`,
      uptime: result.uptime
    };
    isHealthModalOpen.value = true;
    return result;
  } catch (err) {
    const errorMessage = err instanceof ApiError ? err.message : "Неизвестная ошибка";
    healthStatus.value = {
      status: "error",
      title: "Ошибка проверки Gateway",
      message: errorMessage
    };
    isHealthModalOpen.value = true;
    throw err;
  } finally {
    gatewayHealthLoading.value = false;
  }
}

function closeHealthModal() {
  isHealthModalOpen.value = false;
  healthStatus.value = null;
}

// Утилита для форматирования uptime
function formatUptime(seconds: number): string {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  return `${hours}ч ${minutes}м`;
}

// Вернуть из store:
return {
  // ...existing
  isHealthModalOpen,
  healthStatus,
  checkGatewayHealth,
  closeHealthModal
};
```

---

## Шаг 12: Обновить DashboardPage с модальным окном

**Файл:** `frontend/src/pages/DashboardPage.vue`

**Изменения в template:**
```vue
<template>
  <main class="page stack">
    <AppHeader />

    <section class="panel stack">
      <!-- ... existing heading-row ... -->
    </section>

    <!-- Модальное окно -->
    <AppModal
      v-model="readingStore.isHealthModalOpen"
      :header="readingStore.healthStatus?.title || 'Статус Gateway'"
    >
      <div class="health-status" :class="readingStore.healthStatus?.status">
        <div class="health-icon">
          <span v-if="readingStore.healthStatus?.status === 'success'">✓</span>
          <span v-else>!</span>
        </div>
        <p class="health-message">
          {{ readingStore.healthStatus?.message || "Загрузка..." }}
        </p>
        <p v-if="readingStore.healthStatus?.uptime" class="health-uptime">
          Uptime: {{ readingStore.healthStatus.uptime }} сек.
        </p>
      </div>
    </AppModal>
  </main>
</template>
```

**Изменения в script:**
```ts
import AppModal from "../components/AppModal.vue";

// Функция checkGateway теперь не нужна - логика в store
// Оставить только кнопку, которая вызывает readingStore.checkGatewayHealth()
```

**Стили для health-status:**
```vue
<style scoped>
.health-status {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 16px;
  padding: 20px 0;
}

.health-icon {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  font-weight: bold;
}

.health-status.success .health-icon {
  background: rgba(34, 197, 94, 0.1);
  color: rgb(34, 197, 94);
}

.health-status.error .health-icon {
  background: rgba(239, 68, 68, 0.1);
  color: rgb(239, 68, 68);
}

.health-message {
  font-size: 1.1rem;
  margin: 0;
}

.health-uptime {
  font-size: 0.9rem;
  color: var(--text-muted);
  margin: 0;
}
</style>
```

---

## Шаг 13: Добавить ключи локализации для модального окна

**Файл:** `frontend/src/i18n/messages.ts`

**ru:**
```ts
gatewayHealthModalTitle: "Статус Gateway",
gatewayHealthSuccessTitle: "Gateway доступен",
gatewayHealthErrorTitle: "Ошибка проверки Gateway",
gatewayHealthUptime: "Время работы: {value}",
gatewayHealthClose: "Закрыть"
```

**en:**
```ts
gatewayHealthModalTitle: "Gateway Status",
gatewayHealthSuccessTitle: "Gateway is available",
gatewayHealthErrorTitle: "Gateway check failed",
gatewayHealthUptime: "Uptime: {value}",
gatewayHealthClose: "Close"
```

---

## Шаг 14: Финальная проверка

- [ ] Компонент AppModal переиспользуемый (можно использовать на других страницах)
- [ ] Модальное окно закрывается по клику на overlay
- [ ] Модальное окно закрывается по ESC
- [ ] Кнопка-крестик работает
- [ ] Ширина/высота не более 70%
- [ ] Фон затемнён
- [ ] Заголовок по центру
- [ ] Статус Gateway отображается в модальном окне (не в консоли)
- [ ] TypeScript компилируется: `npm run typecheck`
- [ ] Сборка проходит: `npm run build`

---

## Структура файлов после всех изменений

```
frontend/src/
├── types/
│   ├── reading.ts              (добавлены GatewayHealthResponse, GatewayHealthStatus)
│   └── ...
├── services/
│   ├── gateway.api.ts          (новый)
│   ├── reading.api.ts
│   └── ...
├── components/
│   ├── AppModal.vue            (новый, переиспользуемый)
│   ├── CardView.vue
│   └── ...
├── stores/
│   ├── reading.store.ts        (расширен: isHealthModalOpen, healthStatus, closeHealthModal)
│   └── ...
├── pages/
│   └── DashboardPage.vue       (изменён: добавлено модальное окно)
└── i18n/
    └── messages.ts             (дополнен ключами для модального окна)
```

---

## Рекомендации

1. **Предпочтительный вариант:** Шаг 5 (с изменением store) — лучше для переиспользования
2. **Альтернатива:** Шаг 6 (без store) — быстрее, но менее масштабируемо
3. **AppModal:** Универсальный компонент, можно использовать для любых диалогов
4. **Будущее:** Вынести статус Gateway в отдельный store для глобального состояния
5. **Будущее:** Добавить анимацию появления модального окна (уже есть в шаблоне)
