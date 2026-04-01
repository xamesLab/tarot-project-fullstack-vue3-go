<script setup>
defineProps({
  label: {
    type: String,
    default: '',
  },
  hint: {
    type: String,
    default: '',
  },
  error: {
    type: String,
    default: '',
  },
  forId: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <div class="form-field">
    <label v-if="label" class="form-field__label" :for="forId">{{ label }}</label>
    <slot />
    <p v-if="error" class="form-field__error">{{ error }}</p>
    <p v-else-if="hint" class="form-field__hint">{{ hint }}</p>
  </div>
</template>
