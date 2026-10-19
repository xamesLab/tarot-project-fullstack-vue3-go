import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { createMemoryRouter, MemoryRouter, Route, Routes } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { api } from '@/api';
import { TOKEN_STORAGE_KEY } from '@/api/client';
import type { BackendReadingDto, UserDto } from '@/api/types';
import { App } from '@/app/App';
import { ProtectedRoute, routes } from '@/app/router';
import { getAccountingDate } from '@/features/daily-card/bootstrap';
import { useAuthStore } from '@/features/auth/store';
import { usePreferencesStore } from '@/features/preferences/store';

const user: UserDto = {
  id: 'spread-user',
  username: 'alice',
  email: 'alice@example.com',
  role: 'adept',
  status: 'active',
  timezone: 'Europe/Moscow',
  language: 'ru',
  theme: 'dark',
};

const card = (position: number, name: string, nameRu: string) => ({
  key: name.toLowerCase().replace(/ /g, '_'),
  arcana: 'major' as const,
  type: 'Major Arcana',
  type_ru: 'Старший Аркан',
  name,
  name_ru: nameRu,
  interpretation: `${name} offers its lesson.`,
  interpretation_ru: `${nameRu} открывает значение этой позиции.`,
  tags: 'path, clarity',
  tags_ru: 'путь, ясность',
  image_url: `/cards/${position}.png`,
  back_image_url: '/backs/default.png',
  position,
  reversed: false,
  viewed: false,
});

function emptyHistory() {
  return {
    success: true as const,
    readings: [],
    pagination: { page: 1, limit: 20, total: 0, total_pages: 1 },
    retention_info: { role: user.role, retention_days: 90, oldest_date: '2026-01-01' },
  };
}

describe('three-card spreads page', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    usePreferencesStore.setState({ theme: 'dark', language: 'ru' });
    useAuthStore.setState({ user, token: 'jwt', status: 'authenticated', hydrated: true });
  });

  it('protects the route from guests', async () => {
    localStorage.clear();
    useAuthStore.setState({ user: null, token: '', status: 'guest', hydrated: true });
    render(
      <MemoryRouter initialEntries={['/spreads']}>
        <Routes>
          <Route path="/" element={<div>Public page</div>} />
          <Route path="/spreads" element={<ProtectedRoute><div>Private spread</div></ProtectedRoute>} />
        </Routes>
      </MemoryRouter>,
    );

    expect(await screen.findByText('Public page')).toBeInTheDocument();
    expect(screen.queryByText('Private spread')).not.toBeInTheDocument();
  });

  it('creates a spread, reveals all cards through the API, and requests the shared interpretation', async () => {
    vi.spyOn(api.history, 'getHistory').mockResolvedValue(emptyHistory());
    let reading: BackendReadingDto = {
      id: 'spread-today',
      type: 'three_card',
      accounting_date: getAccountingDate(new Date(), user.timezone),
      generated_at: new Date().toISOString(),
      cards: [card(0, 'The Moon', 'Луна'), card(1, 'The Star', 'Звезда'), card(2, 'The Sun', 'Солнце')],
      remaining_today: { count: 1, limit: 3 },
    };
    const generate = vi.spyOn(api.spreads, 'generateThreeCard').mockResolvedValue({ success: true, reading });
    const reveal = vi.spyOn(api.spreads, 'revealCard').mockImplementation(async (_readingId, position) => {
      reading = { ...reading, cards: reading.cards.map((item) => item.position === position ? { ...item, viewed: true } : item) };
      return { success: true, reading };
    });
    const interpretation = vi.spyOn(api.spreads, 'requestInterpretation').mockResolvedValue({
      success: true,
      request: {
        id: 'ask-1',
        reading_id: 'spread-today',
        provider: 'test',
        model: 'oracle-test',
        prompt_version: 'v1',
        status: 'completed',
        response_excerpt: 'Карты ведут из тумана к ясности.',
        requested_at: new Date().toISOString(),
        completed_at: new Date().toISOString(),
      },
    });
    const router = createMemoryRouter(routes, { initialEntries: ['/spreads'] });
    render(<App router={router} />);

    expect(await screen.findByRole('heading', { name: 'Карты ждут вашего вопроса' })).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Сделать расклад' }));
    await waitFor(() => expect(generate).toHaveBeenCalledWith({ timezone: user.timezone }));

    const revealAll = await screen.findByRole('button', { name: 'Раскрыть все карты' });
    fireEvent.click(revealAll);

    await waitFor(() => expect(reveal).toHaveBeenCalledTimes(3));
    expect(reveal.mock.calls.map((call) => call[1])).toEqual([0, 1, 2]);
    await waitFor(() => expect(interpretation).toHaveBeenCalledWith('spread-today'));
    expect(await screen.findByRole('button', { name: 'Развернуть толкование' })).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Развернуть толкование' }));
    expect(await screen.findByText('Карты ведут из тумана к ясности.')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Луна' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Звезда' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Солнце' })).toBeInTheDocument();
  });
});
