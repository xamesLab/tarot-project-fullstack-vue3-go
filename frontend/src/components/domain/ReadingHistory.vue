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
  emptyTitle: {
    type: String,
    default: '',
  },
  emptyCopy: {
    type: String,
    default: '',
  },
})

const { t, locale } = useI18n()
const expandedIds = ref([])

function isExpanded(id) {
  return expandedIds.value.includes(id)
}

function toggleExpanded(id) {
  expandedIds.value = isExpanded(id)
    ? expandedIds.value.filter((item) => item !== id)
    : [...expandedIds.value, id]
}

function formattedDate(value) {
  return formatDateTime(value, locale.value)
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
        </div>

        <BaseButton variant="text" @click="toggleExpanded(reading.id)">
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

.reading-history__cards {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
}
</style>
