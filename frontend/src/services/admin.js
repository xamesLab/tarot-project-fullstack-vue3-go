import { api } from './api'

export function accessCheck() {
  return api.get('/admin/access-check').then((response) => response.data)
}

export function getDashboard() {
  return api.get('/admin/dashboard').then((response) => response.data)
}

export function getUsers(params = {}) {
  return api.get('/admin/users', { params }).then((response) => response.data)
}

export function updateUserRole(userId, payload) {
  return api.patch(`/admin/users/${userId}/role`, payload).then((response) => response.data)
}

export function updateUserStatus(userId, payload) {
  return api.patch(`/admin/users/${userId}/status`, payload).then((response) => response.data)
}

export function getDecks(params = {}) {
  return api.get('/admin/decks', { params }).then((response) => response.data)
}

export function createDeck(formData) {
  return api.post('/admin/decks', formData).then((response) => response.data)
}

export function getCardBacks(params = {}) {
  return api.get('/admin/card-backs', { params }).then((response) => response.data)
}

export function createCardBack(formData) {
  return api.post('/admin/card-backs', formData).then((response) => response.data)
}

export function getDeck(deckId) {
  return api.get(`/admin/decks/${deckId}`).then((response) => response.data)
}

export function updateDeck(deckId, payload) {
  return api.patch(`/admin/decks/${deckId}`, payload).then((response) => response.data)
}

export function getRolePolicies() {
  return api.get('/admin/roles').then((response) => response.data)
}

export function updateRolePolicy(role, payload) {
  return api.patch(`/admin/roles/${role}`, payload).then((response) => response.data)
}

export function getSettings() {
  return api.get('/admin/settings').then((response) => response.data)
}

export function updateSettings(payload) {
  return api.patch('/admin/settings', payload).then((response) => response.data)
}
