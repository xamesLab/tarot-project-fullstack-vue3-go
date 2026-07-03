<script setup>
import { computed } from 'vue'

const props = defineProps({
  tag: {
    type: String,
    default: 'section',
  },
  tone: {
    type: String,
    default: 'default',
  },
  padding: {
    type: String,
    default: 'md',
  },
  hoverable: {
    type: Boolean,
    default: false,
  },
})

const classes = computed(() => [
  'base-surface',
  `base-surface--${props.tone === 'default' ? 'outlined' : props.tone}`,
  `base-surface--padding-${props.padding}`,
  { 'base-surface--hoverable': props.hoverable },
])
</script>

<template>
  <component :is="tag" :class="classes">
    <slot />
  </component>
</template>
