export type UserRole = "user" | "admin";
export type UserStatus = "active" | "banned" | "freeze";

export interface UserProfile {
  id: string;
  login: string;
  role: UserRole;
  status?: UserStatus;
  displayName?: string;
  timezone: string;
  preferredDeckId?: string;
  deckThemeId?: string;
}

export interface AuthResponse {
  token: string;
  user: UserProfile;
}
