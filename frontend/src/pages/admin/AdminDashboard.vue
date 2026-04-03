<script setup>
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import AdminLayout from '../../components/admin/AdminLayout.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseSurface from '../../components/base/BaseSurface.vue'
import BaseTabs from '../../components/base/BaseTabs.vue'
import { useAdminStore } from '../../stores/admin'

const { t } = useI18n()
const adminStore = useAdminStore()
const activeAdminTab = ref('users')

const adminTabs = computed(() => [
  { value: 'users', label: t('admin.usersTitle') },
  { value: 'readings', label: t('admin.dashboard.readingsTitle') },
  { value: 'decks', label: t('admin.decksTitle') },
])

const metricCards = computed(() => {
  const metrics = adminStore.dashboardMetrics

  return [
    {
      key: 'users',
      label: t('admin.dashboard.totalUsers'),
      value: metrics?.users?.total ?? '—',
      meta: t('admin.dashboard.usersMeta', {
        active: metrics?.users?.by_status?.active ?? 0,
        disabled: metrics?.users?.by_status?.disabled ?? 0,
      }),
    },
    {
      key: 'today',
      label: t('admin.dashboard.readingsToday'),
      value: metrics?.readings?.today?.total ?? '—',
      meta: t('admin.dashboard.readingsSplit', {
        cardOfDay: metrics?.readings?.today?.card_of_day ?? 0,
        threeCard: metrics?.readings?.today?.three_card ?? 0,
      }),
    },
    {
      key: 'week',
      label: t('admin.dashboard.readingsWeek'),
      value: metrics?.readings?.this_week?.total ?? '—',
      meta: t('admin.dashboard.readingsSplit', {
        cardOfDay: metrics?.readings?.this_week?.card_of_day ?? 0,
        threeCard: metrics?.readings?.this_week?.three_card ?? 0,
      }),
    },
    {
      key: 'decks',
      label: t('admin.dashboard.totalDecks'),
      value: metrics?.decks?.total ?? '—',
      meta: t('admin.dashboard.decksMeta', {
        done: metrics?.decks?.by_status?.done ?? 0,
        processing: metrics?.decks?.by_status?.processing ?? 0,
        disabled: metrics?.decks?.by_status?.disabled ?? 0,
      }),
    },
  ]
})

onMounted(() => {
  if (adminStore.dashboardStatus === 'idle') {
    adminStore.loadDashboard()
  }
})
</script>

<template>
  <AdminLayout :title="t('admin.title')" :copy="t('admin.copy')">
    <div class="page-stack">
      <div v-if="adminStore.dashboardError" class="admin-message admin-message--error">
        {{ adminStore.dashboardError }}
      </div>

      <section class="admin-metrics">
        <BaseSurface
          v-for="card in metricCards"
          :key="card.key"
          tag="article"
          class="metric-card"
        >
          <p class="text-eyebrow">{{ card.label }}</p>
          <p class="metric-card__value">{{ card.value }}</p>
          <p class="text-body-muted">{{ card.meta }}</p>
        </BaseSurface>
      </section>

      <BaseSurface padding="lg" tag="section">
        <BaseTabs v-model="activeAdminTab" :items="adminTabs">
          <template #panel="{ activeItem }">
            <div v-if="activeItem?.value === 'users'" class="admin-panel">
              <div class="admin-panel__header">
                <div>
                  <p class="text-eyebrow">{{ t('admin.usersTitle') }}</p>
                  <p class="text-body-muted">{{ t('admin.dashboard.usersPanelCopy') }}</p>
                </div>
                <BaseButton to="/admin/users" variant="secondary">{{ t('admin.openUsers') }}</BaseButton>
              </div>

              <div class="summary-grid">
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.users?.by_role?.seeker ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('roles.seeker') }}</p>
                </BaseSurface>
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.users?.by_role?.adept ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('roles.adept') }}</p>
                </BaseSurface>
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.users?.by_role?.master ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('roles.master') }}</p>
                </BaseSurface>
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.users?.by_role?.keeper ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('roles.keeper') }}</p>
                </BaseSurface>
              </div>
            </div>

            <div v-else-if="activeItem?.value === 'readings'" class="admin-panel">
              <div class="admin-panel__header">
                <div>
                  <p class="text-eyebrow">{{ t('admin.dashboard.readingsTitle') }}</p>
                  <p class="text-body-muted">{{ t('admin.dashboard.readingsPanelCopy') }}</p>
                </div>
              </div>

              <div class="summary-grid summary-grid--wide">
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.readings?.today?.card_of_day ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('admin.dashboard.cardOfDayToday') }}</p>
                </BaseSurface>
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.readings?.today?.three_card ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('admin.dashboard.threeCardToday') }}</p>
                </BaseSurface>
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.readings?.this_week?.card_of_day ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('admin.dashboard.cardOfDayWeek') }}</p>
                </BaseSurface>
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.readings?.this_week?.three_card ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('admin.dashboard.threeCardWeek') }}</p>
                </BaseSurface>
              </div>
            </div>

            <div v-else class="admin-panel">
              <div class="admin-panel__header">
                <div>
                  <p class="text-eyebrow">{{ t('admin.decksTitle') }}</p>
                  <p class="text-body-muted">{{ t('admin.dashboard.decksPanelCopy') }}</p>
                </div>
                <BaseButton to="/admin/decks" variant="secondary">{{ t('admin.openDecks') }}</BaseButton>
              </div>

              <div class="summary-grid">
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.decks?.by_status?.done ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('admin.deckStatuses.done') }}</p>
                </BaseSurface>
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.decks?.by_status?.processing ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('admin.deckStatuses.processing') }}</p>
                </BaseSurface>
                <BaseSurface tag="article" class="summary-card" tone="outlined">
                  <p class="summary-card__value">{{ adminStore.dashboardMetrics?.decks?.by_status?.disabled ?? 0 }}</p>
                  <p class="text-body-muted">{{ t('admin.deckStatuses.disabled') }}</p>
                </BaseSurface>
              </div>
            </div>
          </template>
        </BaseTabs>
      </BaseSurface>
    </div>
  </AdminLayout>
</template>

<style scoped>
.admin-metrics {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.metric-card {
  display: grid;
  gap: 10px;
}

.metric-card__value,
.summary-card__value {
  font-family: var(--ds-font-display);
  font-size: clamp(2rem, 4vw, 2.8rem);
  color: var(--ds-color-accent-primary);
  line-height: 1;
}

.admin-panel {
  display: grid;
  gap: 24px;
}

.admin-panel__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.summary-grid--wide {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.summary-card {
  display: grid;
  gap: 8px;
}

.admin-message {
  padding: 14px 16px;
  border: 1px solid var(--ds-color-danger);
  color: var(--ds-color-danger);
  background: rgba(178, 58, 44, 0.12);
}

@media (max-width: 1100px) {
  .admin-metrics,
  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .admin-metrics,
  .summary-grid,
  .summary-grid--wide {
    grid-template-columns: 1fr;
  }

  .admin-panel__header {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
