<template>
  <main class="page auth-page">
    <section class="panel stack auth-panel">
      <h2>{{ t("registerTitle") }}</h2>
      <p class="muted">{{ t("registerHint") }}</p>

      <form class="stack" @submit.prevent="submit">
        <div class="field">
          <label for="login">{{ t("fieldLogin") }}</label>
          <input id="login" v-model.trim="login" autocomplete="username" minlength="3" required />
        </div>

        <div class="field">
          <label for="password">{{ t("fieldPassword") }}</label>
          <input id="password" v-model="password" type="password" autocomplete="new-password" minlength="8" required />
        </div>

        <p v-if="authStore.error" class="error">{{ authStore.error }}</p>
        <button class="btn" type="submit" :disabled="authStore.loading">
          {{ authStore.loading ? t("registerSubmitting") : t("registerSubmit") }}
        </button>
      </form>

      <RouterLink to="/login" class="muted">{{ t("registerToLogin") }}</RouterLink>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

import { useI18n } from "../composables/useI18n";
import { useAuthStore } from "../stores/auth.store";

const authStore = useAuthStore();
const router = useRouter();
const { t } = useI18n();

const login = ref("");
const password = ref("");

async function submit() {
  await authStore.register(login.value, password.value);
  await router.push("/dashboard");
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
    linear-gradient(135deg, rgba(210, 168, 74, 0.14) 0%, transparent 48%),
    linear-gradient(145deg, var(--surface) 0%, var(--surface-strong) 100%);
}

.auth-panel::after {
  content: "";
  position: absolute;
  inset: -30px auto auto -28px;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 1px solid rgba(255, 248, 230, 0.08);
  box-shadow:
    0 0 0 12px rgba(255, 248, 230, 0.03),
    0 0 0 34px rgba(49, 95, 115, 0.05);
}

.auth-panel h2 {
  font-size: clamp(1.6rem, 2vw, 2rem);
  text-transform: uppercase;
  letter-spacing: 0.14em;
}
</style>
