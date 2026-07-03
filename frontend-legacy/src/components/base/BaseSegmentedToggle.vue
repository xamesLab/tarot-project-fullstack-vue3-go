<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    required: true,
  },
  options: {
    type: Array,
    required: true,
  },
  ariaLabel: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue'])

const activeIndex = computed(() => {
  const index = props.options.findIndex((option) => option.value === props.modelValue)
  return index === -1 ? 0 : index
})

const rootStyle = computed(() => ({
  gridTemplateColumns: `repeat(${props.options.length}, minmax(0, 1fr))`,
  '--toggle-columns': props.options.length,
  '--toggle-active-index': activeIndex.value,
}))

const thumbStyle = computed(() => ({
  width: `calc((100% - 8px) / ${props.options.length})`,
  transform: `translateX(calc(${activeIndex.value} * 100%))`,
}))
</script>

<template>
  <div class="base-segmented-toggle" :style="rootStyle" :aria-label="ariaLabel" role="group">
    <span class="base-segmented-toggle__thumb" :style="thumbStyle" aria-hidden="true"></span>
    <button
      v-for="option in options"
      :key="option.value"
      type="button"
      class="base-segmented-toggle__option"
      :class="{ 'is-active': option.value === modelValue }"
      :aria-pressed="option.value === modelValue"
      @click="emit('update:modelValue', option.value)"
    >
      {{ option.label }}
    </button>
  </div>
</template>
