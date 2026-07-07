import { httpApi, httpDailyCardApi } from '@/api/http';
import { mockApi } from '@/api/mock';

const contentApi = import.meta.env.VITE_API_MODE === 'http' ? httpApi : mockApi;
const authApi = import.meta.env.VITE_AUTH_API_MODE === 'http' ? httpApi.auth : mockApi.auth;

export const api = {
  readings: contentApi.readings,
  oracle: contentApi.oracle,
  auth: authApi,
  dailyCard: httpDailyCardApi,
};
