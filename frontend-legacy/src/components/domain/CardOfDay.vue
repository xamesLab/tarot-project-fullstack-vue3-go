<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseBadge from '../base/BaseBadge.vue'
import SectionHeader from '../layout/SectionHeader.vue'
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
      <SectionHeader
        size="sm"
        title-tag="h3"
        :eyebrow="t('readings.cardOfDayLabel')"
        :title="t('readings.cardOfDayReady')"
        :copy="nextAvailableLabel ? t('readings.nextAvailableAt', { value: nextAvailableLabel }) : ''"
      >
        <template #meta>
          <BaseBadge :label="badgeLabel" />
        </template>
      </SectionHeader>
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
  margin-bottom: 24px;
  border-bottom: var(--ds-border-width) solid var(--ds-color-border-subtle);
  padding-bottom: 24px;
}

@media (max-width: 720px) {
  .card-of-day__head {
    margin-bottom: 16px;
    padding-bottom: 16px;
  }
}
</style>
