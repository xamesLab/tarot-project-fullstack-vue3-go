import { ref } from "vue";
import { defineStore } from "pinia";

import type { UserProfile } from "../types/auth";
import type { DeckTheme } from "../types/deck-theme";
import * as settingsApi from "../services/settings.api";

export const useSettingsStore = defineStore("settings", () => {
  const timezone = ref("UTC+00:00");
  const deckThemeId = ref("");
  const themes = ref<DeckTheme[]>([]);
  const loading = ref(false);
  const saving = ref(false);
  const error = ref("");
  const success = ref("");

  function applyUserSettings(user: UserProfile | null) {
    timezone.value = user?.timezone || "UTC+00:00";
    deckThemeId.value = user?.deckThemeId || "";
  }

  async function loadThemes(token: string) {
    loading.value = true;
    error.value = "";
    success.value = "";
    try {
      themes.value = await settingsApi.getDeckThemes(token);
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to load themes";
    } finally {
      loading.value = false;
    }
  }

  async function save(token: string) {
    saving.value = true;
    error.value = "";
    success.value = "";
    try {
      const user = await settingsApi.saveSettings(token, {
        timezone: timezone.value,
        deckThemeId: deckThemeId.value
      });
      applyUserSettings(user);
      success.value = "Settings saved";
      return user;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to save settings";
      throw err;
    } finally {
      saving.value = false;
    }
  }

  return {
    timezone,
    deckThemeId,
    themes,
    loading,
    saving,
    error,
    success,
    applyUserSettings,
    loadThemes,
    save
  };
});
