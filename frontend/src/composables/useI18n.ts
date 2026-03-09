import { messages, type TranslationKey } from "../i18n/messages";
import { useUiPreferencesStore } from "../stores/ui-preferences.store";

export function useI18n() {
  const uiPreferencesStore = useUiPreferencesStore();

  function t(key: TranslationKey): string {
    const active = messages[uiPreferencesStore.locale];
    return active[key] ?? messages.ru[key] ?? messages.en[key] ?? key;
  }

  return {
    t
  };
}
