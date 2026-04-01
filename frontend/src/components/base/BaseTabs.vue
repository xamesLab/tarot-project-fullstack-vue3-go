<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => [],
  },
  modelValue: {
    type: String,
    default: '',
  },
  variant: {
    type: String,
    default: 'underline',
  },
})

const emit = defineEmits(['update:modelValue'])

const activeValue = computed(() => props.modelValue || props.items[0]?.value || '')
const activeItem = computed(() => props.items.find((item) => item.value === activeValue.value) || props.items[0] || null)
</script>

<template>
  <div class="base-tabs" :class="`base-tabs--${variant}`">
    <div class="base-tabs__list" role="tablist">
      <button
        v-for="item in items"
        :key="item.value"
        type="button"
        class="base-tabs__trigger"
        :class="{ 'is-active': item.value === activeValue }"
        :aria-selected="item.value === activeValue"
        @click="emit('update:modelValue', item.value)"
      >
        {{ item.label }}
      </button>
    </div>

    <div class="base-tabs__panel">
      <slot name="panel" :active-item="activeItem">
        <div v-if="activeItem">
          <h3>{{ activeItem.label }}</h3>
          <p>{{ activeItem.content }}</p>
        </div>
      </slot>
    </div>
  </div>
</template>
