<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

const props = defineProps({
  to: {
    type: [String, Object],
    default: null,
  },
  type: {
    type: String,
    default: 'button',
  },
  variant: {
    type: String,
    default: 'primary',
  },
  size: {
    type: String,
    default: 'md',
  },
  fullWidth: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['click'])

const isDisabled = computed(() => props.disabled || props.loading)
const component = computed(() => (props.to ? RouterLink : 'button'))
const componentProps = computed(() => (
  props.to
    ? { to: props.to }
    : { type: props.type, disabled: isDisabled.value || undefined }
))

const classes = computed(() => [
  'base-button',
  `base-button--${props.variant}`,
  `base-button--${props.size}`,
  {
    'base-button--full': props.fullWidth,
    'is-disabled': isDisabled.value,
  },
])

function handleClick(event) {
  if (isDisabled.value) {
    event.preventDefault()
    return
  }

  emit('click', event)
}
</script>

<template>
  <component :is="component" v-bind="componentProps" :class="classes" @click="handleClick">
    <span v-if="loading" class="base-button__loader" aria-hidden="true"></span>
    <slot />
  </component>
</template>
