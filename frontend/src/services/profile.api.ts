import type { UserProfile } from "../types/auth";
import { request } from "./apiClient";

export interface AvailableDeck {
  id: string;
  name: string;
}

export interface SaveProfilePayload {
  displayName?: string;
  timezone?: string;
  preferredDeckId?: string;
}

export function getProfile(token: string): Promise<UserProfile> {
  return request<UserProfile>("/api/me/profile", { token });
}

export function saveProfile(token: string, payload: SaveProfilePayload): Promise<UserProfile> {
  return request<UserProfile>("/api/me/profile", {
    method: "PATCH",
    token,
    body: payload
  });
}

export function getAvailableDecks(token: string): Promise<AvailableDeck[]> {
  return request<AvailableDeck[]>("/api/decks/available", { token });
}
