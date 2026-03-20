<template>
  <section class="mode-switcher panel">
    <ModeToggle type="locale" :model-value="uiPreferencesStore.locale" @update:model-value="onLocaleChange" />
    <ModeToggle type="theme" :model-value="uiPreferencesStore.theme" @update:model-value="onThemeChange" />
  </section>
</template>

<script setup lang="ts">
import type { Locale, ThemeMode } from "../i18n/messages";
import { useI18n } from "../composables/useI18n";
import { useUiPreferencesStore } from "../stores/ui-preferences.store";
import ModeToggle from "./ModeToggle.vue";

const uiPreferencesStore = useUiPreferencesStore();
const { t } = useI18n();

function onLocaleChange(locale: string) {
  uiPreferencesStore.setLocale(locale as Locale);
}

function onThemeChange(theme: string) {
  uiPreferencesStore.setTheme(theme as ThemeMode);
}
</script>

<style scoped>
.mode-switcher {
  display: flex;
  gap: 8px;
  padding: 12px;
  border-radius: 18px;
  background:
    linear-gradient(135deg, rgba(210, 168, 74, 0.14) 0%, transparent 60%),
    linear-gradient(145deg, var(--surface) 0%, var(--surface-strong) 100%);
  box-shadow: var(--shadow-soft);
  backdrop-filter: blur(16px);
  width: fit-content;
}

@media (max-width: 640px) {
  .mode-switcher {
    flex-direction: column;
    width: 100%;
  }
}
</style>
