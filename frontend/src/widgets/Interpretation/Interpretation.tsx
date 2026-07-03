import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { DailyReadingDto } from '@/api/types';
import { ChevronIcon } from '@/shared/icons';
import { AspectBlock, CornerPanel } from '@/shared/ui/Content';
import { Chip } from '@/shared/ui/Chip';
import { Ornament } from '@/shared/ui/Ornament';

export function Interpretation({ reading }: { reading?: DailyReadingDto }) {
  const { t } = useTranslation('home');
  const [expanded, setExpanded] = useState(false);
  return <><section className="brief"><div className="keywords">{reading?.keywords.map((keyword) => <Chip key={keyword}>{keyword}</Chip>)}</div><p>{reading?.brief ?? t('loading')}</p></section><Ornament variant="cross" className="interpretation-ornament"/><div className="expand-wrap"><button type="button" className={`expand-button ${expanded ? 'open' : ''}`} onClick={() => setExpanded((value) => !value)} aria-expanded={expanded}><span>{t(expanded ? 'collapse' : 'expand')}</span><ChevronIcon/></button>{expanded && <div className="spheres-reveal"><div className="spheres">{reading?.spheres.map((sphere) => <AspectBlock key={sphere.id} type={sphere.id} title={sphere.title}>{sphere.text}</AspectBlock>)}</div>{reading && <CornerPanel title={t('messageTitle')}><p>{reading.message.paragraph}</p><p>{reading.message.closing} <em>{reading.message.emphasis}</em></p></CornerPanel>}</div>}</div></>;
}
