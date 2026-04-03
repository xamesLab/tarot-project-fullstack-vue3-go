<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import TarotCard from './TarotCard.vue'
import { formatDateTime } from '../../utils/date'

const props = defineProps({
  reading: {
    type: Object,
    required: true,
  },
  nextAvailableAt: {
    type: String,
    default: '',
  },
  revealLoading: {
    type: Boolean,
    default: false,
  },
})
const emit = defineEmits(['reveal'])

const { t, locale } = useI18n()

const nextAvailableLabel = computed(() => (
  (props.nextAvailableAt || props.reading?.next_available_at)
    ? formatDateTime(props.nextAvailableAt || props.reading?.next_available_at, locale.value)
    : ''
))
const badgeLabel = computed(() => props.reading?.accounting_date || t('dashboard.today'))
</script>

<template>
  <section class="card-of-day">
    <div class="card-of-day__head">
      <div class="card-of-day__summary">
        <p class="text-eyebrow">{{ t('readings.cardOfDayLabel') }}</p>
        <h3 class="card-of-day__title">{{ t('readings.cardOfDayReady') }}</h3>
        <p v-if="nextAvailableLabel" class="card-of-day__meta">
          {{ t('readings.nextAvailableAt', { value: nextAvailableLabel }) }}
        </p>
      </div>
      <span class="card-of-day__badge">{{ badgeLabel }}</span>
    </div>

    <TarotCard
      :card="reading.cards[0]"
      :revealed="Boolean(reading.cards[0]?.viewed)"
      :revealable="true"
      :reveal-loading="revealLoading"
      size="md"
      @reveal="emit('reveal', { readingId: reading.id, position: reading.cards[0]?.position ?? 0 })"
    />
  </section>
</template>

<style scoped>
.card-of-day {
  overflow: visible;
}

.card-of-day__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  margin-bottom: 24px;
  border-bottom: var(--ds-border-width) solid var(--ds-color-border-subtle);
  padding-bottom: 24px;
}

.card-of-day__summary {
  display: grid;
  gap: 8px;
}

.card-of-day__title {
  font-size: 18px;
  color: var(--ds-color-text-primary);
  letter-spacing: 0.05em;
}

.card-of-day__meta {
  font-size: var(--ds-text-body-sm);
  color: var(--ds-color-text-secondary);
}

.card-of-day__badge {
  padding: 5px 14px;
  border: var(--ds-border-width) solid var(--ds-color-border-accent);
  background: rgba(212, 175, 55, 0.1);
  font-family: var(--ds-font-ui);
  font-size: 9px;
  font-weight: 500;
  letter-spacing: 0.09em;
  line-height: 1;
  white-space: nowrap;
  text-transform: uppercase;
  color: var(--ds-color-accent-primary);
}

@media (max-width: 720px) {
  .card-of-day__head {
    gap: 12px;
    margin-bottom: 16px;
    padding-bottom: 16px;
  }

  .card-of-day__summary {
    gap: 6px;
  }

  .card-of-day__title {
    font-size: 16px;
  }

  .card-of-day__meta {
    font-size: var(--ds-text-caption);
  }

  .card-of-day__badge {
    padding: 4px 10px;
    font-size: 8px;
  }
}
</style>
