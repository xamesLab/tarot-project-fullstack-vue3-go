import type {
  AdminApi,
  AdminRolePolicyDto,
  AdminSettingsDto,
  AdminUserDto,
  CardBackDto,
  DeckDto,
  PaginationDto,
} from '@/api/types';

const pause = (milliseconds = 180) => new Promise((resolve) => setTimeout(resolve, milliseconds));

let users: AdminUserDto[] = [
  { id: 'mock-user', username: 'arcana_keeper', email: 'keeper@nebula.local', role: 'keeper', status: 'active', timezone: 'Europe/Moscow', language: 'ru', theme: 'dark', created_at: '2025-10-18T08:30:00Z', last_login_at: '2026-07-19T06:10:00Z' },
  { id: 'mock-seeker', username: 'luna', email: 'luna@example.com', role: 'seeker', status: 'active', timezone: 'Europe/Paris', language: 'ru', theme: 'dark', created_at: '2026-01-07T12:00:00Z', last_login_at: '2026-07-18T20:42:00Z' },
  { id: 'mock-master', username: 'orion', email: 'orion@example.com', role: 'master', status: 'active', timezone: 'America/New_York', language: 'en', theme: 'light', created_at: '2025-12-11T17:45:00Z', last_login_at: '2026-07-17T18:05:00Z' },
  { id: 'mock-disabled', username: 'selene', email: 'selene@example.com', role: 'adept', status: 'disabled', timezone: 'UTC', language: 'en', theme: 'dark', created_at: '2026-03-22T09:15:00Z' },
];

let decks: DeckDto[] = [
  { id: 'rider-waite', name: { ru: 'Таро Райдера — Уэйта', en: 'Rider–Waite Tarot' }, description: { ru: 'Классическая колода.', en: 'The classic tarot deck.' }, author: 'A. E. Waite', year: 1909, tags: ['classic'], status: 'done', card_count: 78, back_image_path: '/uploads/decks/rider-waite/back.png', created_at: '2025-10-18T08:30:00Z' },
  { id: 'celestial', name: { ru: 'Небесное Таро', en: 'Celestial Tarot' }, description: { ru: 'Колода с мотивами ночного неба.', en: 'A deck inspired by the night sky.' }, author: 'Arcana', year: 2026, tags: ['celestial'], status: 'processing', card_count: 78, back_image_path: '/uploads/decks/celestial/back.png', created_at: '2026-07-18T08:30:00Z' },
  { id: 'archived-deck', name: { ru: 'Лунный архив', en: 'Lunar Archive' }, description: { ru: 'Архивная колода.', en: 'An archived deck.' }, author: 'Nebula', year: 2024, tags: ['moon'], status: 'archived', card_count: 78, back_image_path: '/uploads/decks/archived/back.png', created_at: '2024-05-10T08:30:00Z' },
];

let cardBacks: CardBackDto[] = [
  { id: 'midnight-stars', name: { ru: 'Полуночные звёзды', en: 'Midnight Stars' }, description: { ru: 'Золотые звёзды на синем фоне.', en: 'Golden stars on midnight blue.' }, tags: ['stars'], status: 'done', image_path: '/uploads/card-backs/midnight-stars/back.png', created_at: '2025-10-18T08:30:00Z' },
  { id: 'archived-back', name: { ru: 'Старое золото', en: 'Old Gold' }, description: { ru: 'Архивная рубашка.', en: 'An archived card back.' }, tags: ['gold'], status: 'archived', image_path: '/uploads/card-backs/archived/back.png', created_at: '2024-06-14T08:30:00Z' },
];

let policies: AdminRolePolicyDto[] = (['seeker', 'adept', 'master', 'keeper'] as const).map((role, index) => ({
  id: `policy-${role}`,
  role,
  display_name: { ru: role, en: role },
  three_card_limit: [1, 3, 10, 100][index],
  ai_enabled: role !== 'seeker',
  ai_daily_limit: [0, 3, 10, 100][index],
  history_retention_days: [7, 30, 180, 3650][index],
}));

let settings: AdminSettingsDto = {
  id: 'global',
  maintenance_mode: false,
  maintenance_message: { ru: 'Мы скоро вернёмся.', en: 'We will be back shortly.' },
  ai_enabled: true,
  default_role: 'seeker',
  default_deck_id: 'rider-waite',
  default_back_id: 'midnight-stars',
  registration_enabled: true,
};

function paginate<T>(items: T[], page: number, limit: number): { items: T[]; pagination: PaginationDto } {
  const total = items.length;
  return {
    items: items.slice((page - 1) * limit, page * limit),
    pagination: { page, limit, total, total_pages: Math.max(1, Math.ceil(total / limit)) },
  };
}

