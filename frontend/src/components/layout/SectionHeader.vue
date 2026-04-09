<script setup>
import { computed, useSlots } from 'vue'

const props = defineProps({
  tag: {
    type: String,
    default: 'header',
  },
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
  titleTag: {
    type: String,
    default: 'h2',
  },
  size: {
    type: String,
    default: 'md',
  },
  align: {
    type: String,
    default: 'start',
  },
})

const slots = useSlots()

const hasMeta = computed(() => Boolean(slots.meta))
const classes = computed(() => [
  'section-header',
  `section-header--${props.size}`,
  `section-header--${props.align}`,
  { 'section-header--with-meta': hasMeta.value },
])
</script>

<template>
  <component :is="tag" :class="classes">
    <div class="section-header__content">
      <p v-if="eyebrow" class="text-eyebrow">{{ eyebrow }}</p>
      <component :is="titleTag" v-if="title" class="section-header__title">{{ title }}</component>
      <p v-if="copy" class="section-header__copy">{{ copy }}</p>
    </div>

    <div v-if="hasMeta" class="section-header__meta">
      <slot name="meta" />
    </div>
  </component>
</template>

<style scoped>
.section-header {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--ds-space-4);
}

.section-header__content {
  display: grid;
  gap: var(--ds-space-3);
  max-width: 70ch;
}

.section-header__title {
  font-size: clamp(28px, 3vw, 42px);
  color: var(--ds-color-accent-primary);
}

.section-header__copy {
  color: var(--ds-color-text-secondary);
}

.section-header__meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--ds-space-3);
  align-items: center;
  justify-content: flex-end;
}

.section-header--sm .section-header__content {
  gap: var(--ds-space-2);
}

.section-header--sm .section-header__title {
  font-size: var(--ds-text-heading-2);
}

.section-header--lg .section-header__title {
  font-size: clamp(36px, 4vw, 54px);
}

.section-header--center {
  justify-content: center;
  text-align: center;
}

.section-header--center .section-header__meta {
  justify-content: center;
}

@media (max-width: 720px) {
  .section-header,
  .section-header__content {
    gap: var(--ds-space-3);
  }

  .section-header__title {
    font-size: clamp(24px, 8vw, 32px);
  }

  .section-header--sm .section-header__title {
    font-size: var(--ds-text-heading-3);
  }

  .section-header--lg .section-header__title {
    font-size: clamp(30px, 10vw, 40px);
  }

  .section-header__copy {
    max-width: none;
  }
}
</style>
