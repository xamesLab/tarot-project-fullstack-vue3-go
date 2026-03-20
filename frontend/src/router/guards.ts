import type { Pinia } from "pinia";
import type { Router } from "vue-router";

import { useAuthStore } from "../stores/auth.store";

export function installRouterGuards(router: Router, pinia: Pinia) {
  router.beforeEach(async (to) => {
    if (to.meta.devOnly && !import.meta.env.DEV) {
      return { name: "not-found" };
    }

    const authStore = useAuthStore(pinia);

    if (authStore.isAuthenticated && !authStore.user) {
      await authStore.hydrate();
    }

    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
      return { name: "login", query: { redirect: to.fullPath } };
    }

    if (to.meta.publicOnly && authStore.isAuthenticated) {
      return { name: "dashboard" };
    }

    if (to.meta.adminOnly && !authStore.isAdmin) {
      return { name: "dashboard" };
    }

    return true;
  });
}
