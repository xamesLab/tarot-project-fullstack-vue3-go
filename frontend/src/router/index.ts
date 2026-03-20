import type { Pinia } from "pinia";
import { createRouter, createWebHistory } from "vue-router";

import { getDevToolsRoutes } from "./dev-tools.routes";
import { installRouterGuards } from "./guards";

const LoginPage = () => import("../pages/LoginPage.vue");
const RegisterPage = () => import("../pages/RegisterPage.vue");
const DashboardPage = () => import("../pages/DashboardPage.vue");
const ProfilePage = () => import("../pages/ProfilePage.vue");
const SettingsPage = () => import("../pages/SettingsPage.vue");
const AdminShellPage = () => import("../pages/admin/AdminShellPage.vue");
const AdminUsersPage = () => import("../pages/admin/AdminUsersPage.vue");
const AdminDecksPage = () => import("../pages/admin/AdminDecksPage.vue");
const AdminSettingsPage = () => import("../pages/admin/AdminSettingsPage.vue");
const NotFoundPage = () => import("../pages/NotFoundPage.vue");

const baseRoutes = [
  { path: "/", redirect: "/dashboard" },
  {
    path: "/login",
    name: "login",
    component: LoginPage,
    meta: { publicOnly: true }
  },
  {
    path: "/register",
    name: "register",
    component: RegisterPage,
    meta: { publicOnly: true }
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: DashboardPage,
    meta: { requiresAuth: true }
  },
  {
    path: "/profile",
    name: "profile",
    component: ProfilePage,
    meta: { requiresAuth: true }
  },
  {
    path: "/settings",
    name: "settings",
    component: SettingsPage,
    meta: { requiresAuth: true }
  },
  {
    path: "/admin",
    component: AdminShellPage,
    meta: { requiresAuth: true, adminOnly: true },
    children: [
      {
        path: "",
        redirect: { name: "admin-users" }
      },
      {
        path: "users",
        name: "admin-users",
        component: AdminUsersPage
      },
      {
        path: "decks",
        name: "admin-decks",
        component: AdminDecksPage
      },
      {
        path: "settings",
        name: "admin-settings",
        component: AdminSettingsPage
      }
    ]
  },
  {
    path: "/:pathMatch(.*)*",
    name: "not-found",
    component: NotFoundPage
  }
];

const devRoutes = import.meta.env.DEV ? getDevToolsRoutes() : [];

export const router = createRouter({
  history: createWebHistory(),
  routes: [...baseRoutes, ...devRoutes]
});

export function setupRouterGuards(pinia: Pinia) {
  installRouterGuards(router, pinia);
}
