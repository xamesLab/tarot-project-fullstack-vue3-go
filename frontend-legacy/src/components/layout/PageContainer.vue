<script setup>
import { computed } from 'vue'

const props = defineProps({
  wide: {
    type: Boolean,
    default: false,
  },
  tag: {
    type: String,
    default: 'div',
  },
})

const classes = computed(() => ['page-container', { 'page-container--wide': props.wide }])
</script>

<template>
  <component :is="tag" :class="classes">
    <slot />
  </component>
</template>
