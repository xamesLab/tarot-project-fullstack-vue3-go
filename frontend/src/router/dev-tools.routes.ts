import type { RouteRecordRaw } from "vue-router";

const UIKitPage = () => import("../pages/dev/UIKit.vue");

export function getDevToolsRoutes(): RouteRecordRaw[] {
  return [
    {
      path: "/ui-kit",
      name: "ui-kit",
      component: UIKitPage,
      meta: { devOnly: true }
    }
  ];
}
