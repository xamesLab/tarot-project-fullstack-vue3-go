<script setup>
import { onBeforeUnmount, onMounted } from 'vue'

import BaseSurface from './BaseSurface.vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['close'])

function onKeydown(event) {
  if (event.key === 'Escape' && props.isOpen) {
    emit('close')
  }
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <Teleport to="body">
    <div v-if="isOpen" class="base-modal">
      <button type="button" class="base-modal__overlay" aria-label="Close modal" @click="emit('close')"></button>
      <div class="base-modal__dialog">
        <BaseSurface tone="elevated" padding="lg" tag="div">
          <div class="base-modal__head">
            <h3 class="base-modal__title">{{ title }}</h3>
            <button type="button" class="base-modal__close" @click="emit('close')">×</button>
          </div>
          <div class="base-modal__body">
            <slot />
          </div>
          <div v-if="$slots.footer" class="base-modal__footer">
            <slot name="footer" />
          </div>
        </BaseSurface>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.base-modal {
  position: fixed;
  inset: 0;
  z-index: 120;
  display: grid;
  place-items: center;
  padding: 24px;
}

.base-modal__overlay {
  position: absolute;
  inset: 0;
  background: var(--ds-color-overlay);
}

.base-modal__dialog {
  position: relative;
  z-index: 1;
  width: min(560px, 100%);
}

.base-modal__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}

.base-modal__title {
  font-size: var(--ds-text-heading-2);
  color: var(--ds-color-accent-primary);
}

.base-modal__close {
  font-size: 24px;
  color: var(--ds-color-accent-primary);
  cursor: pointer;
}

.base-modal__body {
  color: var(--ds-color-text-secondary);
}

.base-modal__footer {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
