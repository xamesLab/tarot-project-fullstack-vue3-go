import { defineStore } from 'pinia'

import { normalizeApiError } from '../services/api'
import * as askService from '../services/ask'
import * as readingService from '../services/reading'
import { getAccountingDateForTimezone } from '../utils/date'
import { useAuthStore } from './auth'

function resolveRequestTimezone(user) {
  return user?.timezone || Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC'
}

function buildRevealKey(readingId, position) {
  return `${readingId}:${position}`
}

function mergeReadingState(previousReading, nextReading) {
  if (!previousReading) {
    return nextReading
  }

  return {
    ...previousReading,
    ...nextReading,
    cards: nextReading?.cards || previousReading.cards || [],
  }
}

function buildAIStateFromRequest(request) {
  if (!request) {
    return null
  }

  return {
    request_id: request.id || '',
    provider: request.provider || '',
    model: request.model || '',
    prompt_version: request.prompt_version || '',
    language: request.request_payload?.language || '',
    status: request.status || '',
    content: request.response_excerpt || '',
    error_code: request.error_code || '',
    error_message: request.error_message || '',
    completed_at: request.completed_at || null,
    last_requested_at: request.requested_at || null,
    updated_at: request.completed_at || request.requested_at || null,
  }
}

export const useReadingsStore = defineStore('readings', {
  state: () => ({
    cardOfDay: null,
    latestSpread: null,
    cardOfDayNextAvailableAt: '',
    spreadLimit: null,
    history: [],
    historyFilter: 'all',
    historyPagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_pages: 1,
    },
    retentionInfo: null,
    dashboardStatus: 'idle',
    historyStatus: 'idle',
    cardOfDayStatus: 'idle',
    spreadStatus: 'idle',
    dashboardError: '',
    historyError: '',
    cardOfDayError: '',
    spreadError: '',
    revealStatusByKey: {},
    revealErrorByKey: {},
    aiStatusByReadingId: {},
    aiErrorByReadingId: {},
  }),

  actions: {
    clearErrors() {
      this.dashboardError = ''
      this.historyError = ''
      this.cardOfDayError = ''
      this.spreadError = ''
    },

    setReadingAI(readingId, ai) {
      if (!readingId) {
        return
      }

      this.replaceReading({
        id: readingId,
        ai,
      })
    },

    replaceReading(updatedReading) {
      if (!updatedReading?.id) {
        return
      }

      if (this.cardOfDay?.id === updatedReading.id) {
        this.cardOfDay = mergeReadingState(this.cardOfDay, updatedReading)
      }

      if (this.latestSpread?.id === updatedReading.id) {
        this.latestSpread = mergeReadingState(this.latestSpread, updatedReading)
      }

      this.history = this.history.map((reading) => (
        reading.id === updatedReading.id
          ? mergeReadingState(reading, updatedReading)
          : reading
      ))
    },

    prependReading(reading) {
      if (!reading) {
        return
      }

      this.history = [reading, ...this.history.filter((item) => item.id !== reading.id)]
      this.historyPagination.total += 1
    },

    async loadDashboard() {
      const authStore = useAuthStore()
      this.dashboardStatus = 'loading'
      this.dashboardError = ''

      try {
        const response = await readingService.getHistory({
          type: 'all',
          page: 1,
          limit: 20,
        })

        const accountingDate = getAccountingDateForTimezone(resolveRequestTimezone(authStore.user))
        this.history = response.readings || []
        this.cardOfDay = this.history.find((item) => item.type === 'card_of_day' && item.accounting_date === accountingDate) || null
        this.latestSpread = this.history.find((item) => item.type === 'three_card' && item.accounting_date === accountingDate) || null
        this.retentionInfo = response.retention_info || null
        this.dashboardStatus = 'ready'
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.dashboardError = normalized.message
        this.dashboardStatus = 'error'
      }
    },

    async loadHistory({ type = 'all', page = 1, limit = 10 } = {}) {
      this.historyStatus = 'loading'
      this.historyError = ''

      try {
        const response = await readingService.getHistory({ type, page, limit })
        this.history = response.readings || []
        this.historyFilter = type
        this.historyPagination = response.pagination || this.historyPagination
        this.retentionInfo = response.retention_info || null
        this.historyStatus = 'ready'
        return { ok: true }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.historyError = normalized.message
        this.historyStatus = 'error'
        return { ok: false, error: normalized }
      }
    },

    async generateCardOfDay() {
      const authStore = useAuthStore()
      this.cardOfDayStatus = 'loading'
      this.cardOfDayError = ''

      try {
        const response = await readingService.generateCardOfDay({
          timezone: resolveRequestTimezone(authStore.user),
        })

        this.cardOfDay = response.reading
        this.cardOfDayNextAvailableAt = response.reading?.next_available_at || ''
        this.prependReading(response.reading)
        this.cardOfDayStatus = 'ready'
        return { ok: true, reading: response.reading }
      } catch (error) {
        const normalized = normalizeApiError(error)
        const payload = error.response?.data || {}

        if (normalized.code === 'card_of_day_already_generated' && payload.existing_reading) {
          this.cardOfDay = {
            ...payload.existing_reading,
            next_available_at: payload.next_available_at || '',
          }
          this.cardOfDayNextAvailableAt = payload.next_available_at || ''
          this.cardOfDayError = ''
          this.cardOfDayStatus = 'ready'
          return { ok: false, error: normalized, payload }
        }

        this.cardOfDayError = normalized.message
        this.cardOfDayStatus = 'error'
        return { ok: false, error: normalized, payload }
      }
    },

    async generateThreeCardSpread() {
      const authStore = useAuthStore()
      this.spreadStatus = 'loading'
      this.spreadError = ''

      try {
        const response = await readingService.generateThreeCard({
          timezone: resolveRequestTimezone(authStore.user),
        })

        this.latestSpread = response.reading
        this.spreadLimit = response.reading?.remaining_today || null
        this.prependReading(response.reading)
        this.spreadStatus = 'ready'
        return { ok: true, reading: response.reading }
      } catch (error) {
        const normalized = normalizeApiError(error)
        const payload = error.response?.data || {}

        if (normalized.code === 'daily_limit_reached' && payload.current_limit) {
          this.spreadLimit = payload.current_limit
        }

        this.spreadError = normalized.message
        this.spreadStatus = 'error'
        return { ok: false, error: normalized, payload }
      }
    },

    async revealCard(readingId, position) {
      const revealKey = buildRevealKey(readingId, position)
      this.revealStatusByKey = {
        ...this.revealStatusByKey,
        [revealKey]: 'loading',
      }
      this.revealErrorByKey = {
        ...this.revealErrorByKey,
        [revealKey]: '',
      }

      try {
        const response = await readingService.revealReadingCard(readingId, position)
        this.replaceReading(response.reading)
        this.revealStatusByKey = {
          ...this.revealStatusByKey,
          [revealKey]: 'ready',
        }
        return { ok: true, reading: response.reading }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.revealStatusByKey = {
          ...this.revealStatusByKey,
          [revealKey]: 'error',
        }
        this.revealErrorByKey = {
          ...this.revealErrorByKey,
          [revealKey]: normalized.message,
        }
        return { ok: false, error: normalized }
      }
    },

    async requestAIInterpretation(readingId) {
      if (!readingId) {
        return { ok: false }
      }

      const reading = [
        this.cardOfDay,
        this.latestSpread,
        ...this.history,
      ].find((item) => item?.id === readingId)

      if (reading?.ai?.status === 'completed' && reading.ai?.content) {
        return { ok: true, request: null, reused: true }
      }

      if (this.aiStatusByReadingId[readingId] === 'loading') {
        return { ok: false, loading: true }
      }

      this.aiStatusByReadingId = {
        ...this.aiStatusByReadingId,
        [readingId]: 'loading',
      }
      this.aiErrorByReadingId = {
        ...this.aiErrorByReadingId,
        [readingId]: '',
      }

      try {
        const response = await askService.createAskRequest({ reading_id: readingId })
        const request = response.request
        this.setReadingAI(readingId, buildAIStateFromRequest(request))
        this.aiStatusByReadingId = {
          ...this.aiStatusByReadingId,
          [readingId]: 'ready',
        }
        this.aiErrorByReadingId = {
          ...this.aiErrorByReadingId,
          [readingId]: '',
        }
        return { ok: true, request }
      } catch (error) {
        const normalized = normalizeApiError(error)
        this.aiStatusByReadingId = {
          ...this.aiStatusByReadingId,
          [readingId]: 'error',
        }
        this.aiErrorByReadingId = {
          ...this.aiErrorByReadingId,
          [readingId]: normalized.message,
        }
        this.setReadingAI(readingId, {
          ...(reading?.ai || {}),
          status: 'failed',
          error_code: normalized.code,
          error_message: normalized.message,
          updated_at: new Date().toISOString(),
        })
        return { ok: false, error: normalized }
      }
    },
  },
})
