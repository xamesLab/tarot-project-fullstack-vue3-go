import { useEffect, useRef } from 'react';

export function useTilt<T extends HTMLElement>(enabled: boolean) {
  const ref = useRef<T>(null);
  useEffect(() => {
    const element = ref.current;
    if (!element || !enabled || !matchMedia('(hover: hover) and (pointer: fine)').matches || matchMedia('(prefers-reduced-motion: reduce)').matches) return;
    const move = (event: MouseEvent) => {
      const bounds = element.getBoundingClientRect();
      const x = (event.clientX - bounds.left) / bounds.width - 0.5;
      const y = (event.clientY - bounds.top) / bounds.height - 0.5;
      element.style.transition = 'transform .18s ease-out, box-shadow .4s';
      element.style.transform = `rotateY(${x * 10}deg) rotateX(${-y * 10}deg) scale(1.02)`;
    };
    const leave = () => {
      element.style.transition = 'transform .6s cubic-bezier(.23,1,.32,1)';
      element.style.transform = 'rotateY(0) rotateX(0) scale(1)';
    };
    element.addEventListener('mousemove', move);
    element.addEventListener('mouseleave', leave);
    return () => { element.removeEventListener('mousemove', move); element.removeEventListener('mouseleave', leave); };
  }, [enabled]);
  return ref;
}
