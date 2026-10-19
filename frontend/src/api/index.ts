import { httpApi, httpDailyCardApi } from '@/api/http';
import { mockApi } from '@mock/index';
import { mockDailyCardApi } from '@mock/daily-card';

const contentApi = import.meta.env.VITE_API_MODE === 'http' ? httpApi : mockApi;
const accountApi = import.meta.env.VITE_AUTH_API_MODE === 'http' ? httpApi : mockApi;

export const api = {
  readings: contentApi.readings,
  oracle: contentApi.oracle,
  auth: accountApi.auth,
  profile: accountApi.profile,
  catalog: accountApi.catalog,
  history: accountApi.history,
  spreads: accountApi.spreads,
  admin: accountApi.admin,
  dailyCard: import.meta.env.VITE_AUTH_API_MODE === 'http' ? httpDailyCardApi : mockDailyCardApi,
};
