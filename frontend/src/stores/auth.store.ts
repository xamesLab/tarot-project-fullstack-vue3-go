import { computed, ref } from "vue";
import { defineStore } from "pinia";

import type { UserProfile } from "../types/auth";
import * as authApi from "../services/auth.api";

const TOKEN_KEY = "tarot_daily_token";

export const useAuthStore = defineStore("auth", () => {
  const token = ref<string>(localStorage.getItem(TOKEN_KEY) ?? "");
  const user = ref<UserProfile | null>(null);
  const loading = ref(false);
  const error = ref("");

  const isAuthenticated = computed(() => Boolean(token.value));
  const isAdmin = computed(() => user.value?.role === "admin");

  async function hydrate() {
    if (!token.value) {
      return;
    }
    try {
      user.value = await authApi.me(token.value);
    } catch {
      logout();
    }
  }

  async function register(login: string, password: string) {
    loading.value = true;
    error.value = "";
    try {
      const result = await authApi.register({ login, password });
      setSession(result.token, result.user);
    } catch (err) {
      error.value = toErrorMessage(err);
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function login(login: string, password: string) {
    loading.value = true;
    error.value = "";
    try {
      const result = await authApi.login({ login, password });
      setSession(result.token, result.user);
    } catch (err) {
      error.value = toErrorMessage(err);
      throw err;
    } finally {
      loading.value = false;
    }
  }

  function logout() {
    token.value = "";
    user.value = null;
    localStorage.removeItem(TOKEN_KEY);
  }

  function setSession(nextToken: string, nextUser: UserProfile) {
    token.value = nextToken;
    user.value = nextUser;
    localStorage.setItem(TOKEN_KEY, nextToken);
  }

  function setUser(nextUser: UserProfile) {
    user.value = nextUser;
  }

  return {
    token,
    user,
    loading,
    error,
    isAuthenticated,
    isAdmin,
    hydrate,
    register,
    login,
    setUser,
    logout
  };
});

function toErrorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }
  return "Unexpected error";
}
