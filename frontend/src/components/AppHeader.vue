<template>
  <header class="app-header panel">
    <div class="brand">
      <h1>Tarot Daily</h1>
      <p class="muted">{{ t("appSubtitle") }}</p>
    </div>
    <nav class="nav">
      <RouterLink to="/dashboard">{{ t("navDashboard") }}</RouterLink>
      <RouterLink to="/profile">{{ t("navProfile") }}</RouterLink>
      <RouterLink to="/settings">{{ t("navSettings") }}</RouterLink>
      <button class="btn btn-secondary" type="button" @click="logout">{{ t("navLogout") }}</button>
    </nav>
  </header>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";

import { useI18n } from "../composables/useI18n";
import { useAuthStore } from "../stores/auth.store";

const authStore = useAuthStore();
const router = useRouter();
const { t } = useI18n();

function logout() {
  authStore.logout();
  router.push({ name: "login" });
}
</script>

<style scoped>
.app-header {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 18px;
  align-items: center;
  padding: 22px;
  background:
    linear-gradient(135deg, rgba(49, 95, 115, 0.22) 0%, transparent 45%),
    linear-gradient(145deg, var(--surface) 0%, var(--surface-strong) 100%);
}

.app-header::after {
  content: "";
  position: absolute;
  right: -18px;
  bottom: -24px;
  width: 132px;
  height: 132px;
  border-radius: 50%;
  border: 1px solid rgba(255, 248, 230, 0.12);
  box-shadow:
    0 0 0 14px rgba(255, 248, 230, 0.04),
    0 0 0 38px rgba(210, 168, 74, 0.06);
  pointer-events: none;
}

.brand h1 {
  margin: 0;
  font-size: clamp(1.5rem, 2vw, 2rem);
  text-transform: uppercase;
  letter-spacing: 0.14em;
}

.brand p {
  margin: 8px 0 0;
  max-width: 42rem;
  line-height: 1.65;
}

.nav {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
  justify-content: flex-end;
}

.nav a {
  padding: 10px 14px;
  border: 1px solid var(--border-strong);
  border-radius: 999px;
  background: rgba(255, 248, 230, 0.08);
  backdrop-filter: blur(12px);
  font-weight: 700;
  letter-spacing: 0.04em;
}

.nav a:hover,
.nav a.router-link-active {
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-deep) 100%);
  color: var(--accent-contrast);
  border-color: transparent;
}

@media (max-width: 820px) {
  .app-header {
    grid-template-columns: 1fr;
  }

  .nav {
    justify-content: flex-start;
  }
}
</style>
