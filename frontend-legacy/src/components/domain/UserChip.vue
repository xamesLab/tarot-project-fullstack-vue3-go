<script setup>
import BaseBadge from '../base/BaseBadge.vue'

const props = defineProps({
  user: {
    type: Object,
    default: null,
  },
})
</script>

<template>
  <div class="user-chip">
    <div class="user-chip__avatar" aria-hidden="true">
      <svg viewBox="0 0 24 24" class="user-chip__avatar-icon" focusable="false">
        <path
          d="M12 12a4 4 0 1 0-4-4 4 4 0 0 0 4 4Zm0 2c-4.418 0-8 2.239-8 5v1h16v-1c0-2.761-3.582-5-8-5Z"
          fill="currentColor"
        />
      </svg>
    </div>
    <div class="user-chip__meta">
      <span class="user-chip__name">{{ user?.username || 'Guest' }}</span>
      <BaseBadge size="sm" :label="user?.role || 'visitor'" />
    </div>
  </div>
</template>

<style scoped>
.user-chip {
  display: flex;
  align-items: flex-end;
  gap: 10px;
  padding: 3px 0;
  border: 0;
  background: transparent;
}

.user-chip__avatar {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  background: linear-gradient(
    135deg,
    color-mix(in srgb, var(--ds-color-accent-dark) 50%, transparent),
    color-mix(in srgb, var(--ds-color-accent-primary) 50%, transparent)
  );
  color: var(--ds-color-accent-primary);
}

.user-chip__avatar-icon {
  width: 20px;
  height: 20px;
}

.user-chip__meta {
  display: grid;
  gap: 3px;
}

.user-chip__name {
  font-family: var(--ds-font-ui);
  font-size: 12px;
  font-weight: 600;
  color: var(--ds-color-text-primary);
}

</style>
