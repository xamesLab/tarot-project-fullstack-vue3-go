import { defineStore } from 'pinia'

import * as adminService from '../services/admin'
import { normalizeApiError } from '../services/api'

const DEFAULT_PAGINATION = {
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 1,
}

function replaceById(items, nextItem) {
  return items.map((item) => (item.id === nextItem.id ? nextItem : item))
}

function removeById(items, id) {
  return items.filter((item) => item.id !== id)
}

export const useAdminStore = defineStore('admin', {
  state: () => ({
    dashboardMetrics: null,
    dashboardStatus: 'idle',
    dashboardError: '',

    userList: [],
    userFilters: {
      search: '',
      role: 'all',
      status: 'all',
    },
    userPagination: { ...DEFAULT_PAGINATION },
    usersStatus: 'idle',
    usersError: '',
    userMutationStatusById: {},
    userMutationErrorById: {},

    deckList: [],
    deckPagination: { ...DEFAULT_PAGINATION },
    decksStatus: 'idle',
    decksError: '',
    archivedDeckList: [],
    archivedDeckPagination: { ...DEFAULT_PAGINATION },
    archivedDecksStatus: 'idle',
    archivedDecksError: '',
    deckMutationStatusById: {},
    deckMutationErrorById: {},
    deckCreateStatus: 'idle',
    deckCreateError: '',

    cardBackList: [],
    cardBackPagination: { ...DEFAULT_PAGINATION },
    cardBacksStatus: 'idle',
    cardBacksError: '',
    archivedCardBackList: [],
    archivedCardBackPagination: { ...DEFAULT_PAGINATION },
    archivedCardBacksStatus: 'idle',
    archivedCardBacksError: '',
    cardBackMutationStatusById: {},
    cardBackMutationErrorById: {},
    cardBackCreateStatus: 'idle',
    cardBackCreateError: '',

    policies: [],
    policiesStatus: 'idle',
    policiesError: '',
    policyMutationStatusByRole: {},
    policyMutationErrorByRole: {},

    settings: null,
    settingsStatus: 'idle',
    settingsError: '',
    settingsMutationStatus: 'idle',
    settingsMutationError: '',
  }),

  actions: {
    async verifyAccess() {
      try {
        await adminService.accessCheck()
        return { ok: true }
      } catch (error) {
        return { ok: false, error: normalizeApiError(error) }
      }
    },

    async loadDashboard() {
      this.dashboardStatus = 'loading'
      this.dashboardError = ''

      try {
        const response = await adminService.getDashboard()
        this.dashboardMetrics = response.metrics || null
        this.dashboardStatus = 'ready'
        return { ok: true, metrics: this.dashboardMetrics }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.dashboardError = normalized.message
        this.dashboardStatus = 'error'
        return { ok: false, error: normalized }
      }
    },

    async loadUsers(params = {}) {
      const nextFilters = {
        search: params.search ?? this.userFilters.search,
        role: params.role ?? this.userFilters.role,
        status: params.status ?? this.userFilters.status,
      }
      const page = params.page ?? this.userPagination.page ?? 1
      const limit = params.limit ?? this.userPagination.limit ?? 20

      this.usersStatus = 'loading'
      this.usersError = ''

      try {
        const response = await adminService.getUsers({
          search: nextFilters.search || undefined,
          role: nextFilters.role === 'all' ? undefined : nextFilters.role,
          status: nextFilters.status === 'all' ? undefined : nextFilters.status,
          page,
          limit,
        })

        this.userFilters = nextFilters
        this.userList = response.users || []
        this.userPagination = response.pagination || { ...DEFAULT_PAGINATION }
        this.usersStatus = 'ready'
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.usersError = normalized.message
        this.usersStatus = 'error'
        return { ok: false, error: normalized }
      }
    },

    async updateUserRole(userId, role) {
      this.userMutationStatusById[userId] = 'loading'
      this.userMutationErrorById[userId] = ''

      try {
        const response = await adminService.updateUserRole(userId, { role })
        this.userList = replaceById(this.userList, response.user)
        this.userMutationStatusById[userId] = 'ready'
        return { ok: true, user: response.user }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.userMutationStatusById[userId] = 'error'
        this.userMutationErrorById[userId] = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async updateUserStatus(userId, status) {
      this.userMutationStatusById[userId] = 'loading'
      this.userMutationErrorById[userId] = ''

      try {
        const response = await adminService.updateUserStatus(userId, { status })
        this.userList = replaceById(this.userList, response.user)
        this.userMutationStatusById[userId] = 'ready'
        return { ok: true, user: response.user }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.userMutationStatusById[userId] = 'error'
        this.userMutationErrorById[userId] = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async loadDecks(params = {}) {
      const page = params.page ?? this.deckPagination.page ?? 1
      const limit = params.limit ?? this.deckPagination.limit ?? 20

      this.decksStatus = 'loading'
      this.decksError = ''

      try {
        const response = await adminService.getDecks({ page, limit, scope: 'active' })
        this.deckList = response.decks || []
        this.deckPagination = response.pagination || { ...DEFAULT_PAGINATION }
        this.decksStatus = 'ready'
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.decksError = normalized.message
        this.decksStatus = 'error'
        return { ok: false, error: normalized }
      }
    },

    async loadArchivedDecks(params = {}) {
      const page = params.page ?? this.archivedDeckPagination.page ?? 1
      const limit = params.limit ?? this.archivedDeckPagination.limit ?? 20

      this.archivedDecksStatus = 'loading'
      this.archivedDecksError = ''

      try {
        const response = await adminService.getDecks({ page, limit, scope: 'archived' })
        this.archivedDeckList = response.decks || []
        this.archivedDeckPagination = response.pagination || { ...DEFAULT_PAGINATION }
        this.archivedDecksStatus = 'ready'
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.archivedDecksError = normalized.message
        this.archivedDecksStatus = 'error'
        return { ok: false, error: normalized }
      }
    },

    async createDeck(formData) {
      this.deckCreateStatus = 'loading'
      this.deckCreateError = ''

      try {
        const response = await adminService.createDeck(formData)
        const deck = response.deck || null

        if (deck) {
          this.deckList = [deck, ...this.deckList.filter((item) => item.id !== deck.id)]
        }

        this.deckCreateStatus = 'ready'
        await Promise.allSettled([
          this.loadDecks({ page: 1, limit: this.deckPagination.limit || 20 }),
          this.loadDashboard(),
        ])

        return { ok: true, deck, message: response.message || '' }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.deckCreateStatus = 'error'
        this.deckCreateError = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async archiveDeck(deckId) {
      this.deckMutationStatusById[deckId] = 'loading'
      this.deckMutationErrorById[deckId] = ''

      try {
        await adminService.archiveDeck(deckId)
        this.deckMutationStatusById[deckId] = 'ready'
        await Promise.allSettled([
          this.loadDecks({ page: this.deckPagination.page || 1, limit: this.deckPagination.limit || 20 }),
          this.loadArchivedDecks({ page: 1, limit: this.archivedDeckPagination.limit || 20 }),
          this.loadDashboard(),
          this.loadSettings(),
        ])
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.deckMutationStatusById[deckId] = 'error'
        this.deckMutationErrorById[deckId] = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async restoreDeck(deckId) {
      this.deckMutationStatusById[deckId] = 'loading'
      this.deckMutationErrorById[deckId] = ''

      try {
        await adminService.restoreDeck(deckId)
        this.deckMutationStatusById[deckId] = 'ready'
        await Promise.allSettled([
          this.loadDecks({ page: 1, limit: this.deckPagination.limit || 20 }),
          this.loadArchivedDecks({ page: this.archivedDeckPagination.page || 1, limit: this.archivedDeckPagination.limit || 20 }),
          this.loadDashboard(),
          this.loadSettings(),
        ])
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.deckMutationStatusById[deckId] = 'error'
        this.deckMutationErrorById[deckId] = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async deleteDeck(deckId) {
      this.deckMutationStatusById[deckId] = 'loading'
      this.deckMutationErrorById[deckId] = ''

      try {
        await adminService.deleteDeck(deckId)
        this.deckList = removeById(this.deckList, deckId)
        this.archivedDeckList = removeById(this.archivedDeckList, deckId)
        this.deckMutationStatusById[deckId] = 'ready'
        await Promise.allSettled([
          this.loadArchivedDecks({ page: this.archivedDeckPagination.page || 1, limit: this.archivedDeckPagination.limit || 20 }),
          this.loadDashboard(),
          this.loadSettings(),
        ])
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.deckMutationStatusById[deckId] = 'error'
        this.deckMutationErrorById[deckId] = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async loadCardBacks(params = {}) {
      const page = params.page ?? this.cardBackPagination.page ?? 1
      const limit = params.limit ?? this.cardBackPagination.limit ?? 20

      this.cardBacksStatus = 'loading'
      this.cardBacksError = ''

      try {
        const response = await adminService.getCardBacks({ page, limit, scope: 'active' })
        this.cardBackList = response.card_backs || []
        this.cardBackPagination = response.pagination || { ...DEFAULT_PAGINATION }
        this.cardBacksStatus = 'ready'
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.cardBacksError = normalized.message
        this.cardBacksStatus = 'error'
        return { ok: false, error: normalized }
      }
    },

    async loadArchivedCardBacks(params = {}) {
      const page = params.page ?? this.archivedCardBackPagination.page ?? 1
      const limit = params.limit ?? this.archivedCardBackPagination.limit ?? 20

      this.archivedCardBacksStatus = 'loading'
      this.archivedCardBacksError = ''

      try {
        const response = await adminService.getCardBacks({ page, limit, scope: 'archived' })
        this.archivedCardBackList = response.card_backs || []
        this.archivedCardBackPagination = response.pagination || { ...DEFAULT_PAGINATION }
        this.archivedCardBacksStatus = 'ready'
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.archivedCardBacksError = normalized.message
        this.archivedCardBacksStatus = 'error'
        return { ok: false, error: normalized }
      }
    },

    async createCardBack(formData) {
      this.cardBackCreateStatus = 'loading'
      this.cardBackCreateError = ''

      try {
        const response = await adminService.createCardBack(formData)
        const cardBack = response.card_back || null

        if (cardBack) {
          this.cardBackList = [cardBack, ...this.cardBackList.filter((item) => item.id !== cardBack.id)]
        }

        this.cardBackCreateStatus = 'ready'
        await Promise.allSettled([
          this.loadCardBacks({ page: 1, limit: this.cardBackPagination.limit || 20 }),
          this.loadSettings(),
        ])

        return { ok: true, cardBack, message: response.message || '' }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.cardBackCreateStatus = 'error'
        this.cardBackCreateError = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async archiveCardBack(cardBackId) {
      this.cardBackMutationStatusById[cardBackId] = 'loading'
      this.cardBackMutationErrorById[cardBackId] = ''

      try {
        await adminService.archiveCardBack(cardBackId)
        this.cardBackMutationStatusById[cardBackId] = 'ready'
        await Promise.allSettled([
          this.loadCardBacks({ page: this.cardBackPagination.page || 1, limit: this.cardBackPagination.limit || 20 }),
          this.loadArchivedCardBacks({ page: 1, limit: this.archivedCardBackPagination.limit || 20 }),
          this.loadSettings(),
        ])
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.cardBackMutationStatusById[cardBackId] = 'error'
        this.cardBackMutationErrorById[cardBackId] = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async restoreCardBack(cardBackId) {
      this.cardBackMutationStatusById[cardBackId] = 'loading'
      this.cardBackMutationErrorById[cardBackId] = ''

      try {
        await adminService.restoreCardBack(cardBackId)
        this.cardBackMutationStatusById[cardBackId] = 'ready'
        await Promise.allSettled([
          this.loadCardBacks({ page: 1, limit: this.cardBackPagination.limit || 20 }),
          this.loadArchivedCardBacks({ page: this.archivedCardBackPagination.page || 1, limit: this.archivedCardBackPagination.limit || 20 }),
          this.loadSettings(),
        ])
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.cardBackMutationStatusById[cardBackId] = 'error'
        this.cardBackMutationErrorById[cardBackId] = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async deleteCardBack(cardBackId) {
      this.cardBackMutationStatusById[cardBackId] = 'loading'
      this.cardBackMutationErrorById[cardBackId] = ''

      try {
        await adminService.deleteCardBack(cardBackId)
        this.cardBackList = removeById(this.cardBackList, cardBackId)
        this.archivedCardBackList = removeById(this.archivedCardBackList, cardBackId)
        this.cardBackMutationStatusById[cardBackId] = 'ready'
        await Promise.allSettled([
          this.loadArchivedCardBacks({ page: this.archivedCardBackPagination.page || 1, limit: this.archivedCardBackPagination.limit || 20 }),
          this.loadSettings(),
        ])
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.cardBackMutationStatusById[cardBackId] = 'error'
        this.cardBackMutationErrorById[cardBackId] = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async loadPolicies() {
      this.policiesStatus = 'loading'
      this.policiesError = ''

      try {
        const response = await adminService.getRolePolicies()
        this.policies = response.policies || []
        this.policiesStatus = 'ready'
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.policiesError = normalized.message
        this.policiesStatus = 'error'
        return { ok: false, error: normalized }
      }
    },

    async updateRolePolicy(role, payload) {
      this.policyMutationStatusByRole[role] = 'loading'
      this.policyMutationErrorByRole[role] = ''

      try {
        const response = await adminService.updateRolePolicy(role, payload)
        this.policies = this.policies.map((item) => (item.role === role ? response.policy : item))
        this.policyMutationStatusByRole[role] = 'ready'
        return { ok: true, policy: response.policy }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.policyMutationStatusByRole[role] = 'error'
        this.policyMutationErrorByRole[role] = normalized.message
        return { ok: false, error: normalized }
      }
    },

    async loadSettings() {
      this.settingsStatus = 'loading'
      this.settingsError = ''

      try {
        const response = await adminService.getSettings()
        this.settings = response.settings || null
        this.settingsStatus = 'ready'
        return { ok: true, settings: this.settings }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.settingsError = normalized.message
        this.settingsStatus = 'error'
        return { ok: false, error: normalized }
      }
    },

    async updateSettings(payload) {
      this.settingsMutationStatus = 'loading'
      this.settingsMutationError = ''

      try {
        const response = await adminService.updateSettings(payload)
        this.settings = response.settings || null
        this.settingsMutationStatus = 'ready'
        return { ok: true, settings: this.settings }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.settingsMutationStatus = 'error'
        this.settingsMutationError = normalized.message
        return { ok: false, error: normalized }
      }
    },
  },
})
