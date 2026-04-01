<script setup>
import { useI18n } from 'vue-i18n'

import BaseSurface from '../../components/base/BaseSurface.vue'
import PageSection from '../../components/layout/PageSection.vue'

const { t } = useI18n()
</script>

<template>
  <div class="page-stack">
    <PageSection :eyebrow="t('admin.eyebrow')" :title="t('admin.settingsTitle')" :copy="t('admin.settingsCopy')" />
    <BaseSurface padding="lg" tag="section">
      <p class="text-body-muted">{{ t('admin.placeholder') }}</p>
    </BaseSurface>
  </div>
</template>
