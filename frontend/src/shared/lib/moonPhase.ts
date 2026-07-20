import { getMoonIllumination } from 'suncalc';

export const MOON_PHASE_KEYS = [
  'newMoon',
  'waxingCrescent',
  'firstQuarter',
  'waxingGibbous',
  'fullMoon',
  'waningGibbous',
  'lastQuarter',
  'waningCrescent',
] as const;

export type MoonPhaseKey = (typeof MOON_PHASE_KEYS)[number];

export interface MoonSnapshot {
  phase: MoonPhaseKey;
  phaseIndex: number;
  cycleProgress: number;
  illumination: number;
  ageDays: number;
  waxing: boolean;
}

export const SYNODIC_MONTH_DAYS = 29.530588853;

export function normalizeCycleProgress(progress: number) {
  return ((progress % 1) + 1) % 1;
}

export function getMoonPhaseIndex(progress: number) {
  return Math.round(normalizeCycleProgress(progress) * 8) % 8;
}

export function getMoonSnapshot(date: Date): MoonSnapshot {
  const illumination = getMoonIllumination(date);
  const cycleProgress = normalizeCycleProgress(illumination.phase);
  const phaseIndex = getMoonPhaseIndex(cycleProgress);

  return {
    phase: MOON_PHASE_KEYS[phaseIndex],
    phaseIndex,
    cycleProgress,
    illumination: Math.min(1, Math.max(0, illumination.fraction)),
    ageDays: cycleProgress * SYNODIC_MONTH_DAYS,
    waxing: illumination.waxing,
  };
}

export function resolveTimeZone(timeZone?: string | null) {
  const browserTimeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const candidate = timeZone || browserTimeZone || 'UTC';

  try {
    new Intl.DateTimeFormat('en', { timeZone: candidate }).format();
    return candidate;
  } catch {
    return 'UTC';
  }
}

export function formatZonedCalendarDate(date: Date, locale: string, timeZone?: string | null) {
  const resolvedTimeZone = resolveTimeZone(timeZone);

  return {
    date: new Intl.DateTimeFormat(locale, {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
      timeZone: resolvedTimeZone,
    }).format(date),
    weekday: new Intl.DateTimeFormat(locale, {
      weekday: 'long',
      timeZone: resolvedTimeZone,
    }).format(date),
    timeZone: resolvedTimeZone,
  };
}
