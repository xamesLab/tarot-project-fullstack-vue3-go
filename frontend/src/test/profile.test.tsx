import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { createMemoryRouter, MemoryRouter, Route, Routes } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { api } from '@/api';
import type { ReadingHistoryItemDto, UserDto } from '@/api/types';
import { App } from '@/app/App';
import { ProtectedRoute, routes } from '@/app/router';
import { TOKEN_STORAGE_KEY } from '@/api/client';
import { useAuthStore } from '@/features/auth/store';
import { usePreferencesStore } from '@/features/preferences/store';

const sessionUser: UserDto = {
  id: 'user-1',
  username: 'alice',
  email: 'alice@example.com',
  role: 'adept',
  status: 'active',
  timezone: 'Europe/Moscow',
  language: 'ru',
  theme: 'dark',
};

const profileUser: UserDto = {
  ...sessionUser,
  selected_deck_id: 'celestial',
  selected_back_id: 'lunar-gold',
  created_at: '2026-07-03T10:00:00Z',
  last_login_at: '2026-07-18T08:30:00Z',
};

const deckCatalog = {
  success: true as const,
  default_deck_id: 'rider-waite',
  decks: [
    { id: 'rider-waite', name: { ru: 'Таро Райдера — Уэйта', en: 'Rider–Waite Tarot' }, description: { ru: '', en: '' }, author: 'A. E. Waite', year: 1909, tags: [], status: 'done', card_count: 78, back_image_path: '/deck.png' },
    { id: 'celestial', name: { ru: 'Небесное Таро', en: 'Celestial Tarot' }, description: { ru: '', en: '' }, author: 'Arcana', year: 2026, tags: [], status: 'done', card_count: 78, back_image_path: '/celestial.png' },
  ],
};

const cardBackCatalog = {
  success: true as const,
  default_back_id: 'midnight-stars',
  card_backs: [
    { id: 'midnight-stars', name: { ru: 'Полуночные звёзды', en: 'Midnight Stars' }, description: { ru: '', en: '' }, tags: [], status: 'done', image_path: '/stars.png' },
    { id: 'lunar-gold', name: { ru: 'Лунное золото', en: 'Lunar Gold' }, description: { ru: '', en: '' }, tags: [], status: 'done', image_path: '/lunar.png' },
  ],
};

const historyCard = (key: string, name: string, nameRu: string, position: number) => ({
  key,
  arcana: 'major' as const,
  type: 'Major Arcana',
  type_ru: 'Старший Аркан',
  name,
  name_ru: nameRu,
  interpretation: 'A quiet answer is taking shape.',
  interpretation_ru: 'Тихий ответ уже складывается внутри вас.',
  tags: 'intuition',
  tags_ru: 'интуиция',
  image_url: `/cards/${key}.png`,
  back_image_url: '/backs/default.png',
  position,
  reversed: false,
  viewed: true,
});

const historyReadings: ReadingHistoryItemDto[] = [
  {
    id: 'daily-1',
    type: 'card_of_day',
    accounting_date: '2026-07-18',
    generated_at: '2026-07-18T05:30:00Z',
    cards: [historyCard('the_high_priestess', 'The High Priestess', 'Верховная Жрица', 0)],
    ai: { status: 'completed', language: 'ru', content: 'Сегодня особенно важно доверять интуиции.' },
  },
  {
    id: 'spread-1',
    type: 'three_card',
    accounting_date: '2026-07-17',
    generated_at: '2026-07-17T18:10:00Z',
    cards: [
      historyCard('the_moon', 'The Moon', 'Луна', 0),
      historyCard('the_star', 'The Star', 'Звезда', 1),
      historyCard('the_sun', 'The Sun', 'Солнце', 2),
    ],
  },
];

const historyResponse = (readings: ReadingHistoryItemDto[], page = 1, totalPages = 1) => ({
  success: true as const,
  readings,
  pagination: { page, limit: 10, total: readings.length, total_pages: totalPages },
  retention_info: { role: 'adept' as const, retention_days: 90, oldest_date: '2026-04-19' },
});

