<script setup>
import { computed } from 'vue'

const props = defineProps({
  user: {
    type: Object,
    default: null,
  },
})

const initials = computed(() => {
  const source = props.user?.username || props.user?.email || 'TD'

  return source
    .split(/[\s._-]+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() || '')
    .join('')
})
</script>

<template>
  <div class="user-chip">
    <div class="user-chip__avatar">{{ initials }}</div>
    <div class="user-chip__meta">
      <span class="user-chip__name">{{ user?.username || 'Guest' }}</span>
      <span class="user-chip__role">{{ user?.role || 'visitor' }}</span>
    </div>
  </div>
</template>

<style scoped>
.user-chip {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 14px;
  border: 1px solid var(--ds-color-border-accent);
  background: var(--ds-color-bg-surface);
}

.user-chip__avatar {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, var(--ds-color-accent-dark), var(--ds-color-accent-primary));
  color: var(--ds-color-bg-page);
  font-family: var(--ds-font-display);
  font-size: 14px;
  font-weight: 700;
}

.user-chip__meta {
  display: grid;
  gap: 2px;
}

.user-chip__name {
  font-family: var(--ds-font-ui);
  font-size: 12px;
  font-weight: 600;
  color: var(--ds-color-text-primary);
}

.user-chip__role {
  font-family: var(--ds-font-ui);
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--ds-color-accent-primary);
}
</style>
