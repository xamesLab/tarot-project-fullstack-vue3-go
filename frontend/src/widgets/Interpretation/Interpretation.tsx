import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { BackendReadingCardDto, DailyReadingDto, Language } from '@/api/types';
import { getDailyCardPresentation } from '@/features/daily-card/presentation';
import { ChevronIcon } from '@/shared/icons';
import { AspectBlock, CornerPanel } from '@/shared/ui/Content';
import { Chip } from '@/shared/ui/Chip';
import { Ornament } from '@/shared/ui/Ornament';

export function Interpretation({
  reading,
  backendCard,
  language,
  revealed = true,
}: {
  reading?: DailyReadingDto;
  backendCard?: BackendReadingCardDto;
  language: Language;
  revealed?: boolean;
}) {
  const { t } = useTranslation('home');
  const [expanded, setExpanded] = useState(false);
  const card = backendCard ? getDailyCardPresentation(backendCard, language) : null;
  useEffect(() => {
    if (!revealed) setExpanded(false);
  }, [revealed]);
  return (
    <>
      <section className="brief">
        <div className="keywords">{card?.tags.map((tag) => <Chip key={tag}>{tag}</Chip>)}</div>
        <p>{card?.interpretation || t('loading')}</p>
      </section>
      <Ornament variant="cross" className="interpretation-ornament" />
      <div className="expand-wrap">
        <button
          type="button"
          className={`expand-button ${expanded ? 'open' : ''}`}
          onClick={() => setExpanded((value) => !value)}
          aria-expanded={revealed && expanded}
          disabled={!revealed}
          tabIndex={revealed ? 0 : -1}
        >
          <span>{t(expanded ? 'collapse' : 'expand')}</span>
          <ChevronIcon />
        </button>
        {revealed && expanded && (
          <div className="spheres-reveal">
            <div className="spheres">
              {reading?.spheres.map((sphere) => (
                <AspectBlock key={sphere.id} type={sphere.id} title={sphere.title}>
                  {sphere.text}
                </AspectBlock>
              ))}
            </div>
            {reading && (
              <CornerPanel title={t('messageTitle')}>
                <p>
                  {reading.message.paragraph} {reading.message.closing}{' '}
                  <em>{reading.message.emphasis}</em>
                </p>
              </CornerPanel>
            )}
          </div>
        )}
      </div>
    </>
  );
}
