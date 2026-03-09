<template>
  <main class="page">
    <section class="panel stack">
      <h2>{{ t("notFoundTitle") }}</h2>
      <p class="muted">{{ t("notFoundHint") }}</p>
      <RouterLink class="btn" to="/dashboard">{{ t("notFoundBack") }}</RouterLink>
    </section>
  </main>
</template>

<script setup lang="ts">
import { useI18n } from "../composables/useI18n";

const { t } = useI18n();
</script>
