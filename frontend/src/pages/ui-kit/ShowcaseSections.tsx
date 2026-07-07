import { useState, type ReactNode } from 'react';
import {
  CelestialCrossIcon,
  Constellation,
  LogoSigil,
  MoonPhases,
  SphereIcon,
  StarIcon,
  ZodiacWheel,
} from '@/shared/icons';
import { Avatar } from '@/shared/ui/Avatar';
import { Button } from '@/shared/ui/Button';
import { Chip } from '@/shared/ui/Chip';
import { AspectBlock, CornerPanel, Tooltip } from '@/shared/ui/Content';
import { Modal } from '@/shared/ui/Modal';
import { Ornament } from '@/shared/ui/Ornament';
import { TarotCard } from '@/shared/ui/TarotCard';
import { CodeCaption, DemoRow, UiKitSection } from './UiKitSection';

const iconNames = [
  'Всевидящее око',
  'Солнце',
  'Луна',
  'Полнолуние',
  'Звезда',
  'Искра',
  'Пентаграмма',
  'Огонь',
  'Вода',
  'Воздух',
  'Земля',
  'Ключ',
  'Чаша',
  'Меч',
  'Посох',
  'Пентакль',
  'Анк',
  'Свеча',
  'Часы',
  'Молния',
  'Кристалл',
  'Комета',
  'Уроборос',
  'Лотос',
  'Перо',
  'Книга',
  'Портфель',
  'Око',
  'Здоровье',
  'Закрыть',
  'Галочка',
  'Инфо',
  'Талисман',
];

