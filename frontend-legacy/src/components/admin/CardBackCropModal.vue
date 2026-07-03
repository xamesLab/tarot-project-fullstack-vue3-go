<script setup>
import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'

import BaseButton from '../base/BaseButton.vue'
import BaseModal from '../base/BaseModal.vue'

const CROP_WIDTH = 320
const CROP_HEIGHT = 480
const OUTPUT_WIDTH = 800
const OUTPUT_HEIGHT = 1200

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false,
  },
  imageUrl: {
    type: String,
    default: '',
  },
  fileName: {
    type: String,
    default: 'card-back.png',
  },
  title: {
    type: String,
    default: '',
  },
  copy: {
    type: String,
    default: '',
  },
  zoomLabel: {
    type: String,
    default: 'Zoom',
  },
  cancelLabel: {
    type: String,
    default: 'Cancel',
  },
  confirmLabel: {
    type: String,
    default: 'Use crop',
  },
})

const emit = defineEmits(['close', 'confirm'])

const imageElement = ref(null)
const isSaving = ref(false)
const dragState = reactive({
  active: false,
  startX: 0,
  startY: 0,
  startOffsetX: 0,
  startOffsetY: 0,
})

const imageState = reactive({
  naturalWidth: 0,
  naturalHeight: 0,
  minScale: 1,
  scale: 1,
  offsetX: 0,
  offsetY: 0,
})

const stageStyle = computed(() => ({
  width: `${CROP_WIDTH}px`,
  height: `${CROP_HEIGHT}px`,
}))

const imageStyle = computed(() => ({
  width: `${imageState.naturalWidth * imageState.scale}px`,
  height: `${imageState.naturalHeight * imageState.scale}px`,
  transform: `translate(${imageState.offsetX}px, ${imageState.offsetY}px)`,
}))

const zoomValue = computed({
  get() {
    return imageState.scale
  },
  set(value) {
    imageState.scale = clampScale(Number(value))
    clampOffsets()
  },
})

watch(
  () => props.isOpen,
  (isOpen) => {
    if (!isOpen) {
      resetDrag()
      isSaving.value = false
    }
  },
)

function clampScale(value) {
  const min = imageState.minScale || 1
  const max = min * 3
  return Math.min(Math.max(value, min), max)
}

function clampOffsets() {
  const displayedWidth = imageState.naturalWidth * imageState.scale
  const displayedHeight = imageState.naturalHeight * imageState.scale

  const minOffsetX = Math.min(0, CROP_WIDTH - displayedWidth)
  const minOffsetY = Math.min(0, CROP_HEIGHT - displayedHeight)
  const maxOffsetX = 0
  const maxOffsetY = 0

  imageState.offsetX = Math.min(Math.max(imageState.offsetX, minOffsetX), maxOffsetX)
  imageState.offsetY = Math.min(Math.max(imageState.offsetY, minOffsetY), maxOffsetY)
}

function centerImage() {
  imageState.offsetX = (CROP_WIDTH - imageState.naturalWidth * imageState.scale) / 2
  imageState.offsetY = (CROP_HEIGHT - imageState.naturalHeight * imageState.scale) / 2
  clampOffsets()
}

function handleImageLoad() {
  if (!imageElement.value) {
    return
  }

  imageState.naturalWidth = imageElement.value.naturalWidth
  imageState.naturalHeight = imageElement.value.naturalHeight

  const fitScale = Math.max(
    CROP_WIDTH / imageState.naturalWidth,
    CROP_HEIGHT / imageState.naturalHeight,
  )

  imageState.minScale = fitScale
  imageState.scale = fitScale
  centerImage()
}

function resetDrag() {
  dragState.active = false
}

function handlePointerDown(event) {
  dragState.active = true
  dragState.startX = event.clientX
  dragState.startY = event.clientY
  dragState.startOffsetX = imageState.offsetX
  dragState.startOffsetY = imageState.offsetY
}

function handlePointerMove(event) {
  if (!dragState.active) {
    return
  }

  imageState.offsetX = dragState.startOffsetX + (event.clientX - dragState.startX)
  imageState.offsetY = dragState.startOffsetY + (event.clientY - dragState.startY)
  clampOffsets()
}

function handlePointerUp() {
  resetDrag()
}

