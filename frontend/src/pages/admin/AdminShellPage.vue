<template>
  <main class="page stack admin-shell-page">
    <AdminHeader />

    <section class="panel admin-tabs">
      <RouterLink class="tab" to="/admin/users" active-class="tab-active">{{ t("adminNavUsers") }}</RouterLink>
      <RouterLink class="tab" to="/admin/decks" active-class="tab-active">{{ t("adminNavDecks") }}</RouterLink>
      <RouterLink class="tab" to="/admin/settings" active-class="tab-active">
        {{ t("adminNavSettings") }}
      </RouterLink>
    </section>

    <RouterView />
  </main>
</template>

<script setup lang="ts">
import AdminHeader from "../../components/AdminHeader.vue";
import { useI18n } from "../../composables/useI18n";

const { t } = useI18n();
</script>

<style scoped>
.admin-shell-page {
  max-width: 1140px;
}

.admin-tabs {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  padding: 10px;
  background:
    linear-gradient(135deg, rgba(49, 95, 115, 0.1) 0%, transparent 50%),
    linear-gradient(145deg, var(--surface-soft) 0%, var(--surface) 100%);
}

.tab {
  border: 1px solid var(--border-strong);
  border-radius: 999px;
  padding: 10px 16px;
  background: rgba(255, 248, 230, 0.08);
  font-weight: 600;
  letter-spacing: 0.04em;
}

.tab-active {
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-deep) 100%);
  color: var(--accent-contrast);
  border-color: transparent;
}
</style>
