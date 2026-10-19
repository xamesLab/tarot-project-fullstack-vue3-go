import type {
  BackendReadingCardDto,
  BackendReadingDto,
  CardOfDayInterpretationDto,
  DailyCardApi,
  Language,
} from '@/api/types';

interface MockCard {
  key: string;
  arcana: 'major' | 'minor';
  type: string;
  typeRu: string;
  name: string;
  nameRu: string;
  interpretation: Record<Language, string>;
  tags: Record<Language, string>;
}

const cards: MockCard[] = [
  {
    key: 'the_high_priestess',
    arcana: 'major',
    type: 'Major Arcana',
    typeRu: 'Старший Аркан',
    name: 'The High Priestess',
    nameRu: 'Верховная Жрица',
    interpretation: {
      en: 'A quiet answer is taking shape. Make room to listen before choosing your next step.',
      ru: 'Ответ постепенно проясняется. Прислушайтесь к себе, прежде чем выбирать следующий шаг.',
    },
    tags: { en: 'intuition, reflection, patience', ru: 'интуиция, размышление, терпение' },
  },
  {
    key: 'the_star',
    arcana: 'major',
    type: 'Major Arcana',
    typeRu: 'Старший Аркан',
    name: 'The Star',
    nameRu: 'Звезда',
    interpretation: {
      en: 'Hope returns through small, steady actions. Let a clear goal guide your energy.',
      ru: 'Надежда возвращается через небольшие последовательные шаги. Пусть ясная цель направляет вашу энергию.',
    },
    tags: { en: 'hope, renewal, direction', ru: 'надежда, обновление, направление' },
  },
  {
    key: 'the_fool',
    arcana: 'major',
    type: 'Major Arcana',
    typeRu: 'Старший Аркан',
    name: 'The Fool',
    nameRu: 'Шут',
    interpretation: {
      en: 'A fresh beginning is possible. Stay curious while keeping a little room for preparation.',
      ru: 'Возможно новое начало. Сохраняйте любопытство и оставьте время на подготовку.',
    },
    tags: { en: 'beginnings, curiosity, openness', ru: 'начало, любопытство, открытость' },
  },
];

type MockInterpretationCopy = Omit<CardOfDayInterpretationDto, 'general' | 'message'> & {
  general: (cardName: string, cardMeaning: string) => string;
  message: (cardName: string, cardMeaning: string) => string;
};

const interpretationCopy: Record<Language, MockInterpretationCopy> = {
  en: {
    general: (cardName, cardMeaning) => `${cardName}: ${cardMeaning}`,
    message: (cardName, cardMeaning) => `${cardName}: ${cardMeaning} Today, notice what is becoming clear before rushing toward an answer. Let events offer clues, and compare them with what you already know. A quiet pause can separate intuition from a passing worry. Choose one small action that respects your limits, then leave room for new information to change your view.`,
    relationships: 'In close relationships, pay attention to more than words and notice how you both feel together. A calm, direct conversation can clarify expectations when you give the other person time to respond. Do not mistake a pause for rejection or agree only to avoid disappointing someone. Name one boundary that matters to you, then look for a mutual step without pressure.',
    finances: 'For money matters, this is a good day to check details and plan calmly. Review recurring expenses, payment dates, and the terms of an offer before making a commitment. The card does not promise sudden income; it points toward clarity through simple numbers. Delay an impulse purchase if you feel uncertain, and choose one practical step that supports stability.',
    emotional_state: 'Your emotional state today may depend on how much attention you give to other people’s requests. Notice tiredness, tension, or curiosity without judging yourself, and allow a short pause. A comforting activity, a walk, or a conversation with someone you trust can help restore your footing. You do not need to solve everything now; choose a manageable pace and leave room for rest.',
    career_business: 'At work or in business, focus on a task where your experience already makes a difference. Clarify the goal, timeline, and responsibilities before expanding a plan or promising a result. If a new opportunity appears, ask practical questions and consider the resources it requires. Steady action and clear agreements will give you more support today than trying to accelerate everything at once.',
  },
  ru: {
    general: (cardName, cardMeaning) => `${cardName}: ${cardMeaning}`,
    message: (cardName, cardMeaning) => `Карта дня — ${cardName}: ${cardMeaning} Сегодня полезно присмотреться к событиям без спешки и заметить, какая мысль возвращается чаще остальных. Отделите собственные желания от чужих ожиданий, проверьте впечатления фактами и не требуйте от себя мгновенного ответа. Небольшой спокойный шаг поможет сохранить ясность, а новые сведения могут подсказать более точное направление. Сохраняйте любопытство и выбирайте темп, который позволяет вам слышать себя.`,
    relationships: 'Сегодня в отношениях важно замечать не только слова, но и то, как вы оба чувствуете себя рядом. Спокойный прямой разговор поможет прояснить ожидания, если вы дадите собеседнику время ответить. Не принимайте паузу за отказ и не соглашайтесь из страха расстроить другого. Обозначьте одну важную для вас границу, а затем ищите взаимный шаг без давления. Близость крепнет там, где остаётся место уважению и честности.',
    finances: 'В финансовых вопросах день подходит для внимательной проверки деталей и спокойного планирования. Пересмотрите регулярные расходы, сроки платежей и условия предложений, прежде чем принимать обязательства. Карта не обещает внезапного дохода; она напоминает, что ясность начинается с простых цифр. Отложите импульсивную покупку, если сомневаетесь, и выберите один выполнимый шаг для устойчивости. Осторожность сегодня помогает сохранить свободу выбора завтра.',
    emotional_state: 'Душевное состояние сегодня может зависеть от того, сколько внимания вы отдаёте чужим просьбам. Заметьте усталость, напряжение или интерес без критики и позвольте себе сделать короткую паузу. Поддерживающее занятие, прогулка или разговор с близким человеком помогут вернуть ощущение опоры. Не нужно решать всё сразу: выберите посильный ритм и оставьте место для отдыха. Бережное отношение к себе не мешает двигаться вперёд.',
    career_business: 'В карьере и бизнесе полезно сосредоточиться на задаче, где ваш опыт уже приносит заметную пользу. Уточните цель, сроки и ответственность, прежде чем расширять план или обещать результат. Если появилась новая возможность, задайте практические вопросы и оцените необходимые ресурсы. Последовательное действие и ясная договорённость сегодня дадут больше опоры, чем попытка ускорить всё одновременно. Отметьте, какой следующий шаг зависит именно от вас.',
  },
};

