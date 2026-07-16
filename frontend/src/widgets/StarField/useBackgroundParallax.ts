import { useEffect, useRef } from 'react';

type ParallaxVars = {
  x: number;
  y: number;
  scroll: number;
};

function supportsPointerParallax() {
  return window.matchMedia('(hover: hover) and (pointer: fine)').matches;
}

function shouldReduceMotion() {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches;
}

export function useBackgroundParallax() {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const element = ref.current;
    if (!element || shouldReduceMotion()) return;

    const target: ParallaxVars = { x: 0, y: 0, scroll: window.scrollY || 0 };
    const current: ParallaxVars = { ...target };
    const pointerEnabled = supportsPointerParallax();
    let frame = 0;

    const writeStyles = () => {
      const scroll = Math.min(current.scroll, 1600);
      const scrollScale = pointerEnabled
        ? { far: -0.012, mid: -0.024, near: -0.034 }
        : { far: -0.006, mid: -0.011, near: -0.016 };
      element.style.setProperty('--parallax-x', current.x.toFixed(3));
      element.style.setProperty('--parallax-y', current.y.toFixed(3));
      element.style.setProperty('--parallax-scroll', scroll.toFixed(1));
      element.style.setProperty('--parallax-far-x', `${(-current.x * 12).toFixed(2)}px`);
      element.style.setProperty('--parallax-far-y', `${(-current.y * 8 + scroll * scrollScale.far).toFixed(2)}px`);
      element.style.setProperty('--parallax-mid-x', `${(current.x * 24).toFixed(2)}px`);
      element.style.setProperty('--parallax-mid-y', `${(current.y * 16 + scroll * scrollScale.mid).toFixed(2)}px`);
      element.style.setProperty('--parallax-near-x', `${(-current.x * 42).toFixed(2)}px`);
      element.style.setProperty('--parallax-near-y', `${(-current.y * 28 + scroll * scrollScale.near).toFixed(2)}px`);
    };

    const isSettled = () => (
      Math.abs(target.x - current.x) < .001 &&
      Math.abs(target.y - current.y) < .001 &&
      Math.abs(target.scroll - current.scroll) < .2
    );

    const tick = () => {
      current.x += (target.x - current.x) * .1;
      current.y += (target.y - current.y) * .1;
      current.scroll += (target.scroll - current.scroll) * .08;

      if (isSettled()) {
        current.x = target.x;
        current.y = target.y;
        current.scroll = target.scroll;
      }

      writeStyles();
      frame = isSettled() ? 0 : window.requestAnimationFrame(tick);
    };

    const schedule = () => {
      if (!frame) frame = window.requestAnimationFrame(tick);
    };

    const onScroll = () => {
      target.scroll = window.scrollY || 0;
      schedule();
    };

    const onPointerMove = (event: PointerEvent) => {
      target.x = event.clientX / window.innerWidth - 0.5;
      target.y = event.clientY / window.innerHeight - 0.5;
      schedule();
    };

    window.addEventListener('scroll', onScroll, { passive: true });
    if (pointerEnabled) window.addEventListener('pointermove', onPointerMove, { passive: true });
    schedule();

    return () => {
      window.removeEventListener('scroll', onScroll);
      if (pointerEnabled) window.removeEventListener('pointermove', onPointerMove);
      if (frame) window.cancelAnimationFrame(frame);
    };
  }, []);

  return ref;
}
