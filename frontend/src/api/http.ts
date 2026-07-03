import type { NebulaTarotApi } from '@/api/types';

const unavailable = async (): Promise<never> => {
  throw new Error('HTTP API is not implemented. Set VITE_API_MODE=mock.');
};

export const httpApi: NebulaTarotApi = {
  readings: { getDaily: unavailable },
  oracle: { ask: unavailable },
  auth: { register: unavailable },
};
