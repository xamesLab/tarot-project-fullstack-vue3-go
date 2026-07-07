import { fireEvent, render, screen } from '@testing-library/react';
import { I18nextProvider } from 'react-i18next';
import type { ReactNode } from 'react';
import { describe, expect, it, vi } from 'vitest';
import type { BackendReadingCardDto } from '@/api/types';
import { i18n } from '@/features/i18n/i18n';
import { CardOfDayStage } from '@/widgets/CardOfDayStage/CardOfDayStage';

const backendCard = (key: string): BackendReadingCardDto => ({
  key,
  arcana: 'major',
  type: 'II',
  type_ru: 'II',
  name: key,
  name_ru: key,
  interpretation: '',
  interpretation_ru: '',
  tags: '',
  tags_ru: '',
  image_url: `/uploads/${key}.png`,
  back_image_url: '/uploads/back.png',
  position: 0,
  reversed: false,
  viewed: false,
});

const renderStage = (stage: ReactNode) => render(<I18nextProvider i18n={i18n}><div>{stage}</div></I18nextProvider>);

describe('CardOfDayStage backend images', () => {
  it('uses SVG while backend data is absent and never uses mock image URLs', () => {
    const { container } = renderStage(<CardOfDayStage language="en" hint="Flip"/>);
    expect(container.querySelector('img')).not.toBeInTheDocument();
    expect(container.querySelector('[data-card-side="face"] svg')).toBeInTheDocument();
    expect(container.innerHTML).not.toContain('/mock-face.png');
  });

  it('uses backend images and resets the flip for a different card key', () => {
    const onRevealChange = vi.fn();
    const { container, rerender } = renderStage(<CardOfDayStage backendCard={backendCard('first_card')} language="en" hint="Flip" onRevealChange={onRevealChange}/>);
    const card = screen.getByRole('button', { name: 'Flip' });
    expect(screen.getByRole('img', { name: 'first_card' })).toHaveAttribute('src', 'http://localhost:8888/uploads/first_card.png');
    expect(container.querySelector('[data-card-side="back"] img')).toHaveAttribute('src', 'http://localhost:8888/uploads/back.png');
    fireEvent.click(card);
    expect(card).toHaveAttribute('aria-pressed', 'true');
    expect(onRevealChange).toHaveBeenLastCalledWith(true);

    rerender(<I18nextProvider i18n={i18n}><div><CardOfDayStage backendCard={backendCard('second_card')} language="en" hint="Flip"/></div></I18nextProvider>);
    expect(screen.getByRole('button', { name: 'Flip' })).toHaveAttribute('aria-pressed', 'false');
    expect(screen.getByRole('img', { name: 'second_card' })).toBeInTheDocument();
  });

});
