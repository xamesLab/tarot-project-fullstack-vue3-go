import { beforeEach, describe, expect, it, vi } from 'vitest';
import { TOKEN_STORAGE_KEY } from '@/api/client';
import { httpDailyCardApi } from '@/api/http';

const reading = { type: 'card_of_day', accounting_date: '2026-07-03', generated_at: '2026-07-03T10:00:00Z', cards: [] };

describe('daily card HTTP adapter', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it('requests a public guest card', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ success: true, source: 'guest', reading, claim_token: 'claim' }), { status: 200, headers: { 'Content-Type': 'application/json' } }));
    await httpDailyCardApi.generateGuest({ guest_id: '0f12c9cc-1ef2-4e78-a0cd-3a30e7fa4b1a', timezone: 'UTC' });

    expect(fetchMock).toHaveBeenCalledWith('http://localhost:8888/api/reading/guest/card-of-day', expect.objectContaining({ method: 'POST' }));
    const request = fetchMock.mock.calls[0][1] as RequestInit;
    expect(JSON.parse(request.body as string)).toEqual({ guest_id: '0f12c9cc-1ef2-4e78-a0cd-3a30e7fa4b1a', timezone: 'UTC' });
    expect(new Headers(request.headers).has('Authorization')).toBe(false);
  });

  it('sends the JWT and claim token when resolving for a user', async () => {
    localStorage.setItem(TOKEN_STORAGE_KEY, 'jwt');
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ success: true, source: 'guest_claimed', reading: { ...reading, id: 'reading-1' } }), { status: 200, headers: { 'Content-Type': 'application/json' } }));
    await httpDailyCardApi.resolve({ timezone: 'UTC', claim_token: 'claim' });

    const request = fetchMock.mock.calls[0][1] as RequestInit;
    expect(new Headers(request.headers).get('Authorization')).toBe('Bearer jwt');
    expect(JSON.parse(request.body as string)).toEqual({ timezone: 'UTC', claim_token: 'claim' });
  });
});
