import { ref } from "vue";
import { defineStore } from "pinia";

import type { TodayReading } from "../types/reading";
import { generateSpreads, getTodayReading } from "../services/reading.api";
import { ApiError } from "../services/apiClient";

export const useReadingStore = defineStore("reading", () => {
  const today = ref<TodayReading | null>(null);
  const loading = ref(false);
  const generating = ref(false);
  const error = ref("");

  async function loadToday(token: string) {
    loading.value = true;
    error.value = "";
    try {
      today.value = await getTodayReading(token);
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to load reading";
    } finally {
      loading.value = false;
    }
  }

  async function generate(token: string) {
    generating.value = true;
    error.value = "";
    try {
      today.value = await generateSpreads(token);
    } catch (err) {
      if (err instanceof ApiError && err.status === 429) {
        error.value = "Daily limit reached for your role.";
      } else {
        error.value = err instanceof Error ? err.message : "Failed to generate spreads";
      }
      throw err;
    } finally {
      generating.value = false;
    }
  }

  return {
    today,
    loading,
    generating,
    error,
    loadToday,
    generate
  };
});
