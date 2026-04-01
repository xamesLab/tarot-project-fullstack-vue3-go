import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useAppStore } from '../stores/app'

export function useLocale() {
  const appStore = useAppStore()
  const { locale } = useI18n()

  watch(
    () => appStore.language,
    (value) => {
      locale.value = value
      appStore.applyLanguage(value)
    },
    { immediate: true },
  )

  return {
    language: computed(() => appStore.language),
    setLanguage: appStore.setLanguage,
    toggleLanguage: appStore.toggleLanguage,
  }
}
