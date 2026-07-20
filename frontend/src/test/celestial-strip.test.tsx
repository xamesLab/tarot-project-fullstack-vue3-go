import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { i18n } from '@/features/i18n/i18n';
import {
  MOON_PHASE_KEYS,
  formatZonedCalendarDate,
  getMoonPhaseIndex,
  getMoonSnapshot,
  normalizeCycleProgress,
} from '@/shared/lib/moonPhase';
import { CelestialStrip } from '@/widgets/CelestialStrip/CelestialStrip';

describe('moon phase calculations', () => {
  it('maps the lunar cycle to the eight standard phase sectors', () => {
    MOON_PHASE_KEYS.forEach((_, index) => {
      expect(getMoonPhaseIndex(index / 8)).toBe(index);
    });

    expect(getMoonPhaseIndex(0.0624)).toBe(0);
    expect(getMoonPhaseIndex(0.0625)).toBe(1);
    expect(getMoonPhaseIndex(0.9374)).toBe(7);
    expect(getMoonPhaseIndex(0.9375)).toBe(0);
    expect(normalizeCycleProgress(-0.125)).toBe(0.875);
  });

  it('formats one instant as the correct local calendar day', () => {
    const instant = new Date('2026-07-19T21:30:00Z');
    const moscow = formatZonedCalendarDate(instant, 'ru-RU', 'Europe/Moscow');
    const losAngeles = formatZonedCalendarDate(instant, 'en-US', 'America/Los_Angeles');

    expect(moscow.date).toContain('20 июля 2026');
    expect(moscow.weekday).toBe('понедельник');
    expect(losAngeles.date).toContain('July 19, 2026');
    expect(losAngeles.weekday).toBe('Sunday');
  });
});

describe('CelestialStrip', () => {
  const instant = new Date('2026-07-19T21:30:00Z');

  beforeEach(async () => {
    await i18n.changeLanguage('ru');
  });

  it('shows live lunar data and marks exactly one current phase', () => {
    const snapshot = getMoonSnapshot(instant);
    const currentPhaseName = i18n.t(`celestial:phases.${snapshot.phase}.name`);
    render(<CelestialStrip now={instant} timeZone="Europe/Moscow" />);

    expect(screen.getAllByRole('button')).toHaveLength(8);
    expect(screen.getByText(/20 июля 2026/)).toBeInTheDocument();
    expect(screen.getByText('понедельник')).toBeInTheDocument();

    const currentButton = screen.getByRole('button', {
      name: `${currentPhaseName}, текущая фаза`,
    });
    expect(currentButton).toHaveAttribute('aria-current', 'true');

    fireEvent.focus(currentButton);
    expect(screen.getByRole('tooltip')).toHaveTextContent('Освещённость');
    expect(screen.getByRole('tooltip')).toHaveTextContent('Возраст Луны');
  });

  it('opens educational details and supports arrow-key navigation', () => {
    render(<CelestialStrip now={instant} timeZone="Europe/Moscow" />);
    const phaseButtons = screen.getAllByRole('button');

    fireEvent.focus(phaseButtons[0]);
    expect(screen.getByRole('tooltip')).toHaveTextContent('Новолуние');

    fireEvent.keyDown(phaseButtons[0], { key: 'ArrowRight' });
    expect(phaseButtons[1]).toHaveFocus();
    expect(screen.getByRole('tooltip')).toHaveTextContent('Растущий серп');

    fireEvent.keyDown(phaseButtons[1], { key: 'Escape' });
    expect(screen.queryByRole('tooltip')).not.toBeInTheDocument();
  });
});
