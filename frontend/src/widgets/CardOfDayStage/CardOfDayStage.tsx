import type { TarotCardDto } from '@/api/types';
import { ZodiacWheel } from '@/shared/icons';
import { TarotCard } from '@/shared/ui/TarotCard';

export function CardOfDayStage({ card, hint }: { card?: TarotCardDto; hint: string }) {
  return <><ZodiacWheel className="zodiac decor-icon"/><div className="side-rune left"><i/><span>☽ ☉ ♄</span><i/></div><div className="side-rune right"><i/><span>♀ ♃ ♂</span><i/></div><div className="card-c"><div className="card-glow"/><div className="echo left"/><div className="echo right"/><TarotCard roman={card?.roman ?? 'II'} name={card?.illustrationLabel ?? '···'} flippable flipLabel={hint}/><div className="card-hint">{hint}</div></div><div className="card-label"><span className="card-number">{card?.roman ?? 'II'}</span><span className="card-title">{card?.name ?? '···'}</span><span className="card-subtitle">{card?.subtitle ?? '···'}</span></div></>;
}
