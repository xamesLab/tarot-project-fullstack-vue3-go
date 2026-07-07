import { describe, expect, it } from 'vitest';
import type { BackendReadingCardDto } from '@/api/types';
import { getDailyCardPresentation, splitCardTags } from '@/features/daily-card/presentation';

const card: BackendReadingCardDto = {
  key: 'the_high_priestess',
  arcana: 'major',
  type: 'II',
  type_ru: 'II',
  name: 'The High Priestess',
  name_ru: 'Верховная Жрица',
  interpretation: 'Trust your intuition.',
  interpretation_ru: 'Доверьтесь интуиции.',
  tags: 'Intuition, hidden knowledge.',
  tags_ru: 'Интуиция, тайное знание.',
  image_url: '/face.png',
  back_image_url: '/back.png',
  position: 0,
  reversed: false,
  viewed: false,
};

describe('daily card presentation', () => {
  it('localizes major arcana data', () => {
    expect(getDailyCardPresentation(card, 'ru')).toMatchObject({ marker: 'II', name: 'Верховная Жрица', arcana: 'major', reversed: false, tags: ['Интуиция', 'тайное знание'], interpretation: 'Доверьтесь интуиции.' });
    expect(getDailyCardPresentation(card, 'en')).toMatchObject({ marker: 'II', name: 'The High Priestess', tags: ['Intuition', 'hidden knowledge'], interpretation: 'Trust your intuition.' });
  });

  it('uses the localized suit for a reversed minor arcana', () => {
    const minor = { ...card, arcana: 'minor' as const, type: 'Wands', type_ru: 'жезлы', reversed: true };
    expect(getDailyCardPresentation(minor, 'ru')).toMatchObject({ marker: 'Жезлы', arcana: 'minor', reversed: true });
    expect(getDailyCardPresentation(minor, 'en')).toMatchObject({ marker: 'Wands', arcana: 'minor', reversed: true });
  });

  it('removes empty tags and terminal punctuation', () => {
    expect(splitCardTags(' First, , Second!, Third. ')).toEqual(['First', 'Second', 'Third']);
  });
});
