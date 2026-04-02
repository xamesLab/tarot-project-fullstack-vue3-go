<script setup>
import { computed, onMounted, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import AdminLayout from '../../components/admin/AdminLayout.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseInput from '../../components/base/BaseInput.vue'
import BaseSelect from '../../components/base/BaseSelect.vue'
import BaseSurface from '../../components/base/BaseSurface.vue'
import FormField from '../../components/form/FormField.vue'
import { useAuthStore } from '../../stores/auth'
import { useAdminStore } from '../../stores/admin'

const { t, locale } = useI18n()
const adminStore = useAdminStore()
const authStore = useAuthStore()

const filters = reactive({
  search: '',
  role: 'all',
  status: 'all',
})

const rowDrafts = reactive({})

const roleOptions = computed(() => [
  { value: 'all', label: t('admin.filters.allRoles') },
  { value: 'seeker', label: t('roles.seeker') },
  { value: 'adept', label: t('roles.adept') },
  { value: 'master', label: t('roles.master') },
  { value: 'keeper', label: t('roles.keeper') },
])

const statusOptions = computed(() => [
  { value: 'all', label: t('admin.filters.allStatuses') },
  { value: 'active', label: t('admin.statuses.active') },
  { value: 'disabled', label: t('admin.statuses.disabled') },
])

watch(
  () => adminStore.userList,
  (users) => {
    users.forEach((user) => {
      rowDrafts[user.id] = {
        role: user.role,
        status: user.status,
      }
    })
  },
  { immediate: true },
)

function localizeDate(value) {
  if (!value) {
    return '—'
  }

  return new Intl.DateTimeFormat(locale.value, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

function isSaving(userId) {
  return adminStore.userMutationStatusById[userId] === 'loading'
}

async function loadUsers(page = 1) {
  await adminStore.loadUsers({
    page,
    search: filters.search,
    role: filters.role,
    status: filters.status,
  })
}

function applyFilters() {
  loadUsers(1)
}

function resetFilters() {
  filters.search = ''
  filters.role = 'all'
  filters.status = 'all'
  loadUsers(1)
}

async function saveRole(user) {
  const draft = rowDrafts[user.id]
  if (!draft || draft.role === user.role) {
    return
  }

  const result = await adminStore.updateUserRole(user.id, draft.role)
  if (!result.ok) {
    draft.role = user.role
  }
}

async function saveStatus(user) {
  const draft = rowDrafts[user.id]
  if (!draft || draft.status === user.status) {
    return
  }

  const result = await adminStore.updateUserStatus(user.id, draft.status)
  if (!result.ok) {
    draft.status = user.status
  }
}

onMounted(() => {
  filters.search = adminStore.userFilters.search
  filters.role = adminStore.userFilters.role
  filters.status = adminStore.userFilters.status

  if (adminStore.usersStatus === 'idle') {
    loadUsers()
  }
})
</script>

<template>
  <AdminLayout :title="t('admin.usersTitle')" :copy="t('admin.usersCopy')">
    <div class="page-stack">
      <BaseSurface padding="lg" tag="section">
        <form class="filters-grid" @submit.prevent="applyFilters">
          <FormField :label="t('admin.filters.search')" for-id="admin-user-search">
            <BaseInput
              id="admin-user-search"
              v-model="filters.search"
              :placeholder="t('admin.filters.searchPlaceholder')"
            />
          </FormField>

          <FormField :label="t('admin.filters.role')" for-id="admin-user-role-filter">
            <BaseSelect id="admin-user-role-filter" v-model="filters.role" :options="roleOptions" />
          </FormField>

          <FormField :label="t('admin.filters.status')" for-id="admin-user-status-filter">
            <BaseSelect id="admin-user-status-filter" v-model="filters.status" :options="statusOptions" />
          </FormField>

          <div class="filters-actions">
            <BaseButton type="submit" :loading="adminStore.usersStatus === 'loading'">{{ t('admin.applyFilters') }}</BaseButton>
            <BaseButton type="button" variant="secondary" @click="resetFilters">{{ t('admin.resetFilters') }}</BaseButton>
          </div>
        </form>
      </BaseSurface>

      <div v-if="adminStore.usersError" class="admin-message admin-message--error">
        {{ adminStore.usersError }}
      </div>

      <BaseSurface padding="lg" tag="section">
        <div class="list-head">
          <div>
            <p class="text-eyebrow">{{ t('admin.userTable.title') }}</p>
            <p class="text-body-muted">{{ t('admin.userTable.copy', { total: adminStore.userPagination.total }) }}</p>
          </div>
        </div>

        <div v-if="!adminStore.userList.length && adminStore.usersStatus === 'ready'" class="empty-state">
          {{ t('admin.emptyUsers') }}
        </div>

        <div v-else class="user-list">
          <article v-for="user in adminStore.userList" :key="user.id" class="user-card">
            <div class="user-card__identity">
              <div>
                <h3 class="user-card__name">{{ user.username || t('admin.userTable.noUsername') }}</h3>
                <p class="text-body-muted">{{ user.email }}</p>
              </div>
              <span class="user-chip">{{ t(`roles.${user.role}`) }}</span>
            </div>

            <dl class="user-meta">
              <div>
                <dt>{{ t('admin.userTable.status') }}</dt>
                <dd>{{ t(`admin.statuses.${user.status || 'active'}`) }}</dd>
              </div>
              <div>
                <dt>{{ t('admin.userTable.timezone') }}</dt>
                <dd>{{ user.timezone || 'UTC' }}</dd>
              </div>
              <div>
                <dt>{{ t('admin.userTable.createdAt') }}</dt>
                <dd>{{ localizeDate(user.created_at) }}</dd>
              </div>
              <div>
                <dt>{{ t('admin.userTable.updatedAt') }}</dt>
                <dd>{{ localizeDate(user.updated_at) }}</dd>
              </div>
            </dl>

            <div v-if="rowDrafts[user.id]" class="user-actions">
              <FormField :label="t('admin.userTable.role')">
                <BaseSelect v-model="rowDrafts[user.id].role" :options="roleOptions.slice(1)" />
              </FormField>

              <FormField :label="t('admin.userTable.status')">
                <BaseSelect v-model="rowDrafts[user.id].status" :options="statusOptions.slice(1)" />
              </FormField>

              <div class="user-actions__buttons">
                <BaseButton
                  variant="secondary"
                  :disabled="isSaving(user.id)"
                  :loading="isSaving(user.id)"
                  @click="saveRole(user)"
                >
                  {{ t('admin.userTable.saveRole') }}
                </BaseButton>

                <BaseButton
                  variant="secondary"
                  :disabled="isSaving(user.id)"
                  :loading="isSaving(user.id)"
                  @click="saveStatus(user)"
                >
                  {{ t('admin.userTable.saveStatus') }}
                </BaseButton>
              </div>
            </div>

            <p v-if="adminStore.userMutationErrorById[user.id]" class="admin-message admin-message--inline">
              {{ adminStore.userMutationErrorById[user.id] }}
            </p>

            <p v-if="authStore.user?.id === user.id" class="text-body-muted">
              {{ t('admin.userTable.selfGuardHint') }}
            </p>
          </article>
        </div>

        <div class="pagination-bar">
          <BaseButton
            variant="text"
            :disabled="adminStore.userPagination.page <= 1 || adminStore.usersStatus === 'loading'"
            @click="loadUsers(adminStore.userPagination.page - 1)"
          >
            {{ t('admin.pagination.previous') }}
          </BaseButton>
          <span class="text-body-muted">
            {{ t('admin.pagination.page', { page: adminStore.userPagination.page, total: adminStore.userPagination.total_pages }) }}
          </span>
          <BaseButton
            variant="text"
            :disabled="adminStore.userPagination.page >= adminStore.userPagination.total_pages || adminStore.usersStatus === 'loading'"
            @click="loadUsers(adminStore.userPagination.page + 1)"
          >
            {{ t('admin.pagination.next') }}
          </BaseButton>
        </div>
      </BaseSurface>
    </div>
  </AdminLayout>
</template>

<style scoped>
.filters-grid {
  display: grid;
  grid-template-columns: minmax(0, 2fr) repeat(2, minmax(180px, 1fr)) auto;
  gap: 16px;
  align-items: end;
}

.filters-actions,
.user-actions__buttons,
.pagination-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.list-head {
  margin-bottom: 24px;
}

.user-list {
  display: grid;
  gap: 18px;
}

.user-card {
  display: grid;
  gap: 18px;
  padding: 20px;
  border: 1px solid var(--ds-color-border-subtle);
  background: var(--ds-color-surface-secondary);
}

.user-card__identity {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.user-card__name {
  font-size: var(--ds-text-heading-3);
  color: var(--ds-color-accent-primary);
}

.user-chip {
  padding: 6px 12px;
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-caption);
  letter-spacing: 0.08em;
  text-transform: uppercase;
  border: 1px solid var(--ds-color-border-accent);
  color: var(--ds-color-accent-primary);
}

.user-meta {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.user-meta dt {
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-caption);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--ds-color-text-secondary);
  margin-bottom: 6px;
}

.user-meta dd {
  color: var(--ds-color-text-primary);
}

.user-actions {
  display: grid;
  grid-template-columns: minmax(180px, 220px) minmax(180px, 220px) auto;
  gap: 16px;
  align-items: end;
}

.admin-message {
  padding: 14px 16px;
  border: 1px solid var(--ds-color-danger);
  color: var(--ds-color-danger);
  background: rgba(178, 58, 44, 0.12);
}

.admin-message--inline {
  padding: 0;
  border: 0;
  background: transparent;
}

.empty-state {
  padding: 24px 0;
  color: var(--ds-color-text-secondary);
}

.pagination-bar {
  justify-content: space-between;
  margin-top: 24px;
}

@media (max-width: 1100px) {
  .filters-grid,
  .user-actions,
  .user-meta {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .filters-grid,
  .user-actions,
  .user-meta {
    grid-template-columns: 1fr;
  }

  .user-card__identity,
  .pagination-bar {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
