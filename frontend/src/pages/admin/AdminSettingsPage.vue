<template>
  <section class="panel stack">
    <div class="heading-row">
      <h2>{{ t("adminSettingsTitle") }}</h2>
      <button class="btn btn-secondary" type="button" @click="load" :disabled="loading || saving">
        {{ loading ? t("adminRefreshing") : t("adminRefresh") }}
      </button>
    </div>
    <p class="muted">{{ t("adminSettingsHint") }}</p>

    <BaseLoader v-if="loading" :text="t('adminSettingsLoading')" />

    <form v-else class="stack" @submit.prevent="save">
      <div class="field">
        <label for="project-title">{{ t("adminSettingsProjectTitle") }}</label>
        <input id="project-title" v-model.trim="form.projectTitle" minlength="2" maxlength="120" required />
      </div>

      <div class="field">
        <label for="support-email">{{ t("adminSettingsSupportEmail") }}</label>
        <input id="support-email" v-model.trim="form.supportEmail" type="email" />
      </div>

      <label class="field-inline">
        <input v-model="form.maintenanceMode" type="checkbox" />
        <span>{{ t("adminSettingsMaintenanceMode") }}</span>
      </label>

      <div class="settings-grid">
        <div class="field">
          <label for="day-reset-hour">{{ t("adminSettingsDayResetHour") }}</label>
          <input id="day-reset-hour" v-model.number="form.dayResetHour" type="number" min="0" max="23" required />
        </div>
        <div class="field">
          <label for="max-upload-mb">{{ t("adminSettingsMaxUploadMb") }}</label>
          <input id="max-upload-mb" v-model.number="form.maxUploadMb" type="number" min="1" max="100" required />
        </div>
      </div>

      <p v-if="error" class="error">{{ error }}</p>
      <p v-if="success" class="muted">{{ success }}</p>

      <button class="btn" type="submit" :disabled="saving">
        {{ saving ? t("adminSettingsSaving") : t("adminSettingsSave") }}
      </button>
    </form>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";

import { useI18n } from "../../composables/useI18n";
import * as adminApi from "../../services/admin.api";
import { useAuthStore } from "../../stores/auth.store";
import BaseLoader from "../../components/BaseLoader.vue";

const authStore = useAuthStore();
const { t } = useI18n();

const loading = ref(false);
const saving = ref(false);
const error = ref("");
const success = ref("");

const form = reactive({
  projectTitle: "",
  supportEmail: "",
  maintenanceMode: false,
  dayResetHour: 3,
  maxUploadMb: 8
});

async function load() {
  if (!authStore.token) return;
  loading.value = true;
  error.value = "";
  success.value = "";
  try {
    const data = await adminApi.getAdminSettings(authStore.token);
    form.projectTitle = data.projectTitle;
    form.supportEmail = data.supportEmail;
    form.maintenanceMode = data.maintenanceMode;
    form.dayResetHour = data.dayResetHour;
    form.maxUploadMb = data.maxUploadMb;
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to load admin settings";
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!authStore.token) return;
  saving.value = true;
  error.value = "";
  success.value = "";
  try {
    const updated = await adminApi.patchAdminSettings(authStore.token, {
      projectTitle: form.projectTitle,
      supportEmail: form.supportEmail,
      maintenanceMode: form.maintenanceMode,
      dayResetHour: form.dayResetHour,
      maxUploadMb: form.maxUploadMb
    });
    form.projectTitle = updated.projectTitle;
    form.supportEmail = updated.supportEmail;
    form.maintenanceMode = updated.maintenanceMode;
    form.dayResetHour = updated.dayResetHour;
    form.maxUploadMb = updated.maxUploadMb;
    success.value = t("adminSettingsSaved");
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to save admin settings";
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>

<style scoped>
.heading-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.heading-row h2 {
  margin: 0;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 10px;
}
</style>
