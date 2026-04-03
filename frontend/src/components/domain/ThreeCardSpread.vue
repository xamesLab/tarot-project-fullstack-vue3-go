<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import TarotCard from './TarotCard.vue'

const props = defineProps({
  reading: {
    type: Object,
    required: true,
  },
  revealingMap: {
    type: Object,
    default: () => ({}),
  },
})
const emit = defineEmits(['reveal'])

const { t } = useI18n()

const positionKeys = ['past', 'present', 'future']
const badgeLabel = computed(() => {
  if (props.reading?.remaining_today?.limit) {
    return `${props.reading.remaining_today.count}/${props.reading.remaining_today.limit}`
  }

  return t('readings.spreadLabel')
})
</script>

<template>
  <section class="three-card-spread">
    <div class="three-card-spread__head">
      <div class="three-card-spread__summary">
        <p class="text-eyebrow">{{ t('readings.spreadLabel') }}</p>
        <h3 class="three-card-spread__title">{{ t('readings.spreadReady') }}</h3>
        <p v-if="reading.remaining_today" class="three-card-spread__meta">
          {{ t('readings.remainingToday', reading.remaining_today) }}
        </p>
      </div>
      <span class="three-card-spread__badge">{{ badgeLabel }}</span>
    </div>

    <div class="three-card-spread__grid">
      <TarotCard
        v-for="(card, index) in reading.cards"
        :key="`${reading.id}-${card.key}-${index}`"
        :card="card"
        :position-label="t(`readings.positions.${positionKeys[index] || 'future'}`)"
        :revealed="Boolean(card.viewed)"
        :revealable="true"
        :reveal-loading="Boolean(revealingMap[card.position])"
        size="md"
        @reveal="emit('reveal', { readingId: reading.id, position: card.position })"
      />
    </div>
  </section>
</template>

<style scoped>
.three-card-spread {
  overflow: visible;
}

.three-card-spread__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  margin-bottom: 24px;
  border-bottom: var(--ds-border-width) solid var(--ds-color-border-subtle);
  padding-bottom: 24px;
}

.three-card-spread__summary {
  display: grid;
  gap: 8px;
}

.three-card-spread__title {
  font-size: 18px;
  color: var(--ds-color-text-primary);
  letter-spacing: 0.05em;
}

.three-card-spread__meta {
  font-size: var(--ds-text-body-sm);
  color: var(--ds-color-text-secondary);
}

.three-card-spread__badge {
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

.three-card-spread__grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 50px;
}

@media (max-width: 960px) {
  .three-card-spread__grid {
    gap: 32px;
  }
}
</style>
