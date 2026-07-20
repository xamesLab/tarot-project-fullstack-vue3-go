import type {
  DailyCardApi,
  LoginPayload,
  NebulaTarotApi,
  ReadingHistoryParams,
  RegisterPayload,
  UpdateProfilePayload,
  UpdateSettingsPayload,
  AdminAssetsParams,
  AdminUsersParams,
  UpdateAdminRolePolicyPayload,
  UpdateAdminSettingsPayload,
  UserRole,
} from '@/api/types';
import { post, request } from '@/api/client';

const unavailable = async (): Promise<never> => {
  throw new Error('HTTP API is not implemented. Set VITE_API_MODE=mock.');
};

export const httpApi: NebulaTarotApi = {
  readings: { getDaily: unavailable },
  oracle: { ask: unavailable },
  auth: {
    register: (payload: RegisterPayload) => post('/auth/register', payload),
    login: (payload: LoginPayload) => post('/auth/login', payload),
    getCurrentUser: () => request('/me'),
  },
  profile: {
    updateProfile: (payload: UpdateProfilePayload) => request('/me/profile', { method: 'PATCH', body: JSON.stringify(payload) }),
    updateSettings: (payload: UpdateSettingsPayload) => request('/me/settings', { method: 'PATCH', body: JSON.stringify(payload) }),
  },
  catalog: {
    getAvailableDecks: () => request('/decks/available'),
    getAvailableCardBacks: () => request('/card-backs/available'),
  },
  history: {
    getHistory: ({ type, page, limit }: ReadingHistoryParams) => {
      const query = new URLSearchParams({ type, page: String(page), limit: String(limit) });
      return request(`/reading/history?${query.toString()}`);
    },
  },
  spreads: {
    generateThreeCard: (payload) => post('/reading/generate/three-card', payload),
    revealCard: (readingId, position) => post(`/reading/${readingId}/cards/${position}/reveal`, {}),
    requestInterpretation: (readingId) => post('/ask', { reading_id: readingId }),
  },
  admin: {
    getDashboard: () => request('/admin/dashboard'),
    getUsers: ({ search, role, status, page, limit }: AdminUsersParams) => {
      const query = new URLSearchParams({ page: String(page), limit: String(limit) });
      if (search) query.set('search', search);
      if (role && role !== 'all') query.set('role', role);
      if (status && status !== 'all') query.set('status', status);
      return request(`/admin/users?${query.toString()}`);
    },
    updateUserRole: (userId: string, role: UserRole) => request(`/admin/users/${userId}/role`, { method: 'PATCH', body: JSON.stringify({ role }) }),
    updateUserStatus: (userId: string, status: 'active' | 'disabled') => request(`/admin/users/${userId}/status`, { method: 'PATCH', body: JSON.stringify({ status }) }),
    getRolePolicies: () => request('/admin/roles'),
    updateRolePolicy: (role: UserRole, payload: UpdateAdminRolePolicyPayload) => request(`/admin/roles/${role}`, { method: 'PATCH', body: JSON.stringify(payload) }),
    getSettings: () => request('/admin/settings'),
    updateSettings: (payload: UpdateAdminSettingsPayload) => request('/admin/settings', { method: 'PATCH', body: JSON.stringify(payload) }),
    getDecks: ({ scope, page, limit }: AdminAssetsParams) => request(`/admin/decks?${new URLSearchParams({ scope, page: String(page), limit: String(limit) }).toString()}`),
    createDeck: (payload: FormData) => request('/admin/decks', { method: 'POST', body: payload }),
    archiveDeck: (deckId: string) => request(`/admin/decks/${deckId}/archive`, { method: 'POST' }),
    restoreDeck: (deckId: string) => request(`/admin/decks/${deckId}/restore`, { method: 'POST' }),
    deleteDeck: (deckId: string) => request(`/admin/decks/${deckId}`, { method: 'DELETE' }),
    getCardBacks: ({ scope, page, limit }: AdminAssetsParams) => request(`/admin/card-backs?${new URLSearchParams({ scope, page: String(page), limit: String(limit) }).toString()}`),
    createCardBack: (payload: FormData) => request('/admin/card-backs', { method: 'POST', body: payload }),
    archiveCardBack: (cardBackId: string) => request(`/admin/card-backs/${cardBackId}/archive`, { method: 'POST' }),
    restoreCardBack: (cardBackId: string) => request(`/admin/card-backs/${cardBackId}/restore`, { method: 'POST' }),
    deleteCardBack: (cardBackId: string) => request(`/admin/card-backs/${cardBackId}`, { method: 'DELETE' }),
  },
};

export const httpDailyCardApi: DailyCardApi = {
  generateGuest: (payload) => post('/reading/guest/card-of-day', payload, { auth: false }),
  resolve: (payload) => post('/reading/card-of-day/resolve', payload),
  getGuestInterpretation: (payload) => post('/reading/guest/card-of-day/interpretation', payload, { auth: false }),
  getUserInterpretation: (payload) => post('/reading/card-of-day/interpretation', payload),
};
