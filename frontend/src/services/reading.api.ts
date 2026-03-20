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

export function generateCardOfDay(token: string): Promise<TodayReading> {
  return request<TodayReading>("/api/reading/generate/card-of-day", {
    method: "POST",
    token
  });
}

export function generateSpreadA(token: string): Promise<TodayReading> {
  return request<TodayReading>("/api/reading/generate/spread-a", {
    method: "POST",
    token
  });
}

export function generateSpreadB(token: string): Promise<TodayReading> {
  return request<TodayReading>("/api/reading/generate/spread-b", {
    method: "POST",
    token
  });
}
