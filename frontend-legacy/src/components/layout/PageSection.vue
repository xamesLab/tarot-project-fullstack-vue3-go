<script setup>
import { useSlots } from 'vue'

import SectionHeader from './SectionHeader.vue'

const slots = useSlots()

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
    <SectionHeader v-if="eyebrow || title || copy || slots.meta" :eyebrow="eyebrow" :title="title" :copy="copy">
      <template v-if="slots.meta" #meta>
        <slot name="meta" />
      </template>
    </SectionHeader>
    <slot />
  </section>
</template>
