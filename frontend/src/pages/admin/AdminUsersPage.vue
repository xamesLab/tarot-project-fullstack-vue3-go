<template>
  <section class="panel stack">
    <div class="heading-row">
      <h2>{{ t("adminUsersTitle") }}</h2>
      <button class="btn btn-secondary" type="button" @click="reload" :disabled="loading">
        {{ loading ? t("adminRefreshing") : t("adminRefresh") }}
      </button>
    </div>
    <p class="muted">{{ t("adminUsersHint") }}</p>

    <form class="filters-grid" @submit.prevent="applyFilters">
      <div class="field">
        <label for="users-search">{{ t("adminUsersSearch") }}</label>
        <input id="users-search" v-model.trim="searchInput" />
      </div>

      <div class="field">
        <label for="users-role">{{ t("adminUsersRoleFilter") }}</label>
        <select id="users-role" v-model="roleFilter">
          <option value="">all</option>
          <option value="user">user</option>
          <option value="admin">admin</option>
        </select>
      </div>

      <div class="field">
        <label for="users-status">{{ t("adminUsersStatusFilter") }}</label>
        <select id="users-status" v-model="statusFilter">
          <option value="">all</option>
          <option value="active">active</option>
          <option value="banned">banned</option>
          <option value="freeze">freeze</option>
        </select>
      </div>

      <div class="field">
        <label for="users-sort">{{ t("adminUsersSort") }}</label>
        <select id="users-sort" v-model="sortBy">
          <option value="registeredAt:desc">{{ t("adminUsersSortDesc") }}</option>
          <option value="registeredAt:asc">{{ t("adminUsersSortAsc") }}</option>
        </select>
      </div>

      <div class="actions">
        <button class="btn" type="submit" :disabled="loading">{{ t("adminUsersApplyFilters") }}</button>
        <button class="btn btn-secondary" type="button" :disabled="loading" @click="resetFilters">
          {{ t("adminUsersResetFilters") }}
        </button>
      </div>
    </form>

    <p v-if="error" class="error">{{ error }}</p>

    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>{{ t("adminUsersTableName") }}</th>
            <th>{{ t("adminUsersTableRole") }}</th>
            <th>{{ t("adminUsersTableRegistered") }}</th>
            <th>{{ t("adminUsersTableStatus") }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="4" class="muted">{{ t("commonLoading") }}</td>
          </tr>
          <tr v-else-if="users.length === 0">
            <td colspan="4" class="muted">{{ t("adminUsersNoData") }}</td>
          </tr>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.name }}</td>
            <td>{{ user.role }}</td>
            <td>{{ formatDate(user.registeredAt) }}</td>
            <td>{{ user.status }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pager">
      <button class="btn btn-secondary" type="button" :disabled="loading || page <= 1" @click="goPrev">
        {{ t("adminUsersPrev") }}
      </button>
      <span class="muted">{{ page }} / {{ totalPages }}</span>
      <button class="btn btn-secondary" type="button" :disabled="loading || page >= totalPages" @click="goNext">
        {{ t("adminUsersNext") }}
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";

import { useI18n } from "../../composables/useI18n";
import * as adminApi from "../../services/admin.api";
import { useAuthStore } from "../../stores/auth.store";

const authStore = useAuthStore();
const { t } = useI18n();

const loading = ref(false);
const error = ref("");

const users = ref<adminApi.AdminUserListItem[]>([]);
const page = ref(1);
const totalPages = ref(1);

const searchInput = ref("");
const roleFilter = ref("");
const statusFilter = ref<"" | "active" | "banned" | "freeze">("");
const sortBy = ref<"registeredAt:desc" | "registeredAt:asc">("registeredAt:desc");

async function loadUsers() {
  if (!authStore.token) return;
  loading.value = true;
  error.value = "";
  try {
    const data = await adminApi.getAdminUsers(authStore.token, {
      page: page.value,
      pageSize: 20,
      role: roleFilter.value || undefined,
      status: statusFilter.value || undefined,
      sort: sortBy.value,
      name: searchInput.value || undefined
    });
    users.value = data.items;
    totalPages.value = Math.max(data.totalPages, 1);
    page.value = data.page;
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to load users";
  } finally {
    loading.value = false;
  }
}

async function applyFilters() {
  page.value = 1;
  await loadUsers();
}

async function resetFilters() {
  searchInput.value = "";
  roleFilter.value = "";
  statusFilter.value = "";
  sortBy.value = "registeredAt:desc";
  page.value = 1;
  await loadUsers();
}

async function reload() {
  await loadUsers();
}

async function goPrev() {
  if (page.value <= 1) return;
  page.value -= 1;
  await loadUsers();
}

async function goNext() {
  if (page.value >= totalPages.value) return;
  page.value += 1;
  await loadUsers();
}

function formatDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString();
}

onMounted(loadUsers);
</script>

<style scoped>
.heading-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.heading-row h2 {
  margin: 0;
}

.filters-grid {
  display: grid;
  gap: 10px;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
}

.actions {
  display: flex;
  gap: 8px;
  align-items: end;
}

.table-wrap {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  border: 1px solid var(--border);
  text-align: left;
  padding: 10px;
}

thead th {
  background: var(--surface-soft);
}

.pager {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}
</style>
