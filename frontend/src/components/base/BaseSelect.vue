<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  id: {
    type: String,
    default: '',
  },
  options: {
    type: Array,
    default: () => [],
  },
  placeholder: {
    type: String,
    default: '',
  },
  error: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const classes = computed(() => ['base-select', { 'is-error': props.error }])
</script>

<template>
  <div class="base-select-wrap">
    <select
      :id="id"
      :value="modelValue"
      :class="classes"
      :disabled="disabled"
      @change="emit('update:modelValue', $event.target.value)"
    >
      <option v-if="placeholder" disabled value="">
        {{ placeholder }}
      </option>
      <option v-for="option in options" :key="option.value" :value="option.value">
        {{ option.label }}
      </option>
    </select>
  </div>
</template>
