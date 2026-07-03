import { defineStore } from 'pinia'

import * as authService from '../services/auth'
import * as userService from '../services/user'
import { normalizeApiError, TOKEN_STORAGE_KEY } from '../services/api'
import { useAppStore } from './app'

function readStoredToken() {
  if (typeof window === 'undefined') {
    return ''
  }

  return window.localStorage.getItem(TOKEN_STORAGE_KEY) || ''
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: readStoredToken(),
    status: 'idle',
    hydrated: false,
  }),

  getters: {
    isAuthenticated: (state) => Boolean(state.token && state.user),
    hasToken: (state) => Boolean(state.token),
    userRole: (state) => state.user?.role || 'guest',
    isKeeper() {
      return this.userRole === 'keeper'
    },
  },

  actions: {
    persistToken(token) {
      this.token = token || ''

      if (typeof window === 'undefined') {
        return
      }

      if (token) {
        window.localStorage.setItem(TOKEN_STORAGE_KEY, token)
      } else {
        window.localStorage.removeItem(TOKEN_STORAGE_KEY)
      }
    },

    setUser(user) {
      this.user = user
      useAppStore().syncUserPreferences(user)
    },

    clearSession() {
      this.persistToken('')
      this.user = null
      this.status = 'idle'
    },

    async hydrateSession() {
      if (this.hydrated) {
        return this.user
      }

      if (!this.token) {
        this.hydrated = true
        return null
      }

      try {
        await this.loadUser()
      } catch {
        this.clearSession()
      } finally {
        this.hydrated = true
      }

      return this.user
    },

    async loadUser() {
      this.status = 'loading'

      try {
        const response = await userService.getProfile()
        this.setUser(response.user)
        this.status = 'ready'
        return response.user
      } catch (error) {
        this.status = 'error'
        throw normalizeApiError(error)
      }
    },

    async login(credentials) {
      this.status = 'loading'

      try {
        const response = await authService.login(credentials)
        this.persistToken(response.token)
        this.setUser(response.user)
        this.hydrated = true
        this.status = 'ready'

        return { ok: true, user: response.user }
      } catch (error) {
        this.status = 'error'
        return { ok: false, error: normalizeApiError(error) }
      }
    },

    async register(payload) {
      this.status = 'loading'

      try {
        const response = await authService.register(payload)
        this.persistToken(response.token)
        this.setUser(response.user)
        this.hydrated = true
        this.status = 'ready'

        return { ok: true, user: response.user }
      } catch (error) {
        this.status = 'error'
        return { ok: false, error: normalizeApiError(error) }
      }
    },

    async updateProfile(payload) {
      this.status = 'loading'

      try {
        const response = await userService.updateProfile(payload)
        this.setUser(response.user)
        this.status = 'ready'

        return { ok: true, user: response.user }
      } catch (error) {
        this.status = 'error'
        return { ok: false, error: normalizeApiError(error) }
      }
    },

    async updateSettings(payload) {
      this.status = 'loading'

      try {
        const response = await userService.updateSettings(payload)
        this.setUser(response.user)
        this.status = 'ready'

        return { ok: true, user: response.user }
      } catch (error) {
        this.status = 'error'
        return { ok: false, error: normalizeApiError(error) }
      }
    },

    logout() {
      this.clearSession()
      this.hydrated = true
    },
  },
})
