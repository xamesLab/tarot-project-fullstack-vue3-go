<template>
  <AppBackground :is-light="isLight" />
  <div class="app-shell">
    <ModeSwitcher class="mode-switcher" />
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

import AppBackground from "./components/AppBackground.vue";
import ModeSwitcher from "./components/ModeSwitcher.vue";
import { useUiPreferencesStore } from "./stores/ui-preferences.store";

const uiPreferencesStore = useUiPreferencesStore();
const isLight = computed(() => uiPreferencesStore.theme === "light");
</script>

<style scoped>
.app-shell {
  position: relative;
  z-index: 1;
  min-height: 100vh;
}

.mode-switcher {
  position: fixed;
  top: 14px;
  right: 14px;
  z-index: 50;
}
</style>
