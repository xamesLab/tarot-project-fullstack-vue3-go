<template>
  <main class="page stack">
    <AppHeader />

    <section class="panel stack">
      <div class="heading-row">
        <h2>{{ t("dashboardTitle") }}</h2>
        <div class="actions-row">
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
        </div>
      </div>
      <p class="muted">
        {{ t("dashboardCardOfDayNotice") }}
      </p>

      <BaseLoader v-if="readingStore.loading" :text="t('dashboardLoading')" />
      <p v-else-if="readingStore.error" class="error">{{ readingStore.error }}</p>

      <template v-else-if="readingStore.today">
        <p class="muted">
          {{ t("dashboardDate") }}: {{ readingStore.today.date }} | {{ t("dashboardTimezone") }}:
          {{ readingStore.today.timezoneUsed }}
        </p>
        <p class="muted">
          {{ t("dashboardSpreadsUsed") }}: {{ readingStore.today.spreadGenerationsUsed }}
          <template v-if="readingStore.today.spreadGenerationsLimit !== undefined">
            / {{ readingStore.today.spreadGenerationsLimit }} ({{ t("dashboardRemaining") }}:
            {{ readingStore.today.spreadGenerationsRemaining ?? 0 }})
          </template>
          <template v-else>
            / {{ t("dashboardUnlimited") }}
          </template>
        </p>

        <!-- Табы с раскладами -->
        <AppTabs v-model="activeTab" :tabs="tabs">
          <template #default="{ activeTab }">
            <!-- Карта дня -->
            <div v-if="activeTab === 'cardOfDay'" class="tab-content">
              <div class="tab-actions">
                <button
                  class="btn"
                  type="button"
                  @click="generateCardOfDay"
                  :disabled="readingStore.generating || !readingStore.today.canGenerateMore"
                >
                  {{ readingStore.generating ? t("generatingCardOfDay") : t("generateCardOfDay") }}
                </button>
              </div>
              <SpreadView :title="t('spreadCardOfDay')" :cards="[readingStore.today.cardOfDay]" />
            </div>

            <!-- Расклад A -->
            <div v-else-if="activeTab === 'spreadA'" class="tab-content">
              <div class="tab-actions">
                <button
                  class="btn"
                  type="button"
                  @click="generateSpreadA"
                  :disabled="readingStore.generating || !readingStore.today.canGenerateMore"
                >
                  {{ readingStore.generating ? t("generatingSpreadA") : t("generateSpreadA") }}
                </button>
              </div>
              <SpreadView :title="t('spreadA')" :cards="readingStore.today.spreadA" />
            </div>

            <!-- Расклад B -->
            <div v-else-if="activeTab === 'spreadB'" class="tab-content">
              <div class="tab-actions">
                <button
                  class="btn"
                  type="button"
                  @click="generateSpreadB"
                  :disabled="readingStore.generating || !readingStore.today.canGenerateMore"
                >
                  {{ readingStore.generating ? t("generatingSpreadB") : t("generateSpreadB") }}
                </button>
              </div>
              <SpreadView :title="t('spreadB')" :cards="readingStore.today.spreadB" />
            </div>
          </template>
        </AppTabs>
      </template>
      <p v-else class="muted">{{ t("dashboardNoReading") }}</p>
    </section>

    <!-- Модальное окно статуса Gateway -->
    <AppModal
      v-model="readingStore.isHealthModalOpen"
    >
      <template #header>
        <h2>{{ readingStore.healthStatus?.title || t("gatewayHealthModalTitle") }}</h2>
      </template>
      <div class="health-status" :class="readingStore.healthStatus?.status">
        <div class="health-icon">
          <span v-if="readingStore.healthStatus?.status === 'success'">✓</span>
          <span v-else>!</span>
        </div>
        <p class="health-message">
          {{ readingStore.healthStatus?.message || t("commonLoading") }}
        </p>
        <p v-if="readingStore.healthStatus?.uptime" class="health-uptime">
          Uptime: {{ readingStore.healthStatus.uptime }} сек.
        </p>
      </div>
    </AppModal>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import { useI18n } from "../composables/useI18n";
import AppHeader from "../components/AppHeader.vue";
import AppModal from "../components/AppModal.vue";
import AppTabs from "../components/AppTabs.vue";
import BaseLoader from "../components/BaseLoader.vue";
import SpreadView from "../components/SpreadView.vue";
import { useAuthStore } from "../stores/auth.store";
import { useReadingStore } from "../stores/reading.store";

const authStore = useAuthStore();
const readingStore = useReadingStore();
const { t } = useI18n();

const activeTab = ref("cardOfDay");

const tabs = computed(() => [
  { id: "cardOfDay", label: t("tabCardOfDay") },
  { id: "spreadA", label: t("tabSpreadA") },
  { id: "spreadB", label: t("tabSpreadB") }
]);

const generateDisabled = computed(() => {
  if (!authStore.token) return true;
  if (readingStore.loading || readingStore.generating) return true;
  if (!readingStore.today) return true;
  return !readingStore.today.canGenerateMore;
});

async function reload() {
  if (!authStore.token) return;
  await readingStore.loadToday(authStore.token);
}

async function generateCardOfDay() {
  if (!authStore.token) return;
  try {
    await readingStore.generateCardOfDayOnly(authStore.token);
  } catch (err) {
    // Ошибка обрабатывается в store
  }
}

async function generateSpreadA() {
  if (!authStore.token) return;
  try {
    await readingStore.generateSpreadAOnly(authStore.token);
  } catch (err) {
    // Ошибка обрабатывается в store
  }
}

async function generateSpreadB() {
  if (!authStore.token) return;
  try {
    await readingStore.generateSpreadBOnly(authStore.token);
  } catch (err) {
    // Ошибка обрабатывается в store
  }
}

async function checkGateway() {
  if (!authStore.token) return;
  try {
    await readingStore.checkGatewayHealth(authStore.token);
  } catch (err) {
    // Ошибка уже залогирована в store
  }
}

onMounted(async () => {
  await reload();
});
</script>

<style scoped>
.heading-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.heading-row h2 {
  margin: 0;
}

.actions-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

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

.tab-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.tab-actions {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
}

@media (max-width: 640px) {
  .tab-actions {
    justify-content: center;
  }
}
</style>
