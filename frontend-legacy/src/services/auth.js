import { api } from './api'

export function register(data) {
  return api.post('/auth/register', data).then((response) => response.data)
}

export function login(credentials) {
  return api.post('/auth/login', credentials).then((response) => response.data)
}
