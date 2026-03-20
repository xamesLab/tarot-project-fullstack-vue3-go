<template>
  <div
    v-if="props.type === 'locale'"
    class="mode-toggle locale-toggle"
    role="tablist"
    aria-label="Language switcher"
  >
    <button
      v-for="option in localeOptions"
      :key="option.value"
      class="lang-segment"
      :class="{ active: props.modelValue === option.value }"
      type="button"
      role="tab"
      :aria-selected="props.modelValue === option.value"
      :aria-label="option.label"
      @click="select(option.value)"
    >
      <span class="segment-text">{{ option.text }}</span>
      <Icon :name="option.icon" :size="14" class-name="toggle-icon-fallback" />
    </button>
  </div>

  <button
    v-else
    class="mode-toggle theme-toggle-btn"
    type="button"
    :aria-label="themeButtonLabel"
    @click="toggleTheme"
  >
    <span class="theme-symbol">{{ themeSymbol }}</span>
    <Icon :name="themeIconName" :size="15" class-name="toggle-icon-fallback" />
  </button>
</template>

<script setup lang="ts">
import { computed } from "vue";

import Icon from "./Icon.vue";

export type ModeToggleType = "locale" | "theme";

interface LocaleOption {
  value: "en" | "ru";
  label: string;
  text: string;
  icon: "language" | "language_ru";
}

const props = defineProps<{
  type: ModeToggleType;
  modelValue: string;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const localeOptions: LocaleOption[] = [
  { value: "ru", label: "Russian", text: "RU", icon: "language_ru" },
  { value: "en", label: "English", text: "EN", icon: "language" }
];

const isDarkTheme = computed(() => props.modelValue === "dark");
const themeSymbol = computed(() => (isDarkTheme.value ? "\u2600" : "\u263D"));
const themeIconName = computed(() => (isDarkTheme.value ? "sun" : "moon"));
const themeButtonLabel = computed(() =>
  isDarkTheme.value ? "Switch to light theme" : "Switch to dark theme"
);

function select(value: string) {
  emit("update:modelValue", value);
}

function toggleTheme() {
  const nextTheme = props.modelValue === "dark" ? "light" : "dark";
  emit("update:modelValue", nextTheme);
}
</script>

<style scoped>
.mode-toggle {
  transition: background 0.4s, color 0.4s, border-color 0.3s;
}

.locale-toggle {
  display: inline-flex;
  border: 1px solid var(--border-strong);
  border-radius: 20px;
  overflow: hidden;
  background: transparent;
}

.lang-segment {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  min-width: 48px;
  padding: 6px 14px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.06em;
  transition: background 0.2s, color 0.2s, border-color 0.3s;
}

.lang-segment.active {
  background: var(--surface-elevated);
  color: var(--accent);
}

.theme-toggle-btn {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  border: 1px solid var(--border-strong);
  background: transparent;
  color: var(--text-muted);
  font-size: 15px;
  line-height: 1;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: background 0.4s, color 0.4s, border-color 0.3s;
}

.theme-toggle-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.theme-symbol {
  font-family: var(--font-display);
}

:deep(.toggle-icon-fallback) {
  display: none;
}
</style>
