import { api } from '@/api';
import { ApiError } from '@/api/client';
import type { BackendReadingDto, DailyCardApi, UserDto } from '@/api/types';

export const GUEST_DAILY_CARD_STORAGE_KEY = 'nebula-tarot:guest-card-of-day:v3';
const LEGACY_GUEST_DAILY_CARD_STORAGE_KEY = 'nebula-tarot:guest-card-of-day:v1';
const PREVIOUS_GUEST_DAILY_CARD_STORAGE_KEY = 'nebula-tarot:guest-card-of-day:v2';

interface GuestDailyCardCache {
  version: 3;
  guestId: string;
  timezone: string;
  accountingDate: string;
  claimToken: string;
  reading: BackendReadingDto;
}

export interface DailyCardBootstrapResult {
  source: 'guest-cache' | 'guest-api' | 'existing' | 'guest_claimed' | 'generated';
  reading: BackendReadingDto;
  claimToken?: string;
}

interface BootstrapDependencies {
  dailyCardApi?: DailyCardApi;
  storage?: Storage;
  now?: Date;
  browserTimezone?: string;
  createGuestId?: () => string;
}

function dateParts(date: Date, timezone: string) {
  const parts = new Intl.DateTimeFormat('en-US', {
    timeZone: timezone,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    hourCycle: 'h23',
  }).formatToParts(date);
  const value = (type: Intl.DateTimeFormatPartTypes) => Number(parts.find((part) => part.type === type)?.value);
  return { year: value('year'), month: value('month'), day: value('day'), hour: value('hour') };
}

export function getAccountingDate(date: Date, timezone: string): string {
  const { year, month, day, hour } = dateParts(date, timezone);
  const calendarDate = new Date(Date.UTC(year, month - 1, day));
  if (hour < 3) calendarDate.setUTCDate(calendarDate.getUTCDate() - 1);
  return calendarDate.toISOString().slice(0, 10);
}

function readCache(storage: Storage): GuestDailyCardCache | null {
  storage.removeItem(LEGACY_GUEST_DAILY_CARD_STORAGE_KEY);
  storage.removeItem(PREVIOUS_GUEST_DAILY_CARD_STORAGE_KEY);
  const raw = storage.getItem(GUEST_DAILY_CARD_STORAGE_KEY);
  if (!raw) return null;
  try {
    const value = JSON.parse(raw) as Partial<GuestDailyCardCache>;
    const card = value.reading?.cards?.[0];
    if (value.version !== 3 || !value.guestId || !value.timezone || !value.accountingDate || !value.claimToken || !card?.arcana || !card.type || !card.type_ru) {
      throw new Error('Invalid guest daily card cache');
    }
    return value as GuestDailyCardCache;
  } catch {
    storage.removeItem(GUEST_DAILY_CARD_STORAGE_KEY);
    return null;
  }
}

function isCurrent(cache: GuestDailyCardCache, now: Date): boolean {
  try {
    return getAccountingDate(now, cache.timezone) === cache.accountingDate;
  } catch {
    return false;
  }
}

export async function bootstrapDailyCard(user: UserDto | null, dependencies: BootstrapDependencies = {}): Promise<DailyCardBootstrapResult> {
  const dailyCardApi = dependencies.dailyCardApi ?? api.dailyCard;
  const storage = dependencies.storage ?? localStorage;
  const now = dependencies.now ?? new Date();
  let cache = readCache(storage);
  if (cache && !isCurrent(cache, now)) {
    storage.removeItem(GUEST_DAILY_CARD_STORAGE_KEY);
    cache = null;
  }

  if (!user) {
    if (cache) return { source: 'guest-cache', reading: cache.reading, claimToken: cache.claimToken };

    const timezone = dependencies.browserTimezone ?? Intl.DateTimeFormat().resolvedOptions().timeZone ?? 'UTC';
    const guestId = dependencies.createGuestId?.() ?? crypto.randomUUID();
    const response = await dailyCardApi.generateGuest({ guest_id: guestId, timezone });
    const nextCache: GuestDailyCardCache = {
      version: 3,
      guestId,
      timezone,
      accountingDate: response.reading.accounting_date,
      claimToken: response.claim_token,
      reading: response.reading,
    };
    storage.setItem(GUEST_DAILY_CARD_STORAGE_KEY, JSON.stringify(nextCache));
    return { source: 'guest-api', reading: response.reading, claimToken: response.claim_token };
  }

  if (cache) {
    try {
      const response = await dailyCardApi.resolve({ timezone: cache.timezone, claim_token: cache.claimToken });
      storage.removeItem(GUEST_DAILY_CARD_STORAGE_KEY);
      return { source: response.source, reading: response.reading };
    } catch (error) {
      if (!(error instanceof ApiError) || error.code !== 'invalid_guest_claim') throw error;
      storage.removeItem(GUEST_DAILY_CARD_STORAGE_KEY);
    }
  }

  const response = await dailyCardApi.resolve({ timezone: user.timezone });
  storage.removeItem(GUEST_DAILY_CARD_STORAGE_KEY);
  return { source: response.source, reading: response.reading };
}
