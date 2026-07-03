import { api } from './api'

export function generateCardOfDay(payload = {}) {
  return api.post('/reading/generate/card-of-day', payload).then((response) => response.data)
}

export function generateThreeCard(payload = {}) {
  return api.post('/reading/generate/three-card', payload).then((response) => response.data)
}

export function getHistory(params = {}) {
  return api.get('/reading/history', { params }).then((response) => response.data)
}

export function revealReadingCard(readingId, position) {
  return api.post(`/reading/${readingId}/cards/${position}/reveal`).then((response) => response.data)
}
