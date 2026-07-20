import { beforeEach, describe, expect, it, vi } from 'vitest';
import { httpApi } from '@/api/http';
import { ApiError, AUTH_UNAUTHORIZED_EVENT, TOKEN_STORAGE_KEY } from '@/api/client';
import { useAuthStore } from '@/features/auth/store';
import { usePreferencesStore } from '@/features/preferences/store';

const user = { id: '1', username: 'alice', email: 'alice@example.com', role: 'master' as const, status: 'active' as const, timezone: 'Europe/Moscow', language: 'en' as const, theme: 'light' as const };

describe('auth HTTP and session lifecycle', () => {
  beforeEach(() => {
    localStorage.clear();
    usePreferencesStore.setState({ theme: 'dark', language: 'ru' });
    useAuthStore.setState({ user: null, token: '', status: 'guest', hydrated: false });
    vi.restoreAllMocks();
  });

  it('sends login payload and stores the returned session', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ success: true, user, token: 'jwt' }), { status: 200, headers: { 'Content-Type': 'application/json' } }));
    await useAuthStore.getState().login({ email: 'alice@example.com', password: 'Secret123' });
    expect(fetchMock).toHaveBeenCalledWith('http://localhost:8888/api/auth/login', expect.objectContaining({ method: 'POST' }));
    expect(localStorage.getItem(TOKEN_STORAGE_KEY)).toBe('jwt');
    expect(useAuthStore.getState().user?.username).toBe('alice');
    expect(usePreferencesStore.getState()).toMatchObject({ theme: 'light', language: 'en' });
  });

  it('adds the bearer token when restoring the current user', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ success: true, user }), { status: 200, headers: { 'Content-Type': 'application/json' } }));
    await httpApi.auth.getCurrentUser();
    const request = fetchMock.mock.calls[0][1] as RequestInit;
    expect(new Headers(request.headers).get('Authorization')).toBe('Bearer jwt');
  });

  it('uses the protected profile, catalog, and history endpoints', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    const ok = (body: unknown) => new Response(JSON.stringify(body), { status: 200, headers: { 'Content-Type': 'application/json' } });
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(ok({ success: true, user }))
      .mockResolvedValueOnce(ok({ success: true, user }))
      .mockResolvedValueOnce(ok({ success: true, decks: [], default_deck_id: '' }))
      .mockResolvedValueOnce(ok({ success: true, card_backs: [], default_back_id: '' }))
      .mockResolvedValueOnce(ok({ success: true, readings: [], pagination: {}, retention_info: {} }));

    await httpApi.profile.updateProfile({ username: 'alice_2', timezone: 'UTC' });
    await httpApi.profile.updateSettings({ language: 'ru', theme: 'dark' });
    await httpApi.catalog.getAvailableDecks();
    await httpApi.catalog.getAvailableCardBacks();
    await httpApi.history.getHistory({ type: 'three_card', page: 2, limit: 10 });

    expect(fetchMock).toHaveBeenNthCalledWith(1, 'http://localhost:8888/api/me/profile', expect.objectContaining({
      method: 'PATCH',
      body: JSON.stringify({ username: 'alice_2', timezone: 'UTC' }),
    }));
    expect(fetchMock).toHaveBeenNthCalledWith(2, 'http://localhost:8888/api/me/settings', expect.objectContaining({
      method: 'PATCH',
      body: JSON.stringify({ language: 'ru', theme: 'dark' }),
    }));
    expect(fetchMock).toHaveBeenNthCalledWith(3, 'http://localhost:8888/api/decks/available', expect.any(Object));
    expect(fetchMock).toHaveBeenNthCalledWith(4, 'http://localhost:8888/api/card-backs/available', expect.any(Object));
    expect(fetchMock).toHaveBeenNthCalledWith(5, 'http://localhost:8888/api/reading/history?type=three_card&page=2&limit=10', expect.any(Object));
    const headers = new Headers((fetchMock.mock.calls[0][1] as RequestInit).headers);
    expect(headers.get('Authorization')).toBe('Bearer jwt');
  });

  it('uses the protected three-card spread and interpretation endpoints', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    const ok = (body: unknown) => new Response(JSON.stringify(body), { status: 200, headers: { 'Content-Type': 'application/json' } });
    const reading = { id: 'reading-1', type: 'three_card', accounting_date: '2026-07-19', generated_at: '2026-07-19T10:00:00Z', cards: [] };
    const fetchMock = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(ok({ success: true, reading }))
      .mockResolvedValueOnce(ok({ success: true, reading }))
      .mockResolvedValueOnce(ok({ success: true, request: { id: 'ask-1', reading_id: 'reading-1', status: 'completed' } }));

    await httpApi.spreads.generateThreeCard({ timezone: 'Europe/Moscow' });
    await httpApi.spreads.revealCard('reading-1', 2);
    await httpApi.spreads.requestInterpretation('reading-1');

    expect(fetchMock).toHaveBeenNthCalledWith(1, 'http://localhost:8888/api/reading/generate/three-card', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ timezone: 'Europe/Moscow' }),
    }));
    expect(fetchMock).toHaveBeenNthCalledWith(2, 'http://localhost:8888/api/reading/reading-1/cards/2/reveal', expect.objectContaining({ method: 'POST' }));
    expect(fetchMock).toHaveBeenNthCalledWith(3, 'http://localhost:8888/api/ask', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ reading_id: 'reading-1' }),
    }));
  });

  it('parses structured backend errors', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ success: false, error: 'validation_failed', message: 'Email is already registered', details: { field: 'email', reason: 'already_registered' } }), { status: 400, headers: { 'Content-Type': 'application/json' } }));
    await expect(httpApi.auth.register({ email: 'alice@example.com', password: 'Secret123', timezone: 'UTC' })).rejects.toMatchObject({ status: 400, code: 'validation_failed', details: { field: 'email', reason: 'already_registered' } } satisfies Partial<ApiError>);
  });

  it('keeps the session when an authenticated admin request is forbidden', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    const unauthorized = vi.fn();
    window.addEventListener(AUTH_UNAUTHORIZED_EVENT, unauthorized);
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ success: false, error: 'forbidden', message: 'Forbidden' }), { status: 403, headers: { 'Content-Type': 'application/json' } }));
    await expect(httpApi.admin.getDashboard()).rejects.toMatchObject({ status: 403 });
    expect(unauthorized).not.toHaveBeenCalled();
    expect(localStorage.getItem(TOKEN_STORAGE_KEY)).toBe('jwt');
    window.removeEventListener(AUTH_UNAUTHORIZED_EVENT, unauthorized);
  });

  it('lets the browser set the multipart boundary for admin uploads', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ success: true, deck: {} }), { status: 201, headers: { 'Content-Type': 'application/json' } }));
    const payload = new FormData();
    payload.append('name_en', 'Deck');
    payload.append('cards_zip', new File(['zip'], 'deck.zip', { type: 'application/zip' }));
    await httpApi.admin.createDeck(payload);
    const request = fetchMock.mock.calls[0][1] as RequestInit;
    expect(request.body).toBe(payload);
    expect(new Headers(request.headers).has('Content-Type')).toBe(false);
    expect(new Headers(request.headers).get('Authorization')).toBe('Bearer jwt');
  });

  it('keeps the token after a network failure during hydration', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    useAuthStore.setState({ token: 'jwt', hydrated: false, status: 'idle' });
    vi.spyOn(globalThis, 'fetch').mockRejectedValue(new TypeError('Network error'));
    await useAuthStore.getState().hydrateSession();
    expect(localStorage.getItem(TOKEN_STORAGE_KEY)).toBe('jwt');
    expect(useAuthStore.getState().status).toBe('error');
  });

  it('clears the token after unauthorized hydration', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    useAuthStore.setState({ token: 'jwt', hydrated: false, status: 'idle' });
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ success: false, error: 'invalid_token', message: 'Invalid token' }), { status: 401, headers: { 'Content-Type': 'application/json' } }));
    await useAuthStore.getState().hydrateSession();
    expect(localStorage.getItem(TOKEN_STORAGE_KEY)).toBeNull();
    expect(useAuthStore.getState().status).toBe('guest');
  });
});
