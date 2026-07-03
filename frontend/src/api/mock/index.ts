import type { Language, NebulaTarotApi, RegisterPayload } from '@/api/types';
import { getReadingFixture } from './fixtures';

const delay = (milliseconds: number) => new Promise((resolve) => setTimeout(resolve, milliseconds));

const replies: Record<Language, string> = {
  ru: 'Карта советует не спешить. Загляни внутрь: ответ уже зреет в твоей интуиции. Позволь дню стать пространством для созерцания.',
  en: 'The card counsels patience. Look within: the answer already ripens in your intuition. Let this day be a space for contemplation.',
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
      return { id: crypto.randomUUID(), name: payload.name, email: payload.email, birthDate: payload.birthDate };
    },
  },
};
