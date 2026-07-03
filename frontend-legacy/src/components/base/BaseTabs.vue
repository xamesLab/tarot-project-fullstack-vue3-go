<script setup>
import { computed, getCurrentInstance } from 'vue'

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
const instance = getCurrentInstance()

const activeValue = computed(() => props.modelValue || props.items[0]?.value || '')
const activeItem = computed(() => props.items.find((item) => item.value === activeValue.value) || props.items[0] || null)

function getTabId(value) {
  return `base-tabs-${instance?.uid || 'tabs'}-tab-${value}`
}

function getPanelId(value) {
  return `base-tabs-${instance?.uid || 'tabs'}-panel-${value}`
}

function activateByIndex(index) {
  const item = props.items[index]
  if (!item) {
    return
  }

  emit('update:modelValue', item.value)
}

function handleKeydown(event, index) {
  if (!props.items.length) {
    return
  }

  switch (event.key) {
    case 'ArrowRight':
    case 'ArrowDown':
      event.preventDefault()
      activateByIndex((index + 1) % props.items.length)
      break
    case 'ArrowLeft':
    case 'ArrowUp':
      event.preventDefault()
      activateByIndex((index - 1 + props.items.length) % props.items.length)
      break
    case 'Home':
      event.preventDefault()
      activateByIndex(0)
      break
    case 'End':
      event.preventDefault()
      activateByIndex(props.items.length - 1)
      break
    default:
      break
  }
}
</script>

<template>
  <div class="base-tabs" :class="`base-tabs--${variant}`">
    <div class="base-tabs__list" role="tablist">
      <button
        v-for="(item, index) in items"
        :key="item.value"
        type="button"
        class="base-tabs__trigger"
        :class="{ 'is-active': item.value === activeValue }"
        role="tab"
        :id="getTabId(item.value)"
        :aria-selected="item.value === activeValue"
        :aria-controls="getPanelId(item.value)"
        :tabindex="item.value === activeValue ? 0 : -1"
        @click="emit('update:modelValue', item.value)"
        @keydown="handleKeydown($event, index)"
      >
        {{ item.label }}
      </button>
    </div>

    <div
      class="base-tabs__panel"
      role="tabpanel"
      :id="getPanelId(activeItem?.value || 'panel')"
      :aria-labelledby="activeItem ? getTabId(activeItem.value) : undefined"
    >
      <slot name="panel" :active-item="activeItem">
        <div v-if="activeItem">
          <h3>{{ activeItem.label }}</h3>
          <p>{{ activeItem.content }}</p>
        </div>
      </slot>
    </div>
  </div>
</template>
