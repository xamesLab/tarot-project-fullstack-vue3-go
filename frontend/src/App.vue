<template>
  <div class="app-shell">
    <ModeSwitcher class="mode-switcher" />
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import ModeSwitcher from "./components/ModeSwitcher.vue";
</script>

<style scoped>
.app-shell {
  position: relative;
  min-height: 100vh;
}

.app-shell::before,
.app-shell::after {
  content: "";
  position: fixed;
  pointer-events: none;
  z-index: -1;
}

.app-shell::before {
  left: -12vw;
  bottom: 6vh;
  width: 36vw;
  height: 36vw;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(210, 168, 74, 0.16) 0%, transparent 68%);
  filter: blur(12px);
}

.app-shell::after {
  right: 8vw;
  top: 18vh;
  width: 18rem;
  height: 18rem;
  border-radius: 50%;
  border: 1px solid rgba(255, 248, 230, 0.1);
  box-shadow:
    0 0 0 14px rgba(255, 248, 230, 0.04),
    0 0 0 42px rgba(49, 95, 115, 0.03);
}

.mode-switcher {
  position: fixed;
  top: 14px;
  right: 14px;
  z-index: 20;
}
</style>
