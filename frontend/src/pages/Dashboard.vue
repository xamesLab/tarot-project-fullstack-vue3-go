<script setup>
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '../components/base/BaseButton.vue'
import BaseSurface from '../components/base/BaseSurface.vue'
import BaseTabs from '../components/base/BaseTabs.vue'
import CardOfDay from '../components/domain/CardOfDay.vue'
import ThreeCardSpread from '../components/domain/ThreeCardSpread.vue'
import EmptyStateCard from '../components/feedback/EmptyStateCard.vue'
import LoadingCard from '../components/feedback/LoadingCard.vue'
import PageContainer from '../components/layout/PageContainer.vue'
import PageSection from '../components/layout/PageSection.vue'
import { useAuthStore } from '../stores/auth'
import { useReadingsStore } from '../stores/readings'

const { t } = useI18n()
const authStore = useAuthStore()
const readingsStore = useReadingsStore()
const activeReadingTab = ref('card-of-day')

const greeting = computed(() => t('dashboard.greeting', { name: authStore.user?.username || 'Seeker' }))
const isLoading = computed(() => readingsStore.dashboardStatus === 'loading')
const spreadLimitReached = computed(() => (
  readingsStore.spreadLimit && readingsStore.spreadLimit.count <= 0
))
const cardOfDayRevealLoading = computed(() => {
  if (!readingsStore.cardOfDay?.id || !readingsStore.cardOfDay?.cards?.length) {
    return false
  }

  const position = readingsStore.cardOfDay.cards[0]?.position ?? 0
  return readingsStore.revealStatusByKey[`${readingsStore.cardOfDay.id}:${position}`] === 'loading'
})
const spreadRevealMap = computed(() => {
  const reading = readingsStore.latestSpread
  if (!reading?.id) {
    return {}
  }

  return Object.fromEntries(
    (reading.cards || []).map((card) => [
      card.position,
      readingsStore.revealStatusByKey[`${reading.id}:${card.position}`] === 'loading',
    ]),
  )
})
const readingTabs = computed(() => [
  { value: 'card-of-day', label: t('dashboard.cardEyebrow') },
  { value: 'three-card', label: t('dashboard.spreadEyebrow') },
])

async function handleGenerateCardOfDay() {
  await readingsStore.generateCardOfDay()
}

async function handleGenerateSpread() {
  await readingsStore.generateThreeCardSpread()
}

async function handleRevealCard(payload) {
  if (!payload?.readingId || payload?.position == null) {
    return
  }

  await readingsStore.revealCard(payload.readingId, payload.position)
}

onMounted(() => {
  readingsStore.loadDashboard()
})
</script>

<template>
  <PageContainer>
    <div class="page-stack">
      <PageSection :eyebrow="t('dashboard.eyebrow')" :title="greeting" :copy="t('dashboard.copy')" />

      <BaseSurface class="dashboard-tabs-surface" padding="lg" tag="section">
        <BaseTabs v-model="activeReadingTab" :items="readingTabs">
          <template #panel="{ activeItem }">
            <LoadingCard
              v-if="isLoading"
              :title="t('dashboard.loadingTitle')"
              :copy="t('dashboard.loadingCopy')"
            />

            <div v-else-if="activeItem?.value === 'card-of-day'" class="dashboard-reading-panel">
              <div class="section-head">
                <p class="text-eyebrow">{{ t('dashboard.cardEyebrow') }}</p>
                <h3 class="section-head__title section-head__title--small">{{ t('dashboard.cardTitle') }}</h3>
                <p class="section-head__copy">{{ t('dashboard.cardCopy') }}</p>
              </div>

              <p v-if="readingsStore.cardOfDayError" class="status-message status-message--error">
                {{ readingsStore.cardOfDayError }}
              </p>

              <CardOfDay
                v-if="readingsStore.cardOfDay"
                :reading="readingsStore.cardOfDay"
                :next-available-at="readingsStore.cardOfDayNextAvailableAt"
                :reveal-loading="cardOfDayRevealLoading"
                @reveal="handleRevealCard"
              />

              <EmptyStateCard
                v-else
                :title="t('dashboard.cardEmptyTitle')"
                :description="t('dashboard.cardEmptyCopy')"
              >
                <template #actions>
                  <BaseButton :loading="readingsStore.cardOfDayStatus === 'loading'" @click="handleGenerateCardOfDay">
                    {{ t('dashboard.primaryAction') }}
                  </BaseButton>
                </template>
              </EmptyStateCard>
            </div>

            <div v-else class="dashboard-reading-panel">
              <div class="section-head">
                <p class="text-eyebrow">{{ t('dashboard.spreadEyebrow') }}</p>
                <h3 class="section-head__title section-head__title--small">{{ t('dashboard.spreadTitle') }}</h3>
                <p class="section-head__copy">{{ t('dashboard.spreadCopy') }}</p>
              </div>

              <p v-if="readingsStore.spreadError" class="status-message status-message--error">
                {{ readingsStore.spreadError }}
              </p>

              <ThreeCardSpread
                v-if="readingsStore.latestSpread"
                :reading="readingsStore.latestSpread"
                :revealing-map="spreadRevealMap"
                @reveal="handleRevealCard"
              />

              <EmptyStateCard
                v-else
                symbol="☽"
                :title="t('dashboard.spreadEmptyTitle')"
                :description="t('dashboard.spreadEmptyCopy')"
              >
                <template #actions>
                  <BaseButton
                    variant="secondary"
                    :loading="readingsStore.spreadStatus === 'loading'"
                    :disabled="spreadLimitReached"
                    @click="handleGenerateSpread"
                  >
                    {{ t('dashboard.spreadAction') }}
                  </BaseButton>
                </template>
              </EmptyStateCard>
            </div>
          </template>
        </BaseTabs>
      </BaseSurface>

      <BaseSurface class="hero-panel" tone="outlined" padding="lg" tag="section">
        <div class="hero-panel__content">
          <div class="stack-md">
            <p class="text-eyebrow">{{ t('dashboard.today') }}</p>
            <h2 class="hero-panel__title">{{ t('dashboard.heroTitle') }}</h2>
            <p class="hero-panel__copy">{{ t('dashboard.heroCopy') }}</p>
          </div>

          <div class="cluster">
            <BaseButton :loading="readingsStore.cardOfDayStatus === 'loading'" @click="handleGenerateCardOfDay">
              {{ t('dashboard.primaryAction') }}
            </BaseButton>
            <BaseButton
              variant="secondary"
              :loading="readingsStore.spreadStatus === 'loading'"
              :disabled="spreadLimitReached"
              @click="handleGenerateSpread"
            >
              {{ t('dashboard.spreadAction') }}
            </BaseButton>
            <BaseButton variant="ghost" to="/history">{{ t('dashboard.secondaryAction') }}</BaseButton>
          </div>
        </div>
      </BaseSurface>
    </div>
  </PageContainer>
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

.status-message {
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-body-sm);
}

.status-message--error {
  color: var(--ds-color-danger);
}

.dashboard-reading-panel > .section-head {
  margin-bottom: 24px;
}

.dashboard-reading-panel > .status-message {
  margin-bottom: 20px;
}

.dashboard-tabs-surface {
  overflow: hidden;
}
</style>
