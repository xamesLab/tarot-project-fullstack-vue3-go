import type {
  BackendReadingCardDto,
  CardBackDto,
  DeckDto,
  Language,
  LoginPayload,
  NebulaTarotApi,
  ReadingHistoryItemDto,
  RegisterPayload,
  UserDto,
} from '@/api/types';
import { getReadingFixture } from './fixtures';
import { mockAdminApi } from './admin';

const delay = (milliseconds: number) => new Promise((resolve) => setTimeout(resolve, milliseconds));

const replies: Record<Language, string> = {
  ru: 'Карта советует не спешить. Загляни внутрь: ответ уже зреет в твоей интуиции. Позволь дню стать пространством для созерцания.',
  en: 'The card counsels patience. Look within: the answer already ripens in your intuition. Let this day be a space for contemplation.',
};

let mockUser: UserDto = {
  id: 'mock-user',
  username: 'oracle_guest',
  email: 'guest@example.com',
  role: 'keeper',
  status: 'active',
  timezone: 'Europe/Moscow',
  language: 'ru',
  theme: 'dark',
};

const mockDecks: DeckDto[] = [
  {
    id: 'rider-waite',
    name: { ru: 'Таро Райдера — Уэйта', en: 'Rider–Waite Tarot' },
    description: { ru: 'Классическая колода для ежедневной практики.', en: 'A classic deck for daily practice.' },
    author: 'A. E. Waite',
    year: 1909,
    tags: ['classic'],
    status: 'done',
    card_count: 78,
    back_image_path: '/uploads/decks/rider-waite/back.png',
  },
  {
    id: 'celestial',
    name: { ru: 'Небесное Таро', en: 'Celestial Tarot' },
    description: { ru: 'Колода с мотивами ночного неба.', en: 'A deck inspired by the night sky.' },
    author: 'Arcana',
    year: 2026,
    tags: ['celestial'],
    status: 'done',
    card_count: 78,
    back_image_path: '/uploads/decks/celestial/back.png',
  },
];

const mockCardBacks: CardBackDto[] = [
  {
    id: 'midnight-stars',
    name: { ru: 'Полуночные звёзды', en: 'Midnight Stars' },
    description: { ru: 'Золотые звёзды на глубоком синем фоне.', en: 'Golden stars on deep midnight blue.' },
    tags: ['stars'],
    status: 'done',
    image_path: '/uploads/card-backs/midnight-stars/back.png',
  },
  {
    id: 'lunar-gold',
    name: { ru: 'Лунное золото', en: 'Lunar Gold' },
    description: { ru: 'Лунный орнамент в золотых тонах.', en: 'A lunar ornament in warm gold.' },
    tags: ['moon'],
    status: 'done',
    image_path: '/uploads/card-backs/lunar-gold/back.png',
  },
];

function mockHistoryCard(
  key: string,
  name: string,
  nameRu: string,
  position: number,
): BackendReadingCardDto {
  return {
    key,
    arcana: 'major',
    type: 'Major Arcana',
    type_ru: 'Старший Аркан',
    name,
    name_ru: nameRu,
    interpretation: 'Listen to the quiet answer already taking shape within you.',
    interpretation_ru: 'Прислушайтесь к тихому ответу, который уже складывается внутри вас.',
    tags: 'intuition, path',
    tags_ru: 'интуиция, путь',
    image_url: `/uploads/decks/rider-waite/cards/${key}.png`,
    back_image_url: '/uploads/card-backs/midnight-stars/back.png',
    position,
    reversed: false,
    viewed: true,
  };
}

