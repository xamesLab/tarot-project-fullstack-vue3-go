import type { Config } from 'tailwindcss';

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        gold: 'var(--gold)',
        'gold-light': 'var(--gold-light)',
        'gold-dim': 'var(--gold-dim)',
        ink: 'var(--ink)',
        'ink-mid': 'var(--ink-mid)',
        'ink-soft': 'var(--ink-soft)',
        mist: 'var(--mist)',
        'mist-dim': 'var(--mist-dim)',
        cream: 'var(--cream)',
        line: 'var(--line)',
        'line-soft': 'var(--line-soft)',
      },
      fontFamily: {
        display: ['Cinzel', 'Cormorant Garamond', 'serif'],
        serif: ['Cormorant Garamond', 'serif'],
        ui: ['Raleway', 'sans-serif'],
      },
    },
  },
  plugins: [],
} satisfies Config;
