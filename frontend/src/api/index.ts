import { httpApi } from '@/api/http';
import { mockApi } from '@/api/mock';

export const api = import.meta.env.VITE_API_MODE === 'http' ? httpApi : mockApi;
