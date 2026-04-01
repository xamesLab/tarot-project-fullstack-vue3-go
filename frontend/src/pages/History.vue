<script setup>
import { useI18n } from 'vue-i18n'

import BaseSurface from '../components/base/BaseSurface.vue'
import EmptyStateCard from '../components/feedback/EmptyStateCard.vue'
import PageSection from '../components/layout/PageSection.vue'

const { t } = useI18n()
</script>

<template>
  <div class="page-stack">
    <PageSection :eyebrow="t('history.eyebrow')" :title="t('history.title')" :copy="t('history.copy')" />

    <BaseSurface padding="lg" tag="section">
      <EmptyStateCard
        symbol="☾"
        :title="t('history.emptyTitle')"
        :description="t('history.emptyCopy')"
      />
    </BaseSurface>
  </div>
</template>
