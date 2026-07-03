import { api } from './api'

export function getProfile() {
  return api.get('/me').then((response) => response.data)
}

export function updateProfile(data) {
  return api.patch('/me/profile', data).then((response) => response.data)
}

export function updateSettings(data) {
  return api.patch('/me/settings', data).then((response) => response.data)
}
