import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { BackendReadingCardDto, CardOfDayInterpretationDto, Language } from '@/api/types';
import { getDailyCardPresentation } from '@/features/daily-card/presentation';
import { ChevronIcon } from '@/shared/icons';
import { AspectBlock, CornerPanel } from '@/shared/ui/Content';
import { Chip } from '@/shared/ui/Chip';
import { Ornament } from '@/shared/ui/Ornament';

export function Interpretation({
  backendCard,
  interpretation,
  interpretationLoading = false,
  interpretationError = '',
  language,
  revealed = true,
}: {
  backendCard?: BackendReadingCardDto;
  interpretation?: CardOfDayInterpretationDto;
  interpretationLoading?: boolean;
  interpretationError?: string;
  language: Language;
  revealed?: boolean;
}) {
  const { t } = useTranslation('home');
  const [expanded, setExpanded] = useState(false);
  const card = backendCard ? getDailyCardPresentation(backendCard, language) : null;
  const canExpand = revealed && Boolean(interpretation) && !interpretationLoading && !interpretationError;
  const brief = interpretationLoading
    ? t('loading')
    : interpretationError || interpretation?.general || t('loading');

  useEffect(() => {
    if (!revealed || !canExpand) setExpanded(false);
  }, [canExpand, revealed]);

  return (
    <>
      <section className="brief">
        <div className="keywords">{card?.tags.map((tag) => <Chip key={tag}>{tag}</Chip>)}</div>
        <p>{brief}</p>
      </section>
      {revealed && interpretationLoading ? (
        <div className="oracle-thinking" role="status" aria-live="polite">
          <div className="oracle-loader" aria-hidden="true">
            <span />
            <span />
            <span />
          </div>
          <span>{t('oracleThinking')}</span>
        </div>
      ) : null}
      <Ornament variant="cross" className="interpretation-ornament" />
      <div className="expand-wrap">
        <button
          type="button"
          className={`expand-button ${expanded ? 'open' : ''}`}
          onClick={() => setExpanded((value) => !value)}
          aria-expanded={canExpand && expanded}
          disabled={!canExpand}
          tabIndex={canExpand ? 0 : -1}
        >
          <span>{t(expanded ? 'collapse' : 'expand')}</span>
          <ChevronIcon />
        </button>
        {canExpand && expanded && interpretation && (
          <div className="spheres-reveal">
            <div className="spheres">
              <AspectBlock type="career" title={t('interpretation.work')}>
                {interpretation.work}
              </AspectBlock>
              <AspectBlock type="love" title={t('interpretation.family')}>
                {interpretation.family}
              </AspectBlock>
              <AspectBlock type="spirit" title={t('interpretation.advice')}>
                {interpretation.advice}
              </AspectBlock>
            </div>
            <CornerPanel title={t('messageTitle')}>
              <p>{interpretation.general}</p>
            </CornerPanel>
          </div>
        )}
      </div>
    </>
  );
}
