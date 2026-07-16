import { useMemo, type CSSProperties } from 'react';
import { Constellation } from '@/shared/icons';
import { useBackgroundParallax } from './useBackgroundParallax';

interface StarStyle extends CSSProperties { '--dur': string; '--delay': string; }
interface StarConfig {
  count: number;
  prefix: string;
  minSize: number;
  maxSize: number;
  minOpacity: number;
  maxOpacity: number;
  minDuration: number;
  maxDuration: number;
}

function createStars({ count, prefix, minSize, maxSize, minOpacity, maxOpacity, minDuration, maxDuration }: StarConfig) {
  return Array.from({ length: count }, (_, index) => {
    const size = minSize + Math.random() * (maxSize - minSize);
    return {
      key: `${prefix}-${index}`,
      style: {
        width: size,
        height: size,
        top: `${Math.random() * 100}%`,
        left: `${Math.random() * 100}%`,
        '--dur': `${minDuration + Math.random() * (maxDuration - minDuration)}s`,
        '--delay': `${Math.random() * 4}s`,
        opacity: minOpacity + Math.random() * (maxOpacity - minOpacity),
      } as StarStyle,
    };
  });
}

export function BackgroundDecor() {
  const parallaxRef = useBackgroundParallax();
  const farStars = useMemo(() => createStars({ count: 92, prefix: 'far', minSize: .45, maxSize: 1.8, minOpacity: .1, maxOpacity: .32, minDuration: 4, maxDuration: 9 }), []);
  const midStars = useMemo(() => createStars({ count: 54, prefix: 'mid', minSize: .7, maxSize: 2.4, minOpacity: .14, maxOpacity: .42, minDuration: 3.5, maxDuration: 8 }), []);
  const nearStars = useMemo(() => createStars({ count: 22, prefix: 'near', minSize: 1.1, maxSize: 3.2, minOpacity: .18, maxOpacity: .5, minDuration: 3, maxDuration: 6.5 }), []);

  return (
    <div className="background-parallax" ref={parallaxRef} aria-hidden="true">
      <div className="parallax-layer parallax-layer-far">
        <div className="stars stars-far">{farStars.map(({ key, style }) => <i key={key} className="star" style={style}/>)}</div>
      </div>
      <div className="parallax-layer parallax-layer-mid">
        <div className="stars stars-mid">{midStars.map(({ key, style }) => <i key={key} className="star" style={style}/>)}</div>
        <Constellation className="bg-constellation bgc1"/>
        <Constellation variant={2} className="bg-constellation bgc2"/>
      </div>
      <div className="parallax-layer parallax-layer-near">
        <div className="stars stars-near">{nearStars.map(({ key, style }) => <i key={key} className="star star-near" style={style}/>)}</div>
        <div className="orb orb1"/>
        <div className="orb orb2"/>
        <div className="orb orb3"/>
      </div>
    </div>
  );
}