function MysticGlyph({ index }: { index: number }) {
  const glyphs: ReactNode[] = [
    <>
      <path d="M3 16s5-7 13-7 13 7 13 7-5 7-13 7S3 16 3 16Z" />
      <circle cx="16" cy="16" r="4" />
      <circle cx="16" cy="16" r="1.5" fill="currentColor" />
    </>,
    <>
      <circle cx="16" cy="16" r="6" />
      <path d="M16 2v6m0 16v6M2 16h6m16 0h6M6 6l4 4m12 12 4 4m0-20-4 4M10 22l-4 4" />
    </>,
    <path d="M21 4a12 12 0 1 0 7 20A10 10 0 1 1 21 4Z" />,
    <>
      <circle cx="16" cy="16" r="10" />
      <circle cx="13" cy="12" r="2" opacity=".5" />
      <circle cx="20" cy="19" r="1.5" opacity=".5" />
    </>,
    <path d="m16 3 3 9h10l-8 6 3 10-8-6-8 6 3-10-8-6h10l3-9Z" />,
    <path d="M16 3v9m0 8v9M3 16h9m8 0h9M8 8l5 5m6 6 5 5m0-16-5 5m-6 6-5 5" />,
    <path d="m16 3 3 9h10l-8 6 3 10-8-6-8 6 3-10-8-6h10l3-9Z" />,
    <path d="M17 3c2 7-4 8-1 13 1-4 6-5 6-10 5 6 7 12 3 18-4 6-14 6-18 0-3-5-1-10 4-15-1 5 1 7 3 9-1-7 5-9 3-15Z" />,
    <>
      <path d="M16 4c4 6 7 9 7 14a7 7 0 0 1-14 0c0-5 3-8 7-14Z" />
      <path d="M5 27c4-2 7 2 11 0s7 2 11 0" />
    </>,
    <path d="M4 11c5-4 9 3 14-1 3-2 6-2 10 0M4 17c4-3 7 2 11 0s7-3 13 0M4 23c5-3 9 2 14 0 3-2 6-2 10 0" />,
    <>
      <circle cx="16" cy="16" r="11" />
      <path d="M5 16h22M16 5v22" />
    </>,
    <>
      <circle cx="9" cy="22" r="4" />
      <path d="m12 19 14-14m-4 4 4 4m-8 0 4 4" />
    </>,
    <>
      <path d="M8 4h16l-2 10c-1 5-11 5-12 0L8 4Zm8 14v9m-6 1h12" />
      <path d="M9 8h14" />
    </>,
    <>
      <path d="m16 3 3 21-3 5-3-5 3-21ZM12 23h8" />
      <path d="m16 3-3 5h6l-3-5Z" />
    </>,
    <>
      <path d="M16 3v26M12 7h8M10 29h12" />
      <circle cx="16" cy="15" r="3" />
    </>,
    <>
      <circle cx="16" cy="16" r="11" />
      <path d="m16 7 3 6 7 1-5 5 1 7-6-3-6 3 1-7-5-5 7-1 3-6Z" />
    </>,
    <path d="M16 3v26m-7-16h14M16 7c-5 0-7 3-7 6s2 6 7 6 7-3 7-6-2-6-7-6Z" />,
    <>
      <path d="M11 27h10M13 23h6l2-12c-2-5-8-5-10 0l2 12Z" />
      <path d="M16 3v4m-3-3 2 3m4-3-2 3" />
    </>,
    <>
      <circle cx="16" cy="16" r="11" />
      <path d="M16 8v9l6 3" />
    </>,
    <path d="M19 2 8 18h8l-3 12 11-16h-8l3-12Z" />,
    <path d="m16 3 9 8-4 17H11L7 11l9-8Zm-9 8h18M11 28l5-17 5 17" />,
    <>
      <path d="M4 20c7 0 8-12 15-12 5 0 8 4 8 8 0 7-8 10-14 8" />
      <path d="m5 20 5-3 1 6-6-3Z" />
      <circle cx="22" cy="13" r="1" fill="currentColor" />
    </>,
    <>
      <circle cx="16" cy="16" r="11" strokeDasharray="3 2" />
      <path d="m7 22 5 1-1 5" />
      <circle cx="16" cy="16" r="4" />
    </>,
    <path d="M16 27c-8-5-11-11-8-15 2-3 6-2 8 2 2-4 6-5 8-2 3 4 0 10-8 15Zm0-13V5m-5 3 5 6 5-6" />,
    <path d="M8 27c4-8 4-16 14-22-1 10-5 18-14 22Zm3-5 10-13M7 28l4-6" />,
    <>
      <path d="M5 7h10c3 0 5 2 5 5v15H10c-3 0-5-2-5-5V7Zm22 0H17" />
      <path d="M16 10v17" />
    </>,
    <>
      <rect x="4" y="9" width="24" height="17" rx="2" />
      <path d="M11 9V6h10v3M4 16h24M14 16v3h4v-3" />
    </>,
    <>
      <path d="M3 16s5-7 13-7 13 7 13 7-5 7-13 7S3 16 3 16Z" />
      <circle cx="16" cy="16" r="4" />
    </>,
    <>
      <path d="M16 28c-8-5-12-10-9-15 2-4 7-3 9 1 2-4 7-5 9-1 3 5-1 10-9 15Z" />
      <path d="M5 16h6l2-5 5 10 2-5h7" />
    </>,
    <path d="m7 7 18 18M25 7 7 25" />,
    <path d="m5 17 7 7L28 7" />,
    <>
      <circle cx="16" cy="16" r="12" />
      <path d="M16 14v8m0-13v1" />
    </>,
    <>
      <path d="m16 3 11 6v14l-11 6-11-6V9l11-6Z" />
      <path d="m16 8 6 8-6 8-6-8 6-8Z" />
      <circle cx="16" cy="16" r="2" fill="currentColor" />
    </>,
  ];
  return (
    <svg
      viewBox="0 0 32 32"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.1"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      {glyphs[index]}
    </svg>
  );
}

function LogoMark({ variant }: { variant: 1 | 2 | 3 }) {
  if (variant === 1) return <LogoSigil />;
  if (variant === 2)
    return (
      <svg
        viewBox="0 0 100 100"
        fill="none"
        stroke="currentColor"
        strokeLinecap="round"
        strokeLinejoin="round"
      >
        <path
          d="M50 8a16 16 0 1 0 0 32 12 12 0 1 1 0-32Z"
          fill="currentColor"
          stroke="none"
          opacity=".85"
        />
        <rect x="27" y="48" width="9" height="40" rx="1.5" />
        <rect x="64" y="48" width="9" height="40" rx="1.5" />
        <path d="M23 48h17m20 0h17M21 90h58M50 58l3 10-3 10-3-10 3-10Z" />
      </svg>
    );
  return (
    <svg viewBox="0 0 100 100" fill="none" stroke="currentColor" strokeLinejoin="round">
      <path d="m50 4 40 23v46L50 96 10 73V27L50 4Z" opacity=".55" />
      <path d="m50 18 32 32-32 32-32-32 32-32Z" />
      <path d="M26 26h48v48H26z" opacity=".6" />
      <circle cx="50" cy="50" r="6" />
      <circle cx="50" cy="50" r="2" fill="currentColor" stroke="none" />
    </svg>
  );
}

