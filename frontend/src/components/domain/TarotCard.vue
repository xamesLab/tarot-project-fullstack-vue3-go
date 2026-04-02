<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { resolveApiAssetUrl } from '../../services/api'

const props = defineProps({
  card: {
    type: Object,
    required: true,
  },
  positionLabel: {
    type: String,
    default: '',
  },
  compact: {
    type: Boolean,
    default: false,
  },
  size: {
    type: String,
    default: 'lg',
    validator: (value) => ['lg', 'md', 'sm'].includes(value),
  },
})

const { locale } = useI18n()
const imageLoadFailed = ref(false)
const backImageLoadFailed = ref(false)

const title = computed(() => (locale.value === 'ru' ? props.card.name_ru : props.card.name) || props.card.name)
const interpretation = computed(() => (
  locale.value === 'ru' ? props.card.interpretation_ru : props.card.interpretation
) || props.card.interpretation)
const tags = computed(() => ((locale.value === 'ru' ? props.card.tags_ru : props.card.tags) || '')
  .split(',')
  .map((tag) => tag.trim())
  .filter(Boolean))
const ornamentSymbol = computed(() => {
  const symbols = ['✦', '☽', '☉', '⋆']
  return symbols[props.card.position] || '✦'
})
const imageUrl = computed(() => resolveApiAssetUrl(props.card.image_url))
const backImageUrl = computed(() => resolveApiAssetUrl(props.card.back_image_url))

watch(
  () => [props.card?.image_url, props.card?.back_image_url],
  () => {
    imageLoadFailed.value = false
    backImageLoadFailed.value = false
  },
  { immediate: true },
)
</script>

<template>
  <article
    class="tarot-card"
    :class="[
      `tarot-card--${size}`,
      { 'tarot-card--compact': compact },
    ]"
  >
    <div class="tarot-card__body">
      <p v-if="positionLabel" class="tarot-card__position">{{ positionLabel }}</p>
      <h3 class="tarot-card__title">{{ title }}</h3>

      <div class="tarot-card__back">
        <img
          v-if="imageUrl && !imageLoadFailed"
          :src="imageUrl"
          :alt="title"
          class="tarot-card__image"
          loading="lazy"
          @error="imageLoadFailed = true"
        />
        <img
          v-else-if="backImageUrl && !backImageLoadFailed"
          :src="backImageUrl"
          :alt="`${title} card back`"
          class="tarot-card__image tarot-card__image--back"
          loading="lazy"
          @error="backImageLoadFailed = true"
        />
        <div v-else class="tarot-card__ornament">
          <span class="tarot-card__symbol">{{ ornamentSymbol }}</span>
        </div>
      </div>
      <p class="tarot-card__copy">{{ interpretation }}</p>
      <div v-if="tags.length" class="tarot-card__tags">
        <span v-for="tag in tags" :key="tag" class="tarot-card__tag">{{ tag }}</span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.tarot-card {
  --card-frame-max: 320px;
  --card-body-padding: 30px;
  --card-title-size: var(--ds-text-heading-2);
  --card-position-margin: 12px;
  --card-name-margin: 16px;
  --card-copy-size: var(--ds-text-body-md);
  --card-copy-line-height: 1.7;
  --card-copy-margin: 20px;
  --card-tag-gap: 10px;
  --card-back-margin: 20px;
  --card-symbol-size: 48px;
  --card-tag-padding-y: 6px;
  --card-tag-padding-x: 14px;
  --card-tag-font-size: var(--ds-text-caption);
  --card-tag-letter-spacing: 0.05em;
  display: block;
  height: 100%;
  width: min(100%, var(--card-frame-max));
  background: transparent;
  justify-self: center;
}

.tarot-card--md {
  --card-frame-max: 280px;
  --card-body-padding: 24px;
  --card-title-size: 21px;
  --card-position-margin: 10px;
  --card-name-margin: 14px;
  --card-copy-size: 15px;
  --card-copy-line-height: 1.65;
  --card-copy-margin: 18px;
  --card-tag-gap: 8px;
  --card-back-margin: 18px;
  --card-symbol-size: 42px;
  --card-tag-padding-y: 5px;
  --card-tag-padding-x: 12px;
  --card-tag-font-size: 9px;
  --card-tag-letter-spacing: 0.04em;
}

