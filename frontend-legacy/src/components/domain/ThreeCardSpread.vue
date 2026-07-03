<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseBadge from '../base/BaseBadge.vue'
import BaseButton from '../base/BaseButton.vue'
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
  actionLabel: {
    type: String,
    default: '',
  },
  actionLoading: {
    type: Boolean,
    default: false,
  },
  actionDisabled: {
    type: Boolean,
    default: false,
  },
})
const emit = defineEmits(['reveal', 'generate'])

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
          <BaseBadge class="three-card-spread__badge" :label="badgeLabel" />
          <BaseButton
            v-if="actionLabel"
            class="three-card-spread__action"
            size="sm"
            variant="secondary"
            :loading="actionLoading"
            :disabled="actionDisabled"
            @click="emit('generate')"
          >
            {{ actionLabel }}
          </BaseButton>
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

.three-card-spread__action {
  min-height: 34px;
  padding: 7px 12px;
  font-size: var(--ds-text-caption);
  letter-spacing: 0.06em;
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

  .three-card-spread__action {
    min-height: 30px;
    padding: 6px 10px;
    letter-spacing: 0.04em;
  }
}
</style>
