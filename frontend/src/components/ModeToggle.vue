<template>
  <div class="mode-toggle" :class="modeToggleClass">
    <button
      v-for="option in options"
      :key="option.value"
      class="toggle-btn"
      :class="{ 'is-active': modelValue === option.value }"
      type="button"
      :aria-label="option.label"
      :aria-pressed="modelValue === option.value"
      @click="select(option.value)"
    >
      <Icon :name="option.icon" :size="20" />
    </button>
  </div>
</template>

<script setup lang="ts">
import Icon from "./Icon.vue";

export type ModeToggleType = "locale" | "theme";

interface ModeToggleOption {
  value: string;
  label: string;
  icon: "language" | "language_ru" | "sun" | "moon";
}

const props = defineProps<{
  type: ModeToggleType;
  modelValue: string;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const localeOptions: ModeToggleOption[] = [
  { value: "ru", label: "Русский", icon: "language_ru" },
  { value: "en", label: "English", icon: "language" }
];

const themeOptions: ModeToggleOption[] = [
  { value: "light", label: "Light", icon: "sun" },
  { value: "dark", label: "Dark", icon: "moon" }
];

const options = props.type === "locale" ? localeOptions : themeOptions;

const modeToggleClass = props.type === "locale" ? "locale-toggle" : "theme-toggle";

function select(value: string) {
  emit("update:modelValue", value);
}
</script>

<style scoped>
.mode-toggle {
  display: flex;
  background: transparent;
  border-radius: 6px;
  overflow: hidden;
}

.toggle-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--text-muted, #666);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
}

.toggle-btn:hover {
  color: var(--text, #333);
  background: rgba(128, 128, 128, 0.1);
}

.toggle-btn.is-active {
  color: var(--primary, #8b5cf6);
}

.toggle-btn.is-active::after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 2px;
  background: var(--primary, #8b5cf6);
  border-radius: 2px 2px 0 0;
  animation: slideIn 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes slideIn {
  from {
    width: 0;
    opacity: 0;
  }
  to {
    width: 20px;
    opacity: 1;
  }
}

/* Dark mode adjustments */
.dark .toggle-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.dark .toggle-btn.is-active {
  color: var(--primary-light, #a78bfa);
}

.dark .toggle-btn.is-active::after {
  background: var(--primary-light, #a78bfa);
}
</style>
