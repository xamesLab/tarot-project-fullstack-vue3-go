import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { App } from '@/app/App';
import { api } from '@/api';
import { router } from '@/app/router';
import { usePreferencesStore } from '@/features/preferences/store';
import { useAuthStore } from '@/features/auth/store';
import { TOKEN_STORAGE_KEY } from '@/api/client';

describe('NEBULA-TAROT smoke', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    usePreferencesStore.setState({ theme: 'dark', language: 'ru' });
    useAuthStore.setState({ user: null, token: '', status: 'guest', hydrated: true });
    window.history.replaceState({}, '', '/');
    const backendReading = { type: 'card_of_day' as const, accounting_date: '2026-07-03', generated_at: '2026-07-03T10:00:00Z', cards: [{ key: 'seven_of_pentacles', arcana: 'minor' as const, type: 'Pentacles', type_ru: 'пентакли', name: 'Seven of Pentacles', name_ru: 'Семёрка Пентаклей', interpretation: 'Patience brings results.', interpretation_ru: 'Терпение приносит результат.', tags: 'Patience, results.', tags_ru: 'Терпение, результат.', image_url: '/uploads/live-face.png', back_image_url: '/uploads/live-back.png', position: 0, reversed: false, viewed: false }] };
    vi.spyOn(api.dailyCard, 'generateGuest').mockResolvedValue({ success: true, source: 'guest', reading: backendReading, claim_token: 'claim' });
    vi.spyOn(api.dailyCard, 'resolve').mockResolvedValue({ success: true, source: 'existing', reading: { ...backendReading, id: 'reading-1' } });
    vi.spyOn(console, 'info').mockImplementation(() => undefined);
    vi.spyOn(console, 'error').mockImplementation(() => undefined);
  });

  it('keeps the daily page public and opens registration from gated navigation', async () => {
    const { container } = render(<App />);
    expect(screen.getByLabelText('NEBULA-TAROT')).toBeInTheDocument();
    expect(screen.getByText('Карта дня (карта гостя)')).toBeInTheDocument();
    const tooltipButton = screen.getByRole('button', { name: 'войдите чтобы увидеть вашу личную карту' });
    fireEvent.focus(tooltipButton);
    expect(screen.getByRole('tooltip')).toHaveTextContent('войдите чтобы увидеть вашу личную карту');
    const cardFace = await screen.findByRole('img', { name: 'Семёрка Пентаклей' });
    expect(new URL(cardFace.getAttribute('src')!).pathname).toBe('/uploads/live-face.png');
    const detailsShell = container.querySelector<HTMLElement>('.card-details-shell')!;
    const detailsBody = container.querySelector<HTMLElement>('.card-details-body')!;
    const expandButton = container.querySelector<HTMLButtonElement>('.expand-button')!;
    expect(detailsShell).toHaveClass('concealed');
    expect(detailsBody).toHaveAttribute('aria-hidden', 'true');
    expect(expandButton).toBeDisabled();
    expect(expandButton).toHaveAttribute('tabindex', '-1');
    expect(screen.getByText('Переверните карту, чтобы открыть описание')).toBeInTheDocument();

    const flipButton = screen.getByRole('button', { name: 'Нажмите, чтобы перевернуть' });
    fireEvent.click(flipButton);
    expect(detailsShell).toHaveClass('revealed');
    expect(detailsBody).toHaveAttribute('aria-hidden', 'false');
    expect(expandButton).toBeEnabled();
    expect(await screen.findByText('Семёрка Пентаклей')).toBeInTheDocument();
    expect(screen.getByText('Пентакли')).toBeInTheDocument();
    expect(screen.getByText('Младший Аркан · Прямое положение')).toBeInTheDocument();
    expect(screen.getByText('Терпение')).toBeInTheDocument();
    expect(screen.getByText('результат')).toBeInTheDocument();
    expect(screen.getByText('Терпение приносит результат.')).toBeInTheDocument();
    expect(screen.queryByText('Верховная Жрица')).not.toBeInTheDocument();
    fireEvent.click(expandButton);
    expect(await screen.findByText('Любовь & Отношения')).toBeInTheDocument();
    fireEvent.click(flipButton);
    expect(detailsShell).toHaveClass('concealed');
    expect(screen.queryByText('Любовь & Отношения')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Карта дня' }));
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Расклады' }));
    expect(await screen.findByRole('dialog')).toHaveTextContent('Присоединиться к NEBULA-TAROT');
  });

  it('opens login from the guest pill and switches to registration', async () => {
    render(<App />);
    fireEvent.click(screen.getByRole('button', { name: /Гость/ }));
    expect(await screen.findByRole('dialog')).toHaveTextContent('Войти в NEBULA-TAROT');
    fireEvent.click(screen.getByRole('button', { name: 'Нет аккаунта? Создать профиль' }));
    expect(screen.getByRole('dialog')).toHaveTextContent('Присоединиться к NEBULA-TAROT');
    expect(screen.getByLabelText('Имя пользователя (необязательно)')).toBeInTheDocument();
    expect(screen.getByLabelText('Часовой пояс')).toBeInTheDocument();
  });

  it('renders authenticated user data and logs out', () => {
    const user = { id: '1', username: 'alice', email: 'alice@example.com', role: 'adept' as const, status: 'active' as const, timezone: 'UTC', language: 'ru' as const, theme: 'dark' as const };
    localStorage.setItem(TOKEN_STORAGE_KEY, 'token');
    useAuthStore.setState({ user, token: 'token', status: 'authenticated', hydrated: true });
    const { container } = render(<App />);
    expect(container.querySelector('.section-label')).toHaveTextContent(/^Карта дня$/);
    expect(screen.queryByText(/\(карта гостя\)/)).not.toBeInTheDocument();
    expect(screen.getByText('alice')).toBeInTheDocument();
    expect(screen.getByText(/Адепт/)).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Выйти' }));
    expect(useAuthStore.getState().user).toBeNull();
    expect(localStorage.getItem(TOKEN_STORAGE_KEY)).toBeNull();
  });

  it('switches language and theme', () => {
    render(<App />);
    fireEvent.click(screen.getByRole('button', { name: 'EN' }));
    expect(screen.getByRole('button', { name: 'Card of the Day' })).toBeInTheDocument();
    expect(screen.getByText('Card of the Day (guest card)')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'sign in to see your personal card' })).toBeInTheDocument();
    expect(screen.getByText('Flip the card to reveal its meaning')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Switch to light theme' }));
    expect(usePreferencesStore.getState().theme).toBe('light');
  });

  it('registers the UI Kit route in development', async () => {
    await router.navigate('/ui-kit');
    render(<App />);
    expect(await screen.findByText('Библиотека компонентов ARCANA')).toBeInTheDocument();
    expect(screen.getByText('Основы — Палитра')).toBeInTheDocument();
  });
});
