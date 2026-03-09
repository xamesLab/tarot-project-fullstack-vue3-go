<template>
  <section class="mode-switcher panel">
    <label class="field-inline" for="locale-switch">
      <span>{{ t("modeLanguage") }}</span>
      <select id="locale-switch" :value="uiPreferencesStore.locale" @change="onLocaleChange">
        <option value="ru">{{ t("modeLanguageRu") }}</option>
        <option value="en">{{ t("modeLanguageEn") }}</option>
      </select>
    </label>

    <label class="field-inline" for="theme-switch">
      <span>{{ t("modeTheme") }}</span>
      <select id="theme-switch" :value="uiPreferencesStore.theme" @change="onThemeChange">
        <option value="light">{{ t("modeThemeLight") }}</option>
        <option value="dark">{{ t("modeThemeDark") }}</option>
      </select>
    </label>
  </section>
</template>

<script setup lang="ts">
import type { Locale, ThemeMode } from "../i18n/messages";
import { useI18n } from "../composables/useI18n";
import { useUiPreferencesStore } from "../stores/ui-preferences.store";

const uiPreferencesStore = useUiPreferencesStore();
const { t } = useI18n();

function onLocaleChange(event: Event) {
  const locale = (event.target as HTMLSelectElement).value as Locale;
  uiPreferencesStore.setLocale(locale);
}

function onThemeChange(event: Event) {
  const theme = (event.target as HTMLSelectElement).value as ThemeMode;
  uiPreferencesStore.setTheme(theme);
}
</script>

<style scoped>
.mode-switcher {
  display: grid;
  grid-template-columns: repeat(2, minmax(128px, 1fr));
  gap: 10px;
  padding: 12px;
  border-radius: 18px;
  background:
    linear-gradient(135deg, rgba(210, 168, 74, 0.14) 0%, transparent 60%),
    linear-gradient(145deg, var(--surface) 0%, var(--surface-strong) 100%);
  box-shadow: var(--shadow-soft);
  backdrop-filter: blur(16px);
}

.field-inline {
  display: grid;
  gap: 6px;
  font-size: 0.82rem;
}

.field-inline select {
  min-height: 38px;
  border-radius: 12px;
  border: 1px solid var(--border-strong);
  padding: 0 10px;
  background: linear-gradient(180deg, var(--surface-soft) 0%, var(--surface) 100%);
  color: var(--text-main);
}

@media (max-width: 640px) {
  .mode-switcher {
    grid-template-columns: 1fr;
  }
}
</style>
