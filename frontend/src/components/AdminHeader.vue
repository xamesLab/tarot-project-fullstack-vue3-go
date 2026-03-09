<template>
  <header class="admin-header panel">
    <div class="brand">
      <h1>{{ t("adminShellTitle") }}</h1>
      <p class="muted">{{ t("adminShellSubtitle") }}</p>
    </div>
    <nav class="nav">
      <RouterLink to="/admin">{{ t("adminTitle") }}</RouterLink>
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
.admin-header {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 18px;
  align-items: center;
  padding: 22px;
  border-color: var(--border-strong);
  background:
    linear-gradient(135deg, rgba(210, 168, 74, 0.16) 0%, transparent 46%),
    linear-gradient(145deg, var(--surface) 0%, var(--surface-strong) 100%);
}

.admin-header::after {
  content: "";
  position: absolute;
  left: -36px;
  top: -40px;
  width: 140px;
  height: 140px;
  border-radius: 50%;
  border: 1px solid rgba(255, 248, 230, 0.1);
  box-shadow:
    0 0 0 12px rgba(255, 248, 230, 0.04),
    0 0 0 34px rgba(210, 168, 74, 0.05);
  pointer-events: none;
}

.brand h1 {
  margin: 0;
  font-size: clamp(1.45rem, 1.9vw, 1.9rem);
  text-transform: uppercase;
  letter-spacing: 0.12em;
}

.brand p {
  margin: 8px 0 0;
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
  .admin-header {
    grid-template-columns: 1fr;
  }

  .nav {
    justify-content: flex-start;
  }
}
</style>
