<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'

import AppFooter from './components/layout/AppFooter.vue'
import AppHeader from './components/layout/AppHeader.vue'
import BaseSurface from './components/base/BaseSurface.vue'
import { useLocale } from './composables/useLocale'
import { useTheme } from './composables/useTheme'
import { APP_FORBIDDEN_EVENT, APP_UNAUTHORIZED_EVENT } from './services/api'
import { useAppStore } from './stores/app'
import { useAuthStore } from './stores/auth'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

useTheme()
useLocale()

const isReady = ref(false)

function handleUnauthorized() {
  if (!authStore.hasToken) {
    return
  }

  authStore.logout()

  if (!['/login', '/register'].includes(router.currentRoute.value.path)) {
    router.push('/login')
  }
}

function handleForbidden() {
  if (!authStore.isAuthenticated) {
    return
  }

  if (router.currentRoute.value.path !== '/') {
    router.push('/')
  }
}

onMounted(async () => {
  appStore.init()
  await authStore.hydrateSession()
  isReady.value = true

  if (typeof window !== 'undefined') {
    window.addEventListener(APP_UNAUTHORIZED_EVENT, handleUnauthorized)
    window.addEventListener(APP_FORBIDDEN_EVENT, handleForbidden)
  }
})

onBeforeUnmount(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener(APP_UNAUTHORIZED_EVENT, handleUnauthorized)
    window.removeEventListener(APP_FORBIDDEN_EVENT, handleForbidden)
  }
})
</script>

<template>
  <RouterView v-if="route.path === '/ui-kit'" />

  <div v-else-if="isReady" class="app-root">
    <div class="app-background" aria-hidden="true"></div>
    <AppHeader />

    <main class="app-main">
      <RouterView />
    </main>

    <AppFooter />
  </div>

  <div v-else class="app-splash">
    <div class="app-background" aria-hidden="true"></div>
    <BaseSurface class="app-splash__panel" tone="outlined" padding="lg" tag="section">
      <img class="app-splash__logo" src="/branding/nebula-tarot-logo.png" alt="Nebula Tarot" />
      <h1 class="app-splash__title">Nebula Tarot</h1>
      <p class="app-splash__copy">Initializing application shell and restoring session.</p>
      <div class="app-splash__loader"></div>
    </BaseSurface>
  </div>
</template>

<style scoped>
.app-splash__logo {
  width: 120px;
  height: 120px;
  object-fit: contain;
  margin: 0 auto 20px;
  display: block;
  filter: drop-shadow(0 0 18px rgba(212, 175, 55, 0.18));
}
</style>
