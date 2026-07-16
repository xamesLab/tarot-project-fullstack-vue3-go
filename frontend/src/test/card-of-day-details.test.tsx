import { fireEvent, render, screen } from '@testing-library/react';
import { I18nextProvider } from 'react-i18next';
import { describe, expect, it } from 'vitest';
import type { BackendReadingCardDto, CardOfDayInterpretationDto } from '@/api/types';
import { i18n } from '@/features/i18n/i18n';
import { CardOfDayDetails } from '@/widgets/CardOfDayDetails/CardOfDayDetails';

const backendCard: BackendReadingCardDto = {
  key: 'ace_of_wands',
  arcana: 'minor',
  type: 'Wands',
  type_ru: 'Жезлы',
  name: 'Ace of Wands',
  name_ru: 'Туз Жезлов',
  interpretation: 'A new spark.',
  interpretation_ru: 'Новый импульс.',
  tags: 'Inspiration, action.',
  tags_ru: 'Вдохновение, действие.',
  image_url: '/face.png',
  back_image_url: '/back.png',
  position: 0,
  reversed: true,
  viewed: false,
};

const interpretation: CardOfDayInterpretationDto = {
  general: 'Карта зовёт к началу.',
  work: 'В работе выберите первый ясный шаг.',
  family: 'В отношениях говорите прямо.',
  advice: 'Действуйте спокойно.',
};

function view(revealed: boolean, value: CardOfDayInterpretationDto | undefined = interpretation) {
  return <I18nextProvider i18n={i18n}><CardOfDayDetails backendCard={backendCard} interpretation={value} language="ru" revealed={revealed}/></I18nextProvider>;
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
    expect(screen.getByText('Карта зовёт к началу.')).toBeInTheDocument();
    fireEvent.click(expandButton);
    expect(screen.getByText('В работе выберите первый ясный шаг.')).toBeInTheDocument();
    expect(screen.getByText('В отношениях говорите прямо.')).toBeInTheDocument();
    expect(screen.getByText('Действуйте спокойно.')).toBeInTheDocument();

    rerender(view(false));
    expect(screen.queryByText('В работе выберите первый ясный шаг.')).not.toBeInTheDocument();
    rerender(view(true));
    expect(screen.queryByText('В работе выберите первый ясный шаг.')).not.toBeInTheDocument();
  });

  it('shows a compact loading state until AI interpretation arrives', async () => {
    await i18n.changeLanguage('ru');
    render(<I18nextProvider i18n={i18n}><CardOfDayDetails backendCard={backendCard} language="ru" revealed interpretationLoading/></I18nextProvider>);

    expect(screen.getByText('Карты открывают послание…')).toBeInTheDocument();
    expect(screen.getByRole('status')).toHaveTextContent('оракул думает над вашей картой');
    expect(screen.getByRole('button', { name: /Развернуть толкование/i })).toBeDisabled();
  });
});
