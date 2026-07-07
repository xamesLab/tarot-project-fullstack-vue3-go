import type { BackendReadingCardDto, Language } from '@/api/types';
import { resolveApiAssetUrl } from '@/api/client';
import { getDailyCardPresentation } from '@/features/daily-card/presentation';
import { ZodiacWheel } from '@/shared/icons';
import { TarotCard } from '@/shared/ui/TarotCard';

export function CardOfDayStage({ backendCard, language, hint, onRevealChange }: { backendCard?: BackendReadingCardDto; language: Language; hint: string; onRevealChange?: (revealed: boolean) => void }) {
  const card = backendCard ? getDailyCardPresentation(backendCard, language) : null;
  return <><ZodiacWheel className="zodiac decor-icon"/><div className="side-rune left"><i/><span>☽ ☉ ♄</span><i/></div><div className="side-rune right"><i/><span>♀ ♃ ♂</span><i/></div><div className="card-c"><div className="card-glow"/><div className="echo left"/><div className="echo right"/><TarotCard key={backendCard?.key ?? 'loading'} roman={card?.marker ?? '···'} name={card?.name ?? '···'} faceImageUrl={resolveApiAssetUrl(backendCard?.image_url)} faceImageAlt={card?.name} backImageUrl={resolveApiAssetUrl(backendCard?.back_image_url)} flippable flipLabel={hint} onFlipChange={onRevealChange}/><div className="card-hint">{hint}</div></div></>;
}
