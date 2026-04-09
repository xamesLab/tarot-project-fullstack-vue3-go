<script setup>
import { computed, ref, watch } from 'vue'
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
const mobileMenuOpen = ref(false)

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

function toggleMobileMenu() {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

function closeMobileMenu() {
  mobileMenuOpen.value = false
}

async function handleLogout() {
  authStore.logout()
  closeMobileMenu()
  await router.push('/login')
}

watch(
  () => route.fullPath,
  () => {
    closeMobileMenu()
  },
)
</script>

<template>
  <header class="site-header">
    <PageContainer>
      <div class="site-header__inner" :class="{ 'is-menu-open': mobileMenuOpen }">
        <div class="site-header__top-row">
          <div class="site-header__brand">
            <RouterLink class="brand-mark" to="/" @click="closeMobileMenu">
              <img class="brand-mark__symbol" src="/branding/nebula-tarot-logo.png" :alt="t('app.title')" />
              <strong class="brand-mark__name">{{ t('app.title') }}</strong>
            </RouterLink>
          </div>

          <button
            type="button"
            class="site-header__menu-toggle"
            :class="{ 'is-open': mobileMenuOpen }"
            :aria-expanded="mobileMenuOpen"
            aria-controls="site-header-menu"
            aria-label="Toggle navigation menu"
            @click="toggleMobileMenu"
          >
            <span></span>
            <span></span>
            <span></span>
          </button>
        </div>

        <div id="site-header-menu" class="site-header__menu-panel" :class="{ 'is-open': mobileMenuOpen }">
          <nav class="site-header__nav" aria-label="Primary">
            <RouterLink
              v-for="link in primaryLinks"
              :key="link.to"
              class="site-header__nav-link"
              :class="{ 'is-active': isLinkActive(link.to) }"
              :to="link.to"
              @click="closeMobileMenu"
            >
              {{ link.label }}
            </RouterLink>

            <RouterLink v-if="isDevMode" class="site-header__nav-link" to="/ui-kit" @click="closeMobileMenu">
              {{ t('nav.uiKit') }}
            </RouterLink>
          </nav>

          <div class="site-header__controls">
            <div class="site-header__toggles">
              <BaseSegmentedToggle
                class="site-header__toggle"
                :model-value="language"
                :options="languageOptions"
                :aria-label="t('controls.language')"
                @update:model-value="setLanguage"
              />

              <span class="site-header__toggles-separator" aria-hidden="true">|</span>

              <BaseSegmentedToggle
                class="site-header__toggle"
                :model-value="theme"
                :options="themeOptions"
                :aria-label="t('controls.theme')"
                @update:model-value="setTheme"
              />
            </div>

            <div v-if="isAuthenticated" class="site-header__controls-end">
              <UserChip :user="authStore.user" />

              <BaseButton class="site-header__logout-button" variant="text" :aria-label="t('nav.logout')" @click="handleLogout">
                <svg viewBox="0 0 24 24" class="site-header__logout-icon" aria-hidden="true" focusable="false">
                  <path
                    d="M13 4h4a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2h-4"
                    fill="none"
                    stroke="currentColor"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="1.6"
                  />
                  <path
                    d="M10 16l4-4-4-4"
                    fill="none"
                    stroke="currentColor"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="1.6"
                  />
                  <path
                    d="M14 12H5"
                    fill="none"
                    stroke="currentColor"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="1.6"
                  />
                </svg>
              </BaseButton>
            </div>
          </div>
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
  background: linear-gradient(
    180deg,
    color-mix(in srgb, var(--ds-color-accent-primary) 30%, transparent) 0%,
    color-mix(in srgb, var(--ds-color-accent-primary) 14%, transparent) 58%,
    transparent 100%
  );
  backdrop-filter: blur(12px);
}

.site-header__inner {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 18px var(--ds-space-5);
}

.site-header__top-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
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

.brand-mark__name {
  font-family: var(--ds-font-display);
  font-size: 20px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--ds-color-accent-primary);
}

.site-header__menu-toggle {
  display: none;
  width: 42px;
  height: 42px;
  align-items: center;
  justify-content: center;
  padding: 0;
  border: 0;
  background: transparent;
  color: var(--ds-color-accent-primary);
  cursor: pointer;
}

.site-header__menu-toggle span {
  position: absolute;
  width: 18px;
  height: 1px;
  background: currentColor;
  transition:
    transform var(--ds-motion-base),
    opacity var(--ds-motion-base);
}

.site-header__menu-toggle span:nth-child(1) {
  transform: translateY(-6px);
}

.site-header__menu-toggle span:nth-child(2) {
  transform: translateY(0);
}

