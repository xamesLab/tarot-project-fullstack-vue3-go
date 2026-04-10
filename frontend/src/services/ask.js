import { api } from './api'

export function createAskRequest(payload) {
  return api.post('/ask', payload).then((response) => response.data)
}

export function getAskRequest(requestId) {
  return api.get(`/ask/${requestId}`).then((response) => response.data)
}
