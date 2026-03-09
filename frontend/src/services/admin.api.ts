import type { DeckTheme } from "../types/deck-theme";
import { request } from "./apiClient";

export interface CreateThemePayload {
  name: string;
}

export interface RolePolicy {
  role: string;
  dailySpreadLimit: number;
  isUnlimited: boolean;
}

export interface UpsertRolePolicyPayload {
  role: string;
  dailySpreadLimit: number;
}

export interface AssignUserRolePayload {
  login: string;
  role: string;
}

export interface RoleAssignmentResult {
  login: string;
  role: string;
  dailySpreadLimit: number;
  isUnlimited: boolean;
}

export interface AdminUsersQuery {
  page?: number;
  pageSize?: number;
  role?: string;
  status?: "active" | "banned" | "freeze";
  sort?: "registeredAt:desc" | "registeredAt:asc";
  name?: string;
}

export interface AdminUserListItem {
  id: string;
  login: string;
  name: string;
  role: string;
  status: "active" | "banned" | "freeze";
  registeredAt: string;
}

export interface AdminUsersListResponse {
  items: AdminUserListItem[];
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface AdminDeck {
  id: string;
  name: string;
  status: "process" | "done";
  createdByUserId: string;
  createdByName: string;
  backImageUrl?: string;
  cards?: Record<string, string>;
  cardsUploadedCount: number;
  requiredCardsCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface AdminGeneralSettings {
  projectTitle: string;
  supportEmail: string;
  maintenanceMode: boolean;
  dayResetHour: number;
  maxUploadMb: number;
  updatedAt: string;
}

export interface PatchAdminGeneralSettingsPayload {
  projectTitle?: string;
  supportEmail?: string;
  maintenanceMode?: boolean;
  dayResetHour?: number;
  maxUploadMb?: number;
}

export function getAdminDeckThemes(token: string): Promise<DeckTheme[]> {
  return request<DeckTheme[]>("/api/admin/deck-themes", { token });
}

export function createDeckTheme(token: string, payload: CreateThemePayload): Promise<DeckTheme> {
  return request<DeckTheme>("/api/admin/deck-themes", {
    method: "POST",
    token,
    body: payload
  });
}

export function getRolePolicies(token: string): Promise<RolePolicy[]> {
  return request<RolePolicy[]>("/api/admin/roles", { token });
}

export function upsertRolePolicy(token: string, payload: UpsertRolePolicyPayload): Promise<RolePolicy> {
  return request<RolePolicy>("/api/admin/roles", {
    method: "POST",
    token,
    body: payload
  });
}

export function assignUserRole(token: string, payload: AssignUserRolePayload): Promise<RoleAssignmentResult> {
  return request<RoleAssignmentResult>("/api/admin/users/role", {
    method: "PATCH",
    token,
    body: payload
  });
}

export function getAdminUsers(token: string, query: AdminUsersQuery = {}): Promise<AdminUsersListResponse> {
  const search = new URLSearchParams();
  if (query.page !== undefined) search.set("page", String(query.page));
  if (query.pageSize !== undefined) search.set("pageSize", String(query.pageSize));
  if (query.role) search.set("role", query.role);
  if (query.status) search.set("status", query.status);
  if (query.sort) search.set("sort", query.sort);
  if (query.name) search.set("name", query.name);

  const suffix = search.toString();
  const path = suffix ? `/api/admin/users?${suffix}` : "/api/admin/users";
  return request<AdminUsersListResponse>(path, { token });
}

export function getAdminDecks(token: string): Promise<AdminDeck[]> {
  return request<AdminDeck[]>("/api/admin/decks", { token });
}

export function createAdminDeck(token: string, payload: { name: string }): Promise<AdminDeck> {
  return request<AdminDeck>("/api/admin/decks", {
    method: "POST",
    token,
    body: payload
  });
}

export function getAdminDeckById(token: string, deckId: string): Promise<AdminDeck> {
  return request<AdminDeck>(`/api/admin/decks/${encodeURIComponent(deckId)}`, { token });
}

export function uploadAdminDeckBack(token: string, deckId: string, file: File): Promise<AdminDeck> {
  const formData = new FormData();
  formData.append("file", file);
  return request<AdminDeck>(`/api/admin/decks/${encodeURIComponent(deckId)}/assets/back`, {
    method: "POST",
    token,
    body: formData
  });
}

export function uploadAdminDeckCard(token: string, deckId: string, cardKey: string, file: File): Promise<AdminDeck> {
  const formData = new FormData();
  formData.append("file", file);
  return request<AdminDeck>(
    `/api/admin/decks/${encodeURIComponent(deckId)}/assets/cards/${encodeURIComponent(cardKey)}`,
    {
      method: "POST",
      token,
      body: formData
    }
  );
}

export function getAdminSettings(token: string): Promise<AdminGeneralSettings> {
  return request<AdminGeneralSettings>("/api/admin/settings", { token });
}

export function patchAdminSettings(
  token: string,
  payload: PatchAdminGeneralSettingsPayload
): Promise<AdminGeneralSettings> {
  return request<AdminGeneralSettings>("/api/admin/settings", {
    method: "PATCH",
    token,
    body: payload
  });
}
