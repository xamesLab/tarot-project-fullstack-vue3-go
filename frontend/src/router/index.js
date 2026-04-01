import { createRouter, createWebHistory } from 'vue-router'

import { useAuthStore } from '../stores/auth'
import Dashboard from '../pages/Dashboard.vue'
import History from '../pages/History.vue'
import Login from '../pages/Login.vue'
import Profile from '../pages/Profile.vue'
import Register from '../pages/Register.vue'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: {
      titleKey: 'auth.loginTitle',
      guestOnly: true,
    },
  },
  {
    path: '/register',
    name: 'register',
    component: Register,
    meta: {
      titleKey: 'auth.registerTitle',
      guestOnly: true,
    },
  },
  {
    path: '/',
    name: 'dashboard',
    component: Dashboard,
    meta: {
      titleKey: 'dashboard.pageTitle',
      requiresAuth: true,
    },
  },
  {
    path: '/profile',
    name: 'profile',
    component: Profile,
    meta: {
      titleKey: 'profile.title',
      requiresAuth: true,
    },
  },
  {
    path: '/history',
    name: 'history',
    component: History,
    meta: {
      titleKey: 'history.title',
      requiresAuth: true,
    },
  },
  {
    path: '/admin',
    name: 'admin',
    component: () => import('../pages/admin/AdminDashboard.vue'),
    meta: {
      titleKey: 'admin.title',
      requiresAuth: true,
      requiresKeeper: true,
    },
  },
  {
    path: '/admin/users',
    name: 'admin-users',
    component: () => import('../pages/admin/Users.vue'),
    meta: {
      titleKey: 'admin.usersTitle',
      requiresAuth: true,
      requiresKeeper: true,
    },
  },
  {
    path: '/admin/decks',
    name: 'admin-decks',
    component: () => import('../pages/admin/Decks.vue'),
    meta: {
      titleKey: 'admin.decksTitle',
      requiresAuth: true,
      requiresKeeper: true,
    },
  },
  {
    path: '/admin/role-policies',
    name: 'admin-role-policies',
    component: () => import('../pages/admin/RolePolicies.vue'),
    meta: {
      titleKey: 'admin.policiesTitle',
      requiresAuth: true,
      requiresKeeper: true,
    },
  },
  {
    path: '/admin/settings',
    name: 'admin-settings',
    component: () => import('../pages/admin/Settings.vue'),
    meta: {
      titleKey: 'admin.settingsTitle',
      requiresAuth: true,
      requiresKeeper: true,
    },
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
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

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  if (!authStore.hydrated) {
    await authStore.hydrateSession()
  }

  if (to.meta.guestOnly && authStore.isAuthenticated) {
    return { name: 'dashboard' }
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return {
      name: 'login',
      query: {
        redirect: to.fullPath,
      },
    }
  }

  if (to.meta.requiresKeeper && !authStore.isKeeper) {
    return { name: 'dashboard' }
  }

  if (to.name === 'ui-kit' && !import.meta.env.DEV) {
    return { name: 'dashboard' }
  }

  return true
})

export default router
