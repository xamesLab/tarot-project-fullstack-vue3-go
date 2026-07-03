import type { DailyReadingDto, Language } from '@/api/types';

const readings: Record<Language, DailyReadingDto> = {
  ru: {
    date: '2026-05-28',
    card: { id: 'high-priestess', arcana: 'major', roman: 'II', name: 'Верховная Жрица', illustrationLabel: 'The High Priestess', subtitle: 'Старший Аркан · Прямое положение' },
    keywords: ['Интуиция', 'Тайное знание', 'Внутренний голос', 'Цикличность'],
    brief: 'Сегодня вселенная приглашает вас замолчать и прислушаться. Ответы, которые вы ищете, живут не снаружи — они ждут в глубине вашего собственного существа.',
    spheres: [
      { id: 'love', title: 'Любовь & Отношения', text: 'Сейчас не время форсировать события. Карта указывает на скрытые чувства — возможно, важный человек ещё не готов открыться. Наблюдайте, не давите.' },
      { id: 'career', title: 'Карьера & Финансы', text: 'Информация, которой вы не располагаете, может изменить ход дел. Соберите все данные, прежде чем принимать решения. Это день анализа, а не действия.' },
      { id: 'health', title: 'Здоровье & Тело', text: 'Обратите внимание на сигналы тела. Жрица символизирует интуитивное здоровье — медитация и покой сегодня особенно целительны.' },
      { id: 'spirit', title: 'Духовность & Рост', text: 'Пик духовной восприимчивости. Сновидения сегодня насыщенны — ведите дневник. Открывается доступ к коллективному бессознательному.' }
    ],
    message: { paragraph: 'Верховная Жрица — хранительница тайн, сидящая между двумя колоннами: светом и тьмой. Её появление сегодня не случайно.', closing: 'Доверьтесь тому, что чувствуете, а не только тому, что можете объяснить.', emphasis: 'Ваша интуиция сейчас острее разума.' }
  },
  en: {
    date: '2026-05-28',
    card: { id: 'high-priestess', arcana: 'major', roman: 'II', name: 'The High Priestess', illustrationLabel: 'The High Priestess', subtitle: 'Major Arcana · Upright' },
    keywords: ['Intuition', 'Secret Knowledge', 'Inner Voice', 'Cycles'],
    brief: 'Today the universe invites you to fall silent and listen. The answers you seek do not live outside — they wait in the depths of your own being.',
    spheres: [
      { id: 'love', title: 'Love & Relations', text: 'This is not the time to force things. The card points to hidden feelings — perhaps someone important is not yet ready to open up. Observe, do not press.' },
      { id: 'career', title: 'Career & Finance', text: 'Information you lack could change everything. Gather all the facts before deciding. This is a day for analysis, not action.' },
      { id: 'health', title: 'Health & Body', text: 'Pay attention to the signals of your body. The Priestess embodies intuitive health — meditation and rest are especially healing today.' },
      { id: 'spirit', title: 'Spirit & Growth', text: 'A peak of spiritual receptivity. Dreams run vivid tonight — keep a journal. The gate to the collective unconscious stands open.' }
    ],
    message: { paragraph: 'The High Priestess is the keeper of mysteries, seated between two pillars: light and dark. Her appearance today is no accident.', closing: 'Trust what you feel, not only what you can explain.', emphasis: 'Your intuition is sharper than reason now.' }
  }
};

export const getReadingFixture = (language: Language) => structuredClone(readings[language]);
