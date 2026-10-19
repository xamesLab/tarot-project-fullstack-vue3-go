import { render, screen } from '@testing-library/react';
import { createMemoryRouter, MemoryRouter, Route, Routes } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { api } from '@/api';
import { TOKEN_STORAGE_KEY } from '@/api/client';
import { App } from '@/app/App';
import { KeeperRoute, routes } from '@/app/router';
import { useAuthStore } from '@/features/auth/store';
import { usePreferencesStore } from '@/features/preferences/store';

const keeper = { id: 'keeper-1', username: 'keeper', email: 'keeper@example.com', role: 'keeper' as const, status: 'active' as const, timezone: 'Europe/Moscow', language: 'ru' as const, theme: 'dark' as const };
describe('admin route', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    usePreferencesStore.setState({ theme: 'dark', language: 'ru' });
    useAuthStore.setState({ user: keeper, token: 'jwt', status: 'authenticated', hydrated: true });
  });

  it('renders the keeper-only page, exposes all tabs, and reads the selected tab from the URL', async () => {
    const getUsers = vi.spyOn(api.admin, 'getUsers').mockResolvedValue({
      success: true,
      users: [{ ...keeper, created_at: '2026-01-01T00:00:00Z' }],
      pagination: { page: 1, limit: 20, total: 1, total_pages: 1 },
    });
    const router = createMemoryRouter(routes, { initialEntries: ['/admin?tab=users'] });
    render(<App router={router} />);

    expect(await screen.findByRole('heading', { name: 'Панель управления' })).toBeInTheDocument();
    expect(screen.getAllByRole('tab')).toHaveLength(5);
    expect(screen.getByRole('link', { name: 'Админка' })).toHaveAttribute('aria-current', 'page');

    expect(screen.getByRole('tab', { name: 'Пользователи' })).toHaveAttribute('aria-selected', 'true');
    expect(await screen.findByText('keeper@example.com')).toBeInTheDocument();
    expect(getUsers).toHaveBeenCalledWith({ search: '', role: 'all', status: 'all', page: 1, limit: 20 });
  });

  it('redirects an authenticated non-keeper away from protected admin content', async () => {
    useAuthStore.setState({ user: { ...keeper, role: 'master' }, token: 'jwt', status: 'authenticated', hydrated: true });
    render(
      <MemoryRouter initialEntries={['/admin']}>
        <Routes>
          <Route path="/profile" element={<div>Safe profile</div>} />
          <Route path="/admin" element={<KeeperRoute><div>Secret admin</div></KeeperRoute>} />
        </Routes>
      </MemoryRouter>,
    );
    expect(await screen.findByText('Safe profile')).toBeInTheDocument();
    expect(screen.queryByText('Secret admin')).not.toBeInTheDocument();
  });
});
