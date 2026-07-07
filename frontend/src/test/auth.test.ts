import { beforeEach, describe, expect, it, vi } from 'vitest';
import { httpApi } from '@/api/http';
import { ApiError, TOKEN_STORAGE_KEY } from '@/api/client';
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

  it('parses structured backend errors', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ success: false, error: 'validation_failed', message: 'Email is already registered', details: { field: 'email', reason: 'already_registered' } }), { status: 400, headers: { 'Content-Type': 'application/json' } }));
    await expect(httpApi.auth.register({ email: 'alice@example.com', password: 'Secret123', timezone: 'UTC' })).rejects.toMatchObject({ status: 400, code: 'validation_failed', details: { field: 'email', reason: 'already_registered' } } satisfies Partial<ApiError>);
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
