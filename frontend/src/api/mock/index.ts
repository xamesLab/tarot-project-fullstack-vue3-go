import type { Language, LoginPayload, NebulaTarotApi, RegisterPayload, UserDto } from '@/api/types';
import { getReadingFixture } from './fixtures';

const delay = (milliseconds: number) => new Promise((resolve) => setTimeout(resolve, milliseconds));

const replies: Record<Language, string> = {
  ru: 'Карта советует не спешить. Загляни внутрь: ответ уже зреет в твоей интуиции. Позволь дню стать пространством для созерцания.',
  en: 'The card counsels patience. Look within: the answer already ripens in your intuition. Let this day be a space for contemplation.',
};

let mockUser: UserDto = {
  id: 'mock-user',
  username: 'oracle_guest',
  email: 'guest@example.com',
  role: 'seeker',
  status: 'active',
  timezone: 'Europe/Moscow',
  language: 'ru',
  theme: 'dark',
};

export const mockApi: NebulaTarotApi = {
  readings: {
    async getDaily(language) {
      await delay(350);
      return getReadingFixture(language);
    },
  },
  oracle: {
    async ask(_question, language) {
      await delay(600);
      return { id: crypto.randomUUID(), text: replies[language] };
    },
  },
  auth: {
    async register(payload: RegisterPayload) {
      await delay(650);
      mockUser = { ...mockUser, id: crypto.randomUUID(), username: payload.username || payload.email.split('@')[0], email: payload.email, timezone: payload.timezone };
      return { success: true, user: mockUser, token: 'mock-token' };
    },
    async login(payload: LoginPayload) {
      await delay(450);
      mockUser = { ...mockUser, email: payload.email };
      return { success: true, user: mockUser, token: 'mock-token' };
    },
    async getCurrentUser() {
      await delay(150);
      return { success: true, user: mockUser };
    },
  },
};
