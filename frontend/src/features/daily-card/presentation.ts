import type { BackendReadingCardDto, Language } from '@/api/types';

export interface DailyCardPresentation {
  marker: string;
  name: string;
  arcana: 'major' | 'minor';
  reversed: boolean;
  tags: string[];
  interpretation: string;
}

export function splitCardTags(value: string): string[] {
  return value
    .split(',')
    .map((tag) => tag.trim().replace(/[.!?]+$/, ''))
    .filter(Boolean);
}

export function getDailyCardPresentation(card: BackendReadingCardDto, language: Language): DailyCardPresentation {
  const isRussian = language === 'ru';
  const rawMarker = isRussian ? card.type_ru : card.type;
  return {
    marker: card.arcana === 'minor' ? rawMarker.charAt(0).toUpperCase() + rawMarker.slice(1) : rawMarker,
    name: isRussian ? card.name_ru : card.name,
    arcana: card.arcana,
    reversed: card.reversed,
    tags: splitCardTags(isRussian ? card.tags_ru : card.tags),
    interpretation: isRussian ? card.interpretation_ru : card.interpretation,
  };
}
