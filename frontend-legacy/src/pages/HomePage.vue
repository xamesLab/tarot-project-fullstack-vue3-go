<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import PageContainer from '../components/layout/PageContainer.vue'

const { t } = useI18n()

const checklist = computed(() => [
  t('home.points.backend'),
  t('home.points.frontend'),
  t('home.points.themes'),
  t('home.points.i18n'),
  t('home.points.seed'),
])
</script>

<template>
  <PageContainer tag="section">
    <section class="hero-card">
      <div class="hero-card__body">
        <p class="hero-card__eyebrow">{{ t('home.eyebrow') }}</p>
        <h2 class="hero-card__title">{{ t('home.heading') }}</h2>
        <p class="hero-card__copy">{{ t('home.copy') }}</p>

        <ul class="hero-card__list">
          <li v-for="item in checklist" :key="item">{{ item }}</li>
        </ul>
      </div>
    </section>
  </PageContainer>
</template>
