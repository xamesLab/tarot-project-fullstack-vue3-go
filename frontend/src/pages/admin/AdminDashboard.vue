<script setup>
import { useI18n } from 'vue-i18n'

import BaseSurface from '../../components/base/BaseSurface.vue'
import PageSection from '../../components/layout/PageSection.vue'

const { t } = useI18n()
</script>

<template>
  <div class="page-stack">
    <PageSection :eyebrow="t('admin.eyebrow')" :title="t('admin.title')" :copy="t('admin.copy')" />

    <div class="page-grid page-grid--two">
      <BaseSurface padding="lg" tag="section">
        <div class="stack-sm">
          <p class="text-eyebrow">{{ t('admin.usersTitle') }}</p>
          <p class="text-body-muted">{{ t('admin.usersCopy') }}</p>
        </div>
      </BaseSurface>

      <BaseSurface padding="lg" tag="section">
        <div class="stack-sm">
          <p class="text-eyebrow">{{ t('admin.decksTitle') }}</p>
          <p class="text-body-muted">{{ t('admin.decksCopy') }}</p>
        </div>
      </BaseSurface>
    </div>
  </div>
</template>
