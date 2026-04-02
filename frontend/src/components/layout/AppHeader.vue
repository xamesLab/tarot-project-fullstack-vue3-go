<script setup>
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

import BaseButton from '../base/BaseButton.vue'
import BaseSegmentedToggle from '../base/BaseSegmentedToggle.vue'
import PageContainer from './PageContainer.vue'
import UserChip from '../domain/UserChip.vue'
import { useTheme } from '../../composables/useTheme'
import { useLocale } from '../../composables/useLocale'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const { theme, setTheme } = useTheme()
const { language, setLanguage } = useLocale()
const isDevMode = import.meta.env.DEV

const isAuthenticated = computed(() => authStore.isAuthenticated)

const languageOptions = computed(() => [
  { label: 'EN', value: 'en' },
  { label: 'RU', value: 'ru' },
])

const themeOptions = computed(() => [
  { label: t('controls.light'), value: 'light' },
  { label: t('controls.dark'), value: 'dark' },
])

const primaryLinks = computed(() => {
  if (isAuthenticated.value) {
    const links = [
      { to: '/', label: t('nav.dashboard') },
      { to: '/history', label: t('nav.history') },
      { to: '/profile', label: t('nav.profile') },
    ]

    if (authStore.isKeeper) {
      links.push({ to: '/admin', label: t('nav.admin') })
    }

    return links
  }

  return [
    { to: '/login', label: t('nav.login') },
    { to: '/register', label: t('nav.register') },
  ]
})

function isLinkActive(linkTo) {
  if (linkTo === '/') {
    return route.path === '/'
  }

  return route.path === linkTo || route.path.startsWith(`${linkTo}/`)
}

async function handleLogout() {
  authStore.logout()
  await router.push('/login')
}
</script>

<template>
  <header class="site-header">
    <PageContainer>
      <div class="site-header__inner">
        <div class="site-header__brand">
          <RouterLink class="brand-mark" to="/">
            <img class="brand-mark__symbol" src="/branding/nebula-tarot-logo.png" :alt="t('app.title')" />
            <span class="brand-mark__copy">
              <span class="text-eyebrow">{{ t('app.badge') }}</span>
              <strong class="brand-mark__name">{{ t('app.title') }}</strong>
            </span>
          </RouterLink>
        </div>

        <nav class="site-header__nav" aria-label="Primary">
          <RouterLink
            v-for="link in primaryLinks"
            :key="link.to"
            class="site-header__nav-link"
            :class="{ 'is-active': isLinkActive(link.to) }"
            :to="link.to"
          >
            {{ link.label }}
          </RouterLink>

          <RouterLink v-if="isDevMode" class="site-header__nav-link" to="/ui-kit">
            {{ t('nav.uiKit') }}
          </RouterLink>
        </nav>

        <div class="site-header__controls">
          <BaseSegmentedToggle
            :model-value="language"
            :options="languageOptions"
            :aria-label="t('controls.language')"
            @update:model-value="setLanguage"
          />

          <BaseSegmentedToggle
            :model-value="theme"
            :options="themeOptions"
            :aria-label="t('controls.theme')"
            @update:model-value="setTheme"
          />

          <UserChip v-if="isAuthenticated" :user="authStore.user" />

          <BaseButton v-if="isAuthenticated" variant="text" @click="handleLogout">
            {{ t('nav.logout') }}
          </BaseButton>
        </div>
      </div>
    </PageContainer>
  </header>
</template>

<style scoped>
.site-header {
  position: sticky;
  top: 0;
  z-index: 50;
  border-bottom: 1px solid var(--ds-color-border-accent);
  background: rgba(0, 0, 0, 0.9);
  backdrop-filter: blur(12px);
}

.site-header__inner {
  min-height: var(--ds-header-height);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 18px 0;
}

.site-header__brand {
  flex-shrink: 0;
}

.brand-mark {
  display: flex;
  align-items: center;
  gap: 14px;
}

.brand-mark__symbol {
  width: 56px;
  height: 56px;
  object-fit: contain;
  flex-shrink: 0;
  filter: drop-shadow(0 0 10px rgba(212, 175, 55, 0.18));
}

.brand-mark__copy {
  display: grid;
  gap: 2px;
}

.brand-mark__name {
  font-family: var(--ds-font-display);
  font-size: 20px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--ds-color-accent-primary);
}

.site-header__nav {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.site-header__nav-link {
  padding: 10px 12px;
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-caption);
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--ds-color-text-secondary);
  transition: color var(--ds-motion-fast), background var(--ds-motion-fast);
}

.site-header__nav-link:hover,
.site-header__nav-link.is-active {
  color: var(--ds-color-accent-primary);
  background: rgba(212, 175, 55, 0.08);
}

.site-header__controls {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  margin-left: auto;
}

@media (max-width: 1080px) {
  .site-header__inner {
    flex-direction: column;
    align-items: stretch;
  }

  .site-header__nav,
  .site-header__controls {
    justify-content: center;
    margin-left: 0;
  }
}
</style>
