import {
  useEffect,
  useId,
  useMemo,
  useRef,
  useState,
  type KeyboardEvent,
} from 'react';
import { useTranslation } from 'react-i18next';
import {
  MOON_PHASE_KEYS,
  formatZonedCalendarDate,
  getMoonSnapshot,
  type MoonPhaseKey,
} from '@/shared/lib/moonPhase';
import './celestial-strip.css';

const REFRESH_INTERVAL_MS = 15 * 60 * 1000;

interface CelestialStripProps {
  timeZone?: string | null;
  now?: Date;
  className?: string;
}

function MoonPhaseIcon({ phase }: { phase: MoonPhaseKey }) {
  const fill = {
    newMoon: null,
    waxingCrescent: <path d="M10 2.8a7.2 7.2 0 0 1 0 14.4 10.5 8.2 0 0 0 0-14.4Z" />,
    firstQuarter: <path d="M10 2.8a7.2 7.2 0 0 1 0 14.4Z" />,
    waxingGibbous: <path d="M10 2.8a7.2 7.2 0 1 1 0 14.4 4 7.2 0 0 0 0-14.4Z" />,
    fullMoon: <circle cx="10" cy="10" r="7.2" />,
    waningGibbous: <path d="M10 2.8a7.2 7.2 0 1 0 0 14.4 4 7.2 0 0 1 0-14.4Z" />,
    lastQuarter: <path d="M10 2.8a7.2 7.2 0 0 0 0 14.4Z" />,
    waningCrescent: <path d="M10 2.8a7.2 7.2 0 0 0 0 14.4 10.5 8.2 0 0 1 0-14.4Z" />,
  }[phase];

  return (
    <svg viewBox="0 0 20 20" aria-hidden="true">
      <circle className="celestial-moon-outline" cx="10" cy="10" r="7.2" />
      <g className="celestial-moon-fill">{fill}</g>
    </svg>
  );
}

export function CelestialStrip({ timeZone, now, className = '' }: CelestialStripProps) {
  const { t, i18n } = useTranslation('celestial');
  const fixedNow = now?.getTime();
  const [instant, setInstant] = useState(() => now ?? new Date());
  const [openPhaseIndex, setOpenPhaseIndex] = useState<number | null>(null);
  const phaseButtons = useRef<Array<HTMLButtonElement | null>>([]);
  const popoverId = `${useId()}-moon-phase-details`;

  useEffect(() => {
    if (fixedNow !== undefined) {
      setInstant(new Date(fixedNow));
      return undefined;
    }

    const refresh = () => setInstant(new Date());
    const interval = window.setInterval(refresh, REFRESH_INTERVAL_MS);
    const refreshWhenVisible = () => {
      if (!document.hidden) refresh();
    };
    document.addEventListener('visibilitychange', refreshWhenVisible);

    return () => {
      window.clearInterval(interval);
      document.removeEventListener('visibilitychange', refreshWhenVisible);
    };
  }, [fixedNow]);

  const locale = i18n.resolvedLanguage === 'en' ? 'en-US' : 'ru-RU';
  const snapshot = useMemo(() => getMoonSnapshot(instant), [instant]);
  const calendar = useMemo(
    () => formatZonedCalendarDate(instant, locale, timeZone),
    [instant, locale, timeZone],
  );
  const illuminationPercent = Math.round(snapshot.illumination * 100);
  const lunarAge = new Intl.NumberFormat(locale, { maximumFractionDigits: 1 }).format(snapshot.ageDays);
  const selectedPhase = openPhaseIndex === null ? null : MOON_PHASE_KEYS[openPhaseIndex];
  const currentPhaseName = t(`phases.${snapshot.phase}.name`);

  const handlePhaseKeyDown = (event: KeyboardEvent<HTMLButtonElement>, index: number) => {
    let targetIndex: number | null = null;
    if (event.key === 'ArrowRight' || event.key === 'ArrowDown') targetIndex = (index + 1) % MOON_PHASE_KEYS.length;
    if (event.key === 'ArrowLeft' || event.key === 'ArrowUp') targetIndex = (index - 1 + MOON_PHASE_KEYS.length) % MOON_PHASE_KEYS.length;
    if (event.key === 'Home') targetIndex = 0;
    if (event.key === 'End') targetIndex = MOON_PHASE_KEYS.length - 1;
    if (event.key === 'Escape') {
      setOpenPhaseIndex(null);
      return;
    }
    if (targetIndex === null) return;

    event.preventDefault();
    setOpenPhaseIndex(targetIndex);
    phaseButtons.current[targetIndex]?.focus();
  };

  return (
    <section
      className={`celestial-strip ${className}`.trim()}
      aria-label={t('summaryLabel', { phase: currentPhaseName })}
      onBlur={(event) => {
        if (!event.currentTarget.contains(event.relatedTarget)) setOpenPhaseIndex(null);
      }}
    >
      <div className="celestial-phases" role="group" aria-label={t('phaseGroupLabel')}>
        {MOON_PHASE_KEYS.map((phase, index) => {
          const phaseName = t(`phases.${phase}.name`);
          const isCurrent = index === snapshot.phaseIndex;
          const isOpen = index === openPhaseIndex;

          return (
            <button
              className={`celestial-phase-button ${isCurrent ? 'is-current' : ''} ${isOpen ? 'is-open' : ''}`.trim()}
              type="button"
              key={phase}
              ref={(node) => { phaseButtons.current[index] = node; }}
              aria-label={isCurrent ? `${phaseName}, ${t('currentMarker')}` : phaseName}
              aria-current={isCurrent ? 'true' : undefined}
              aria-expanded={isOpen}
              aria-controls={isOpen ? popoverId : undefined}
              onClick={() => setOpenPhaseIndex(index)}
              onFocus={() => setOpenPhaseIndex(index)}
              onKeyDown={(event) => handlePhaseKeyDown(event, index)}
            >
              <MoonPhaseIcon phase={phase} />
            </button>
          );
        })}
      </div>

      <div className="celestial-date-line">
        <span>{calendar.date}</span>
        <i aria-hidden="true">·</i>
        <span>{calendar.weekday}</span>
        <i aria-hidden="true">·</i>
        <span>{currentPhaseName}</span>
      </div>

      {selectedPhase && (
        <div className="celestial-popover" id={popoverId} role="tooltip">
          <strong>{t(`phases.${selectedPhase}.name`)}</strong>
          <p>{t(`phases.${selectedPhase}.description`)}</p>
          {openPhaseIndex === snapshot.phaseIndex ? (
            <dl className="celestial-metrics">
              <div><dt>{t('details.illumination')}</dt><dd>{illuminationPercent}%</dd></div>
              <div><dt>{t('details.age')}</dt><dd>{t('details.days', { value: lunarAge })}</dd></div>
              <div><dt>{t('details.direction')}</dt><dd>{snapshot.waxing ? t('waxing') : t('waning')}</dd></div>
            </dl>
          ) : (
            <span className="celestial-current-note">{t('details.currentPhase', { phase: currentPhaseName })}</span>
          )}
        </div>
      )}
    </section>
  );
}
