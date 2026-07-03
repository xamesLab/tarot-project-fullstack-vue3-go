<script setup>
import BaseSurface from '../base/BaseSurface.vue'

defineProps({
  title: {
    type: String,
    default: '',
  },
  description: {
    type: String,
    default: '',
  },
  symbol: {
    type: String,
    default: '✦',
  },
})
</script>

<template>
  <BaseSurface class="empty-card" tone="ghost" padding="lg" tag="section">
    <div class="empty-card__symbol">{{ symbol }}</div>
    <h3 class="empty-card__title">{{ title }}</h3>
    <p class="empty-card__copy">{{ description }}</p>
    <div v-if="$slots.actions" class="empty-card__actions">
      <slot name="actions" />
    </div>
  </BaseSurface>
</template>

<style scoped>
.empty-card {
  text-align: center;
}

.empty-card__symbol {
  margin-bottom: 18px;
  font-size: 32px;
  color: var(--ds-color-accent-primary);
}

.empty-card__title {
  margin-bottom: 10px;
  font-size: var(--ds-text-heading-2);
  color: var(--ds-color-accent-primary);
}

.empty-card__copy {
  max-width: 42ch;
  margin: 0 auto;
  color: var(--ds-color-text-secondary);
}

.empty-card__actions {
  margin-top: 24px;
  display: flex;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
}
</style>
