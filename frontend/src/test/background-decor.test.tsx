import { render } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { BackgroundDecor } from '@/widgets/StarField/StarField';

const defaultMatchMedia = window.matchMedia;

function mockMatchMedia(reducedMotion = false) {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: (query: string) => ({
      matches: query.includes('prefers-reduced-motion') ? reducedMotion : query.includes('hover: hover'),
      media: query,
      onchange: null,
      addEventListener: () => undefined,
      removeEventListener: () => undefined,
      addListener: () => undefined,
      removeListener: () => undefined,
      dispatchEvent: () => false,
    }),
  });
}

afterEach(() => {
  Object.defineProperty(window, 'matchMedia', { writable: true, value: defaultMatchMedia });
  vi.restoreAllMocks();
});

describe('BackgroundDecor', () => {
  it('renders far, middle, and near parallax layers', () => {
    mockMatchMedia();
    const { container } = render(<BackgroundDecor />);

    expect(container.querySelector('.background-parallax')).toBeInTheDocument();
    expect(container.querySelectorAll('.parallax-layer')).toHaveLength(3);
    expect(container.querySelector('.parallax-layer-far .stars-far')).toBeInTheDocument();
    expect(container.querySelector('.parallax-layer-mid .stars-mid')).toBeInTheDocument();
    expect(container.querySelector('.parallax-layer-near .stars-near')).toBeInTheDocument();
    expect(container.querySelector('.parallax-layer-near .orb1')).toBeInTheDocument();
    expect(container.querySelectorAll('.star')).toHaveLength(168);
  });

  it('does not write parallax variables when reduced motion is preferred', () => {
    mockMatchMedia(true);
    const { container } = render(<BackgroundDecor />);
    const root = container.querySelector<HTMLElement>('.background-parallax')!;

    window.dispatchEvent(new Event('scroll'));
    window.dispatchEvent(new Event('pointermove'));

    expect(root.style.getPropertyValue('--parallax-far-x')).toBe('');
    expect(root.style.getPropertyValue('--parallax-mid-y')).toBe('');
    expect(root.style.getPropertyValue('--parallax-near-y')).toBe('');
  });
});
