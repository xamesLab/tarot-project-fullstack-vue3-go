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
  general: 'A concise overall reading.',
  message: 'A separate expanded message for the day.',
  relationships: 'Relationship insight for today.',
  finances: 'Financial insight for today.',
  emotional_state: 'Emotional state insight for today.',
  career_business: 'Career and business insight for today.',
};

function view(revealed: boolean, value: CardOfDayInterpretationDto | undefined = interpretation) {
  return <I18nextProvider i18n={i18n}><CardOfDayDetails backendCard={backendCard} interpretation={value} language="ru" revealed={revealed}/></I18nextProvider>;
}

describe('CardOfDayDetails', () => {
  it('shows the message before four spheres and resets expanded content when closed', async () => {
    await i18n.changeLanguage('ru');
    const { container, rerender } = render(view(false));
    const shell = container.querySelector<HTMLElement>('.card-details-shell')!;
    const expandButton = container.querySelector<HTMLButtonElement>('.expand-button')!;

    expect(shell).toHaveClass('concealed');
    expect(expandButton).toBeDisabled();
    rerender(view(true));
    expect(shell).toHaveClass('revealed');
    expect(expandButton).toBeEnabled();
    expect(expandButton).toHaveAttribute('aria-expanded', 'true');
    expect(screen.getByText('A concise overall reading.')).toBeInTheDocument();

    expect(container.querySelector('.corner-panel')).toHaveTextContent('A separate expanded message for the day.');
    expect(screen.getByText('Relationship insight for today.')).toBeInTheDocument();
    expect(screen.getByText('Financial insight for today.')).toBeInTheDocument();
    expect(screen.getByText('Emotional state insight for today.')).toBeInTheDocument();
    expect(screen.getByText('Career and business insight for today.')).toBeInTheDocument();

    const reveal = container.querySelector<HTMLElement>('.spheres-reveal')!;
    expect(reveal.firstElementChild).toHaveClass('corner-panel');
    expect(reveal.lastElementChild).toHaveClass('spheres');
    expect([...reveal.querySelectorAll('.aspect-block h3')].map((element) => element.textContent)).toEqual([
      i18n.t('interpretation.relationships', { ns: 'home' }),
      i18n.t('interpretation.finances', { ns: 'home' }),
      i18n.t('interpretation.emotionalState', { ns: 'home' }),
      i18n.t('interpretation.careerBusiness', { ns: 'home' }),
    ]);

    fireEvent.click(expandButton);
    expect(expandButton).toHaveAttribute('aria-expanded', 'false');
    expect(screen.queryByText('Relationship insight for today.')).not.toBeInTheDocument();
    rerender(view(false));
    expect(screen.queryByText('Relationship insight for today.')).not.toBeInTheDocument();
    rerender(view(true));
    expect(screen.getByText('Relationship insight for today.')).toBeInTheDocument();
  });

  it('shows a compact loading state until AI interpretation arrives', async () => {
    await i18n.changeLanguage('ru');
    render(<I18nextProvider i18n={i18n}><CardOfDayDetails backendCard={backendCard} language="ru" revealed interpretationLoading/></I18nextProvider>);

    expect(screen.getByText('Карты открывают послание…')).toBeInTheDocument();
    expect(screen.getByRole('status')).toHaveTextContent('оракул думает над вашей картой');
    expect(screen.getByRole('button', { name: /Развернуть толкование/i })).toBeDisabled();
  });
});