async function confirmCrop() {
  if (!imageElement.value) {
    return
  }

  isSaving.value = true

  const canvas = document.createElement('canvas')
  canvas.width = OUTPUT_WIDTH
  canvas.height = OUTPUT_HEIGHT
  const context = canvas.getContext('2d')
  if (!context) {
    isSaving.value = false
    return
  }

  const sourceX = -imageState.offsetX / imageState.scale
  const sourceY = -imageState.offsetY / imageState.scale
  const sourceWidth = CROP_WIDTH / imageState.scale
  const sourceHeight = CROP_HEIGHT / imageState.scale

  context.drawImage(
    imageElement.value,
    sourceX,
    sourceY,
    sourceWidth,
    sourceHeight,
    0,
    0,
    OUTPUT_WIDTH,
    OUTPUT_HEIGHT,
  )

  canvas.toBlob((blob) => {
    isSaving.value = false
    if (!blob) {
      return
    }

    const baseName = props.fileName.replace(/\.[^.]+$/, '') || 'card-back'
    const croppedFile = new File([blob], `${baseName}-cropped.png`, { type: 'image/png' })
    emit('confirm', croppedFile)
  }, 'image/png')
}

onBeforeUnmount(() => {
  resetDrag()
})
</script>

<template>
  <BaseModal :is-open="isOpen" :title="title" @close="emit('close')">
    <div class="crop-modal">
      <p v-if="copy" class="crop-modal__copy">{{ copy }}</p>

      <div
        class="crop-stage"
        :style="stageStyle"
        @pointermove="handlePointerMove"
        @pointerup="handlePointerUp"
        @pointerleave="handlePointerUp"
      >
        <img
          v-if="imageUrl"
          ref="imageElement"
          :src="imageUrl"
          :alt="fileName"
          class="crop-stage__image"
          :style="imageStyle"
          draggable="false"
          @load="handleImageLoad"
          @pointerdown.prevent="handlePointerDown"
        />
        <div class="crop-stage__frame" aria-hidden="true"></div>
      </div>

      <div class="crop-controls">
        <label class="crop-controls__label" for="card-back-crop-zoom">{{ zoomLabel }}</label>
        <input
          id="card-back-crop-zoom"
          v-model="zoomValue"
          class="crop-controls__range"
          type="range"
          :min="imageState.minScale || 1"
          :max="(imageState.minScale || 1) * 3"
          step="0.01"
        />
      </div>
    </div>

    <template #footer>
      <BaseButton variant="text" @click="emit('close')">{{ cancelLabel }}</BaseButton>
      <BaseButton :loading="isSaving" @click="confirmCrop">{{ confirmLabel }}</BaseButton>
    </template>
  </BaseModal>
</template>

<style scoped>
.crop-modal {
  display: grid;
  gap: 18px;
}

.crop-modal__copy {
  margin: 0;
  color: var(--ds-color-text-secondary);
}

.crop-stage {
  position: relative;
  margin: 0 auto;
  overflow: hidden;
  background:
    radial-gradient(circle at 20% 20%, rgba(201, 162, 86, 0.12), transparent 35%),
    linear-gradient(180deg, rgba(22, 19, 16, 0.96), rgba(12, 10, 8, 0.98));
  touch-action: none;
  user-select: none;
  cursor: grab;
}

.crop-stage:active {
  cursor: grabbing;
}

.crop-stage__image {
  position: absolute;
  top: 0;
  left: 0;
  max-width: none;
  transform-origin: top left;
}

.crop-stage__frame {
  position: absolute;
  inset: 0;
  border: 1px solid rgba(201, 162, 86, 0.8);
  box-shadow:
    inset 0 0 0 999px rgba(0, 0, 0, 0.18),
    0 0 0 1px rgba(201, 162, 86, 0.25);
  pointer-events: none;
}

.crop-stage__frame::before,
.crop-stage__frame::after {
  content: '';
  position: absolute;
  inset: 18px;
  border-top: 1px solid rgba(201, 162, 86, 0.25);
  border-bottom: 1px solid rgba(201, 162, 86, 0.25);
}

.crop-stage__frame::after {
  inset: 0 18px;
  border-top: 0;
  border-bottom: 0;
  border-left: 1px solid rgba(201, 162, 86, 0.25);
  border-right: 1px solid rgba(201, 162, 86, 0.25);
}

.crop-controls {
  display: grid;
  gap: 10px;
}

.crop-controls__label {
  color: var(--ds-color-text-secondary);
  font-size: var(--ds-text-label-sm);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.crop-controls__range {
  width: 100%;
}
</style>
