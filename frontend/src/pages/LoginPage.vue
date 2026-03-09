<template>
  <main class="page auth-page">
    <section class="panel stack auth-panel">
      <h2>{{ t("loginTitle") }}</h2>
      <p class="muted">{{ t("loginHint") }}</p>

      <form class="stack" @submit.prevent="submit">
        <div class="field">
          <label for="login">{{ t("fieldLogin") }}</label>
          <input id="login" v-model.trim="login" autocomplete="username" required />
        </div>

        <div class="field">
          <label for="password">{{ t("fieldPassword") }}</label>
          <input id="password" v-model="password" type="password" autocomplete="current-password" required />
        </div>

        <p v-if="authStore.error" class="error">{{ authStore.error }}</p>
        <button class="btn" type="submit" :disabled="authStore.loading">
          {{ authStore.loading ? t("loginSubmitting") : t("loginSubmit") }}
        </button>
      </form>

      <RouterLink to="/register" class="muted">{{ t("loginToRegister") }}</RouterLink>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";

import { useI18n } from "../composables/useI18n";
import { useAuthStore } from "../stores/auth.store";

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();
const { t } = useI18n();

const login = ref("");
const password = ref("");

async function submit() {
  await authStore.login(login.value, password.value);
  const redirectTo = typeof route.query.redirect === "string" ? route.query.redirect : "/dashboard";
  await router.push(redirectTo);
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: grid;
  align-content: start;
}

.auth-panel {
  max-width: 460px;
  margin: 0 auto;
  padding: 26px;
  background:
    linear-gradient(135deg, rgba(49, 95, 115, 0.14) 0%, transparent 48%),
    linear-gradient(145deg, var(--surface) 0%, var(--surface-strong) 100%);
}

.auth-panel::after {
  content: "";
  position: absolute;
  inset: auto -26px -38px auto;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 1px solid rgba(255, 248, 230, 0.08);
  box-shadow:
    0 0 0 12px rgba(255, 248, 230, 0.03),
    0 0 0 34px rgba(210, 168, 74, 0.05);
}

.auth-panel h2 {
  font-size: clamp(1.6rem, 2vw, 2rem);
  text-transform: uppercase;
  letter-spacing: 0.14em;
}
</style>
