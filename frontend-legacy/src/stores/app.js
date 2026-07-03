import { defineStore } from 'pinia'

const THEME_STORAGE_KEY = 'tarot_daily_theme'
const LANGUAGE_STORAGE_KEY = 'tarot_daily_language'

function resolveInitialLanguage() {
  if (typeof window === 'undefined') {
    return 'en'
  }

  const stored = window.localStorage.getItem(LANGUAGE_STORAGE_KEY)
  if (stored === 'en' || stored === 'ru') {
    return stored
  }

  return navigator.language.toLowerCase().startsWith('ru') ? 'ru' : 'en'
}

function resolveInitialTheme() {
  if (typeof window === 'undefined') {
    return 'light'
  }

  const stored = window.localStorage.getItem(THEME_STORAGE_KEY)
  if (stored === 'light' || stored === 'dark') {
    return stored
  }

  return 'light'
}

export const useAppStore = defineStore('app', {
  state: () => ({
    theme: resolveInitialTheme(),
    language: resolveInitialLanguage(),
    initialized: false,
  }),

  actions: {
    init() {
      this.applyTheme(this.theme)
      this.applyLanguage(this.language)
      this.initialized = true
    },

    applyTheme(theme) {
      if (typeof document !== 'undefined') {
        document.documentElement.setAttribute('data-theme', theme)
      }

      if (typeof window !== 'undefined') {
        window.localStorage.setItem(THEME_STORAGE_KEY, theme)
      }
    },

    applyLanguage(language) {
      if (typeof document !== 'undefined') {
        document.documentElement.setAttribute('lang', language)
      }

      if (typeof window !== 'undefined') {
        window.localStorage.setItem(LANGUAGE_STORAGE_KEY, language)
      }
    },

    setTheme(theme) {
      this.theme = theme
      this.applyTheme(theme)
    },

    toggleTheme() {
      this.setTheme(this.theme === 'light' ? 'dark' : 'light')
    },

    setLanguage(language) {
      this.language = language
      this.applyLanguage(language)
    },

    toggleLanguage() {
      this.setLanguage(this.language === 'en' ? 'ru' : 'en')
    },

    syncUserPreferences(user) {
      if (!user) {
        return
      }

      if (user.theme === 'light' || user.theme === 'dark') {
        this.setTheme(user.theme)
      }

      if (user.language === 'en' || user.language === 'ru') {
        this.setLanguage(user.language)
      }
    },
  },
})
