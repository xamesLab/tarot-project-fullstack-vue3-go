<template>
  <div class="loader" aria-live="polite">{{ props.text ?? t("commonLoading") }}</div>
</template>

<script setup lang="ts">
import { useI18n } from "../composables/useI18n";

const props = defineProps<{
  text?: string;
}>();

const { t } = useI18n();
</script>

<style scoped>
.loader {
  min-height: 54px;
  display: grid;
  place-items: center;
  border: 1px solid var(--border-strong);
  border-radius: 18px;
  background:
    linear-gradient(135deg, rgba(210, 168, 74, 0.1) 0%, transparent 60%),
    linear-gradient(145deg, var(--surface-soft) 0%, var(--surface) 100%);
  color: var(--text-muted);
  letter-spacing: 0.12em;
  text-transform: uppercase;
}
</style>
