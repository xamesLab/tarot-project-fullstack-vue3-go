<script setup>
import BaseSurface from '../base/BaseSurface.vue'

defineProps({
  eyebrow: {
    type: String,
    default: '',
  },
  title: {
    type: String,
    default: '',
  },
  copy: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <BaseSurface class="auth-card" tone="outlined" padding="lg" tag="section">
    <div class="auth-card__head">
      <p v-if="eyebrow" class="text-eyebrow">{{ eyebrow }}</p>
      <h1 class="auth-card__title">{{ title }}</h1>
      <p v-if="copy" class="auth-card__copy">{{ copy }}</p>
    </div>

    <div class="auth-card__body">
      <slot />
    </div>

    <div v-if="$slots.footer" class="auth-card__footer">
      <slot name="footer" />
    </div>
  </BaseSurface>
</template>

<style scoped>
.auth-card {
  width: min(560px, 100%);
  margin: 0 auto;
}

.auth-card__head {
  display: grid;
  gap: 12px;
  margin-bottom: 32px;
}

.auth-card__title {
  font-size: clamp(32px, 4vw, 46px);
  color: var(--ds-color-accent-primary);
}

.auth-card__copy {
  color: var(--ds-color-text-secondary);
}

.auth-card__body {
  display: grid;
  gap: 24px;
}

.auth-card__footer {
  margin-top: 28px;
  color: var(--ds-color-text-secondary);
}
</style>
