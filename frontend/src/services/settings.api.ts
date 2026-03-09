import type { UserProfile } from "../types/auth";
import type { DeckTheme } from "../types/deck-theme";
import { request } from "./apiClient";

export interface SaveSettingsPayload {
  displayName?: string;
  timezone?: string;
  preferredDeckId?: string;
  deckThemeId?: string;
}

export function getDeckThemes(token: string): Promise<DeckTheme[]> {
  return request<DeckTheme[]>("/api/deck-themes", { token });
}

export function saveSettings(token: string, payload: SaveSettingsPayload): Promise<UserProfile> {
  return request<UserProfile>("/api/me/settings", {
    method: "PATCH",
    token,
    body: payload
  });
}
