<template>
  <article class="card-view panel" :class="{ 'is-revealed': revealed }">
    <button class="card-button" type="button" @click="toggleReveal">
      <img :src="activeImageUrl" :alt="revealed ? card.name : t('cardBackAlt')" class="card-image" />
    </button>
    <div class="card-copy stack">
      <p v-if="!revealed" class="muted reveal-hint">{{ t("cardRevealHint") }}</p>
      <template v-else>
        <strong class="card-title">{{ card.name }}</strong>
        <p class="card-interpretation">{{ card.interpretation }}</p>
      </template>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue";

import { useI18n } from "../composables/useI18n";
import type { ReadingCard } from "../types/reading";

const props = defineProps<{
  card: ReadingCard;
}>();

const { t } = useI18n();
const revealed = ref(false);

const activeImageUrl = computed(() => {
  return revealed.value ? props.card.imageUrl : props.card.backImageUrl;
});

function toggleReveal() {
  revealed.value = !revealed.value;
}

watch(
  () => [props.card.id, props.card.state, props.card.imageUrl, props.card.backImageUrl],
  () => {
    revealed.value = false;
  }
);
</script>

<style scoped>
.card-view {
  display: grid;
  grid-template-columns: minmax(92px, 108px) 1fr;
  gap: 18px;
  align-items: start;
  padding: 18px;
  background:
    linear-gradient(135deg, rgba(210, 168, 74, 0.08) 0%, transparent 52%),
    linear-gradient(145deg, var(--surface) 0%, var(--surface-elevated) 100%);
  transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;
}

.card-view:hover {
  transform: translateY(-2px);
  box-shadow: 0 18px 36px rgba(25, 18, 31, 0.2);
}

.card-view.is-revealed {
  border-color: var(--border-strong);
}

.card-button {
  position: relative;
  border: 0;
  background: transparent;
  padding: 0;
  cursor: pointer;
}

.card-button::before {
  content: "";
  position: absolute;
  inset: -10px;
  border-radius: 18px;
  background: radial-gradient(circle, rgba(210, 168, 74, 0.18) 0%, transparent 70%);
  opacity: 0;
  transition: opacity 0.18s ease;
  pointer-events: none;
}

.card-view:hover .card-button::before,
.card-view.is-revealed .card-button::before {
  opacity: 1;
}

.card-image {
  width: 100%;
  aspect-ratio: 2 / 3;
  object-fit: cover;
  border-radius: 14px;
  border: 1px solid var(--border-strong);
  background: var(--surface-soft);
  box-shadow: 0 18px 28px rgba(25, 18, 31, 0.24);
}

.card-copy {
  min-height: 100%;
  align-content: center;
}

.reveal-hint {
  align-self: center;
  font-size: 0.96rem;
}

.card-title {
  font-family: var(--font-display);
  font-size: 1.16rem;
  letter-spacing: 0.05em;
}

.card-interpretation {
  line-height: 1.65;
}

@media (max-width: 640px) {
  .card-view {
    grid-template-columns: 1fr;
    justify-items: center;
    text-align: center;
  }

  .card-button {
    width: min(132px, 100%);
  }
}
</style>
