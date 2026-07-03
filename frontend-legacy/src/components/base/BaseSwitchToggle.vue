<script setup>
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  label: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div class="base-switch">
    <span v-if="label" class="base-switch__label">{{ label }}</span>
    <button
      type="button"
      class="base-switch__control"
      :class="{ 'is-checked': modelValue }"
      :aria-pressed="modelValue"
      @click="emit('update:modelValue', !modelValue)"
    ></button>
  </div>
</template>
