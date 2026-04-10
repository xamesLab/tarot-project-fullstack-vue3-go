<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseBadge from '../base/BaseBadge.vue'
import BaseButton from '../base/BaseButton.vue'
import BaseSurface from '../base/BaseSurface.vue'
import EmptyStateCard from '../feedback/EmptyStateCard.vue'
import TarotCard from './TarotCard.vue'
import { formatDateTime } from '../../utils/date'

const props = defineProps({
  readings: {
    type: Array,
    default: () => [],
  },
  aiStatusByReadingId: {
    type: Object,
    default: () => ({}),
  },
  aiErrorByReadingId: {
    type: Object,
    default: () => ({}),
  },
  emptyTitle: {
    type: String,
    default: '',
  },
  emptyCopy: {
    type: String,
    default: '',
  },
})
const emit = defineEmits(['request-ai'])

const { t, locale } = useI18n()
const expandedIds = ref([])

function isExpanded(id) {
  return expandedIds.value.includes(id)
}

function toggleExpanded(reading) {
  if (isExpanded(reading.id)) {
    expandedIds.value = expandedIds.value.filter((item) => item !== reading.id)
    return
  }

  expandedIds.value = [...expandedIds.value, reading.id]

  if (isAIRequestAllowed(reading) && !reading.ai?.content && aiStatus(reading.id) !== 'loading') {
    emit('request-ai', reading.id)
  }
}

function formattedDate(value) {
  return formatDateTime(value, locale.value)
}

function aiStatus(readingId) {
  return props.aiStatusByReadingId[readingId] || ''
}

function aiError(reading) {
  return props.aiErrorByReadingId[reading.id] || reading.ai?.error_message || ''
}

function isAIRequestAllowed(reading) {
  return Array.isArray(reading?.cards) && reading.cards.every((card) => Boolean(card?.viewed))
}

function hasAIContent(reading) {
  return Boolean(reading?.ai?.content)
}

function shouldShowMeaning(reading) {
  return isAIRequestAllowed(reading) || hasAIContent(reading) || aiStatus(reading.id) === 'loading' || Boolean(aiError(reading))
}
</script>

<template>
  <div v-if="!readings.length">
    <EmptyStateCard symbol="☾" :title="emptyTitle" :description="emptyCopy" />
  </div>

  <div v-else class="reading-history">
    <BaseSurface
      v-for="reading in readings"
      :key="reading.id"
      class="reading-history__item"
      padding="lg"
      tag="article"
    >
      <div class="reading-history__summary">
        <div class="stack-sm">
          <p class="text-eyebrow">{{ t(`readings.types.${reading.type}`) }}</p>
          <h3 class="reading-history__title">{{ formattedDate(reading.generated_at) }}</h3>
          <p class="reading-history__copy">
            {{ t('readings.accountingDateLabel', { value: reading.accounting_date }) }}
          </p>
          <div v-if="hasAIContent(reading)" class="reading-history__status">
            <BaseBadge variant="success" size="sm">{{ t('history.aiReadyBadge') }}</BaseBadge>
          </div>
        </div>

        <BaseButton variant="text" @click="toggleExpanded(reading)">
          {{ isExpanded(reading.id) ? t('readings.hideDetails') : t('readings.showDetails') }}
        </BaseButton>
      </div>

      <div class="reading-history__chips">
        <BaseBadge
          v-for="card in reading.cards"
          :key="`${reading.id}-${card.key}`"
          class="reading-history__chip"
          variant="subtle"
          size="sm"
        >
          {{ locale === 'ru' ? card.name_ru : card.name }}
        </BaseBadge>
      </div>

      <div v-if="isExpanded(reading.id)" class="reading-history__cards">
        <TarotCard
          v-for="(card, index) in reading.cards"
          :key="`${reading.id}-${card.key}-${index}`"
          :card="card"
          compact
        />
      </div>

      <div v-if="isExpanded(reading.id) && shouldShowMeaning(reading)" class="reading-history__meaning">
        <p class="text-eyebrow">{{ t('dashboard.meaningTitle') }}</p>

        <p v-if="aiStatus(reading.id) === 'loading'" class="reading-history__meaning-copy">
          {{ t('dashboard.meaningLoadingCopy') }}
        </p>

        <template v-else-if="hasAIContent(reading)">
          <p class="reading-history__meaning-copy">{{ reading.ai.content }}</p>
        </template>

        <template v-else-if="aiError(reading)">
          <p class="reading-history__meaning-copy reading-history__meaning-copy--error">{{ aiError(reading) }}</p>
          <div class="reading-history__meaning-actions">
            <BaseButton size="sm" variant="secondary" @click="emit('request-ai', reading.id)">
              {{ t('dashboard.meaningRetryAction') }}
            </BaseButton>
          </div>
        </template>

        <p v-else class="reading-history__meaning-copy">
          {{ t('history.aiWaitingCopy') }}
        </p>
      </div>
    </BaseSurface>
  </div>
</template>

<style scoped>
.reading-history {
  display: grid;
  gap: 16px;
}

.reading-history__item {
  display: grid;
  gap: 16px;
}

.reading-history__summary {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 16px;
}

.reading-history__title {
  font-size: var(--ds-text-heading-3);
  color: var(--ds-color-accent-primary);
}

.reading-history__copy {
  color: var(--ds-color-text-secondary);
}

.reading-history__chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.reading-history__status {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.reading-history__cards {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
}

.reading-history__meaning {
  display: grid;
  gap: 8px;
  padding-top: 16px;
  border-top: var(--ds-border-width) solid var(--ds-color-border-subtle);
}

.reading-history__meaning-copy {
  color: var(--ds-color-text-secondary);
  line-height: 1.6;
}

.reading-history__meaning-copy--error {
  color: var(--ds-color-danger);
}

.reading-history__meaning-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
</style>
