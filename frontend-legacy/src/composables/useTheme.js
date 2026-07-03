import { computed } from 'vue'

import { useAppStore } from '../stores/app'

export function useTheme() {
  const appStore = useAppStore()

  return {
    theme: computed(() => appStore.theme),
    setTheme: appStore.setTheme,
    toggleTheme: appStore.toggleTheme,
  }
}