const readingCards = new Map<string, MockCard>();

function getAccountingDate(timezone: string, now = new Date()): string {
  const parts = new Intl.DateTimeFormat('en-US', {
    timeZone: timezone,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    hourCycle: 'h23',
  }).formatToParts(now);
  const value = (type: Intl.DateTimeFormatPartTypes) => Number(parts.find((part) => part.type === type)?.value);
  const date = new Date(Date.UTC(value('year'), value('month') - 1, value('day')));
  if (value('hour') < 3) date.setUTCDate(date.getUTCDate() - 1);
  return date.toISOString().slice(0, 10);
}

function cardFor(seed: string): MockCard {
  const hash = Array.from(seed).reduce((value, character) => (value * 31 + character.charCodeAt(0)) >>> 0, 7);
  return cards[hash % cards.length];
}

function toReading(card: MockCard, accountingDate: string, id?: string): BackendReadingDto {
  if (id) readingCards.set(id, card);
  const readingCard: BackendReadingCardDto = {
    key: card.key,
    arcana: card.arcana,
    type: card.type,
    type_ru: card.typeRu,
    name: card.name,
    name_ru: card.nameRu,
    interpretation: card.interpretation.en,
    interpretation_ru: card.interpretation.ru,
    tags: card.tags.en,
    tags_ru: card.tags.ru,
    image_url: '',
    back_image_url: '',
    position: 0,
    reversed: false,
    viewed: false,
  };
  return {
    id,
    type: 'card_of_day',
    accounting_date: accountingDate,
    generated_at: new Date().toISOString(),
    cards: [readingCard],
    deck_id: 'rider-waite',
    back_id: 'midnight-stars',
  };
}

function makeClaimToken(guestId: string, accountingDate: string): string {
  return `mock-claim:${encodeURIComponent(guestId)}:${accountingDate}`;
}

function cardFromClaimToken(token: string): { card: MockCard; date: string; guestId: string } | undefined {
  const match = /^mock-claim:([^:]+):(\d{4}-\d{2}-\d{2})$/.exec(token);
  if (!match) return undefined;
  const guestId = decodeURIComponent(match[1]);
  return { card: cardFor(`${guestId}:${match[2]}`), date: match[2], guestId };
}

function interpretationFor(card: MockCard, language: Language): CardOfDayInterpretationDto {
  const copy = interpretationCopy[language];
  const cardName = language === 'en' ? card.name : card.nameRu;
  return {
    general: copy.general(cardName, card.interpretation[language]),
    message: copy.message(cardName, card.interpretation[language]),
    relationships: copy.relationships,
    finances: copy.finances,
    emotional_state: copy.emotional_state,
    career_business: copy.career_business,
  };
}

export const mockDailyCardApi: DailyCardApi = {
  async generateGuest({ guest_id, timezone }) {
    const accountingDate = getAccountingDate(timezone);
    const card = cardFor(`${guest_id}:${accountingDate}`);
    return {
      success: true,
      source: 'guest',
      reading: toReading(card, accountingDate),
      claim_token: makeClaimToken(guest_id, accountingDate),
    };
  },

  async resolve({ timezone = 'UTC', claim_token }) {
    const claim = claim_token ? cardFromClaimToken(claim_token) : undefined;
    if (claim) {
      const readingId = `mock-reading-${encodeURIComponent(claim.guestId)}-${claim.date}`;
      return {
        success: true,
        source: 'guest_claimed',
        reading: toReading(claim.card, claim.date, readingId),
      };
    }

    const accountingDate = getAccountingDate(timezone);
    const readingId = `mock-reading-user-${accountingDate}`;
    return {
      success: true,
      source: 'existing',
      reading: toReading(cardFor(`user:${accountingDate}`), accountingDate, readingId),
    };
  },

  async getGuestInterpretation({ claim_token, language }) {
    const claim = cardFromClaimToken(claim_token);
    return {
      success: true,
      source: 'generated',
      interpretation: interpretationFor(claim?.card ?? cards[0], language),
    };
  },

  async getUserInterpretation({ reading_id, language }) {
    return {
      success: true,
      source: 'generated',
      interpretation: interpretationFor(readingCards.get(reading_id) ?? cards[0], language),
    };
  },
};
