<template>
  <main class="page stack">
    <AppHeader />

    <section class="panel stack">
      <div class="heading-row">
        <h2>{{ t("dashboardTitle") }}</h2>
        <div class="actions-row">
          <button class="btn btn-secondary" type="button" @click="reload" :disabled="readingStore.loading">
            {{ readingStore.loading ? t("dashboardReloading") : t("dashboardReload") }}
          </button>
          <button class="btn" type="button" @click="generate" :disabled="generateDisabled">
            {{ readingStore.generating ? t("dashboardGenerating") : t("dashboardGenerate") }}
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
        <SpreadView :title="t('spreadCardOfDay')" :cards="[readingStore.today.cardOfDay]" />
        <SpreadView :title="t('spreadA')" :cards="readingStore.today.spreadA" />
        <SpreadView :title="t('spreadB')" :cards="readingStore.today.spreadB" />
      </template>
      <p v-else class="muted">{{ t("dashboardNoReading") }}</p>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted } from "vue";

import { useI18n } from "../composables/useI18n";
import AppHeader from "../components/AppHeader.vue";
import BaseLoader from "../components/BaseLoader.vue";
import SpreadView from "../components/SpreadView.vue";
import { useAuthStore } from "../stores/auth.store";
import { useReadingStore } from "../stores/reading.store";

const authStore = useAuthStore();
const readingStore = useReadingStore();
const { t } = useI18n();
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

async function generate() {
  if (!authStore.token) return;
  await readingStore.generate(authStore.token);
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
</style>