.tarot-card--sm {
  --card-frame-max: 220px;
  --card-body-padding: 18px;
  --card-title-size: 18px;
  --card-position-margin: 8px;
  --card-name-margin: 12px;
  --card-copy-size: 14px;
  --card-copy-line-height: 1.55;
  --card-copy-margin: 16px;
  --card-tag-gap: 6px;
  --card-back-margin: 16px;
  --card-symbol-size: 34px;
  --card-tag-padding-y: 4px;
  --card-tag-padding-x: 10px;
  --card-tag-font-size: 8px;
  --card-tag-letter-spacing: 0.03em;
}

.tarot-card__body {
  padding: var(--card-body-padding);
}

.tarot-card__back {
  width: 100%;
  aspect-ratio: 2 / 3;
  display: flex;
  align-items: center;
  justify-content: center;
  background:
    radial-gradient(circle at top, rgba(212, 175, 55, 0.1), transparent 58%),
    var(--ds-color-bg-elevated);
  border: 2px solid var(--ds-color-accent-primary);
  margin-bottom: var(--card-back-margin);
  transition:
    transform var(--ds-motion-base),
    box-shadow var(--ds-motion-base),
    border-color var(--ds-motion-base);
  overflow: hidden;
}

.tarot-card__back:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.6);
  border-color: color-mix(in srgb, var(--ds-color-accent-primary) 65%, white);
}

.tarot-card__ornament {
  position: relative;
  width: 80%;
  height: 80%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: var(--ds-border-width) solid var(--ds-color-accent-primary);
}

.tarot-card__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.tarot-card__image--back {
  opacity: 0.98;
}

.tarot-card__ornament::before,
.tarot-card__ornament::after {
  content: '';
  position: absolute;
  background: var(--ds-color-accent-primary);
}

.tarot-card__ornament::before {
  top: 50%;
  left: 20%;
  width: 60%;
  height: 1px;
}

.tarot-card__ornament::after {
  top: 20%;
  left: 50%;
  width: 1px;
  height: 60%;
}

.tarot-card__symbol {
  position: relative;
  z-index: 1;
  font-size: var(--card-symbol-size);
  color: var(--ds-color-accent-primary);
}

.tarot-card__position {
  font-family: var(--ds-font-ui);
  font-size: var(--ds-text-caption);
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--ds-color-text-muted);
  text-align: center;
  margin-bottom: var(--card-position-margin);
}

.tarot-card__title {
  font-size: var(--card-title-size);
  font-weight: 400;
  letter-spacing: 0.05em;
  color: var(--ds-color-text-primary);
  text-align: center;
  line-height: 1.2;
  margin-bottom: var(--card-name-margin);
}

.tarot-card__copy {
  font-size: var(--card-copy-size);
  line-height: var(--card-copy-line-height);
  color: var(--ds-color-text-secondary);
  text-align: center;
  margin-bottom: var(--card-copy-margin);
}

.tarot-card__tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--card-tag-gap);
  justify-content: center;
}

.tarot-card__tag {
  padding: var(--card-tag-padding-y) var(--card-tag-padding-x);
  border: var(--ds-border-width) solid var(--ds-color-border-accent);
  background: rgba(212, 175, 55, 0.1);
  font-family: var(--ds-font-ui);
  font-size: var(--card-tag-font-size);
  line-height: 1;
  letter-spacing: var(--card-tag-letter-spacing);
  color: var(--ds-color-text-secondary);
}

.tarot-card--compact {
  min-height: 100%;
}

.tarot-card--compact .tarot-card__title {
  font-size: var(--ds-text-heading-3);
}

.tarot-card--compact .tarot-card__copy {
  font-size: var(--ds-text-body-sm);
}
</style>
