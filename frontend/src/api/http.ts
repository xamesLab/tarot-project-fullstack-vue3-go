import type { DailyCardApi, LoginPayload, NebulaTarotApi, RegisterPayload } from '@/api/types';
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
};

export const httpDailyCardApi: DailyCardApi = {
  generateGuest: (payload) => post('/reading/guest/card-of-day', payload, { auth: false }),
  resolve: (payload) => post('/reading/card-of-day/resolve', payload),
  getGuestInterpretation: (payload) => post('/reading/guest/card-of-day/interpretation', payload, { auth: false }),
  getUserInterpretation: (payload) => post('/reading/card-of-day/interpretation', payload),
};