.site-header__menu-toggle span:nth-child(3) {
  transform: translateY(6px);
}

.site-header__menu-toggle.is-open span:nth-child(1) {
  transform: rotate(45deg);
}

.site-header__menu-toggle.is-open span:nth-child(2) {
  opacity: 0;
}

.site-header__menu-toggle.is-open span:nth-child(3) {
  transform: rotate(-45deg);
}

.site-header__menu-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex: 1;
  min-width: 0;
}

.site-header__nav {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-left: 0;
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
}

.site-header__toggles {
  display: flex;
  align-items: center;
  gap: 10px;
}

.site-header__toggles-separator {
  display: none;
}

.site-header__controls-end {
  display: flex;
  align-items: center;
  gap: 10px;
}

.site-header__logout-button {
  min-width: 36px;
  min-height: 36px;
  padding: 0;
}

.site-header__logout-icon {
  width: 36px;
  height: 36px;
}

@media (max-width: 860px) {
  .site-header__inner {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
    padding: 14px var(--ds-space-5);
  }

  .site-header__top-row {
    position: relative;
    width: 100%;
    min-height: 48px;
    justify-content: center;
  }

  .site-header__menu-toggle {
    position: absolute;
    top: 50%;
    right: 4px;
    transform: translateY(-50%);
    display: inline-flex;
    flex-shrink: 0;
  }

  .site-header__brand {
    width: 100%;
    min-width: 0;
  }

  .brand-mark {
    display: grid;
    grid-template-columns: 42px minmax(0, 1fr) 42px;
    align-items: center;
    width: 100%;
    gap: 0;
  }

  .brand-mark::after {
    content: "";
    width: 42px;
    height: 42px;
    grid-column: 3;
  }

  .brand-mark__symbol {
    grid-column: 1;
    justify-self: start;
  }

  .brand-mark__name {
    grid-column: 2;
    justify-self: center;
    text-align: center;
  }

  .site-header__menu-panel {
    display: none;
    width: 100%;
    border-top: 1px solid rgba(212, 175, 55, 0.14);
    padding-top: 14px;
  }

.site-header__menu-panel.is-open {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    gap: 7px;
    width: 100%;
  }

  .site-header__nav {
    order: 2;
    width: 100%;
    justify-self: stretch;
    margin-left: 0;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: stretch;
    justify-content: flex-start;
    gap: 6px;
  }

  .site-header__controls {
    order: 1;
    width: 100%;
    justify-self: stretch;
    display: flex;
    flex-wrap: nowrap;
    align-items: center;
    justify-content: space-between;
    gap: 6px;
  }

  .site-header__toggles {
    gap: 0;
    margin-right: auto;
  }

  .site-header__toggles-separator {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    color: var(--ds-color-text-secondary);
    line-height: 1.4;
  }

  .site-header__controls :deep(.base-segmented-toggle) {
    width: fit-content;
  }

  .site-header__controls-end {
    margin-left: auto;
    gap: 6px;
  }

  .site-header__controls :deep(.site-header__toggle.base-segmented-toggle) {
    display: flex;
    padding: 0 6px;
    border: 0;
    background: transparent;
    box-shadow: none;
  }

  .site-header__controls :deep(.site-header__toggle .base-segmented-toggle__thumb) {
    display: none;
  }

  .site-header__controls :deep(.site-header__toggle .base-segmented-toggle__option) {
    min-height: auto;
    padding: 0;
    color: var(--ds-color-text-secondary);
  }

  .site-header__controls :deep(.site-header__toggle .base-segmented-toggle__option.is-active) {
    display: none;
  }

  .site-header__controls :deep(.site-header__toggle .base-segmented-toggle__option:hover) {
    color: var(--ds-color-accent-primary);
  }

  .site-header__controls :deep(.user-chip) {
    justify-content: flex-end;
  }

  .site-header__controls :deep(.user-chip__meta) {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .site-header__controls :deep(.user-chip__avatar) {
    display: none;
  }

  .site-header__logout-icon {
    width: 18px;
    height: 18px;
  }

  .site-header__nav-link {
    flex: 1 1 0;
    display: flex;
    align-items: center;
    justify-content: center;
    text-align: center;
  }
}

@media (max-width: 720px) {
  .site-header__inner {
    padding: 10px var(--ds-space-4);
  }

  .brand-mark {
    gap: 10px;
  }

  .brand-mark__symbol {
    width: 42px;
    height: 42px;
  }

  .brand-mark__name {
    font-size: 16px;
    letter-spacing: 0.1em;
  }

  .site-header__nav-link {
    padding: 8px 10px;
  }
}
</style>
