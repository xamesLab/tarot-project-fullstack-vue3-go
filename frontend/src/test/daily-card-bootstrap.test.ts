import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ApiError } from '@/api/client';
import type { BackendReadingDto, DailyCardApi, UserDto } from '@/api/types';
import { bootstrapDailyCard, getAccountingDate, GUEST_DAILY_CARD_STORAGE_KEY } from '@/features/daily-card/bootstrap';

const reading = (date = '2026-07-03'): BackendReadingDto => ({
  type: 'card_of_day',
  accounting_date: date,
  generated_at: `${date}T10:00:00Z`,
  cards: [{ key: 'the_high_priestess', arcana: 'major', type: 'II', type_ru: 'II', name: 'The High Priestess', name_ru: 'Верховная Жрица', interpretation: '', interpretation_ru: '', tags: '', tags_ru: '', image_url: '/face.png', back_image_url: '/back.png', position: 0, reversed: false, viewed: false }],
});

const user: UserDto = { id: 'user-1', username: 'alice', email: 'alice@example.com', role: 'seeker', status: 'active', timezone: 'UTC', language: 'ru', theme: 'dark' };

function apiMock(): DailyCardApi {
  return { generateGuest: vi.fn(), resolve: vi.fn(), getGuestInterpretation: vi.fn(), getUserInterpretation: vi.fn() };
}

describe('daily card bootstrap', () => {
  beforeEach(() => localStorage.clear());

  it('creates a guest card once and restores it from localStorage', async () => {
    const dailyCardApi = apiMock();
    vi.mocked(dailyCardApi.generateGuest).mockResolvedValue({ success: true, source: 'guest', reading: reading(), claim_token: 'claim' });
    const dependencies = { dailyCardApi, now: new Date('2026-07-03T12:00:00Z'), browserTimezone: 'UTC', createGuestId: () => '0f12c9cc-1ef2-4e78-a0cd-3a30e7fa4b1a' };

    expect((await bootstrapDailyCard(null, dependencies)).source).toBe('guest-api');
    const cached = await bootstrapDailyCard(null, dependencies);
    expect(cached.source).toBe('guest-cache');
    expect(cached.claimToken).toBe('claim');
    expect(dailyCardApi.generateGuest).toHaveBeenCalledTimes(1);
    expect(localStorage.getItem(GUEST_DAILY_CARD_STORAGE_KEY)).toContain('the_high_priestess');
  });

  it('discards the legacy cache that does not contain classification fields', async () => {
    localStorage.setItem('nebula-tarot:guest-card-of-day:v1', JSON.stringify({ version: 1, reading: { cards: [{ key: 'legacy' }] } }));
    const dailyCardApi = apiMock();
    vi.mocked(dailyCardApi.generateGuest).mockResolvedValue({ success: true, source: 'guest', reading: reading(), claim_token: 'claim' });

    await bootstrapDailyCard(null, { dailyCardApi, now: new Date('2026-07-03T12:00:00Z'), browserTimezone: 'UTC' });
    expect(localStorage.getItem('nebula-tarot:guest-card-of-day:v1')).toBeNull();
    expect(localStorage.getItem(GUEST_DAILY_CARD_STORAGE_KEY)).toContain('"version":3');
  });

  it('replaces a stale guest card after the accounting day changes', async () => {
    const dailyCardApi = apiMock();
    vi.mocked(dailyCardApi.generateGuest)
      .mockResolvedValueOnce({ success: true, source: 'guest', reading: reading(), claim_token: 'claim-1' })
      .mockResolvedValueOnce({ success: true, source: 'guest', reading: reading('2026-07-04'), claim_token: 'claim-2' });
    const base = { dailyCardApi, browserTimezone: 'UTC', createGuestId: () => crypto.randomUUID() };

    await bootstrapDailyCard(null, { ...base, now: new Date('2026-07-03T12:00:00Z') });
    await bootstrapDailyCard(null, { ...base, now: new Date('2026-07-04T12:00:00Z') });
    expect(dailyCardApi.generateGuest).toHaveBeenCalledTimes(2);
    expect(localStorage.getItem(GUEST_DAILY_CARD_STORAGE_KEY)).toContain('claim-2');
  });

  it('claims the cached guest card after authentication and clears the cache', async () => {
    const dailyCardApi = apiMock();
    vi.mocked(dailyCardApi.generateGuest).mockResolvedValue({ success: true, source: 'guest', reading: reading(), claim_token: 'claim' });
    vi.mocked(dailyCardApi.resolve).mockResolvedValue({ success: true, source: 'guest_claimed', reading: { ...reading(), id: 'reading-1' } });
    const dependencies = { dailyCardApi, now: new Date('2026-07-03T12:00:00Z'), browserTimezone: 'UTC', createGuestId: () => '0f12c9cc-1ef2-4e78-a0cd-3a30e7fa4b1a' };

    await bootstrapDailyCard(null, dependencies);
    expect((await bootstrapDailyCard(user, dependencies)).source).toBe('guest_claimed');
    expect(dailyCardApi.resolve).toHaveBeenCalledWith({ timezone: 'UTC', claim_token: 'claim' });
    expect(localStorage.getItem(GUEST_DAILY_CARD_STORAGE_KEY)).toBeNull();
  });

  it('keeps the guest cache after a network failure', async () => {
    const dailyCardApi = apiMock();
    vi.mocked(dailyCardApi.generateGuest).mockResolvedValue({ success: true, source: 'guest', reading: reading(), claim_token: 'claim' });
    vi.mocked(dailyCardApi.resolve).mockRejectedValue(new TypeError('Network error'));
    const dependencies = { dailyCardApi, now: new Date('2026-07-03T12:00:00Z'), browserTimezone: 'UTC' };

    await bootstrapDailyCard(null, dependencies);
    await expect(bootstrapDailyCard(user, dependencies)).rejects.toThrow('Network error');
    expect(localStorage.getItem(GUEST_DAILY_CARD_STORAGE_KEY)).not.toBeNull();
  });

  it('falls back to normal resolve after an invalid claim', async () => {
    const dailyCardApi = apiMock();
    vi.mocked(dailyCardApi.generateGuest).mockResolvedValue({ success: true, source: 'guest', reading: reading(), claim_token: 'claim' });
    vi.mocked(dailyCardApi.resolve)
      .mockRejectedValueOnce(new ApiError(400, 'invalid_guest_claim', 'Invalid claim'))
      .mockResolvedValueOnce({ success: true, source: 'generated', reading: { ...reading(), id: 'reading-2' } });
    const dependencies = { dailyCardApi, now: new Date('2026-07-03T12:00:00Z'), browserTimezone: 'UTC' };

    await bootstrapDailyCard(null, dependencies);
    expect((await bootstrapDailyCard(user, dependencies)).source).toBe('generated');
    expect(dailyCardApi.resolve).toHaveBeenLastCalledWith({ timezone: 'UTC' });
    expect(localStorage.getItem(GUEST_DAILY_CARD_STORAGE_KEY)).toBeNull();
  });

  it('uses the backend 03:00 boundary for accounting dates', () => {
    expect(getAccountingDate(new Date('2026-07-03T02:59:00Z'), 'UTC')).toBe('2026-07-02');
    expect(getAccountingDate(new Date('2026-07-03T03:00:00Z'), 'UTC')).toBe('2026-07-03');
  });
});
