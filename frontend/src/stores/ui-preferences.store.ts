import { ref } from "vue";
import { defineStore } from "pinia";

import type { Locale, ThemeMode } from "../i18n/messages";

const LOCALE_KEY = "tarot_daily_locale";
const THEME_KEY = "tarot_daily_theme";
const DEFAULT_LOCALE: Locale = "ru";
const DEFAULT_THEME: ThemeMode = "light";

export const useUiPreferencesStore = defineStore("uiPreferences", () => {
  const locale = ref<Locale>(DEFAULT_LOCALE);
  const theme = ref<ThemeMode>(DEFAULT_THEME);

  function initialize() {
    const storedLocale = localStorage.getItem(LOCALE_KEY);
    const storedTheme = localStorage.getItem(THEME_KEY);

    if (isLocale(storedLocale)) {
      locale.value = storedLocale;
    } else {
      locale.value = DEFAULT_LOCALE;
      localStorage.setItem(LOCALE_KEY, DEFAULT_LOCALE);
    }

    if (isTheme(storedTheme)) {
      theme.value = storedTheme;
    } else {
      theme.value = DEFAULT_THEME;
      localStorage.setItem(THEME_KEY, DEFAULT_THEME);
    }

    applyTheme(theme.value);
  }

  function setLocale(nextLocale: Locale) {
    locale.value = nextLocale;
    localStorage.setItem(LOCALE_KEY, nextLocale);
  }

  function setTheme(nextTheme: ThemeMode) {
    theme.value = nextTheme;
    localStorage.setItem(THEME_KEY, nextTheme);
    applyTheme(nextTheme);
  }

  return {
    locale,
    theme,
    initialize,
    setLocale,
    setTheme
  };
});

function applyTheme(theme: ThemeMode) {
  document.documentElement.setAttribute("data-theme", theme);
}

function isLocale(value: string | null): value is Locale {
  return value === "ru" || value === "en";
}

function isTheme(value: string | null): value is ThemeMode {
  return value === "light" || value === "dark";
}