let mockHistory: ReadingHistoryItemDto[] = [
  {
    id: 'history-daily-1',
    type: 'card_of_day',
    accounting_date: '2026-07-18',
    generated_at: '2026-07-18T05:30:00Z',
    cards: [mockHistoryCard('the_high_priestess', 'The High Priestess', 'Верховная Жрица', 0)],
    deck_id: 'rider-waite',
    back_id: 'midnight-stars',
    ai: { status: 'completed', language: 'ru', content: 'Сегодня интуиция говорит тише обычного, но точнее любых доводов.' },
  },
  {
    id: 'history-spread-1',
    type: 'three_card',
    accounting_date: '2026-07-17',
    generated_at: '2026-07-17T18:10:00Z',
    cards: [
      mockHistoryCard('the_moon', 'The Moon', 'Луна', 0),
      mockHistoryCard('the_star', 'The Star', 'Звезда', 1),
      mockHistoryCard('the_sun', 'The Sun', 'Солнце', 2),
    ],
    deck_id: 'rider-waite',
    back_id: 'midnight-stars',
  },
  {
    id: 'history-daily-2',
    type: 'card_of_day',
    accounting_date: '2026-07-16',
    generated_at: '2026-07-16T04:45:00Z',
    cards: [mockHistoryCard('the_fool', 'The Fool', 'Шут', 0)],
    deck_id: 'rider-waite',
    back_id: 'midnight-stars',
  },
];

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
  profile: {
    async updateProfile(payload) {
      await delay(250);
      mockUser = { ...mockUser, ...payload };
      return { success: true, user: mockUser };
    },
    async updateSettings(payload) {
      await delay(250);
      mockUser = { ...mockUser, ...payload };
      return { success: true, user: mockUser };
    },
  },
  catalog: {
    async getAvailableDecks() {
      await delay(150);
      return { success: true, decks: structuredClone(mockDecks), default_deck_id: mockDecks[0].id };
    },
    async getAvailableCardBacks() {
      await delay(150);
      return { success: true, card_backs: structuredClone(mockCardBacks), default_back_id: mockCardBacks[0].id };
    },
  },
  history: {
    async getHistory({ type, page, limit }) {
      await delay(180);
      const filtered = type === 'all' ? mockHistory : mockHistory.filter((reading) => reading.type === type);
      const start = (page - 1) * limit;
      return {
        success: true,
        readings: structuredClone(filtered.slice(start, start + limit)),
        pagination: {
          page,
          limit,
          total: filtered.length,
          total_pages: Math.max(1, Math.ceil(filtered.length / limit)),
        },
        retention_info: {
          role: mockUser.role,
          retention_days: 30,
          oldest_date: '2026-06-18',
        },
      };
    },
  },
  spreads: {
    async generateThreeCard(payload) {
      await delay(420);
      const now = new Date();
      const timezone = payload.timezone || mockUser.timezone || 'UTC';
      const parts = new Intl.DateTimeFormat('en-US', {
        timeZone: timezone,
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        hourCycle: 'h23',
      }).formatToParts(now);
      const part = (type: Intl.DateTimeFormatPartTypes) => parts.find((item) => item.type === type)?.value || '01';
      const accountingDate = new Date(Date.UTC(Number(part('year')), Number(part('month')) - 1, Number(part('day'))));
      if (Number(part('hour')) < 3) accountingDate.setUTCDate(accountingDate.getUTCDate() - 1);
      const reading: ReadingHistoryItemDto = {
        id: crypto.randomUUID(),
        type: 'three_card',
        accounting_date: accountingDate.toISOString().slice(0, 10),
        generated_at: now.toISOString(),
        cards: [
          mockHistoryCard('the_moon', 'The Moon', 'Луна', 0),
          mockHistoryCard('the_star', 'The Star', 'Звезда', 1),
          mockHistoryCard('the_sun', 'The Sun', 'Солнце', 2),
        ].map((card) => ({ ...card, viewed: false })),
        deck_id: mockUser.selected_deck_id || 'rider-waite',
        back_id: mockUser.selected_back_id || 'midnight-stars',
        remaining_today: { count: 2, limit: 3 },
      };
      mockHistory = [reading, ...mockHistory];
      return { success: true, reading: structuredClone(reading) };
    },
    async revealCard(readingId, position) {
      await delay(260);
      const reading = mockHistory.find((item) => item.id === readingId);
      if (!reading) throw new Error('Reading not found');
      reading.cards = reading.cards.map((card) => card.position === position ? { ...card, viewed: true } : card);
      return { success: true, reading: structuredClone(reading) };
    },
    async requestInterpretation(readingId) {
      await delay(480);
      const reading = mockHistory.find((item) => item.id === readingId);
      if (!reading) throw new Error('Reading not found');
      const content = mockUser.language === 'en'
        ? 'The past asks to be understood, the present restores hope, and the future invites you to act openly and with confidence.'
        : 'Прошлое просит быть понятым, настоящее возвращает надежду, а будущее приглашает действовать открыто и уверенно.';
      const completedAt = new Date().toISOString();
      reading.ai = { status: 'completed', language: mockUser.language, content, completed_at: completedAt };
      return {
        success: true,
        request: {
          id: crypto.randomUUID(),
          reading_id: readingId,
          provider: 'mock',
          model: 'mock-oracle',
          prompt_version: 'mock-v1',
          status: 'completed',
          response_excerpt: content,
          requested_at: completedAt,
          completed_at: completedAt,
        },
      };
    },
  },
  admin: mockAdminApi,
};
