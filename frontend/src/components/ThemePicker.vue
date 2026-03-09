<template>
  <div class="field">
    <label for="deck-theme">{{ t("settingsDeckTheme") }}</label>
    <select id="deck-theme" :value="modelValue" @change="onChange">
      <option value="">{{ t("settingsDeckThemeDefault") }}</option>
      <option v-for="theme in themes" :key="theme.id" :value="theme.id">
        {{ theme.name }}
      </option>
    </select>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "../composables/useI18n";
import type { DeckTheme } from "../types/deck-theme";

const props = defineProps<{
  modelValue: string;
  themes: DeckTheme[];
}>();

const emit = defineEmits<{
  (event: "update:modelValue", value: string): void;
}>();

const { t } = useI18n();

function onChange(event: Event) {
  const value = (event.target as HTMLSelectElement).value;
  emit("update:modelValue", value);
}
</script>
