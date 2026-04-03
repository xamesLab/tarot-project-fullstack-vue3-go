<script setup>
import { computed, getCurrentInstance, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

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
const { t } = useI18n()
const instance = getCurrentInstance()
const dialogRef = ref(null)
const closeButtonRef = ref(null)
const lastFocusedElement = ref(null)
const modalTitleId = computed(() => `base-modal-title-${instance?.uid || 'dialog'}`)

function getFocusableElements() {
  if (!dialogRef.value) {
    return []
  }

  return Array.from(
    dialogRef.value.querySelectorAll(
      'button:not([disabled]), [href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
    ),
  )
}

function onKeydown(event) {
  if (!props.isOpen) {
    return
  }

  if (event.key === 'Escape') {
    emit('close')
    return
  }

  if (event.key === 'Tab') {
    const focusableElements = getFocusableElements()
    if (!focusableElements.length) {
      return
    }

    const firstElement = focusableElements[0]
    const lastElement = focusableElements[focusableElements.length - 1]

    if (event.shiftKey && document.activeElement === firstElement) {
      event.preventDefault()
      lastElement.focus()
    } else if (!event.shiftKey && document.activeElement === lastElement) {
      event.preventDefault()
      firstElement.focus()
    }
  }
}

function applyOpenState() {
  if (typeof document === 'undefined') {
    return
  }

  document.body.style.overflow = props.isOpen ? 'hidden' : ''
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  if (typeof document !== 'undefined') {
    document.body.style.overflow = ''
  }
  window.removeEventListener('keydown', onKeydown)
})

watch(
  () => props.isOpen,
  async (isOpen) => {
    if (typeof document === 'undefined') {
      return
    }

    if (isOpen) {
      lastFocusedElement.value = document.activeElement
      applyOpenState()
      await nextTick()
      closeButtonRef.value?.focus()
      return
    }

    applyOpenState()
    if (lastFocusedElement.value?.focus) {
      lastFocusedElement.value.focus()
    }
  },
  { immediate: true },
)
</script>

<template>
  <Teleport to="body">
    <div v-if="isOpen" class="base-modal">
      <button
        type="button"
        class="base-modal__overlay"
        :aria-label="t('controls.close')"
        @click="emit('close')"
      ></button>
      <div
        ref="dialogRef"
        class="base-modal__dialog"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="title ? modalTitleId : undefined"
      >
        <BaseSurface tone="elevated" padding="lg" tag="div">
          <div class="base-modal__head">
            <h3 :id="modalTitleId" class="base-modal__title">{{ title }}</h3>
            <button
              ref="closeButtonRef"
              type="button"
              class="base-modal__close"
              :aria-label="t('controls.close')"
              :title="t('controls.close')"
              @click="emit('close')"
            >
              ×
            </button>
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
