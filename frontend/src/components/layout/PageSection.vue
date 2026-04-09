<script setup>
import SectionHeader from './SectionHeader.vue'

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
    <SectionHeader v-if="eyebrow || title || copy" :eyebrow="eyebrow" :title="title" :copy="copy" />
    <slot />
  </section>
</template>
