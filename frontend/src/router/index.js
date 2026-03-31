import { createRouter, createWebHistory } from 'vue-router'

import HomePage from '../pages/HomePage.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomePage,
    meta: {
      titleKey: 'home.title',
    },
  },
]

if (import.meta.env.DEV) {
  routes.push({
    path: '/ui-kit',
    name: 'ui-kit',
    component: () => import('../pages/dev/UIKit.vue'),
    meta: {
      titleKey: 'uiKit.title',
    },
  })
}

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
