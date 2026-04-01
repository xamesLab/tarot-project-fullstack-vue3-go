<script setup>
defineProps({
  eyebrow: {
    type: String,
    default: '',
  },
  title: {
    type: String,
    default: '',
  },
  copy: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <section class="page-stack">
    <header v-if="eyebrow || title || copy" class="section-head">
      <p v-if="eyebrow" class="text-eyebrow">{{ eyebrow }}</p>
      <h2 v-if="title" class="section-head__title">{{ title }}</h2>
      <p v-if="copy" class="section-head__copy">{{ copy }}</p>
    </header>
    <slot />
  </section>
</template>
