<template>
  <section class="spread-view stack">
    <div class="spread-heading">
      <h3>{{ title }}</h3>
    </div>
    <p v-if="!cards.length" class="muted">{{ t("spreadNoCards") }}</p>
    <CardView v-for="card in cards" :key="card.id + card.state" :card="card" />
  </section>
</template>

<script setup lang="ts">
import { useI18n } from "../composables/useI18n";
import type { ReadingCard } from "../types/reading";
import CardView from "./CardView.vue";

defineProps<{
  title: string;
  cards: ReadingCard[];
}>();

const { t } = useI18n();
</script>

<style scoped>
.spread-view {
  position: relative;
  padding: 18px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background:
    linear-gradient(135deg, rgba(49, 95, 115, 0.1) 0%, transparent 45%),
    linear-gradient(145deg, var(--surface-soft) 0%, var(--surface) 100%);
  box-shadow: var(--shadow-soft);
}

.spread-heading {
  display: flex;
  align-items: center;
  gap: 14px;
}

.spread-heading::before,
.spread-heading::after {
  content: "";
  height: 1px;
  flex: 1;
  background: linear-gradient(90deg, transparent 0%, var(--border-strong) 50%, transparent 100%);
}

.spread-heading h3 {
  font-size: 1.05rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  white-space: nowrap;
}

@media (max-width: 640px) {
  .spread-heading {
    gap: 10px;
  }

  .spread-heading h3 {
    font-size: 0.95rem;
    white-space: normal;
    text-align: center;
  }
}
</style>
