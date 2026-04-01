<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '../components/base/BaseButton.vue'
import BaseSurface from '../components/base/BaseSurface.vue'
import EmptyStateCard from '../components/feedback/EmptyStateCard.vue'
import LoadingCard from '../components/feedback/LoadingCard.vue'
import PageSection from '../components/layout/PageSection.vue'
import { useAuthStore } from '../stores/auth'

const { t } = useI18n()
const authStore = useAuthStore()

const greeting = computed(() => t('dashboard.greeting', { name: authStore.user?.username || 'Seeker' }))
</script>

<template>
  <div class="page-stack">
    <PageSection :eyebrow="t('dashboard.eyebrow')" :title="greeting" :copy="t('dashboard.copy')">
      <BaseSurface class="hero-panel" tone="outlined" padding="lg" tag="section">
        <div class="hero-panel__content">
          <div class="stack-md">
            <p class="text-eyebrow">{{ t('dashboard.today') }}</p>
            <h2 class="hero-panel__title">{{ t('dashboard.heroTitle') }}</h2>
            <p class="hero-panel__copy">{{ t('dashboard.heroCopy') }}</p>
          </div>

          <div class="cluster">
            <BaseButton>{{ t('dashboard.primaryAction') }}</BaseButton>
            <BaseButton variant="secondary" to="/history">{{ t('dashboard.secondaryAction') }}</BaseButton>
          </div>
        </div>
      </BaseSurface>
    </PageSection>

    <div v-if="!authStore.hydrated || authStore.status === 'loading'" class="page-grid page-grid--two">
      <LoadingCard :title="t('dashboard.loadingTitle')" :copy="t('dashboard.loadingCopy')" />
      <LoadingCard :title="t('dashboard.loadingSecondaryTitle')" :copy="t('dashboard.loadingSecondaryCopy')" />
    </div>

    <div v-else class="page-grid page-grid--two">
      <BaseSurface padding="lg" tag="section">
        <div class="stack-md">
          <div class="section-head">
            <p class="text-eyebrow">{{ t('dashboard.cardEyebrow') }}</p>
            <h3 class="section-head__title section-head__title--small">{{ t('dashboard.cardTitle') }}</h3>
            <p class="section-head__copy">{{ t('dashboard.cardCopy') }}</p>
          </div>

          <EmptyStateCard
            :title="t('dashboard.cardEmptyTitle')"
            :description="t('dashboard.cardEmptyCopy')"
          >
            <template #actions>
              <BaseButton>{{ t('dashboard.primaryAction') }}</BaseButton>
            </template>
          </EmptyStateCard>
        </div>
      </BaseSurface>

      <BaseSurface padding="lg" tag="section">
        <div class="stack-md">
          <div class="section-head">
            <p class="text-eyebrow">{{ t('dashboard.spreadEyebrow') }}</p>
            <h3 class="section-head__title section-head__title--small">{{ t('dashboard.spreadTitle') }}</h3>
            <p class="section-head__copy">{{ t('dashboard.spreadCopy') }}</p>
          </div>

          <EmptyStateCard
            symbol="☽"
            :title="t('dashboard.spreadEmptyTitle')"
            :description="t('dashboard.spreadEmptyCopy')"
          >
            <template #actions>
              <BaseButton variant="secondary">{{ t('dashboard.spreadAction') }}</BaseButton>
            </template>
          </EmptyStateCard>
        </div>
      </BaseSurface>
    </div>
  </div>
</template>

<style scoped>
.hero-panel {
  overflow: hidden;
}

.hero-panel__content {
  display: grid;
  gap: 28px;
}

.hero-panel__title {
  font-size: clamp(32px, 4vw, 52px);
  color: var(--ds-color-accent-primary);
}

.hero-panel__copy {
  max-width: 60ch;
  color: var(--ds-color-text-secondary);
}

.section-head__title--small {
  font-size: var(--ds-text-heading-2);
}
</style>
