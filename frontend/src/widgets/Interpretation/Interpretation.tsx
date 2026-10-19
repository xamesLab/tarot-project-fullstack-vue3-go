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
    setExpanded(revealed && canExpand);
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
            <CornerPanel title={t('messageTitle')}>
              <p>{interpretation.message}</p>
            </CornerPanel>
            <div className="spheres">
              <AspectBlock type="love" title={t('interpretation.relationships')}>
                {interpretation.relationships}
              </AspectBlock>
              <AspectBlock type="finance" title={t('interpretation.finances')}>
                {interpretation.finances}
              </AspectBlock>
              <AspectBlock type="spirit" title={t('interpretation.emotionalState')}>
                {interpretation.emotional_state}
              </AspectBlock>
              <AspectBlock type="career" title={t('interpretation.careerBusiness')}>
                {interpretation.career_business}
              </AspectBlock>
            </div>
          </div>
        )}
      </div>
    </>
  );
}
