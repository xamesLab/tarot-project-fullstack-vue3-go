import axios from 'axios'

export const TOKEN_STORAGE_KEY = 'tarot_daily_token'
export const APP_UNAUTHORIZED_EVENT = 'tarot:unauthorized'
export const APP_FORBIDDEN_EVENT = 'tarot:forbidden'

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

export const api = axios.create({
  baseURL,
})

export function resolveApiAssetUrl(path) {
  if (!path) {
    return ''
  }

  if (/^https?:\/\//i.test(path)) {
    return path
  }

  const apiRoot = baseURL.replace(/\/api\/?$/, '/')

  try {
    return new URL(path, apiRoot).toString()
  } catch {
    return path
  }
}

api.interceptors.request.use((config) => {
  const token = typeof window !== 'undefined' ? window.localStorage.getItem(TOKEN_STORAGE_KEY) : null

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  if (typeof FormData !== 'undefined' && config.data instanceof FormData) {
    delete config.headers['Content-Type']
  } else if (!config.headers['Content-Type']) {
    config.headers['Content-Type'] = 'application/json'
  }

  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status
    const token = typeof window !== 'undefined' ? window.localStorage.getItem(TOKEN_STORAGE_KEY) : null

    if (typeof window !== 'undefined') {
      if (status === 401 && token) {
        window.dispatchEvent(new CustomEvent(APP_UNAUTHORIZED_EVENT, { detail: normalizeApiError(error) }))
      }

      if (status === 403 && token) {
        window.dispatchEvent(new CustomEvent(APP_FORBIDDEN_EVENT, { detail: normalizeApiError(error) }))
      }
    }

    return Promise.reject(error)
  },
)

export function normalizeApiError(error) {
  if (!error) {
    return {
      status: 0,
      code: 'unknown_error',
      message: 'Unknown error',
    }
  }

  if (error.response) {
    return {
      status: error.response.status,
      code: error.response.data?.error || 'request_failed',
      message: error.response.data?.message || 'Request failed',
    }
  }

  if (error.request) {
    return {
      status: 0,
      code: 'network_error',
      message: 'Network error. Please check your connection and backend availability.',
    }
  }

  return {
    status: 0,
    code: 'client_error',
    message: error.message || 'Client error',
  }
}
