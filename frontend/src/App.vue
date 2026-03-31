<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, RouterView, useRoute } from 'vue-router'

const { locale, t } = useI18n()
const route = useRoute()

const storedTheme = typeof window !== 'undefined' ? window.localStorage.getItem('theme') : null
const storedLanguage = typeof window !== 'undefined' ? window.localStorage.getItem('language') : null

const theme = ref(storedTheme || 'light')

if (storedLanguage && ['en', 'ru'].includes(storedLanguage)) {
  locale.value = storedLanguage
}

const pageTitle = computed(() => route.meta.titleKey ? t(route.meta.titleKey) : t('app.title'))

const isUIKitRoute = computed(() => route.path === '/ui-kit')
const isDevMode = import.meta.env.DEV

function applyTheme(value) {
  document.documentElement.setAttribute('data-theme', value)
  window.localStorage.setItem('theme', value)
}

function toggleTheme() {
  theme.value = theme.value === 'light' ? 'dark' : 'light'
}

function toggleLanguage() {
  locale.value = locale.value === 'en' ? 'ru' : 'en'
}

onMounted(() => {
  applyTheme(theme.value)
})

watch(theme, (value) => {
  applyTheme(value)
})

watch(locale, (value) => {
  window.localStorage.setItem('language', value)
  document.documentElement.setAttribute('lang', value)
})
</script>

<template>
  <div :class="['app-shell', { 'app-shell--ui-kit': isUIKitRoute }]">
    <header class="app-shell__header">
      <div>
        <p class="app-shell__eyebrow">{{ t('app.badge') }}</p>
        <h1 class="app-shell__title">{{ pageTitle }}</h1>
      </div>

      <div class="app-shell__actions">
        <button class="shell-button shell-button--ghost" type="button" @click="toggleLanguage">
          {{ locale.toUpperCase() }}
        </button>
        <button class="shell-button" type="button" @click="toggleTheme">
          {{ theme === 'light' ? t('app.darkMode') : t('app.lightMode') }}
        </button>
      </div>
    </header>

    <nav class="app-shell__nav" aria-label="Primary">
      <RouterLink class="nav-link" to="/">{{ t('nav.home') }}</RouterLink>
      <RouterLink v-if="isDevMode" class="nav-link" to="/ui-kit">{{ t('nav.uiKit') }}</RouterLink>
    </nav>

    <main class="app-shell__content">
      <RouterView />
    </main>

    <footer v-if="!isUIKitRoute" class="app-shell__footer">
      <p>{{ t('app.footer') }}</p>
    </footer>
  </div>
</template>