describe('profile route', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    usePreferencesStore.setState({ theme: 'dark', language: 'ru' });
    useAuthStore.setState({ user: null, token: '', status: 'guest', hydrated: true });
  });

  it('loads current API data and renders both profile tabs', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    useAuthStore.setState({ user: sessionUser, token: 'jwt', status: 'authenticated', hydrated: true });
    const getCurrentUser = vi.spyOn(api.auth, 'getCurrentUser').mockResolvedValue({ success: true, user: profileUser });
    const getDecks = vi.spyOn(api.catalog, 'getAvailableDecks').mockResolvedValue(deckCatalog);
    const getCardBacks = vi.spyOn(api.catalog, 'getAvailableCardBacks').mockResolvedValue(cardBackCatalog);
    const getHistory = vi.spyOn(api.history, 'getHistory').mockResolvedValue(historyResponse([]));
    const updateProfile = vi.spyOn(api.profile, 'updateProfile').mockImplementation(async (payload) => ({
      success: true,
      user: { ...profileUser, ...payload },
    }));
    const router = createMemoryRouter(routes, { initialEntries: ['/profile'] });

    render(<App router={router} />);

    expect(await screen.findByRole('heading', { name: 'Личный кабинет' })).toBeInTheDocument();
    await waitFor(() => expect(getCurrentUser).toHaveBeenCalledTimes(1));
    expect(screen.getAllByText('alice').length).toBeGreaterThanOrEqual(2);
    expect(screen.getByText('alice@example.com')).toBeInTheDocument();
    expect(screen.getByText('28 мая 1994')).toBeInTheDocument();
    expect(screen.getByText('Санкт-Петербург')).toBeInTheDocument();
    expect(screen.getByText(/с июля 2026/)).toBeInTheDocument();
    expect(screen.getByText(/Адепт · Europe\/Moscow/)).toBeInTheDocument();
    expect(screen.getByText('214')).toBeInTheDocument();
    expect(screen.getByText(/Верховная Жрица — ваша карта души/)).toBeInTheDocument();

    const profileLink = screen.getByRole('link', { name: 'Профиль' });
    expect(profileLink).toHaveAttribute('href', '/profile');
    expect(profileLink).toHaveAttribute('aria-current', 'page');

    const settingsTab = screen.getByRole('tab', { name: 'Настройки' });
    const historyTab = screen.getByRole('tab', { name: 'История раскладов' });
    expect(settingsTab).toHaveAttribute('aria-selected', 'true');
    expect(await screen.findByRole('heading', { name: 'Аккаунт' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Язык и оформление' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Колода и рубашка' })).toBeInTheDocument();
    expect(screen.getByRole('tabpanel')).not.toBeEmptyDOMElement();
    await waitFor(() => {
      expect(getDecks).toHaveBeenCalledTimes(1);
      expect(getCardBacks).toHaveBeenCalledTimes(1);
    });
    expect(screen.getByLabelText('Колода')).toHaveValue('celestial');
    expect(screen.getByLabelText('Рубашка колоды')).toHaveValue('lunar-gold');

    fireEvent.change(screen.getByLabelText('Имя пользователя'), { target: { value: 'alice_2' } });
    fireEvent.click(screen.getByRole('button', { name: 'Сохранить профиль' }));
    await waitFor(() => expect(updateProfile).toHaveBeenCalledWith({ username: 'alice_2', timezone: 'Europe/Moscow' }));
    expect(useAuthStore.getState().user?.username).toBe('alice_2');
    expect(await screen.findByText('Профиль сохранён — карты помнят.')).toBeInTheDocument();

    fireEvent.keyDown(settingsTab, { key: 'ArrowRight' });
    expect(historyTab).toHaveAttribute('aria-selected', 'true');
    expect(historyTab).toHaveFocus();
    expect(await screen.findByRole('heading', { name: 'Архив раскладов' })).toBeInTheDocument();
    await waitFor(() => expect(getHistory).toHaveBeenCalledWith({ type: 'all', page: 1, limit: 10 }));
    expect(screen.getByText('История пока пуста')).toBeInTheDocument();
  });

  it('filters daily cards and other readings with independent checkboxes and expands details', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    useAuthStore.setState({ user: sessionUser, token: 'jwt', status: 'authenticated', hydrated: true });
    vi.spyOn(api.auth, 'getCurrentUser').mockResolvedValue({ success: true, user: profileUser });
    vi.spyOn(api.catalog, 'getAvailableDecks').mockResolvedValue(deckCatalog);
    vi.spyOn(api.catalog, 'getAvailableCardBacks').mockResolvedValue(cardBackCatalog);
    const getHistory = vi.spyOn(api.history, 'getHistory').mockImplementation(async ({ type }) => {
      const readings = type === 'all' ? historyReadings : historyReadings.filter((reading) => reading.type === type);
      return historyResponse(readings);
    });
    const router = createMemoryRouter(routes, { initialEntries: ['/profile'] });

    render(<App router={router} />);

    fireEvent.click(await screen.findByRole('tab', { name: 'История раскладов' }));
    const panel = screen.getByRole('tabpanel');
    expect(await within(panel).findByText('Верховная Жрица')).toBeInTheDocument();
    expect(within(panel).getByText('Луна')).toBeInTheDocument();

    const dailyFilter = within(panel).getByRole('checkbox', { name: 'Карты дня' });
    const otherFilter = within(panel).getByRole('checkbox', { name: 'Другие расклады' });
    expect(dailyFilter).toHaveAttribute('aria-checked', 'true');
    expect(otherFilter).toHaveAttribute('aria-checked', 'true');

    fireEvent.click(dailyFilter);
    await waitFor(() => expect(getHistory).toHaveBeenCalledWith({ type: 'three_card', page: 1, limit: 10 }));
    await waitFor(() => expect(within(panel).queryByText('Верховная Жрица')).not.toBeInTheDocument());
    expect(within(panel).getByText('Луна')).toBeInTheDocument();

    fireEvent.click(otherFilter);
    expect(await within(panel).findByText('Все фильтры выключены')).toBeInTheDocument();

    fireEvent.click(dailyFilter);
    await waitFor(() => expect(getHistory).toHaveBeenCalledWith({ type: 'card_of_day', page: 1, limit: 10 }));
    fireEvent.click(await within(panel).findByRole('button', { name: 'Показать детали' }));
    expect(await within(panel).findByText('Тихий ответ уже складывается внутри вас.')).toBeInTheDocument();
    expect(within(panel).getByText('Сегодня особенно важно доверять интуиции.')).toBeInTheDocument();
  });

  it('persists deck and appearance settings and synchronizes preferences', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    useAuthStore.setState({ user: sessionUser, token: 'jwt', status: 'authenticated', hydrated: true });
    vi.spyOn(api.auth, 'getCurrentUser').mockResolvedValue({ success: true, user: profileUser });
    vi.spyOn(api.catalog, 'getAvailableDecks').mockResolvedValue(deckCatalog);
    vi.spyOn(api.catalog, 'getAvailableCardBacks').mockResolvedValue(cardBackCatalog);
    const updateSettings = vi.spyOn(api.profile, 'updateSettings').mockImplementation(async (payload) => ({
      success: true,
      user: { ...profileUser, ...payload },
    }));
    const router = createMemoryRouter(routes, { initialEntries: ['/profile'] });

    render(<App router={router} />);

    await screen.findByRole('heading', { name: 'Колода и рубашка' });
    await waitFor(() => expect(screen.getByLabelText('Колода')).not.toBeDisabled());
    fireEvent.change(screen.getByLabelText('Колода'), { target: { value: 'rider-waite' } });
    fireEvent.change(screen.getByLabelText('Рубашка колоды'), { target: { value: 'midnight-stars' } });
    fireEvent.click(screen.getByRole('button', { name: 'Сохранить колоду' }));
    await waitFor(() => expect(updateSettings).toHaveBeenCalledWith({
      selected_deck_id: 'rider-waite',
      selected_back_id: 'midnight-stars',
    }));

    fireEvent.click(within(screen.getByRole('group', { name: 'Язык интерфейса' })).getByRole('button', { name: 'English' }));
    fireEvent.click(within(screen.getByRole('group', { name: 'Тема' })).getByRole('button', { name: 'Светлая' }));
    fireEvent.click(screen.getByRole('button', { name: 'Сохранить оформление' }));
    await waitFor(() => expect(updateSettings).toHaveBeenCalledWith({ language: 'en', theme: 'light' }));
    expect(usePreferencesStore.getState()).toMatchObject({ language: 'en', theme: 'light' });
    expect(document.documentElement).toHaveClass('light');
  });

  it('redirects guests to the public page', async () => {
    render(
      <MemoryRouter initialEntries={['/profile']}>
        <Routes>
          <Route path="/" element={<div>Public page</div>} />
          <Route path="/profile" element={<ProtectedRoute><div>Private profile</div></ProtectedRoute>} />
        </Routes>
      </MemoryRouter>,
    );

    expect(await screen.findByText('Public page')).toBeInTheDocument();
    expect(screen.queryByText('Private profile')).not.toBeInTheDocument();
  });
});
