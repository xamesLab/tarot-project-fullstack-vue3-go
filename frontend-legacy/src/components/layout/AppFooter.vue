<script setup>
import { useI18n } from 'vue-i18n'

import PageContainer from './PageContainer.vue'

const { t } = useI18n()
</script>

<template>
  <footer class="site-footer">
    <PageContainer>
      <div class="site-footer__inner">
        <div class="site-footer__divider">
          <span>✦</span>
        </div>
        <p class="site-footer__copy">{{ t('app.footer') }}</p>
      </div>
    </PageContainer>
  </footer>
</template>

<style scoped>
.site-footer {
  padding: 0 20px 48px;
}

.site-footer__inner {
  display: grid;
  gap: 18px;
  text-align: center;
}

.site-footer__divider {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  color: var(--ds-color-accent-primary);
}

.site-footer__divider::before,
.site-footer__divider::after {
  content: "";
  height: 1px;
  flex: 1;
  background: linear-gradient(to right, transparent, var(--ds-color-border-accent), transparent);
}

.site-footer__copy {
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-caption);
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--ds-color-text-muted);
}
</style>
