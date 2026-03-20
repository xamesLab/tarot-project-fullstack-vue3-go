<template>
  <div class="tabs-container">
    <div class="tabs-list" role="tablist" :aria-label="ariaLabel">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :id="`tab-${tab.id}`"
        class="tab-button"
        :class="{ 'is-active': modelValue === tab.id }"
        role="tab"
        :aria-selected="modelValue === tab.id"
        :aria-controls="`panel-${tab.id}`"
        :tabindex="modelValue === tab.id ? 0 : -1"
        @click="selectTab(tab.id)"
        @keydown="handleKeydown($event, tab.id)"
      >
        {{ tab.label }}
      </button>
    </div>
    <div class="tabs-content">
      <slot :active-tab="modelValue" />
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Tab {
  id: string;
  label: string;
}

defineProps<{
  tabs: Tab[];
  modelValue: string;
  ariaLabel?: string;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

function selectTab(tabId: string) {
  emit("update:modelValue", tabId);
}

function handleKeydown(event: KeyboardEvent, currentTabId: string) {
  const tabIds = defineProps<{ tabs: Tab[] }>().tabs.map((t) => t.id);
  const currentIndex = tabIds.indexOf(currentTabId);

  if (event.key === "ArrowRight") {
    event.preventDefault();
    const nextIndex = (currentIndex + 1) % tabIds.length;
    selectTab(tabIds[nextIndex]);
    focusTab(tabIds[nextIndex]);
  } else if (event.key === "ArrowLeft") {
    event.preventDefault();
    const prevIndex = (currentIndex - 1 + tabIds.length) % tabIds.length;
    selectTab(tabIds[prevIndex]);
    focusTab(tabIds[prevIndex]);
  }
}

function focusTab(tabId: string) {
  const tabElement = document.getElementById(`tab-${tabId}`);
  if (tabElement) {
    tabElement.focus();
  }
}
</script>

<style scoped>
.tabs-container {
  width: 100%;
}

.tabs-list {
  display: flex;
  gap: 4px;
  border-bottom: 2px solid var(--border, rgba(0, 0, 0, 0.1));
  padding-bottom: 0;
  overflow-x: auto;
}

.tab-button {
  position: relative;
  padding: 12px 20px;
  background: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  cursor: pointer;
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--text-muted, #666);
  transition: color 0.2s ease, border-color 0.2s ease, background 0.2s ease;
  white-space: nowrap;
}

.tab-button:hover {
  color: var(--text, #333);
  background: var(--surface-soft, rgba(0, 0, 0, 0.02));
}

.tab-button.is-active {
  color: var(--primary, #8b5cf6);
  border-bottom-color: var(--primary, #8b5cf6);
}

.tab-button:focus-visible {
  outline: 2px solid var(--primary, #8b5cf6);
  outline-offset: -2px;
  border-radius: 4px 4px 0 0;
}

.tabs-content {
  padding: 24px 0 0;
}

@media (max-width: 640px) {
  .tab-button {
    padding: 10px 14px;
    font-size: 0.88rem;
  }

  .tabs-list {
    gap: 2px;
  }
}
</style>
