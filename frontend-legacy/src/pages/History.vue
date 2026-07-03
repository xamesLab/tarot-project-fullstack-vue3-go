<script setup>
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseButton from '../components/base/BaseButton.vue'
import BaseSurface from '../components/base/BaseSurface.vue'
import BaseTabs from '../components/base/BaseTabs.vue'
import ReadingHistory from '../components/domain/ReadingHistory.vue'
import LoadingCard from '../components/feedback/LoadingCard.vue'
import PageContainer from '../components/layout/PageContainer.vue'
import PageSection from '../components/layout/PageSection.vue'
import { useReadingsStore } from '../stores/readings'

const { t } = useI18n()
const readingsStore = useReadingsStore()
const activeFilter = ref('all')

const filterOptions = computed(() => [
  { value: 'all', label: t('history.filters.all') },
  { value: 'card_of_day', label: t('history.filters.cardOfDay') },
  { value: 'three_card', label: t('history.filters.threeCard') },
])

async function loadPage(page = 1) {
  await readingsStore.loadHistory({
    type: activeFilter.value,
    page,
    limit: 10,
  })
}

async function handleFilterChange(value) {
  activeFilter.value = value
  await loadPage(1)
}

async function goToPage(page) {
  if (page < 1 || page > readingsStore.historyPagination.total_pages) {
    return
  }
  await loadPage(page)
}

async function handleRequestAI(readingId) {
  await readingsStore.requestAIInterpretation(readingId)
}

onMounted(() => {
  loadPage(1)
})
</script>

<template>
  <PageContainer>
    <div class="page-stack">
      <PageSection :eyebrow="t('history.eyebrow')" :title="t('history.title')" :copy="t('history.copy')" />

      <BaseSurface padding="lg" tag="section">
        <BaseTabs
          :model-value="activeFilter"
          :items="filterOptions"
          @update:model-value="handleFilterChange"
        >
          <template #panel>
            <p v-if="readingsStore.retentionInfo" class="history-toolbar__meta">
              {{ t('history.retentionInfo', readingsStore.retentionInfo) }}
            </p>
          </template>
        </BaseTabs>
      </BaseSurface>

      <LoadingCard
        v-if="readingsStore.historyStatus === 'loading'"
        :title="t('history.loadingTitle')"
        :copy="t('history.loadingCopy')"
      />

      <BaseSurface v-else-if="readingsStore.historyError" padding="lg" tag="section">
        <p class="history-error">{{ readingsStore.historyError }}</p>
      </BaseSurface>

      <ReadingHistory
        v-else
        :readings="readingsStore.history"
        :ai-status-by-reading-id="readingsStore.aiStatusByReadingId"
        :ai-error-by-reading-id="readingsStore.aiErrorByReadingId"
        :empty-title="t('history.emptyTitle')"
        :empty-copy="t('history.emptyCopy')"
        @request-ai="handleRequestAI"
      />

      <div
        v-if="readingsStore.historyPagination.total_pages > 1 && readingsStore.historyStatus !== 'loading'"
        class="history-pagination"
      >
        <BaseButton
          variant="ghost"
          :disabled="readingsStore.historyPagination.page <= 1"
          @click="goToPage(readingsStore.historyPagination.page - 1)"
        >
          {{ t('history.previousPage') }}
        </BaseButton>

        <span class="history-pagination__label">
          {{ t('history.pageIndicator', readingsStore.historyPagination) }}
        </span>

        <BaseButton
          variant="ghost"
          :disabled="readingsStore.historyPagination.page >= readingsStore.historyPagination.total_pages"
          @click="goToPage(readingsStore.historyPagination.page + 1)"
        >
          {{ t('history.nextPage') }}
        </BaseButton>
      </div>
    </div>
  </PageContainer>
</template>

<style scoped>
.history-toolbar__meta {
  color: var(--ds-color-text-secondary);
}

.history-error {
  color: var(--ds-color-danger);
  font-family: var(--ds-font-ui);
}

.history-pagination {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.history-pagination__label {
  font-family: var(--ds-font-ui);
  color: var(--ds-color-text-secondary);
}
</style>
