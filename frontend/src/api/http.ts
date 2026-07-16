import type { DailyCardApi, LoginPayload, NebulaTarotApi, RegisterPayload } from '@/api/types';
import { request } from '@/api/client';

const unavailable = async (): Promise<never> => {
  throw new Error('HTTP API is not implemented. Set VITE_API_MODE=mock.');
};

export const httpApi: NebulaTarotApi = {
  readings: { getDaily: unavailable },
  oracle: { ask: unavailable },
  auth: {
    register: (payload: RegisterPayload) => request('/auth/register', { method: 'POST', body: JSON.stringify(payload) }),
    login: (payload: LoginPayload) => request('/auth/login', { method: 'POST', body: JSON.stringify(payload) }),
    getCurrentUser: () => request('/me'),
  },
};

export const httpDailyCardApi: DailyCardApi = {
  generateGuest: (payload) => request('/reading/guest/card-of-day', { method: 'POST', body: JSON.stringify(payload) }, { auth: false }),
  resolve: (payload) => request('/reading/card-of-day/resolve', { method: 'POST', body: JSON.stringify(payload) }),
  getGuestInterpretation: (payload) => request('/reading/guest/card-of-day/interpretation', { method: 'POST', body: JSON.stringify(payload) }, { auth: false }),
  getUserInterpretation: (payload) => request('/reading/card-of-day/interpretation', { method: 'POST', body: JSON.stringify(payload) }),
};
