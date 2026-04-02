import { api } from './api'

export function getAvailableDecks() {
  return api.get('/decks/available').then((response) => response.data)
}

export function getAvailableCardBacks() {
  return api.get('/card-backs/available').then((response) => response.data)
}
