import type { AuthResponse, UserProfile } from "../types/auth";
import { request } from "./apiClient";

interface AuthPayload {
  login: string;
  password: string;
}

export function register(payload: AuthPayload): Promise<AuthResponse> {
  return request<AuthResponse>("/api/auth/register", {
    method: "POST",
    body: payload
  });
}

export function login(payload: AuthPayload): Promise<AuthResponse> {
  return request<AuthResponse>("/api/auth/login", {
    method: "POST",
    body: payload
  });
}

export function me(token: string): Promise<UserProfile> {
  return request<UserProfile>("/api/me", { token });
}
