<template>
  <main class="page stack">
    <AppHeader />

    <section class="panel stack">
      <h2>{{ t("settingsTitle") }}</h2>
      <p class="muted">{{ t("settingsHint") }}</p>

      <BaseLoader v-if="settingsStore.loading" :text="t('settingsLoading')" />

      <form v-else class="stack" @submit.prevent="save">
        <div class="field">
          <label for="timezone">{{ t("settingsTimezone") }}</label>
          <input id="timezone" v-model="settingsStore.timezone" placeholder="UTC+00:00 or Europe/Moscow" />
          <p class="muted">{{ t("settingsTimezoneHint") }}</p>
        </div>

        <ThemePicker v-model="settingsStore.deckThemeId" :themes="settingsStore.themes" />

        <p v-if="settingsStore.error" class="error">{{ settingsStore.error }}</p>
        <p v-if="settingsStore.success" class="muted">{{ settingsStore.success }}</p>
        <button class="btn" type="submit" :disabled="settingsStore.saving">
          {{ settingsStore.saving ? t("settingsSaving") : t("settingsSave") }}
        </button>
      </form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { onMounted } from "vue";

import { useI18n } from "../composables/useI18n";
import AppHeader from "../components/AppHeader.vue";
import BaseLoader from "../components/BaseLoader.vue";
import ThemePicker from "../components/ThemePicker.vue";
import { useAuthStore } from "../stores/auth.store";
import { useSettingsStore } from "../stores/settings.store";

const authStore = useAuthStore();
const settingsStore = useSettingsStore();
const { t } = useI18n();

onMounted(async () => {
  settingsStore.applyUserSettings(authStore.user);
  if (!authStore.token) return;
  await settingsStore.loadThemes(authStore.token);
});

async function save() {
  if (!authStore.token) return;
  const updated = await settingsStore.save(authStore.token);
  if (updated) {
    authStore.setUser(updated);
  }
}
</script>
