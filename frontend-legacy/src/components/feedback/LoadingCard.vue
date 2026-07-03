<script setup>
import BaseSurface from '../base/BaseSurface.vue'

defineProps({
  title: {
    type: String,
    default: 'Loading',
  },
  copy: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <BaseSurface class="loading-card" tone="ghost" padding="lg" tag="section">
    <div class="loading-card__pulse"></div>
    <h3 class="loading-card__title">{{ title }}</h3>
    <p v-if="copy" class="loading-card__copy">{{ copy }}</p>
  </BaseSurface>
</template>

<style scoped>
.loading-card {
  text-align: center;
}

.loading-card__pulse {
  width: 54px;
  height: 54px;
  margin: 0 auto 18px;
  border: 1px solid var(--ds-color-border-accent);
  animation: ds-pulse 1.4s ease-in-out infinite;
}

.loading-card__title {
  margin-bottom: 10px;
  font-size: var(--ds-text-heading-2);
  color: var(--ds-color-accent-primary);
}

.loading-card__copy {
  color: var(--ds-color-text-secondary);
}
</style>
