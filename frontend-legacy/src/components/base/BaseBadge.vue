<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: {
    type: String,
    default: '',
  },
  variant: {
    type: String,
    default: 'accent',
  },
  size: {
    type: String,
    default: 'md',
  },
})

const classes = computed(() => [
  'base-badge',
  `base-badge--${props.variant}`,
  `base-badge--${props.size}`,
])
</script>

<template>
  <span :class="classes">
    <slot>{{ label }}</slot>
  </span>
</template>

<style scoped>
.base-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: fit-content;
  min-height: 28px;
  padding: 6px 14px;
  border: var(--ds-border-width) solid var(--ds-color-border-accent);
  background: rgba(212, 175, 55, 0.1);
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-caption);
  font-weight: 600;
  letter-spacing: 0.1em;
  line-height: 1;
  text-transform: uppercase;
  white-space: nowrap;
}

.base-badge--sm {
  min-height: 24px;
  padding: 5px 10px;
  font-size: 9px;
}

.base-badge--accent {
  border-color: var(--ds-color-border-accent);
  background: rgba(212, 175, 55, 0.1);
  color: var(--ds-color-accent-primary);
}

.base-badge--subtle {
  border-color: var(--ds-color-border-subtle);
  background: rgba(255, 255, 255, 0.04);
  color: var(--ds-color-text-secondary);
}

.base-badge--danger {
  border-color: color-mix(in srgb, var(--ds-color-danger) 55%, transparent);
  background: color-mix(in srgb, var(--ds-color-danger) 12%, transparent);
  color: var(--ds-color-danger);
}

.base-badge--success {
  border-color: color-mix(in srgb, var(--ds-color-accent-primary) 45%, transparent);
  background: color-mix(in srgb, var(--ds-color-accent-primary) 14%, transparent);
  color: var(--ds-color-accent-primary);
}

@media (max-width: 768px) {
  .base-badge {
    min-height: 24px;
    padding: 5px 10px;
  }
}
</style>
