import { fireEvent, render, screen } from '@testing-library/react';
import { I18nextProvider } from 'react-i18next';
import { describe, expect, it } from 'vitest';
import type { BackendReadingCardDto, DailyReadingDto } from '@/api/types';
import { i18n } from '@/features/i18n/i18n';
import { CardOfDayDetails } from '@/widgets/CardOfDayDetails/CardOfDayDetails';

const backendCard: BackendReadingCardDto = {
  key: 'ace_of_wands', arcana: 'minor', type: 'Wands', type_ru: 'жезлы', name: 'Ace of Wands', name_ru: 'Туз Жезлов', interpretation: 'A new spark.', interpretation_ru: 'Новый импульс.', tags: 'Inspiration, action.', tags_ru: 'Вдохновение, действие.', image_url: '/face.png', back_image_url: '/back.png', position: 0, reversed: true, viewed: false,
};

const reading: DailyReadingDto = {
  date: '2026-07-05',
  card: { id: 'mock', arcana: 'major', roman: 'II', name: 'Мок', illustrationLabel: 'Mock', subtitle: 'Mock', imageUrl: '', backImageUrl: '' },
  keywords: ['mock'],
  brief: 'mock',
  spheres: [{ id: 'love', title: 'Любовь', text: 'Скрытый mock-текст' }],
  message: { paragraph: 'Mock message', closing: 'Mock closing', emphasis: 'Mock emphasis' },
};

function view(revealed: boolean) {
  return <I18nextProvider i18n={i18n}><CardOfDayDetails backendCard={backendCard} reading={reading} language="ru" revealed={revealed}/></I18nextProvider>;
}

describe('CardOfDayDetails', () => {
  it('conceals all details until revealed and resets expanded content when closed', async () => {
    await i18n.changeLanguage('ru');
    const { container, rerender } = render(view(false));
    const shell = container.querySelector<HTMLElement>('.card-details-shell')!;
    const body = container.querySelector<HTMLElement>('.card-details-body')!;
    const expandButton = container.querySelector<HTMLButtonElement>('.expand-button')!;

    expect(shell).toHaveClass('concealed');
    expect(body).toHaveAttribute('aria-hidden', 'true');
    expect(expandButton).toBeDisabled();
    expect(expandButton).toHaveAttribute('tabindex', '-1');
    expect(screen.getByText('Переверните карту, чтобы открыть описание')).toBeInTheDocument();

    rerender(view(true));
    expect(shell).toHaveClass('revealed');
    expect(body).toHaveAttribute('aria-hidden', 'false');
    expect(expandButton).toBeEnabled();
    expect(screen.getByText('Жезлы')).toHaveClass('card-number-minor');
    expect(screen.getByText('Младший Аркан · Перевёрнутое положение')).toBeInTheDocument();
    fireEvent.click(expandButton);
    expect(screen.getByText('Скрытый mock-текст')).toBeInTheDocument();

    rerender(view(false));
    expect(screen.queryByText('Скрытый mock-текст')).not.toBeInTheDocument();
    rerender(view(true));
    expect(screen.queryByText('Скрытый mock-текст')).not.toBeInTheDocument();
  });
});
