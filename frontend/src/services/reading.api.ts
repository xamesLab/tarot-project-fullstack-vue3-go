import type { TodayReading } from "../types/reading";
import { request } from "./apiClient";

export function getTodayReading(token: string): Promise<TodayReading> {
  return request<TodayReading>("/api/reading/today", { token });
}

export function generateSpreads(token: string): Promise<TodayReading> {
  return request<TodayReading>("/api/reading/generate", {
    method: "POST",
    token
  });
}