export const mockAdminApi: AdminApi = {
  async getDashboard() {
    await pause();
    const activeDecks = decks.filter((deck) => deck.status !== 'archived');
    return {
      success: true,
      metrics: {
        users: {
          total: users.length,
          by_role: {
            seeker: users.filter((user) => user.role === 'seeker').length,
            adept: users.filter((user) => user.role === 'adept').length,
            master: users.filter((user) => user.role === 'master').length,
            keeper: users.filter((user) => user.role === 'keeper').length,
          },
          by_status: {
            active: users.filter((user) => user.status === 'active').length,
            disabled: users.filter((user) => user.status === 'disabled').length,
          },
        },
        readings: {
          today: { total: 34, card_of_day: 27, three_card: 7 },
          this_week: { total: 186, card_of_day: 142, three_card: 44 },
        },
        decks: {
          total: decks.length,
          by_status: {
            done: activeDecks.filter((deck) => deck.status === 'done').length,
            processing: activeDecks.filter((deck) => deck.status === 'processing').length,
            disabled: activeDecks.filter((deck) => deck.status === 'disabled').length,
            archived: decks.filter((deck) => deck.status === 'archived').length,
          },
        },
      },
    };
  },
  async getUsers({ search = '', role = 'all', status = 'all', page, limit }) {
    await pause();
    const needle = search.trim().toLowerCase();
    const filtered = users.filter((user) => (!needle || `${user.username} ${user.email}`.toLowerCase().includes(needle)) && (role === 'all' || user.role === role) && (status === 'all' || user.status === status));
    const result = paginate(filtered, page, limit);
    return { success: true, users: structuredClone(result.items), pagination: result.pagination };
  },
  async updateUserRole(userId, role) {
    await pause();
    users = users.map((user) => user.id === userId ? { ...user, role } : user);
    return { success: true, user: structuredClone(users.find((user) => user.id === userId)!) };
  },
  async updateUserStatus(userId, status) {
    await pause();
    users = users.map((user) => user.id === userId ? { ...user, status } : user);
    return { success: true, user: structuredClone(users.find((user) => user.id === userId)!) };
  },
  async getRolePolicies() {
    await pause();
    return { success: true, policies: structuredClone(policies) };
  },
  async updateRolePolicy(role, payload) {
    await pause();
    policies = policies.map((policy) => policy.role === role ? { ...policy, ...payload } : policy);
    return { success: true, policy: structuredClone(policies.find((policy) => policy.role === role)!) };
  },
  async getSettings() {
    await pause();
    return { success: true, settings: structuredClone(settings) };
  },
  async updateSettings(payload) {
    await pause();
    settings = { ...settings, ...payload };
    return { success: true, settings: structuredClone(settings) };
  },
  async getDecks({ scope, page, limit }) {
    await pause();
    const filtered = decks.filter((deck) => scope === 'archived' ? deck.status === 'archived' : deck.status !== 'archived');
    const result = paginate(filtered, page, limit);
    return { success: true, decks: structuredClone(result.items), pagination: result.pagination };
  },
  async createDeck(payload) {
    await pause(350);
    const deck: DeckDto = { id: crypto.randomUUID(), name: { en: String(payload.get('name_en') || ''), ru: String(payload.get('name_ru') || '') }, description: { en: String(payload.get('description_en') || ''), ru: String(payload.get('description_ru') || '') }, author: String(payload.get('author') || ''), year: Number(payload.get('year')) || 0, tags: String(payload.get('tags') || '').split(',').map((tag) => tag.trim()).filter(Boolean), status: 'done', card_count: 78, back_image_path: '', created_at: new Date().toISOString() };
    decks = [deck, ...decks];
    return { success: true, deck: structuredClone(deck), message: 'Deck created' };
  },
  async archiveDeck(deckId) { await pause(); decks = decks.map((deck) => deck.id === deckId ? { ...deck, status: 'archived' } : deck); return { success: true, deck: structuredClone(decks.find((deck) => deck.id === deckId)!) }; },
  async restoreDeck(deckId) { await pause(); decks = decks.map((deck) => deck.id === deckId ? { ...deck, status: 'done' } : deck); return { success: true, deck: structuredClone(decks.find((deck) => deck.id === deckId)!) }; },
  async deleteDeck(deckId) { await pause(); decks = decks.filter((deck) => deck.id !== deckId); return { success: true }; },
  async getCardBacks({ scope, page, limit }) {
    await pause();
    const filtered = cardBacks.filter((back) => scope === 'archived' ? back.status === 'archived' : back.status !== 'archived');
    const result = paginate(filtered, page, limit);
    return { success: true, card_backs: structuredClone(result.items), pagination: result.pagination };
  },
  async createCardBack(payload) {
    await pause(300);
    const cardBack: CardBackDto = { id: crypto.randomUUID(), name: { en: String(payload.get('name_en') || ''), ru: String(payload.get('name_ru') || '') }, description: { en: String(payload.get('description_en') || ''), ru: String(payload.get('description_ru') || '') }, tags: String(payload.get('tags') || '').split(',').map((tag) => tag.trim()).filter(Boolean), status: 'done', image_path: '', created_at: new Date().toISOString() };
    cardBacks = [cardBack, ...cardBacks];
    return { success: true, card_back: structuredClone(cardBack), message: 'Card back created' };
  },
  async archiveCardBack(cardBackId) { await pause(); cardBacks = cardBacks.map((back) => back.id === cardBackId ? { ...back, status: 'archived' } : back); return { success: true, card_back: structuredClone(cardBacks.find((back) => back.id === cardBackId)!) }; },
  async restoreCardBack(cardBackId) { await pause(); cardBacks = cardBacks.map((back) => back.id === cardBackId ? { ...back, status: 'done' } : back); return { success: true, card_back: structuredClone(cardBacks.find((back) => back.id === cardBackId)!) }; },
  async deleteCardBack(cardBackId) { await pause(); cardBacks = cardBacks.filter((back) => back.id !== cardBackId); return { success: true }; },
};
