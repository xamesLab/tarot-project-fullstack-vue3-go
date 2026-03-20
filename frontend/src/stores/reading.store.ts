import { ref } from "vue";
import { defineStore } from "pinia";

import type { TodayReading, GatewayHealthResponse, GatewayHealthStatus } from "../types/reading";
import { generateSpreads, generateCardOfDay, generateSpreadA, generateSpreadB, getTodayReading } from "../services/reading.api";
import { getGatewayHealth } from "../services/gateway.api";
import { ApiError } from "../services/apiClient";

export const useReadingStore = defineStore("reading", () => {
  const today = ref<TodayReading | null>(null);
  const loading = ref(false);
  const generating = ref(false);
  const error = ref("");
  const gatewayHealthLoading = ref(false);
  const isHealthModalOpen = ref(false);
  const healthStatus = ref<GatewayHealthStatus | null>(null);

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

  async function generateCardOfDayOnly(token: string) {
    generating.value = true;
    error.value = "";
    try {
      today.value = await generateCardOfDay(token);
    } catch (err) {
      if (err instanceof ApiError && err.status === 429) {
        error.value = "Daily limit reached for your role.";
      } else {
        error.value = err instanceof Error ? err.message : "Failed to generate card of the day";
      }
      throw err;
    } finally {
      generating.value = false;
    }
  }

  async function generateSpreadAOnly(token: string) {
    generating.value = true;
    error.value = "";
    try {
      today.value = await generateSpreadA(token);
    } catch (err) {
      if (err instanceof ApiError && err.status === 429) {
        error.value = "Daily limit reached for your role.";
      } else {
        error.value = err instanceof Error ? err.message : "Failed to generate spread A";
      }
      throw err;
    } finally {
      generating.value = false;
    }
  }

  async function generateSpreadBOnly(token: string) {
    generating.value = true;
    error.value = "";
    try {
      today.value = await generateSpreadB(token);
    } catch (err) {
      if (err instanceof ApiError && err.status === 429) {
        error.value = "Daily limit reached for your role.";
      } else {
        error.value = err instanceof Error ? err.message : "Failed to generate spread B";
      }
      throw err;
    } finally {
      generating.value = false;
    }
  }

  async function checkGatewayHealth(token: string): Promise<GatewayHealthResponse> {
    gatewayHealthLoading.value = true;
    try {
      const result = await getGatewayHealth(token);
      healthStatus.value = {
        status: "success",
        title: "Gateway доступен",
        message: `Время работы: ${formatUptime(result.uptime)}`,
        uptime: result.uptime
      };
      isHealthModalOpen.value = true;
      return result;
    } catch (err) {
      const errorMessage = err instanceof ApiError ? err.message : "Неизвестная ошибка";
      healthStatus.value = {
        status: "error",
        title: "Ошибка проверки Gateway",
        message: errorMessage
      };
      isHealthModalOpen.value = true;
      throw err;
    } finally {
      gatewayHealthLoading.value = false;
    }
  }

  function closeHealthModal() {
    isHealthModalOpen.value = false;
    healthStatus.value = null;
  }

  function formatUptime(seconds: number): string {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}ч ${minutes}м`;
  }

  return {
    today,
    loading,
    generating,
    error,
    gatewayHealthLoading,
    isHealthModalOpen,
    healthStatus,
    loadToday,
    generate,
    generateCardOfDayOnly,
    generateSpreadAOnly,
    generateSpreadBOnly,
    checkGatewayHealth,
    closeHealthModal
  };
});
