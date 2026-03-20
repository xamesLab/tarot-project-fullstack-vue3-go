<template>
  <section class="mode-switcher panel">
    <ModeToggle
      type="locale"
      :model-value="uiPreferencesStore.locale"
      @update:model-value="onLocaleChange"
    />
    <ModeToggle
      type="theme"
      :model-value="uiPreferencesStore.theme"
      @update:model-value="onThemeChange"
    />
  </section>
</template>

<script setup lang="ts">
import type { Locale, ThemeMode } from "../i18n/messages";
import { useUiPreferencesStore } from "../stores/ui-preferences.store";
import ModeToggle from "./ModeToggle.vue";

const uiPreferencesStore = useUiPreferencesStore();

function onLocaleChange(locale: string) {
  uiPreferencesStore.setLocale(locale as Locale);
}

function onThemeChange(theme: string) {
  uiPreferencesStore.setTheme(theme as ThemeMode);
}
</script>

<style scoped>
.mode-switcher.panel {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: 20px;
  border: 1px solid var(--border-strong);
  background: var(--surface);
  box-shadow: var(--shadow-soft);
  backdrop-filter: blur(12px);
  width: fit-content;
}

@media (max-width: 640px) {
  .mode-switcher.panel {
    padding: 8px;
    gap: 8px;
  }
}
</style>
