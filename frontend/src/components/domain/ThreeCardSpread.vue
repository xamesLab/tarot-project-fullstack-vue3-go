<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseBadge from '../base/BaseBadge.vue'
import SectionHeader from '../layout/SectionHeader.vue'
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
      <SectionHeader
        size="sm"
        title-tag="h3"
        :eyebrow="t('readings.spreadLabel')"
        :title="t('readings.spreadReady')"
        :copy="reading.remaining_today ? t('readings.remainingToday', reading.remaining_today) : ''"
      >
        <template #meta>
          <BaseBadge :label="badgeLabel" />
        </template>
      </SectionHeader>
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
  margin-bottom: 24px;
  border-bottom: var(--ds-border-width) solid var(--ds-color-border-subtle);
  padding-bottom: 24px;
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

@media (max-width: 720px) {
  .three-card-spread__head {
    margin-bottom: 16px;
    padding-bottom: 16px;
  }

  .three-card-spread__grid {
    gap: 20px;
  }
}
</style>
