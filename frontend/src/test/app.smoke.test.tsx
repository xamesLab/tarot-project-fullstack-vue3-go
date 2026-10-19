import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { createMemoryRouter } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { api } from '@/api';
import { App } from '@/app/App';
import { routes } from '@/app/router';
import { TOKEN_STORAGE_KEY } from '@/api/client';
import { useAuthStore } from '@/features/auth/store';
import { usePreferencesStore } from '@/features/preferences/store';

const backendReading = {
  type: 'card_of_day' as const,
  accounting_date: '2026-07-03',
  generated_at: '2026-07-03T10:00:00Z',
  cards: [{
    key: 'seven_of_pentacles',
    arcana: 'minor' as const,
    type: 'Pentacles',
    type_ru: 'Пентакли',
    name: 'Seven of Pentacles',
    name_ru: 'Семёрка Пентаклей',
    interpretation: 'Patience brings results.',
    interpretation_ru: 'Терпение приносит результат.',
    tags: 'Patience, results.',
    tags_ru: 'Терпение, результат.',
    image_url: '/uploads/live-face.png',
    back_image_url: '/uploads/live-back.png',
    position: 0,
    reversed: false,
    viewed: false,
  }],
};

describe('NEBULA-TAROT smoke', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    usePreferencesStore.setState({ theme: 'dark', language: 'ru' });
    useAuthStore.setState({ user: null, token: '', status: 'guest', hydrated: true });
    window.history.replaceState({}, '', '/');
    vi.spyOn(api.dailyCard, 'generateGuest').mockResolvedValue({ success: true, source: 'guest', reading: backendReading, claim_token: 'claim' });
    vi.spyOn(api.dailyCard, 'resolve').mockResolvedValue({ success: true, source: 'existing', reading: { ...backendReading, id: 'reading-1' } });
    vi.spyOn(api.dailyCard, 'getGuestInterpretation').mockResolvedValue({ success: true, source: 'cache', interpretation: { general: 'AI general', message: 'AI day message', relationships: 'AI relationships', finances: 'AI finances', emotional_state: 'AI emotional state', career_business: 'AI career business' } });
    vi.spyOn(api.dailyCard, 'getUserInterpretation').mockResolvedValue({ success: true, source: 'reading', interpretation: { general: 'AI user general', message: 'AI user day message', relationships: 'AI user relationships', finances: 'AI user finances', emotional_state: 'AI user emotional state', career_business: 'AI user career business' } });
    vi.spyOn(console, 'info').mockImplementation(() => undefined);
    vi.spyOn(console, 'error').mockImplementation(() => undefined);
  });

  it('keeps the daily page public, requests AI only after reveal, and opens registration from gated navigation', async () => {
    const { container } = render(<App />);
    expect(screen.getByLabelText('NEBULA-TAROT')).toBeInTheDocument();
    expect(screen.getByText('Карта дня (карта гостя)')).toBeInTheDocument();
    const tooltipButton = screen.getByRole('button', { name: 'Войдите, чтобы увидеть вашу личную карту' });
    fireEvent.focus(tooltipButton);
    expect(screen.getByRole('tooltip')).toHaveTextContent('Войдите, чтобы увидеть вашу личную карту');
    const cardFace = await screen.findByRole('img', { name: 'Семёрка Пентаклей' });
    expect(new URL(cardFace.getAttribute('src')!).pathname).toBe('/uploads/live-face.png');

    const detailsShell = container.querySelector<HTMLElement>('.card-details-shell')!;
    const detailsBody = container.querySelector<HTMLElement>('.card-details-body')!;
    const expandButton = container.querySelector<HTMLButtonElement>('.expand-button')!;
    expect(detailsShell).toHaveClass('concealed');
    expect(detailsBody).toHaveAttribute('aria-hidden', 'true');
    expect(expandButton).toBeDisabled();
    expect(screen.getByText('Переверните карту, чтобы открыть описание')).toBeInTheDocument();
    expect(api.dailyCard.getGuestInterpretation).not.toHaveBeenCalled();

    const flipButton = screen.getByRole('button', { name: 'Нажмите, чтобы перевернуть' });
    fireEvent.click(flipButton);
    expect(detailsShell).toHaveClass('revealed');
    expect(detailsBody).toHaveAttribute('aria-hidden', 'false');
    expect(await screen.findByText('AI general')).toBeInTheDocument();
    expect(api.dailyCard.getGuestInterpretation).toHaveBeenCalledWith({ claim_token: 'claim', language: 'ru' });
    expect(screen.getByText('Пентакли')).toBeInTheDocument();
    expect(screen.getByText('Младший Аркан · Прямое положение')).toBeInTheDocument();
    expect(screen.getByText('Терпение')).toBeInTheDocument();
    expect(screen.getByText('результат')).toBeInTheDocument();
    await waitFor(() => expect(expandButton).toBeEnabled());
    expect(expandButton).toHaveAttribute('aria-expanded', 'true');
    expect(screen.getByText('AI day message')).toBeInTheDocument();
    expect(await screen.findByText('AI relationships')).toBeInTheDocument();
    expect(screen.getByText('AI finances')).toBeInTheDocument();
    expect(screen.getByText('AI emotional state')).toBeInTheDocument();
    expect(screen.getByText('AI career business')).toBeInTheDocument();

    fireEvent.click(flipButton);
    expect(detailsShell).toHaveClass('concealed');
    expect(screen.queryByText('AI relationships')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Карта дня' }));
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Расклады' }));
    expect(await screen.findByRole('dialog')).toBeInTheDocument();
  });

  it('replies with the oracle silence placeholder to every question without calling the API', async () => {
    const askOracle = vi.spyOn(api.oracle, 'ask');
    const { container } = render(<App />);
    const input = container.querySelector<HTMLInputElement>('.chat-input input')!;
    const sendButton = container.querySelector<HTMLButtonElement>('.chat-input button')!;
    const silenceReply = 'Оракул сейчас молчит';

    fireEvent.change(input, { target: { value: 'What should I focus on today?' } });
    fireEvent.click(sendButton);
    expect(await screen.findByText(silenceReply)).toBeInTheDocument();

    fireEvent.change(input, { target: { value: 'Should I change my plans?' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    expect(screen.getAllByText(silenceReply)).toHaveLength(2);
    expect(askOracle).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole('button', { name: 'EN' }));
    const englishInput = await screen.findByPlaceholderText('Ask the cards a question…');
    fireEvent.change(englishInput, { target: { value: 'What should I do next?' } });
    fireEvent.keyDown(englishInput, { key: 'Enter' });
    expect(await screen.findByText('The Oracle is silent right now.')).toBeInTheDocument();
    expect(askOracle).not.toHaveBeenCalled();
  });

  it('opens login from the guest pill and switches to registration', async () => {
    render(<App />);
    fireEvent.click(screen.getByRole('button', { name: /Гость/ }));
    expect(await screen.findByRole('dialog')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: /Создать профиль/ }));
    expect(screen.getByRole('dialog')).toBeInTheDocument();
    expect(screen.getByLabelText(/Имя пользователя/)).toBeInTheDocument();
    expect(screen.getByLabelText('Часовой пояс')).toBeInTheDocument();
  });

  it('opens the responsive site menu and keeps gated navigation available', async () => {
    render(<App />);
    const menuButton = screen.getByRole('button', { name: 'Открыть меню' });
    fireEvent.click(menuButton);

    const drawer = screen.getByRole('dialog', { name: 'Меню' });
    const navigation = within(drawer).getByRole('navigation', { name: 'Основная навигация' });
    expect(within(navigation).getByRole('button', { name: 'Карта дня' })).toHaveAttribute('aria-current', 'page');
    expect(menuButton).toHaveAttribute('aria-expanded', 'true');
    const themeControl = within(drawer).getByRole('group', { name: 'Тема' });
    expect(within(themeControl).getByRole('button', { name: 'Тёмная тема' })).toHaveAttribute('aria-pressed', 'true');
    fireEvent.click(within(themeControl).getByRole('button', { name: 'Светлая тема' }));
    expect(usePreferencesStore.getState().theme).toBe('light');

    fireEvent.keyDown(drawer, { key: 'Escape' });
    expect(screen.queryByRole('dialog', { name: 'Меню' })).not.toBeInTheDocument();
    expect(menuButton).toHaveAttribute('aria-expanded', 'false');

    fireEvent.click(menuButton);
    fireEvent.click(within(screen.getByRole('dialog', { name: 'Меню' })).getByRole('button', { name: 'Расклады' }));
    expect(screen.queryByRole('dialog', { name: 'Меню' })).not.toBeInTheDocument();
    expect(await screen.findByRole('dialog')).toBeInTheDocument();
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
    const router = createMemoryRouter(routes, {
      initialEntries: ['/ui-kit'],
      future: { v7_relativeSplatPath: true },
    });
    render(<App router={router} />);
    expect(await screen.findByRole('heading', { name: /Библиотека компонентов ARCANA/ })).toBeInTheDocument();
    const sectionNavigation = screen.getByRole('navigation', { name: 'Навигация по UI Kit' });
    expect(within(sectionNavigation).getAllByRole('link')).toHaveLength(16);
    expect(screen.getByRole('button', { name: 'Включить светлую тему' })).toBeInTheDocument();
    expect(document.querySelectorAll('.kit-section')).toHaveLength(16);

    const menuButton = document.querySelector<HTMLButtonElement>('.kit-menu-button');
    expect(menuButton).toHaveAttribute('aria-label', 'Открыть меню компонентов');
    if (!menuButton) throw new Error('UI Kit menu button was not rendered');
    fireEvent.click(menuButton);
    const drawer = screen.getByRole('dialog', { name: 'Компоненты' });
    expect(within(drawer).getAllByRole('link')).toHaveLength(16);
    expect(menuButton).toHaveAttribute('aria-expanded', 'true');
    fireEvent.keyDown(drawer, { key: 'Escape' });
    expect(screen.queryByRole('dialog', { name: 'Компоненты' })).not.toBeInTheDocument();
    expect(menuButton).toHaveAttribute('aria-expanded', 'false');
  });
});
