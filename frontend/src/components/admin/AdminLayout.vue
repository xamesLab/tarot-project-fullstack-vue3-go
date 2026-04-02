<script setup>
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

import PageContainer from '../layout/PageContainer.vue'
import PageSection from '../layout/PageSection.vue'

const props = defineProps({
  title: {
    type: String,
    default: '',
  },
  copy: {
    type: String,
    default: '',
  },
})

const route = useRoute()
const { t } = useI18n()

const navigationItems = computed(() => [
  { to: '/admin', label: t('admin.title') },
  { to: '/admin/users', label: t('admin.usersTitle') },
  { to: '/admin/role-policies', label: t('admin.policiesTitle') },
  { to: '/admin/settings', label: t('admin.settingsTitle') },
  { to: '/admin/decks', label: t('admin.decksTitle') },
])

const breadcrumb = computed(() => `${t('admin.title')} / ${props.title || t('admin.title')}`)

function isActive(path) {
  if (path === '/admin') {
    return route.path === path
  }

  return route.path === path || route.path.startsWith(`${path}/`)
}
</script>

<template>
  <PageContainer>
    <div class="page-stack">
      <PageSection :eyebrow="breadcrumb" :title="title" :copy="copy" />

      <div class="admin-layout">
        <aside class="admin-layout__sidebar">
          <nav class="admin-nav" aria-label="Admin">
            <RouterLink
              v-for="item in navigationItems"
              :key="item.to"
              class="admin-nav__link"
              :class="{ 'is-active': isActive(item.to) }"
              :to="item.to"
            >
              {{ item.label }}
            </RouterLink>
          </nav>
        </aside>

        <div class="admin-layout__content">
          <slot />
        </div>
      </div>
    </div>
  </PageContainer>
</template>

<style scoped>
.admin-layout {
  display: grid;
  grid-template-columns: minmax(220px, 260px) minmax(0, 1fr);
  gap: 24px;
  align-items: start;
}

.admin-layout__sidebar {
  position: sticky;
  top: calc(var(--ds-header-height) + 24px);
}

.admin-nav {
  display: grid;
  gap: 8px;
  padding: 18px;
  border: 1px solid var(--ds-color-border-subtle);
  background: var(--ds-color-surface-secondary);
  box-shadow: var(--ds-shadow-soft);
}

.admin-nav__link {
  padding: 12px 14px;
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-caption);
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--ds-color-text-secondary);
  border: 1px solid transparent;
  transition: color var(--ds-motion-fast), border-color var(--ds-motion-fast), background var(--ds-motion-fast);
}

.admin-nav__link:hover,
.admin-nav__link.is-active {
  color: var(--ds-color-accent-primary);
  border-color: var(--ds-color-border-accent);
  background: rgba(212, 175, 55, 0.08);
}

.admin-layout__content {
  min-width: 0;
}

@media (max-width: 960px) {
  .admin-layout {
    grid-template-columns: 1fr;
  }

  .admin-layout__sidebar {
    position: static;
  }

  .admin-nav {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }
}
</style>
