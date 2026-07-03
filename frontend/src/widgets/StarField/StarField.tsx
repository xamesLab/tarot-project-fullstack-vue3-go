import { useMemo, type CSSProperties } from 'react';
import { Constellation } from '@/shared/icons';

interface StarStyle extends CSSProperties { '--dur': string; '--delay': string; }

export function BackgroundDecor() {
  const stars = useMemo(() => Array.from({ length: 70 }, (_, index) => { const size = Math.random() * 2 + .5; return { key: index, style: { width: size, height: size, top: `${Math.random() * 100}%`, left: `${Math.random() * 100}%`, '--dur': `${3 + Math.random() * 5}s`, '--delay': `${Math.random() * 4}s`, opacity: .1 + Math.random() * .3 } as StarStyle }; }), []);
  return <><div className="stars" aria-hidden="true">{stars.map(({ key, style }) => <i key={key} className="star" style={style}/>)}</div><div className="orb orb1"/><div className="orb orb2"/><div className="orb orb3"/><Constellation className="bg-constellation bgc1"/><Constellation variant={2} className="bg-constellation bgc2"/></>;
}
