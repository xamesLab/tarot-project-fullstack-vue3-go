<script setup>
import BaseSurface from '../base/BaseSurface.vue'
import SectionHeader from '../layout/SectionHeader.vue'

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
    <SectionHeader class="auth-card__head" size="lg" title-tag="h1" :eyebrow="eyebrow" :title="title" :copy="copy" />

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
  margin-bottom: 32px;
}

.auth-card__body {
  display: grid;
  gap: 24px;
}

.auth-card__footer {
  margin-top: 28px;
  color: var(--ds-color-text-secondary);
}

@media (max-width: 720px) {
  .auth-card.base-surface--padding-lg {
    padding: var(--ds-space-5);
  }

  .auth-card__head {
    margin-bottom: 20px;
  }

  .auth-card__body {
    gap: 16px;
  }

  .auth-card__footer {
    margin-top: 18px;
    font-size: var(--ds-text-body-sm);
  }

  .auth-card :deep(.section-header--lg .section-header__content) {
    gap: var(--ds-space-2);
  }

  .auth-card :deep(.section-header--lg .section-header__title) {
    font-size: clamp(26px, 7vw, 34px);
  }

  .auth-card :deep(.section-header__copy) {
    font-size: var(--ds-text-body-sm);
    line-height: 1.5;
  }
}
</style>
