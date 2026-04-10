<script setup>
import { computed, onMounted, ref, watch } from 'vue'
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
import SectionHeader from '../components/layout/SectionHeader.vue'
import { useAuthStore } from '../stores/auth'
import { useReadingsStore } from '../stores/readings'

const { t } = useI18n()
const authStore = useAuthStore()
const readingsStore = useReadingsStore()
const activeReadingTab = ref('card-of-day')
const currentSpreadLimit = computed(() => readingsStore.spreadLimit || readingsStore.latestSpread?.remaining_today || null)

const greeting = computed(() => t('dashboard.greeting', { name: authStore.user?.username || 'Seeker' }))
const isLoading = computed(() => readingsStore.dashboardStatus === 'loading')
const spreadLimitReached = computed(() => (
  currentSpreadLimit.value && currentSpreadLimit.value.count <= 0
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
const cardOfDayMeaningReady = computed(() => Boolean(
  readingsStore.cardOfDay?.cards?.length
  && readingsStore.cardOfDay.cards.every((card) => Boolean(card?.viewed)),
))
const spreadMeaningReady = computed(() => Boolean(
  readingsStore.latestSpread?.cards?.length
  && readingsStore.latestSpread.cards.every((card) => Boolean(card?.viewed)),
))
const activeMeaningReading = computed(() => (
  activeReadingTab.value === 'card-of-day'
    ? readingsStore.cardOfDay
    : readingsStore.latestSpread
))
const activeMeaningVisible = computed(() => (
  activeReadingTab.value === 'card-of-day'
    ? cardOfDayMeaningReady.value
    : spreadMeaningReady.value
))
const activeMeaningLoading = computed(() => {
  const readingId = activeMeaningReading.value?.id
  return Boolean(
    readingId
    && (
      readingsStore.aiStatusByReadingId[readingId] === 'loading'
      || activeMeaningReading.value?.ai?.status === 'processing'
    )
  )
})
const activeMeaningError = computed(() => {
  const readingId = activeMeaningReading.value?.id
  if (!readingId) {
    return ''
  }

  return readingsStore.aiErrorByReadingId[readingId]
    || activeMeaningReading.value?.ai?.error_message
    || ''
})
const activeMeaningText = computed(() => activeMeaningReading.value?.ai?.content || '')

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

async function handleRetryMeaning() {
  if (!activeMeaningReading.value?.id) {
    return
  }

  await readingsStore.requestAIInterpretation(activeMeaningReading.value.id)
}

async function maybeRequestActiveMeaning() {
  if (!activeMeaningVisible.value || !activeMeaningReading.value?.id) {
    return
  }

  if (activeMeaningReading.value?.ai?.content) {
    return
  }

  if (activeMeaningLoading.value || activeMeaningError.value) {
    return
  }

  await readingsStore.requestAIInterpretation(activeMeaningReading.value.id)
}

onMounted(() => {
  readingsStore.loadDashboard()
})

watch(
  [
    activeReadingTab,
    () => activeMeaningVisible.value,
    () => activeMeaningReading.value?.id,
    () => activeMeaningReading.value?.ai?.status,
    () => activeMeaningReading.value?.ai?.content,
  ],
  () => {
    void maybeRequestActiveMeaning()
  },
  { immediate: true },
)
</script>

<template>
  <PageContainer>
    <div class="page-stack">
      <PageSection class="dashboard-intro" :title="greeting" :copy="t('dashboard.copy')">
        <template #meta>
          <BaseButton variant="ghost" to="/history">{{ t('dashboard.secondaryAction') }}</BaseButton>
        </template>
      </PageSection>

      <BaseSurface class="dashboard-tabs-surface" padding="lg" tag="section">
        <BaseTabs v-model="activeReadingTab" class="dashboard-reading-tabs" :items="readingTabs">
          <template #panel="{ activeItem }">
            <LoadingCard
              v-if="isLoading"
              :title="t('dashboard.loadingTitle')"
              :copy="t('dashboard.loadingCopy')"
            />

            <div v-else-if="activeItem?.value === 'card-of-day'" class="dashboard-reading-panel">
              <p v-if="readingsStore.cardOfDayError" class="status-message status-message--error">
                {{ readingsStore.cardOfDayError }}
              </p>

              <CardOfDay
                v-if="readingsStore.cardOfDay"
                class="dashboard-card-of-day"
                :reading="readingsStore.cardOfDay"
                :next-available-at="readingsStore.cardOfDayNextAvailableAt"
                :reveal-loading="cardOfDayRevealLoading"
                @reveal="handleRevealCard"
              />

              <EmptyStateCard
                v-else
                class="dashboard-empty-card"
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
              <p v-if="readingsStore.spreadError" class="status-message status-message--error">
                {{ readingsStore.spreadError }}
              </p>

              <ThreeCardSpread
                v-if="readingsStore.latestSpread"
                class="dashboard-three-card-spread"
                :reading="readingsStore.latestSpread"
                :revealing-map="spreadRevealMap"
                :action-label="t('dashboard.newSpreadAction')"
                :action-loading="readingsStore.spreadStatus === 'loading'"
                :action-disabled="spreadLimitReached"
                @generate="handleGenerateSpread"
                @reveal="handleRevealCard"
              />

              <EmptyStateCard
                v-else
                class="dashboard-empty-card"
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

      <BaseSurface v-if="activeMeaningVisible" class="dashboard-meaning-surface" padding="md" tag="section">
        <SectionHeader size="sm" title-tag="h3" :title="t('dashboard.meaningTitle')" />

        <div class="dashboard-meaning">
          <div v-if="activeMeaningLoading" class="dashboard-meaning__state">
            <p class="text-eyebrow">{{ t('dashboard.meaningLoadingLabel') }}</p>
            <p class="dashboard-meaning__copy">{{ t('dashboard.meaningLoadingCopy') }}</p>
          </div>

          <div v-else-if="activeMeaningText" class="dashboard-meaning__state">
            <p class="dashboard-meaning__copy">{{ activeMeaningText }}</p>
          </div>

          <div v-else-if="activeMeaningError" class="dashboard-meaning__state">
            <p class="text-eyebrow">{{ t('dashboard.meaningErrorLabel') }}</p>
            <p class="dashboard-meaning__copy dashboard-meaning__copy--error">{{ activeMeaningError }}</p>
            <div class="dashboard-meaning__actions">
              <BaseButton size="sm" variant="secondary" @click="handleRetryMeaning">
                {{ t('dashboard.meaningRetryAction') }}
              </BaseButton>
            </div>
          </div>
        </div>
      </BaseSurface>
    </div>
  </PageContainer>
</template>

<style scoped>
.status-message {
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-body-sm);
}

.status-message--error {
  color: var(--ds-color-danger);
}

.dashboard-reading-panel > .status-message {
  margin-bottom: 20px;
}

.dashboard-tabs-surface {
  overflow: hidden;
}

.dashboard-meaning-surface {
  display: grid;
  gap: var(--ds-space-5);
}

.dashboard-meaning {
  display: grid;
  gap: var(--ds-space-4);
}

.dashboard-meaning__state {
  display: grid;
  gap: var(--ds-space-2);
}

.dashboard-meaning__copy {
  color: var(--ds-color-text-secondary);
  line-height: 1.65;
}

.dashboard-meaning__copy--error {
  color: var(--ds-color-danger);
}

.dashboard-meaning__actions {
  display: flex;
  align-items: center;
  gap: var(--ds-space-3);
}

@media (max-width: 720px) {
  .dashboard-intro {
    gap: var(--ds-space-2);
  }

  .dashboard-intro :deep(.section-header) {
    gap: var(--ds-space-2);
  }

  .dashboard-intro :deep(.section-header__content) {
    gap: var(--ds-space-2);
  }

  .dashboard-intro :deep(.section-header__title) {
    font-size: clamp(22px, 7vw, 28px);
    line-height: 1.08;
  }

  .dashboard-intro :deep(.section-header__copy) {
    max-width: 32ch;
    font-size: var(--ds-text-body-sm);
    line-height: 1.55;
  }

  .dashboard-intro :deep(.section-header__meta) {
    gap: var(--ds-space-2);
  }

  .dashboard-tabs-surface {
    padding: var(--ds-space-5);
    padding-bottom: 0;
  }

  .dashboard-reading-tabs :deep(.base-tabs) {
    gap: var(--ds-space-3);
  }

  .dashboard-reading-tabs :deep(.base-tabs__list) {
    margin-inline: calc(var(--ds-space-5) * -1);
    padding-inline: var(--ds-space-5);
    scroll-padding-inline: var(--ds-space-5);
  }

  .dashboard-reading-tabs :deep(.base-tabs__trigger) {
    padding: 10px 14px;
    font-size: 11px;
    letter-spacing: 0.05em;
  }

  .dashboard-reading-tabs :deep(.base-tabs--underline .base-tabs__trigger.is-active::before) {
    left: 5px;
    font-size: 9px;
  }

  .dashboard-reading-panel > .status-message {
    margin-bottom: 14px;
    font-size: var(--ds-text-body-sm);
    line-height: 1.5;
  }

  .dashboard-card-of-day :deep(.card-of-day__head),
  .dashboard-three-card-spread :deep(.three-card-spread__head) {
    margin-bottom: 14px;
    padding-bottom: 14px;
  }

  .dashboard-card-of-day :deep(.section-header),
  .dashboard-three-card-spread :deep(.section-header) {
    gap: var(--ds-space-2);
  }

  .dashboard-card-of-day :deep(.section-header__content),
  .dashboard-three-card-spread :deep(.section-header__content) {
    gap: 6px;
  }

  .dashboard-card-of-day :deep(.section-header__title),
  .dashboard-three-card-spread :deep(.section-header__title) {
    font-size: 18px;
    line-height: 1.15;
  }

  .dashboard-card-of-day :deep(.section-header__copy),
  .dashboard-three-card-spread :deep(.section-header__copy) {
    font-size: var(--ds-text-body-sm);
    line-height: 1.5;
  }

  .dashboard-card-of-day :deep(.section-header__meta),
  .dashboard-three-card-spread :deep(.section-header__meta) {
    gap: var(--ds-space-2);
  }

  .dashboard-card-of-day :deep(.base-badge),
  .dashboard-three-card-spread :deep(.base-badge) {
    min-height: 22px;
    padding: 4px 8px;
    font-size: 9px;
  }

  .dashboard-card-of-day :deep(.tarot-card__tags),
  .dashboard-three-card-spread :deep(.tarot-card__tags) {
    display: none;
  }

  .dashboard-three-card-spread :deep(.three-card-spread__grid) {
    flex-wrap: nowrap;
    justify-content: space-between;
    align-items: flex-start;
    gap: 10px;
  }

  .dashboard-three-card-spread :deep(.three-card-spread__badge) {
    display: none;
  }

  .dashboard-three-card-spread :deep(.three-card-spread__grid > .tarot-card) {
    width: calc((100% - 20px) / 3);
    min-width: 0;
    flex: 0 1 calc((100% - 20px) / 3);
    --card-frame-max: none;
    --card-body-padding: 10px;
    --card-title-size: 13px;
    --card-position-margin: 5px;
    --card-name-margin: 7px;
    --card-copy-size: 10px;
    --card-copy-line-height: 1.45;
    --card-copy-margin: 8px;
    --card-tag-gap: 4px;
    --card-back-margin: 8px;
    --card-symbol-size: 22px;
    --card-tag-padding-y: 2px;
    --card-tag-padding-x: 6px;
    --card-tag-font-size: 6px;
    --card-tag-letter-spacing: 0.01em;
  }

  .dashboard-three-card-spread :deep(.three-card-spread__grid .tarot-card__position) {
    font-size: 9px;
    letter-spacing: 0.04em;
    line-height: 1.3;
  }

  .dashboard-three-card-spread :deep(.three-card-spread__grid .tarot-card__title) {
    min-height: calc(13px * 1.2 * 2);
    display: flex;
    align-items: flex-start;
    justify-content: center;
    letter-spacing: 0.01em;
    line-height: 1.2;
    text-wrap: balance;
  }

  .dashboard-three-card-spread :deep(.three-card-spread__grid .tarot-card__copy) {
    line-height: 1.45;
  }

  .dashboard-three-card-spread :deep(.three-card-spread__grid .tarot-card__tag) {
    line-height: 1.15;
  }

  .dashboard-empty-card {
    padding: var(--ds-space-5);
  }

  .dashboard-empty-card :deep(.empty-card__symbol) {
    margin-bottom: 12px;
    font-size: 26px;
  }

  .dashboard-empty-card :deep(.empty-card__title) {
    margin-bottom: 8px;
    font-size: var(--ds-text-heading-3);
  }

  .dashboard-empty-card :deep(.empty-card__copy) {
    max-width: 30ch;
    font-size: var(--ds-text-body-sm);
    line-height: 1.55;
  }

  .dashboard-empty-card :deep(.empty-card__actions) {
    margin-top: 18px;
  }

  .dashboard-meaning-surface {
    gap: var(--ds-space-4);
  }

  .dashboard-meaning {
    gap: var(--ds-space-3);
  }

  .dashboard-meaning__copy {
    font-size: var(--ds-text-body-sm);
    line-height: 1.55;
  }

  .dashboard-meaning__actions {
    gap: var(--ds-space-2);
  }
}
</style>