function LogoVariant({ variant, title }: { variant: 1 | 2 | 3; title: string }) {
  return (
    <div className="kit-logo-card">
      <div className="kit-logo-title">{title}</div>
      <div className="kit-logo-lockup">
        <LogoMark variant={variant} />
        <div>
          ARCANA<small>Карты & Судьба</small>
        </div>
      </div>
      <div className="kit-logo-sizes">
        {[64, 40, 26].map((size) => (
          <div key={size}>
            <span style={{ width: size, height: size }}>
              <LogoMark variant={variant} />
            </span>
            <small>{size}</small>
          </div>
        ))}
      </div>
    </div>
  );
}

function DecorCell({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="kit-decor-cell">
      {children}
      <small>{label}</small>
    </div>
  );
}

export function ShowcaseSections() {
  const [modalOpen, setModalOpen] = useState(false);
  return (
    <>
      <UiKitSection
        id="card"
        number="10"
        title="Карта Таро"
        description="Сигнатурный компонент: иллюстрация в двойной рамке с римским номером и подписью, плюс блок названия."
      >
        <div className="tarot-demo">
          <TarotCard roman="II" name="The High Priestess" tilt />
          <div className="kit-card-label">
            <span>II</span>
            <b>Верховная Жрица</b>
            <small>Старший Аркан · Прямое положение</small>
          </div>
        </div>
        <CodeCaption>.tarot · .cardlbl (наведите на карту)</CodeCaption>
      </UiKitSection>
      <UiKitSection
        id="aspects"
        number="11"
        title="Аспект-блоки"
        description="Иконка + заголовок + интерпретация. Поднимаются при наведении."
      >
        <div className="aspect-demo">
          <AspectBlock type="love" title="Любовь">
            Не время форсировать. Карта указывает на скрытые чувства — наблюдайте, не давите.
          </AspectBlock>
          <AspectBlock type="career" title="Карьера">
            Соберите все данные прежде, чем решать. День анализа, а не действия.
          </AspectBlock>
          <AspectBlock type="spirit" title="Духовность">
            Пик восприимчивости. Ведите дневник снов — Жрица открывает бессознательное.
          </AspectBlock>
        </div>
        <CodeCaption>.aspect</CodeCaption>
      </UiKitSection>
      <UiKitSection
        id="panel"
        number="12"
        title="Панель с угловыми орнаментами"
        description="Текстовый блок «Послание дня» в рамке с золотыми уголками."
      >
        <CornerPanel title="Послание дня">
          <p>
            Верховная Жрица — хранительница тайн, сидящая между двумя колоннами: светом и тьмой.
            Сегодня откажитесь от соблазна немедленных ответов.{' '}
            <em>Ваша интуиция сейчас острее разума.</em>
          </p>
        </CornerPanel>
        <CodeCaption>.panel · .corner</CodeCaption>
      </UiKitSection>
      <UiKitSection
        id="tooltips"
        number="13"
        title="Тултипы & Подсказки"
        description="Всплывающие подсказки на наведение: короткий тултип со стрелкой, инфо-иконка и развёрнутый поповер."
      >
        <DemoRow>
          <Tooltip label="Одна карта на сегодня">
            <Button variant="outline">Вытянуть карту</Button>
          </Tooltip>
          <Tooltip below label="Ключевое значение аркана">
            <Chip>Интуиция</Chip>
          </Tooltip>
          <Tooltip label="Мария А. · Скорпион">
            <Avatar initials="МА" />
          </Tooltip>
        </DemoRow>
        <CodeCaption>.tip › .tip-bubble ( .below для нижнего )</CodeCaption>
        <div className="kit-ornament-gap" />
        <p className="kit-inline-note">
          Карта выпала в положении{' '}
          <Tooltip label="Перевёрнутая карта — ослабленное или обратное значение">
            <span className="kit-dashed">реверс</span>
          </Tooltip>{' '}
          — прислушайтесь к тени смысла.
        </p>
        <CodeCaption>.dashed › .tip-bubble (инлайн-подсказка)</CodeCaption>
        <div className="kit-ornament-gap" />
        <DemoRow>
          <span className="kit-popover-title">Верховная Жрица</span>
          <Tooltip
            below
            label={
              <>
                <b>Аркан II</b>
                <br />
                Символ интуиции и тайного знания. Приглашает замолчать и прислушаться к внутреннему
                голосу, а не форсировать события.
              </>
            }
          >
            <button type="button" className="info-dot">
              ?
            </button>
          </Tooltip>
        </DemoRow>
        <CodeCaption>.info-dot › .tip-bubble.pop-bubble (поповер)</CodeCaption>
      </UiKitSection>
      <UiKitSection
        id="modal"
        number="14"
        title="Модальное окно"
        description="Кнопка открывает окно по центру: ширина 50%, минимальная высота 50%, заголовок с крестиком в шапке и кнопка внизу."
      >
        <Button onClick={() => setModalOpen(true)}>Открыть расклад</Button>
        <CodeCaption>.modal-overlay › .modal (.modal-head · .modal-body · .modal-foot)</CodeCaption>
        <Modal
          open={modalOpen}
          onClose={() => setModalOpen(false)}
          title="Полный расклад дня"
          panelClassName="ui-kit-modal"
          footer={
            <Button variant="outline" onClick={() => setModalOpen(false)}>
              Принять послание
            </Button>
          }
        >
          <div className="kit-modal-card">
            <span>II</span>
            <b>Верховная Жрица</b>
            <small>Старший Аркан · Прямое положение</small>
          </div>
          <p className="modal-sample">
            Вселенная приглашает вас замолчать и прислушаться. Ответы, которые вы ищете, живут не
            снаружи — они ждут в глубине вашего существа.
          </p>
          <DemoRow>
            <Chip>Интуиция</Chip>
            <Chip>Тайное знание</Chip>
            <Chip>Цикличность</Chip>
          </DemoRow>
        </Modal>
      </UiKitSection>
      <UiKitSection
        id="icons"
        number="15"
        title="Иконки & Декор"
        description="Единая коллекция линейных иконок в мистическом ключе — используемые в библиотеке и дополнительные по теме таро и астрологии. Все на currentColor, наследуют золото и тему."
      >
        <div className="icon-grid">
          {iconNames.map((name, index) => (
            <div className="kit-icon-cell" key={name}>
              <MysticGlyph index={index} />
              <small>{name}</small>
            </div>
          ))}
        </div>
        <CodeCaption>.icon-grid › .icon-cell (currentColor, наследует .ico color)</CodeCaption>
        <p className="kit-decor-description">
          Декоративные SVG-элементы фона из карты дня — созвездия, лунный серп, зодиакальное колесо,
          фазы луны, орнаменты.
        </p>
        <div className="kit-decor-grid">
          <DecorCell label=".bg-constellation">
            <Constellation />
          </DecorCell>
          <DecorCell label=".bg-crescent">
            <CelestialCrossIcon />
          </DecorCell>
          <DecorCell label=".zodiac-wheel">
            <ZodiacWheel />
          </DecorCell>
          <DecorCell label="Фазы луны">
            <MoonPhases />
          </DecorCell>
          <DecorCell label=".card-flank">
            <SphereIcon type="spirit" />
          </DecorCell>
          <DecorCell label=".ornament (звезда)">
            <Ornament />
          </DecorCell>
          <DecorCell label=".corner-ornament">
            <StarIcon />
          </DecorCell>
          <DecorCell label=".footer-divider">
            <MoonPhases />
          </DecorCell>
        </div>
        <CodeCaption>.deco-grid › .deco-cell</CodeCaption>
      </UiKitSection>
      <UiKitSection
        id="logo"
        number="16"
        title="Логотип"
        description="Три варианта итогового знака ARCANA, каждый в трёх размерах: полный локап (знак + вордмарк), затем чистый знак — крупный, средний и мелкий (для аватара / фавикона)."
      >
        <div className="logo-grid">
          <LogoVariant variant={1} title="Вариант I · Всевидящее око" />
          <LogoVariant variant={2} title="Вариант II · Луна и колонны" />
          <LogoVariant variant={3} title="Вариант III · Сигил-звезда" />
        </div>
        <CodeCaption>
          .logo-grid › .logo-card · .logo-lockup · .logo-sizes (.logo-lg / .logo-md / .logo-sm)
        </CodeCaption>
      </UiKitSection>
    </>
  );
}
