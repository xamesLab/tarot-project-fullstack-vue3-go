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
  required: {
    type: Boolean,
    default: false,
  },
})
</script>

<template>
  <div class="form-field">
    <label v-if="label" class="form-field__label" :for="forId">
      {{ label }}
      <span v-if="required" class="form-field__required" aria-hidden="true">*</span>
    </label>
    <slot />
    <p v-if="error" class="form-field__error">{{ error }}</p>
    <p v-else-if="hint" class="form-field__hint">{{ hint }}</p>
  </div>
</template>

<style scoped>
.form-field__required {
  color: var(--ds-color-accent-primary);
}
</style>
