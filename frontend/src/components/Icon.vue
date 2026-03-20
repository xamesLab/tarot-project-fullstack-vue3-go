<template>
  <svg
    v-if="svgContent"
    :class="className"
    :style="{ width: size + 'px', height: size + 'px' }"
    v-html="svgContent"
  ></svg>
</template>

<script setup lang="ts">
import { computed } from "vue";

import languageSvg from "../assets/icons/language.svg?raw";
import languageRuSvg from "../assets/icons/language-ru.svg?raw";
import sunSvg from "../assets/icons/sun.svg?raw";
import moonSvg from "../assets/icons/moon.svg?raw";

export type IconName = "language" | "language_ru" | "sun" | "moon";

const props = withDefaults(
  defineProps<{
    name: IconName;
    size?: number;
    className?: string;
  }>(),
  {
    size: 20,
    className: ""
  }
);

const svgContent = computed(() => {
  switch (props.name) {
    case "language":
      return languageSvg;
    case "language_ru":
      return languageRuSvg;
    case "sun":
      return sunSvg;
    case "moon":
      return moonSvg;
    default:
      return "";
  }
});
</script>

<style scoped>
svg {
  display: block;
}
</style>
